package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/DerAndereAndi/eebus-go/ship"
)

// process the ship handshake and return an error if the handshake failed
func (c *ConnectionHandler) shipHandshake() error {
	if err := c.handshakeInit(); err != nil {
		return err
	}

	if err := c.handshakeHello(); err != nil {
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

	return nil
}

// process the ship handshake and return an error if the handshake failed
func (c *ConnectionHandler) handshakeInit() error {
	var data []byte
	var msgType byte
	var err error

	shipInit := []byte{ship.MsgTypeInit, 0x00}

	if c.Role == ShipRoleClient {
		// CMI_STATE_CLIENT_SEND
		if err := c.writeWebsocketMessage(shipInit); err != nil {
			return err
		}
	}

	// CMI_STATE_SERVER_WAIT
	// CMI_STATE_CLIENT_WAIT
	data, err = c.readNextMessage(cmiTimeout)
	if err != nil {
		return err
	}

	if c.Role == ShipRoleServer {
		if err := c.writeWebsocketMessage(shipInit); err != nil {
			return err
		}
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

	return nil
}

func (c *ConnectionHandler) handshakeHello() error {
	waitingDuration := uint(tHelloInit.Milliseconds())

	// SME_HELLO_STATE_READY_INIT
	helloMsg := ship.ConnectionHello{
		ConnectionHello: ship.ConnectionHelloType{
			Phase:   ship.ConnectionHelloPhaseTypeReady,
			Waiting: &waitingDuration,
		},
	}

	if err := c.sendModel(helloMsg); err != nil {
		return err
	}

	// SME_HELLO_STATE_READY_LISTEN
	data, err := c.readNextMessage(cmiTimeout)
	if err != nil {
		return err
	}
	_, data = c.parseMessage(data, true)

	var helloReturnMsg ship.ConnectionHello
	if err := json.Unmarshal(data, &helloReturnMsg); err != nil {
		return err
	}

	switch helloReturnMsg.ConnectionHello.Phase {
	case ship.ConnectionHelloPhaseTypeReady:
		fmt.Println("Got ready message")
	case ship.ConnectionHelloPhaseTypeAborted:
		return errors.New("Connection aborted")
	default:
		return errors.New(fmt.Sprintf("Unexpected connection hello phase: %s", helloReturnMsg.ConnectionHello.Phase))
	}

	// HELLO_OK
	return nil
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

	if c.Role == ShipRoleServer {
		// SME_PROT_H_STATE_SERVER_INIT
		data, err = c.readNextMessage(cmiTimeout)
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

		if err = c.sendModel(protocolHandshake); err != nil {
			return err
		}

		// SME_PROT_H_STATE_SERVER_LISTEN_CONFIRM
		data, err = c.readNextMessage(cmiTimeout)
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

		if err = c.sendModel(protocolHandshake); err != nil {
			return err
		}

		// SME_PROT_H_STATE_CLIENT_LISTEN_CHOICE
		data, err = c.readNextMessage(cmiTimeout)
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
			},
		}
		if err = c.sendModel(protocolHandshake); err != nil {
			return err
		}
		// SME_PROT_H_STATE_CLIENT_OK
	}

	fmt.Println("Got protocol handshake")
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

	if err = c.sendModel(pinState); err != nil {
		return err
	}

	// SME_PIN_STATE_CHECK_LISTEN
	data, err = c.readNextMessage(cmiTimeout)
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
		fmt.Println("Got pin state: none")
	case ship.PinStateTypeRequired:
		return errors.New("Got pin state: required (unsupported)")
	case ship.PinStateTypeOptional:
		return errors.New("Got pin state: optional (unsupported)")
	case ship.PinStateTypePinOk:
		fmt.Println("Got pin state: ok (unsupported)")
	default:
		return errors.New("Got invalid pin state")
	}
	return nil
}

func (c *ConnectionHandler) handshakeAccessMethods() error {
	var data []byte
	var err error

	// Access Methods
	accessMethodsRequest := ship.AccessMethodsRequest{
		AccessMethodsRequest: ship.AccessMethodsRequestType{},
	}

	if err = c.sendModel(accessMethodsRequest); err != nil {
		return err
	}

	for {
		data, err = c.readNextMessage(cmiTimeout)
		if err != nil {
			return err
		}
		_, data = c.parseMessage(data, true)

		dataString := string(data)

		if strings.Contains(dataString, "\"accessMethodsRequest\":{") {
			fmt.Println("Got access methods request")
			methodsId := "Test"

			accessMethods := ship.AccessMethods{
				AccessMethods: ship.AccessMethodsType{
					Id: &methodsId,
				},
			}
			if err = c.sendModel(accessMethods); err != nil {
				return err
			}
		} else if strings.Contains(dataString, "\"accessMethods\":{") {
			fmt.Println("Got access methods")
			break
		} else {
			return errors.New(fmt.Sprintf("access methods: invalid response: %s", dataString))
		}
	}

	return nil
}
