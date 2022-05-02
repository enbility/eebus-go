package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DerAndereAndi/eebus-go/ship"
)

const (
	cmiTimeout              = 10 * time.Second // SHIP 4.2
	cmiCloseTimeout         = 100 * time.Millisecond
	tHelloInit              = 60 * time.Second
	tHelloProlongThrInc     = 30 * time.Second
	tHelloProlongWaitingGap = 15 * time.Second
	tHellogProlongMin       = 1 * time.Second
)

type shipMessageExchangeState uint

const (
	// Connection Mode Initialisation (CMI) SHIP 13.4.3
	cmiStatInitStart shipMessageExchangeState = iota
	cmiStateClientSend
	cmiStateClientWait
	cmiStateClientEvaluate
	cmiStateServerWait
	cmiStateServerEvaluate
	// Connection Data Preparation SHIP 13.4.4
	smeHelloState
	smeHelloStateReady
	smeHelloStateReadyInit
	smeHelloStateReadyListen
	smeHelloStateReadyTimeout
	smeHelloStatePending
	smeHelloStatePendingInit
	smeHelloStatePendingListen
	smeHelloStatePendingTimeout
	smeHelloStateOk
	smeHelloStateAbort
	// Connection State Protocol Handhsake SHIP 13.4.4.2
	smeProtHStateServerInit
	smeProtHStateClientInit
	smeProtHStateServerListenProposal
	smeProtHStateClientListenChoice
	smeProtHStateListenConfirm
	smeProtHStateTimeout
	smeProtHStateClientOk
	smeProtHStateServerOk
	// Connection PIN State 13.4.5
	sneOubStateCheckInit
	smePinStateCheckListen
	smePinStateCheckError
	smePinStateCheckBusyInit
	smePinStateCheckBusyWait
	smePinStateCheckOk
	smePinStateAskInit
	smePinStateAskProcess
	smePinStateAskRestricted
	smePinStateAskOk
	// ConnectionAccess Methods Identification 13.4.6

	// Everything is done
	smeComplete
)

// process the ship handshake and return an error if the handshake failed
func (c *ConnectionHandler) shipHandshake(trusted bool) error {
	if err := c.handshakeInit(); err != nil {
		return err
	}

	if err := c.handshakeHello(trusted); err != nil {
		return err
	}

	if err := c.handshakeProtocol(); err != nil {
		return err
	}

	if err := c.handshakePin(); err != nil {
		return err
	}

	if err := c.handshakeAccessMethods(); err != nil {
		return err
	}

	c.smeState = smeComplete

	return nil
}

// process the ship handshake and return an error if the handshake failed
func (c *ConnectionHandler) handshakeInit() error {
	var data []byte
	var msgType byte
	var err error

	c.smeState = cmiStatInitStart

	shipInit := []byte{ship.MsgTypeInit, 0x00}

	if c.role == ShipRoleClient {
		// CMI_STATE_CLIENT_SEND
		c.smeState = cmiStateClientSend
		if err := c.writeWebsocketMessage(shipInit); err != nil {
			return err
		}
		c.smeState = cmiStateClientWait
	} else {
		c.smeState = cmiStateServerWait
	}

	// CMI_STATE_SERVER_WAIT
	// CMI_STATE_CLIENT_WAIT
	data, _, err = c.readNextShipMessage(cmiTimeout)
	if err != nil {
		return err
	}

	if data == nil || c.isConnectionClosed {
		return nil
	}

	if c.role == ShipRoleServer {
		c.smeState = cmiStateServerEvaluate
		if err := c.writeWebsocketMessage(shipInit); err != nil {
			return err
		}
	} else {
		c.smeState = cmiStateClientEvaluate
	}

	// CMI_STATE_SERVER_EVALUATE
	// CMI_STATE_CLIENT_EVALUATE
	msgType, data = c.parseMessage(data, false)

	if msgType != ship.MsgTypeInit {
		return errors.New("Invalid SHIP MessageType, expected 0 and got " + string(msgType))
	}
	if data[0] != byte(0) {
		return errors.New("Invalid SHIP MessageValue, expected 0 and got " + string(data))
	}

	c.smeState = smeHelloState

	return nil
}

