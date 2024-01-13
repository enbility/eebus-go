package ship

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/ship/model"
)

// handle incoming SHIP messages and coordinate Handshake States
func (c *ShipConnectionImpl) handleShipMessage(timeout bool, message []byte) {
	if len(message) > 2 {
		var closeMsg model.ConnectionClose
		err := c.processShipJsonMessage(message, &closeMsg)
		if err == nil && closeMsg.ConnectionClose.Phase != "" {
			switch closeMsg.ConnectionClose.Phase {
			case model.ConnectionClosePhaseTypeAnnounce:
				// SHIP 13.4.7: Connection Termination Confirm
				closeMessage := model.ConnectionClose{
					ConnectionClose: model.ConnectionCloseType{
						Phase: model.ConnectionClosePhaseTypeConfirm,
					},
				}

				_ = c.sendShipModel(model.MsgTypeEnd, closeMessage)

				// wait a bit to let it send
				<-time.After(500 * time.Millisecond)

				//
				c.dataHandler.CloseDataConnection(4001, "close")
				c.serviceDataProvider.HandleConnectionClosed(c, c.getState() == SmeStateComplete)
			case model.ConnectionClosePhaseTypeConfirm:
				// we got a confirmation so close this connection
				c.dataHandler.CloseDataConnection(4001, "close")
				c.serviceDataProvider.HandleConnectionClosed(c, c.getState() == SmeStateComplete)
			}

			return
		}
	}

	c.handleState(timeout, message)
}

// set a new handshake state and handle timers if needed
func (c *ShipConnectionImpl) setState(newState ShipMessageExchangeState, err error) {
	c.mux.Lock()

	oldState := c.smeState

	c.smeState = newState

	switch newState {
	case SmeHelloStateReadyInit:
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, tHelloInit)
	case SmeHelloStatePendingInit:
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, tHelloInit)
	case SmeHelloStateOk:
		c.stopHandshakeTimer()
	case SmeHelloStateAbort, SmeHelloStateAbortDone, SmeHelloStateRemoteAbortDone, SmeHelloStateRejected:
		c.stopHandshakeTimer()
	case SmeProtHStateClientListenChoice:
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, cmiTimeout)
	case SmeProtHStateClientOk:
		c.stopHandshakeTimer()
	}

	c.smeError = nil
	if oldState != newState {
		c.smeError = err
		state := ShipState{
			State: newState,
			Error: err,
		}
		c.mux.Unlock()
		c.serviceDataProvider.HandleShipHandshakeStateUpdate(c.remoteSKI, state)
		return
	}
	c.mux.Unlock()
}

func (c *ShipConnectionImpl) getState() ShipMessageExchangeState {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.smeState
}

// handle handshake state transitions
func (c *ShipConnectionImpl) handleState(timeout bool, message []byte) {
	switch c.getState() {
	case SmeStateError:
		logging.Log().Debug(c.RemoteSKI, "connection is in error state")
		return

	// cmiStateInit
	case CmiStateInitStart:
		// triggered without a message received
		c.handshakeInit_cmiStateInitStart()

	case CmiStateClientWait:
		if timeout {
			c.endHandshakeWithError(errors.New("ship client handshake timeout"))
			return
		}

		c.handshakeInit_cmiStateClientWait(message)

	case CmiStateServerWait:
		if timeout {
			c.endHandshakeWithError(errors.New("ship server handshake timeout"))
			return
		}
		c.handshakeInit_cmiStateServerWait(message)

	// smeHello

	case SmeHelloState:
		// check if the service is already trusted or the role is client,
		// which means it was initiated from this service usually by triggering the
		// pairing service
		// go to substate ready if so, otherwise to substate pending

		if c.serviceDataProvider.IsRemoteServiceForSKIPaired(c.remoteSKI) || c.role == ShipRoleClient {
			c.setState(SmeHelloStateReadyInit, nil)
		} else {
			c.setState(SmeHelloStatePendingInit, nil)
		}
		c.handleState(timeout, message)

	case SmeHelloStateReadyInit:
		c.handshakeHello_Init()

	case SmeHelloStateReadyListen:
		c.handshakeHello_ReadyListen(timeout, message)

	case SmeHelloStatePendingInit:
		c.handshakeHello_PendingInit()

	case SmeHelloStatePendingListen:
		c.handshakeHello_PendingListen(timeout, message)

	case SmeHelloStateOk:
		c.handshakeProtocol_Init()

	case SmeHelloStateAbort:
		c.handshakeHello_Abort()

	case SmeHelloStateAbortDone, SmeHelloStateRemoteAbortDone:
		go func() {
			time.Sleep(time.Second)
			c.CloseConnection(false, 4452, "Node rejected by application")
		}()

	// smeProtocol

	case SmeProtHStateServerListenProposal:
		c.handshakeProtocol_smeProtHStateServerListenProposal(message)

	case SmeProtHStateServerListenConfirm:
		c.handshakeProtocol_smeProtHStateServerListenConfirm(message)

	case SmeProtHStateClientListenChoice:
		c.stopHandshakeTimer()
		c.handshakeProtocol_smeProtHStateClientListenChoice(message)

	case SmeProtHStateClientOk:
		c.setAndHandleState(SmePinStateCheckInit)

	case SmeProtHStateServerOk:
		c.setAndHandleState(SmePinStateCheckInit)

	// smePinState

	case SmePinStateCheckInit:
		c.handshakePin_Init()

	case SmePinStateCheckListen:
		c.handshakePin_smePinStateCheckListen(message)

	case SmePinStateCheckOk:
		c.handshakeAccessMethods_Init()

	// smeAccessMethods

	case SmeAccessMethodsRequest:
		c.handshakeAccessMethods_Request(message)
	}
}

