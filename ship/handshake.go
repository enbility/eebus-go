package ship

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/ship/model"
)

// handle incoming SHIP messages and coordinate Handshake States
func (c *ShipConnection) handleShipMessage(timeout bool, message []byte) {
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
				c.DataHandler.CloseDataConnection(4495, "timeout")
				c.serviceDataProvider.HandleConnectionClosed(c, c.getState() == SmeComplete)
			case model.ConnectionClosePhaseTypeConfirm:
				// we got a confirmation so close this connection
				c.DataHandler.CloseDataConnection(4496, "close")
				c.serviceDataProvider.HandleConnectionClosed(c, c.getState() == SmeComplete)
			}

			return
		}
	}

	c.handleState(timeout, message)
}

// set a new handshake state and handle timers if needed
func (c *ShipConnection) setState(newState ShipMessageExchangeState, err error) {
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
	case SmeHelloStateAbort:
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
		c.serviceDataProvider.HandleShipHandshakeStateUpdate(c.RemoteSKI, state)
		return
	}
	c.mux.Unlock()
}

func (c *ShipConnection) getState() ShipMessageExchangeState {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.smeState
}

// handle handshake state transitions
func (c *ShipConnection) handleState(timeout bool, message []byte) {
	switch c.getState() {
	case SmeError:
		logging.Log.Debug(c.RemoteSKI, "connection is in error state")
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

		if c.serviceDataProvider.IsRemoteServiceForSKIPaired(c.RemoteSKI) || c.role == ShipRoleClient {
			c.setState(SmeHelloStateReadyInit, nil)
		} else {
			c.setState(SmeHelloStatePendingInit, nil)
		}
		c.handleState(timeout, message)

	case SmeHelloStateReadyInit:
		c.handshakeHello_Init()

	case SmeHelloStateReadyListen:
		if timeout {
			c.setState(SmeHelloStateAbort, nil)
			c.handleState(false, nil)
			return
		}

		c.handshakeHello_ReadyListen(message)

	case SmeHelloStatePendingInit:
		c.handshakeHello_PendingInit()

	case SmeHelloStatePendingListen:
		if timeout {
			// The device needs to be in a state for the user to allow trusting the device
			// e.g. either the web UI or by other means
			if !c.serviceDataProvider.AllowWaitingForTrust(c.remoteShipID) {
				c.handshakeHello_PendingTimeout()
				return
			}

			c.handshakeHello_PendingProlongationRequest()
			return
		}

		c.handshakeHello_PendingListen(message)

	case SmeHelloStateOk:
		c.handshakeProtocol_Init()

	case SmeHelloStateAbort:
		c.handshakeHello_Abort()

	// smeProtocol

	case SmeProtHStateServerListenProposal:
		c.handshakeProtocol_smeProtHStateServerListenProposal(message)

	case SmeProtHStateServerListenConfirm:
		c.handshakeProtocol_smeProtHStateServerListenConfirm(message)

	case SmeProtHStateClientListenChoice:
		c.stopHandshakeTimer()
		c.handshakeProtocol_smeProtHStateClientListenChoice(message)

	case SmeProtHStateClientOk:
		c.setState(SmePinStateCheckInit, nil)
		c.handleState(false, nil)

	case SmeProtHStateServerOk:
		c.setState(SmePinStateCheckInit, nil)
		c.handleState(false, nil)

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

// SHIP handshake is approved, now set the new state and the SPINE read handler
func (c *ShipConnection) approveHandshake() {
	// Report to SPINE local device about this remote device connection
	c.spineDataProcessing = c.deviceLocalCon.AddRemoteDevice(c.RemoteSKI, c)
	c.stopHandshakeTimer()
	c.setState(SmeComplete, nil)
}

// end the handshake process because of an error
func (c *ShipConnection) endHandshakeWithError(err error) {
	c.stopHandshakeTimer()

	c.setState(SmeError, err)

	logging.Log.Debug(c.RemoteSKI, "SHIP handshake error:", err)

	c.CloseConnection(true, err.Error())
}

// set the handshake timer to a new duration and start the channel
func (c *ShipConnection) setHandshakeTimer(timerType timeoutTimerType, duration time.Duration) {
	c.stopHandshakeTimer()

	c.setHandshakeTimerRunning(true)
	c.handshakeTimerType = timerType

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
func (c *ShipConnection) stopHandshakeTimer() {
	if !c.getHandshakeTimerRunnging() {
		return
	}

	select {
	case c.handshakeTimerStopChan <- struct{}{}:
	default:
	}
	c.setHandshakeTimerRunning(false)
}

func (c *ShipConnection) setHandshakeTimerRunning(value bool) {
	c.handshakeTimerMux.Lock()
	defer c.handshakeTimerMux.Unlock()

	c.handshakeTimerRunning = value
}

func (c *ShipConnection) getHandshakeTimerRunnging() bool {
	c.handshakeTimerMux.Lock()
	defer c.handshakeTimerMux.Unlock()

	return c.handshakeTimerRunning
}