func (c *ConnectionHandler) handshakeHelloSend(phase ship.ConnectionHelloPhaseType, waitingDuration time.Duration, prolongation bool) error {
	helloMsg := ship.ConnectionHello{
		ConnectionHello: ship.ConnectionHelloType{
			Phase: phase,
		},
	}

	if waitingDuration > 0 {
		waiting := uint(waitingDuration.Milliseconds())
		helloMsg.ConnectionHello.Waiting = &waiting
	}
	if prolongation {
		helloMsg.ConnectionHello.ProlongationRequest = &prolongation
	}

	if err := c.sendShipModel(ship.MsgTypeControl, helloMsg); err != nil {
		return err
	}
	return nil
}

func (c *ConnectionHandler) handshakeHello(trusted bool) error {
	remoteSmeState := smeHelloState

	// handling local and remote trust states separately
	// unless one aborts or both trust

	if trusted {
		// SME_HELLO_STATE_READY_INIT
		c.smeState = smeHelloStateReadyInit

		if err := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypeReady, tHelloInit, false); err != nil {
			c.smeState = smeHelloStateAbort

			return err
		}

		// SME_HELLO_STATE_READY_LISTEN
		c.smeState = smeHelloStateReadyListen
	} else {
		// SME_HELLO_STATE_PENDING_INIT
		c.smeState = smeHelloStatePendingInit

		if err := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
			c.smeState = smeHelloStateAbort

			return err
		}

		c.connectionDelegate.requestUserTrustForService(c.remoteService)

		// SME_HELLO_STATE_PENDING_LISTEN
		c.smeState = smeHelloStatePendingListen
	}

	currentTimeout := cmiTimeout

	for {
		data, trustState, err := c.readNextShipMessage(currentTimeout)
		// an error is returned on a timeout
		if err != nil {
			if trusted {
				c.smeState = smeHelloStateReadyTimeout

				if sendErr := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypeAborted, 0, false); sendErr != nil {
					return sendErr
				}

				return err
			} else {
				// Timeout and we didn't receive a user input, so ask for prolongation
				c.smeState = smeHelloStatePendingTimeout

				if sendErr := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypePending, 0, true); sendErr != nil {
					return sendErr
				}

				c.smeState = smeHelloStatePendingListen

				currentTimeout = tHelloProlongThrInc

				continue
			}
		}
		// if data is nil, we got a local trust update
		if data == nil {
			if !trustState {
				c.smeState = smeHelloStateAbort

				return errors.New("Trust denied. Connection aborted")
			} else {
				if sendErr := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypeReady, tHelloInit, false); sendErr != nil {
					c.smeState = smeHelloStateAbort

					return sendErr
				}

				// HELLO_OK
				c.smeState = smeHelloStateReadyListen

				if remoteSmeState == smeHelloStateOk {
					c.smeState = smeHelloStateOk
					return nil
				}

				continue
			}
		}
		// we got a new message
		_, data = c.parseMessage(data, true)

		var helloReturnMsg ship.ConnectionHello
		if err := json.Unmarshal(data, &helloReturnMsg); err != nil {
			c.smeState = smeHelloStateAbort

			return err
		}

		switch helloReturnMsg.ConnectionHello.Phase {
		case ship.ConnectionHelloPhaseTypeReady:
			// HELLO_OK
			remoteSmeState = smeHelloStateOk

			if c.smeState == smeHelloStateReadyListen {
				c.smeState = smeHelloStateOk
			}
		case ship.ConnectionHelloPhaseTypePending:
			// if we got a prolongation request, accept it
			if helloReturnMsg.ConnectionHello.ProlongationRequest != nil && *helloReturnMsg.ConnectionHello.ProlongationRequest {
				if sendErr := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypePending, tHelloInit, false); sendErr != nil {
					c.smeState = smeHelloStateAbort

					return sendErr
				}
			}

		case ship.ConnectionHelloPhaseTypeAborted:
			c.smeState = smeHelloStateAbort

			return errors.New("Connection aborted")
		default:
			c.smeState = smeHelloStateAbort

			return errors.New(fmt.Sprintf("Unexpected connection hello phase: %s", helloReturnMsg.ConnectionHello.Phase))
		}

		if c.smeState == smeHelloStateOk && remoteSmeState == smeHelloStateOk {
			return nil
		}
	}
}