// set a state and trigger handling it
func (c *ShipConnectionImpl) setAndHandleState(state ShipMessageExchangeState) {
	c.setState(state, nil)
	c.handleState(false, nil)
}

// SHIP handshake is approved, now set the new state and the SPINE read handler
func (c *ShipConnectionImpl) approveHandshake() {
	// Report to SPINE local device about this remote device connection
	c.spineDataProcessing = c.serviceDataProvider.SetupRemoteDevice(c.remoteSKI, c)
	c.stopHandshakeTimer()
	c.setState(SmeStateComplete, nil)
}

// end the handshake process because of an error
func (c *ShipConnectionImpl) endHandshakeWithError(err error) {
	c.stopHandshakeTimer()

	c.setState(SmeStateError, err)

	logging.Log().Debug(c.RemoteSKI, "SHIP handshake error:", err)

	c.CloseConnection(true, 0, err.Error())

	state := ShipState{
		State: SmeStateError,
		Error: err,
	}
	c.serviceDataProvider.HandleShipHandshakeStateUpdate(c.remoteSKI, state)
}

// set the handshake timer to a new duration and start the channel
func (c *ShipConnectionImpl) setHandshakeTimer(timerType timeoutTimerType, duration time.Duration) {
	c.stopHandshakeTimer()

	c.setHandshakeTimerRunning(true)
	c.setHandshakeTimerType(timerType)

	go func() {
		select {
		case <-c.handshakeTimerStopChan:
			return
		case <-time.After(duration):
			c.setHandshakeTimerRunning(false)
			c.handleState(true, nil)
			return
		}
	}()
}

// stop the handshake timer and close the channel
func (c *ShipConnectionImpl) stopHandshakeTimer() {
	if !c.getHandshakeTimerRunnging() {
		return
	}

	select {
	case c.handshakeTimerStopChan <- struct{}{}:
	default:
	}
	c.setHandshakeTimerRunning(false)
}

func (c *ShipConnectionImpl) setHandshakeTimerRunning(value bool) {
	c.handshakeTimerMux.Lock()
	defer c.handshakeTimerMux.Unlock()

	c.handshakeTimerRunning = value
}

func (c *ShipConnectionImpl) getHandshakeTimerRunnging() bool {
	c.handshakeTimerMux.Lock()
	defer c.handshakeTimerMux.Unlock()

	return c.handshakeTimerRunning
}

func (c *ShipConnectionImpl) setHandshakeTimerType(timerType timeoutTimerType) {
	c.handshakeTimerMux.Lock()
	defer c.handshakeTimerMux.Unlock()

	c.handshakeTimerType = timerType
}

func (c *ShipConnectionImpl) getHandshakeTimerType() timeoutTimerType {
	c.handshakeTimerMux.Lock()
	defer c.handshakeTimerMux.Unlock()

	return c.handshakeTimerType
}
