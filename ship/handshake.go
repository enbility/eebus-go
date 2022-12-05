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
				c.DataHandler.CloseDataConnection()
				c.serviceDataProvider.HandleConnectionClosed(c, c.smeState == smeComplete)
			case model.ConnectionClosePhaseTypeConfirm:
				// we got a confirmation so close this connection
				c.DataHandler.CloseDataConnection()
				c.serviceDataProvider.HandleConnectionClosed(c, c.smeState == smeComplete)
			}

			return
		}
	}

	c.handleState(timeout, message)
}

// set a new handshake state and handle timers if needed
func (c *ShipConnection) setState(newState shipMessageExchangeState) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.smeState = newState

	switch newState {
	case smeHelloStateReadyInit:
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, tHelloInit)
	case smeHelloStatePendingInit:
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, tHelloInit)
	case smeHelloStateOk:
		c.stopHandshakeTimer()
	case smeHelloStateAbort:
		c.stopHandshakeTimer()
	case smeProtHStateClientListenChoice:
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, cmiTimeout)
	case smeProtHStateClientOk:
		c.stopHandshakeTimer()
	}
}

func (c *ShipConnection) getState() shipMessageExchangeState {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.smeState
}

// handle handshake state transitions
func (c *ShipConnection) handleState(timeout bool, message []byte) {
	switch c.getState() {
	// cmiStateInit
	case cmiStateInitStart:
		// triggered without a message received
		c.handshakeInit_cmiStateInitStart()

	case cmiStateClientWait:
		if timeout {
			logging.Log.Trace("timeout")
			c.endHandshakeWithError(errors.New("ship handshake timeout"))
			return
		}

		c.handshakeInit_cmiStateClientWait(message)

	case cmiStateServerWait:
		if timeout {
			logging.Log.Trace("timeout")
			c.endHandshakeWithError(errors.New("ship handshake timeout"))
			return
		}
		c.handshakeInit_cmiStateServerWait(message)

	// smeHello

	case smeHelloState:
		// go into the 1st  substate right away
		c.setState(smeHelloStateReadyInit)
		c.handleState(timeout, message)

	case smeHelloStateReadyInit:
		c.handshakeHello_Init()

	case smeHelloStateReadyListen:
		if timeout {
			logging.Log.Trace("timeout")
			c.setState(smeHelloStateAbort)
			c.handleState(false, nil)
			return
		}

		c.handshakeHello_ReadyListen(message)

	case smeHelloStatePendingInit:
		c.handshakeHello_PendingInit()

	case smeHelloStatePendingListen:
		if timeout {
			logging.Log.Trace("timeout")
			c.handshakeHello_PendingTimeout()
			return
		}

		c.handshakeHello_PendingListen(message)

	case smeHelloStateOk:
		c.handshakeProtocol_Init()

	case smeHelloStateAbort:
		c.handshakeHello_Abort()

	// smeProtocol

	case smeProtHStateServerListenProposal:
		c.handshakeProtocol_smeProtHStateServerListenProposal(message)

	case smeProtHStateServerListenConfirm:
		c.handshakeProtocol_smeProtHStateServerListenConfirm(message)

	case smeProtHStateClientListenChoice:
		c.stopHandshakeTimer()
		c.handshakeProtocol_smeProtHStateClientListenChoice(message)

	case smeProtHStateClientOk:
		c.setState(smePinStateCheckInit)
		c.handleState(false, nil)

	case smeProtHStateServerOk:
		c.setState(smePinStateCheckInit)
		c.handleState(false, nil)

		// smePinState

	case smePinStateCheckInit:
		c.handshakePin_Init()

	case smePinStateCheckListen:
		c.handshakePin_smePinStateCheckListen(message)

	// smeAccessMethods

	case smeAccessMethodsRequest:
		c.handshakeAccessMethods_Request(message)
	}
}

// SHIP handshake is approved, now set the new state and the SPINE read handler
func (c *ShipConnection) approveHandshake() {
	// Report to SPINE local device about this remote device connection
	c.spineDataProcessing = c.spineLocalDevice.AddRemoteDevice(c.RemoteSKI, c)
	c.stopHandshakeTimer()
	c.setState(smeComplete)
}

// end the handshake process because of an error
func (c *ShipConnection) endHandshakeWithError(err error) {
	c.stopHandshakeTimer()

	c.setState(smeError)

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
