package ship

import (
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/ship/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

// Handshake Hello covers the states smeHello...

// SME_HELLO_STATE_READY_INIT
func (c *ShipConnection) handshakeHello_Init() {
	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypeReady, tHelloInit, false); err != nil {
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return
	}
	fmt.Println("sent Hello", c.role)

	c.setState(smeHelloStateReadyListen)
}

// SME_HELLO_STATE_READY_LISTEN
func (c *ShipConnection) handshakeHello_ReadyListen(message []byte) {
	fmt.Println("received Hello", c.role)
	var helloReturnMsg model.ConnectionHello
	if err := c.processShipJsonMessage(message, &helloReturnMsg); err != nil {
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return
	}

	hello := helloReturnMsg.ConnectionHello

	switch hello.Phase {
	case model.ConnectionHelloPhaseTypeReady:
		// HELLO_OK
		c.setState(smeHelloStateOk)

	case model.ConnectionHelloPhaseTypePending:
		// the phase is still pending an no prolongationRequest is set, ignore the message
		if hello.ProlongationRequest == nil {
			return
		}

		// if we got a prolongation request, accept it
		if *hello.ProlongationRequest {
			if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
				c.endHandshakeWithError(err)
				return
			}
		}

	case model.ConnectionHelloPhaseTypeAborted:
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return

	default:
		// don't accept any other responses
		logging.Log.Errorf("Unexpected connection hello phase: %s", hello.Phase)
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return
	}

	c.handleState(false, nil)
}

// SME_HELLO_ABORT
func (c *ShipConnection) handshakeHello_Abort() {
	c.setState(smeHelloStateAbort)
	c.stopHandshakeTimer()

	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypeAborted, 0, false); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.CloseConnection(false)
}

// SME_HELLO_PENDING_INIT
// TODO: clarify in which scenario and how we need to support getting into this state
func (c *ShipConnection) handshakeHello_PendingInit() {
	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setState(smeHelloStatePendingListen)
}

// SME_HELLO_PENDING_LISTEN
func (c *ShipConnection) handshakeHello_PendingListen(message []byte) {
	var helloReturnMsg model.ConnectionHello
	if err := c.processShipJsonMessage(message, &helloReturnMsg); err != nil {
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return
	}

	hello := helloReturnMsg.ConnectionHello

	switch hello.Phase {
	case model.ConnectionHelloPhaseTypeReady:
		if hello.Waiting == nil {
			c.setState(smeHelloStateAbort)
			c.handleState(false, nil)
			return
		}

		c.stopHandshakeTimer()

		newDuration := time.Duration(*hello.Waiting) * time.Millisecond
		duration := tHelloProlongThrInc
		if newDuration >= duration {
			// the duration has to be reduced
			duration = newDuration - duration

			// check if it is less than T_hello_prolong_min
			if newDuration >= tHelloProlongMin {
				c.setHandshakeTimer(timeoutTimerTypeSendProlongationRequest, duration)
				return
			}
		}

		if newDuration < tHelloProlongMin {
			// I interpret 13.4.4.1.3 Page 64 Line 1550-1553 as this resulting in a timeout state
			// TODO: verify this
			c.setState(smeHelloStateAbort)
			c.handleState(false, nil)
		}

	case model.ConnectionHelloPhaseTypePending:
		if hello.Waiting != nil && hello.ProlongationRequest == nil {
			c.stopHandshakeTimer()

			newDuration := time.Duration(*hello.Waiting) * time.Millisecond
			c.lastReceivedWaitingValue = newDuration
			duration := tHelloProlongThrInc
			if newDuration >= duration {
				// the duration has to be reduced
				duration = newDuration - duration

				// check if it is less than T_hello_prolong_min
				if newDuration >= tHelloProlongMin {
					c.setHandshakeTimer(timeoutTimerTypeSendProlongationRequest, duration)
					return
				}
			}

			if newDuration < tHelloProlongMin {
				// I interpret 13.4.4.1.3 Page 64 Line 1557-1560 as this resulting in a timeout state
				// TODO: verify this
				c.setState(smeHelloStateAbort)
				c.handleState(false, nil)
			}

			return
		}

		if hello.Waiting == nil && hello.ProlongationRequest != nil && *hello.ProlongationRequest {
			// if we got a prolongation request, accept it
			if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
				c.endHandshakeWithError(err)
			}

			return
		}

		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)

	case model.ConnectionHelloPhaseTypeAborted:
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return

	default:
		// don't accept any other responses
		logging.Log.Errorf("Unexpected connection hello phase: %s", hello.Phase)
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return
	}

	c.handleState(false, nil)
}

func (c *ShipConnection) handshakeHello_PendingTimeout() {
	if c.handshakeTimerType != timeoutTimerTypeSendProlongationRequest {
		c.setState(smeHelloStateAbort)
		c.handleState(false, nil)
		return
	}

	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, 0, true); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	if c.lastReceivedWaitingValue == 0 {
		newValue := float64(tHelloInit.Milliseconds()) * 1.1
		c.lastReceivedWaitingValue = time.Duration(newValue)
	}
	c.setHandshakeTimer(timeoutTimerTypeProlongRequestReply, c.lastReceivedWaitingValue)
}

func (c *ShipConnection) handshakeHelloSend(phase model.ConnectionHelloPhaseType, waitingDuration time.Duration, prolongation bool) error {
	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: phase,
		},
	}

	if waitingDuration > 0 {
		helloMsg.ConnectionHello.Waiting = util.Ptr(uint(waitingDuration.Milliseconds()))
	}
	if prolongation {
		helloMsg.ConnectionHello.ProlongationRequest = &prolongation
	}

	if err := c.sendShipModel(model.MsgTypeControl, helloMsg); err != nil {
		return err
	}
	return nil
}
