package ship

import (
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/util"
)

// handle incoming SHIP messages and coordinate Handshake States
func (c *ShipConnection) handleShipMessage(timeout bool, message []byte) {
	if len(message) > 2 {
		logging.Log.Trace("Recv:", c.RemoteSKI, string(message[1:]))
	}

	// TODO: first check if this is a close message

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
		c.stopHandshakeTimer()
		c.handshakeProtocol_Init()

	case smeHelloStateAbort:
		c.stopHandshakeTimer()
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
	c.setState(smeComplete)
}

// end the handshake process because of an error
func (c *ShipConnection) endHandshakeWithError(err error) {
	c.stopHandshakeTimer()

	c.setState(smeError)

	if c.handshakeError != nil {
		logging.Log.Error(c.RemoteSKI, "SHIP handshake error:", c.handshakeError)
	} else {
		logging.Log.Error(c.RemoteSKI, "SHIP handshake error unknown")
	}

	c.CloseConnection(false)
}

// set the handshake timer to a new duration and start the channel
func (c *ShipConnection) setHandshakeTimer(timerType timeoutTimerType, duration time.Duration) {
	c.stopHandshakeTimer()

	c.handshakeTimer.Reset(cmiTimeout)
	c.handshakeTimerRunning = true
	c.handshakeTimerType = timerType

	go func() {
		for {
			select {
			case <-c.handshakeTimerStopChan:
				fmt.Println("EXIT 1")
				return
			case <-c.handshakeTimer.C:
				c.handleState(true, nil)
				fmt.Println("EXIT 2")
				return
			}
		}
	}()
}

// stop the handshake timer and close the channel
func (c *ShipConnection) stopHandshakeTimer() {
	if !c.handshakeTimerRunning {
		return
	}

	c.handshakeTimer.Stop()
	if !util.IsChannelClosed(c.handshakeTimerStopChan) {
		close(c.handshakeTimerStopChan)
	}
	c.handshakeTimerRunning = false
}