func (c *ConnectionHandler) handshakeProtocol() error {
	var data []byte
	var err error

	protocolHandshake := ship.MessageProtocolHandshake{
		MessageProtocolHandshake: ship.MessageProtocolHandshakeType{
			Version: ship.Version{Major: 1, Minor: 0},
			Formats: ship.MessageProtocolFormatsType{
				Format: []ship.MessageProtocolFormatType{ship.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	if c.role == ShipRoleServer {
		// SME_PROT_H_STATE_SERVER_INIT
		data, _, err = c.readNextShipMessage(cmiTimeout)
		if err != nil {
			return err
		}
		_, data = c.parseMessage(data, true)

		// SME_PROT_H_STATE_SERVER_LISTEN_PROPOSAL
		messageProtocolHandshake := ship.MessageProtocolHandshake{}
		if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
			return err
		}

		if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeAnnounceMax {
			return errors.New("Invalid protocol handshake request")
		}

		protocolHandshake.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeSelect

		if err = c.sendShipModel(ship.MsgTypeControl, protocolHandshake); err != nil {
			return err
		}

		// SME_PROT_H_STATE_SERVER_LISTEN_CONFIRM
		data, _, err = c.readNextShipMessage(cmiTimeout)
		if err != nil {
			return err
		}
		_, data = c.parseMessage(data, true)

		messageProtocolHandshake = ship.MessageProtocolHandshake{}
		if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
			return err
		}

		if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeSelect {
			return errors.New("Invalid protocol handshake response")
		}

	} else {
		// SME_PROT_H_STATE_CLIENT_INIT

		protocolHandshake.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeAnnounceMax

		if err = c.sendShipModel(ship.MsgTypeControl, protocolHandshake); err != nil {
			return err
		}

		// SME_PROT_H_STATE_CLIENT_LISTEN_CHOICE
		data, _, err = c.readNextShipMessage(cmiTimeout)
		if err != nil {
			return err
		}
		_, data = c.parseMessage(data, true)

		messageProtocolHandshake := ship.MessageProtocolHandshake{}
		if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
			return err
		}

		if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeSelect {
			return errors.New("Invalid protocol handshake response")
		}

		protocolHandshake = ship.MessageProtocolHandshake{
			MessageProtocolHandshake: ship.MessageProtocolHandshakeType{
				HandshakeType: ship.ProtocolHandshakeTypeTypeSelect,
				Version:       ship.Version{Major: 1, Minor: 0},
				Formats: ship.MessageProtocolFormatsType{
					Format: []ship.MessageProtocolFormatType{ship.MessageProtocolFormatTypeUTF8},
				},
			},
		}
		if err = c.sendShipModel(ship.MsgTypeControl, protocolHandshake); err != nil {
			return err
		}
		// SME_PROT_H_STATE_CLIENT_OK
	}

	return nil
}

func (c *ConnectionHandler) handshakePin() error {
	var data []byte
	var err error

	// PIN State
	// SME_PIN_STATE_CHECK_INIT
	pinState := ship.ConnectionPinState{
		ConnectionPinState: ship.ConnectionPinStateType{
			PinState: ship.PinStateTypeNone,
		},
	}

	if err = c.sendShipModel(ship.MsgTypeControl, pinState); err != nil {
		return err
	}

	// SME_PIN_STATE_CHECK_LISTEN
	data, _, err = c.readNextShipMessage(cmiTimeout)
	if err != nil {
		return err
	}
	_, data = c.parseMessage(data, true)

	var connectionPinState ship.ConnectionPinState
	if err := json.Unmarshal([]byte(data), &connectionPinState); err != nil {
		return err
	}

	switch connectionPinState.ConnectionPinState.PinState {
	case ship.PinStateTypeNone:
		return nil
	case ship.PinStateTypeRequired:
		return errors.New("Got pin state: required (unsupported)")
	case ship.PinStateTypeOptional:
		return errors.New("Got pin state: optional (unsupported)")
	case ship.PinStateTypePinOk:
		return errors.New("Got pin state: ok (unsupported)")
	default:
		return errors.New("Got invalid pin state")
	}
}

func (c *ConnectionHandler) handshakeAccessMethods() error {
	var data []byte
	var err error

	// Access Methods
	accessMethodsRequest := ship.AccessMethodsRequest{
		AccessMethodsRequest: ship.AccessMethodsRequestType{},
	}

	if err = c.sendShipModel(ship.MsgTypeControl, accessMethodsRequest); err != nil {
		return err
	}

	for {
		data, _, err = c.readNextShipMessage(cmiTimeout)
		if err != nil {
			return err
		}
		if data == nil {
			continue
		}

		_, data = c.parseMessage(data, true)

		dataString := string(data)

		if strings.Contains(dataString, "\"accessMethodsRequest\":{") {
			methodsId := c.localService.ShipID

			accessMethods := ship.AccessMethods{
				AccessMethods: ship.AccessMethodsType{
					Id: &methodsId,
				},
			}
			if err = c.sendShipModel(ship.MsgTypeControl, accessMethods); err != nil {
				return err
			}
		} else if strings.Contains(dataString, "\"accessMethods\":{") {
			// compare SHIP ID to stored value on pairing. SKI + SHIP ID should be verified on connection
			// otherwise close connection with error "close 4450: SHIP id mismatch"

			var accessMethods ship.AccessMethods
			if err := json.Unmarshal([]byte(data), &accessMethods); err != nil {
				return err
			}

			if accessMethods.AccessMethods.Id == nil {
				return errors.New("Access methods response does not contain SHIP ID")
			}

			if len(c.remoteService.ShipID) > 0 && c.remoteService.ShipID != *accessMethods.AccessMethods.Id {
				return errors.New("SHIP id mismatch")
			}

			c.remoteService.ShipID = *accessMethods.AccessMethods.Id
			c.connectionDelegate.shipIDUpdateForService(c.remoteService)

			return nil
		} else {
			return errors.New(fmt.Sprintf("access methods: invalid response: %s", dataString))
		}
	}
}

func (c *ConnectionHandler) shipClose() {
	if c.conn == nil {
		return
	}

	// SHIP 13.4.7: Connection Termination
	closeMessage := ship.ConnectionClose{
		ConnectionClose: ship.ConnectionCloseType{
			Phase: ship.ConnectionClosePhaseTypeAnnounce,
		},
	}

	_ = c.sendShipModel(ship.MsgTypeControl, closeMessage)
}

// read the next message from the websocket connection
// read trust update
// return an error if the provided timeout is reached
func (c *ConnectionHandler) readNextShipMessage(duration time.Duration) ([]byte, bool, error) {
	timeout := time.NewTimer(duration)
	select {
	case <-timeout.C:
		if c.isConnectionClosed {
			return nil, false, nil
		}
		return nil, false, errors.New("Timeout waiting for message")
	case trust := <-c.shipTrustChannel:
		// Attention: we need to make sure the channel is only filled if we are in the Hello State!
		return nil, trust, nil
	case msg := <-c.shipReadChannel:
		timeout.Stop()
		return msg, false, nil
	}
}
