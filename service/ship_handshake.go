package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/ship"
)

// handle incoming SHIP messages and coordinate Handshake States
func (c *shipConnection) handleShipState(timeout bool, message []byte) {
	if c.handshakeTimerRunning {
		if !c.handshakeTimer.Stop() {
			<-c.handshakeTimer.C
		}
		c.handshakeTimerRunning = false
	}

	if len(message) > 2 {
		logging.Log.Trace("Recv:", c.remoteService.SKI, string(message[1:]))
	}

	// TODO: first check if this is a close message

	switch c.smeState {
	// cmiStateInit
	case cmiStateInitStart:
		// triggered without a message received
		c.handshakeInit_cmiStateInitStart()

	case cmiStateClientWait:
		if timeout {
			c.endHandshakeWithError(errors.New("ship handshake timeout"))
			return
		}

		c.handshakeInit_cmiStateClientWait(message)

	case cmiStateServerWait:
		if timeout {
			c.endHandshakeWithError(errors.New("ship handshake timeout"))
			return
		}
		c.handshakeInit_cmiStateServerWait(message)

	// smeHello

	case smeHelloStateReadyListen:
		if timeout {
			c.endHandshakeWithError(errors.New("ship handshake timeout"))
			return
		}

		c.handshakeHello_ReadyListen(message)

	case smeHelloStatePendingListen:
		if timeout {
			c.handshakeHello_AskProlongation()
			return
		}

		if !c.connectionIsTrusted {
			c.smeState = smeHelloStateAbort
			c.endHandshakeWithError(errors.New("trust denied. close connection"))
			return
		}

	// smeProtocol

	case smeProtHStateServerListenProposal:
		c.handshakeProtocol_smeProtHStateServerListenProposal(message)

	case smeProtHStateServerListenConfirm:
		c.handshakeProtocol_smeProtHStateServerListenConfirm(message)

	case smeProtHStateClientListenChoice:
		c.handshakeProtocol_smeProtHStateClientListenChoice(message)

	// smePinState

	case smePinStateCheckListen:
		c.handshakePin_smePinStateCheckListen(message)

	// smeAccessMethods

	case smeAccessMethodsRequest:
		c.handshakeAccessMethods_Request(message)

	// Generic

	case smeHelloStateAbort, smeHandshakeError:
		c.shutdown(true)
	case smeComplete:
		// TODO: handshake is completed, we now probably get close messages
		return
	}
}

func (c *shipConnection) approveHandshake() {
	// Report to SPINE local device about this remote device connection
	c.readHandler = c.connectionDelegate.addRemoteDeviceConnection(c.remoteService.SKI, c)
	c.smeState = smeComplete
}

// end the handshake process because of an error
func (c *shipConnection) endHandshakeWithError(err error) {
	c.smeState = smeHandshakeError

	if c.handshakeError != nil {
		logging.Log.Error(c.remoteService.SKI, "SHIP handshake error:", c.handshakeError)
	} else {
		logging.Log.Error(c.remoteService.SKI, "SHIP handshake error unknown")
	}
	c.shutdown(false)
}

// Handshake initialization

// cmiStateInitStart
func (c *shipConnection) handshakeInit_cmiStateInitStart() {
	switch c.role {
	case ShipRoleClient:
		// CMI_STATE_CLIENT_SEND
		c.smeState = cmiStateClientSend
		if err := c.writeWebsocketMessage(shipInit); err != nil {
			c.endHandshakeWithError(err)
			return
		}
		c.smeState = cmiStateClientWait
	case ShipRoleServer:
		c.smeState = cmiStateServerWait
	}

	c.resetHandshakeTimer(cmiTimeout)
}

// CMI_STATE_SERVER_WAIT
func (c *shipConnection) handshakeInit_cmiStateServerWait(message []byte) {
	c.smeState = cmiStateServerEvaluate

	c.handshakeInit_cmiStateEvaluate(message)

	if err := c.writeWebsocketMessage(shipInit); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.smeState = smeHelloState
	c.handshakeHello_Init()
}

// CMI_STATE_CLIENT_WAIT
func (c *shipConnection) handshakeInit_cmiStateClientWait(message []byte) {
	c.smeState = cmiStateClientEvaluate

	c.handshakeInit_cmiStateEvaluate(message)

	c.smeState = smeHelloState
	c.handshakeHello_Init()
}

// CMI_STATE_SERVER_EVALUATE
// CMI_STATE_CLIENT_EVALUATE
func (c *shipConnection) handshakeInit_cmiStateEvaluate(message []byte) {
	msgType, data := c.parseMessage(message, false)

	if msgType != ship.MsgTypeInit {
		c.endHandshakeWithError(fmt.Errorf("Invalid SHIP MessageType, expected 0 and got %s" + string(msgType)))
		return
	}
	if data[0] != byte(0) {
		c.endHandshakeWithError(fmt.Errorf("Invalid SHIP MessageValue, expected 0 and got %s" + string(data)))
	}
}

// Handshake Hello

// smeHelloState
func (c *shipConnection) handshakeHello_Init() {
	switch c.connectionIsTrusted {
	case true:
		c.smeState = smeHelloStateReadyInit

		c.handshakeHello_Trust()
	case false:
		c.smeState = smeHelloStatePendingInit

		if err := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
			c.smeState = smeHelloStateAbort

			c.endHandshakeWithError(err)
			return
		}

		c.smeState = smeHelloStatePendingListen

		c.connectionDelegate.requestUserTrustForService(c.remoteService)
	}

	c.resetHandshakeTimer(cmiTimeout)
}

// Timeout reached we didn't receive a user input, so ask for prolongation
func (c *shipConnection) handshakeHello_AskProlongation() {
	c.smeState = smeHelloStatePendingTimeout

	if err := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypePending, 0, true); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.smeState = smeHelloStatePendingListen

	c.resetHandshakeTimer(tHelloProlongThrInc)
}

// A user trust was received
func (c *shipConnection) handshakeHello_Trust() {
	if err := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypeReady, tHelloInit, false); err != nil {
		c.smeState = smeHelloStateAbort
		c.endHandshakeWithError(err)
		return
	}

	c.smeState = smeHelloStateReadyListen
}

func (c *shipConnection) handshakeHello_ReadyListen(message []byte) {
	_, data := c.parseMessage(message, true)

	var helloReturnMsg ship.ConnectionHello
	if err := json.Unmarshal(data, &helloReturnMsg); err != nil {
		c.smeState = smeHelloStateAbort
		c.endHandshakeWithError(err)
		return
	}

	switch helloReturnMsg.ConnectionHello.Phase {
	case ship.ConnectionHelloPhaseTypeReady:
		// HELLO_OK
		c.smeState = smeHelloStateOk
	case ship.ConnectionHelloPhaseTypePending:
		// if we got a prolongation request, accept it
		if helloReturnMsg.ConnectionHello.ProlongationRequest != nil && *helloReturnMsg.ConnectionHello.ProlongationRequest {
			if err := c.handshakeHelloSend(ship.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
				c.smeState = smeHelloStateAbort
				c.endHandshakeWithError(err)
				return
			}
		}

	case ship.ConnectionHelloPhaseTypeAborted:
		c.smeState = smeHelloStateAbort
		c.endHandshakeWithError(errors.New("Connection aborted"))
		return
	default:
		c.smeState = smeHelloStateAbort
		c.endHandshakeWithError(fmt.Errorf("Unexpected connection hello phase: %s", helloReturnMsg.ConnectionHello.Phase))
		return
	}

	if c.smeState == smeHelloStateOk {
		c.handshakeProtocol_Init()
	}
}

func (c *shipConnection) handshakeHelloSend(phase ship.ConnectionHelloPhaseType, waitingDuration time.Duration, prolongation bool) error {
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

// smeHelloStateOk

func (c *shipConnection) handshakeProtocol_Init() {
	switch c.role {
	case ShipRoleServer:
		c.smeState = smeProtHStateServerInit
		c.resetHandshakeTimer(cmiTimeout)
		c.smeState = smeProtHStateServerListenProposal
	case ShipRoleClient:
		c.smeState = smeProtHStateClientInit
		c.handshakeProtocol_smeProtHStateClientInit()
	}
}

// provide a ship.MessageProtocolHandshake struct
func (c *shipConnection) protocolHandshake() ship.MessageProtocolHandshake {
	protocolHandshake := ship.MessageProtocolHandshake{
		MessageProtocolHandshake: ship.MessageProtocolHandshakeType{
			Version: ship.Version{Major: 1, Minor: 0},
			Formats: ship.MessageProtocolFormatsType{
				Format: []ship.MessageProtocolFormatType{ship.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	return protocolHandshake
}

func (c *shipConnection) handshakeProtocol_smeProtHStateServerListenProposal(message []byte) {
	_, data := c.parseMessage(message, true)

	messageProtocolHandshake := ship.MessageProtocolHandshake{}
	if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeAnnounceMax {
		c.endHandshakeWithError(errors.New("Invalid protocol handshake request"))
		return
	}

	protocolHandshake := c.protocolHandshake()
	protocolHandshake.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeSelect

	if err := c.sendShipModel(ship.MsgTypeControl, protocolHandshake); err != nil {
		c.endHandshakeWithError(err)
	}

	c.smeState = smeProtHStateServerListenConfirm
}

func (c *shipConnection) handshakeProtocol_smeProtHStateServerListenConfirm(message []byte) {
	_, data := c.parseMessage(message, true)

	var messageProtocolHandshake ship.MessageProtocolHandshake
	if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeSelect {
		c.endHandshakeWithError(errors.New("Invalid protocol handshake response"))
		return
	}

	c.smeState = smeProtHStateServerOk
	c.handshakePin_Init()
}

func (c *shipConnection) handshakeProtocol_smeProtHStateClientInit() {
	c.smeState = smeProtHStateClientInit

	protocolHandshake := c.protocolHandshake()
	protocolHandshake.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeAnnounceMax

	if err := c.sendShipModel(ship.MsgTypeControl, protocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.smeState = smeProtHStateClientListenChoice
}

func (c *shipConnection) handshakeProtocol_smeProtHStateClientListenChoice(message []byte) {
	_, data := c.parseMessage(message, true)

	messageProtocolHandshake := ship.MessageProtocolHandshake{}
	if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeSelect {
		c.endHandshakeWithError(errors.New("Invalid protocol handshake response"))
		return
	}

	protocolHandshake := c.protocolHandshake()
	protocolHandshake.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeSelect

	if err := c.sendShipModel(ship.MsgTypeControl, protocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.smeState = smeProtHStateClientOk
	c.handshakePin_Init()
}

// handshake PIN

func (c *shipConnection) handshakePin_Init() {
	c.smeState = smePinStateCheckInit

	pinState := ship.ConnectionPinState{
		ConnectionPinState: ship.ConnectionPinStateType{
			PinState: ship.PinStateTypeNone,
		},
	}

	if err := c.sendShipModel(ship.MsgTypeControl, pinState); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.smeState = smePinStateCheckListen
}

func (c *shipConnection) handshakePin_smePinStateCheckListen(message []byte) {
	_, data := c.parseMessage(message, true)

	var connectionPinState ship.ConnectionPinState
	if err := json.Unmarshal([]byte(data), &connectionPinState); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	switch connectionPinState.ConnectionPinState.PinState {
	case ship.PinStateTypeNone:
		c.smeState = smePinStateCheckOk
		c.handshakeAccessMethods_Init()
	case ship.PinStateTypeRequired:
		c.endHandshakeWithError(errors.New("Got pin state: required (unsupported)"))
	case ship.PinStateTypeOptional:
		c.endHandshakeWithError(errors.New("Got pin state: optional (unsupported)"))
	case ship.PinStateTypePinOk:
		c.endHandshakeWithError(errors.New("Got pin state: ok (unsupported)"))
	default:
		c.endHandshakeWithError(errors.New("Got invalid pin state"))
	}
}

func (c *shipConnection) handshakeAccessMethods_Init() {
	// Access Methods
	accessMethodsRequest := ship.AccessMethodsRequest{
		AccessMethodsRequest: ship.AccessMethodsRequestType{},
	}

	if err := c.sendShipModel(ship.MsgTypeControl, accessMethodsRequest); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.resetHandshakeTimer(cmiTimeout)
	c.smeState = smeAccessMethodsRequest
}

func (c *shipConnection) handshakeAccessMethods_Request(message []byte) {
	_, data := c.parseMessage(message, true)

	dataString := string(data)

	if strings.Contains(dataString, "\"accessMethodsRequest\":{") {
		methodsId := c.localService.ShipID

		accessMethods := ship.AccessMethods{
			AccessMethods: ship.AccessMethodsType{
				Id: &methodsId,
			},
		}
		if err := c.sendShipModel(ship.MsgTypeControl, accessMethods); err != nil {
			c.endHandshakeWithError(err)
		}
		return
	} else if strings.Contains(dataString, "\"accessMethods\":{") {
		// compare SHIP ID to stored value on pairing. SKI + SHIP ID should be verified on connection
		// otherwise close connection with error "close 4450: SHIP id mismatch"

		var accessMethods ship.AccessMethods
		if err := json.Unmarshal([]byte(data), &accessMethods); err != nil {
			c.endHandshakeWithError(err)
			return
		}

		if accessMethods.AccessMethods.Id == nil {
			c.endHandshakeWithError(errors.New("Access methods response does not contain SHIP ID"))
			return
		}

		if len(c.remoteService.ShipID) > 0 && c.remoteService.ShipID != *accessMethods.AccessMethods.Id {
			c.endHandshakeWithError(errors.New("SHIP id mismatch"))
			return
		}

		c.remoteService.ShipID = *accessMethods.AccessMethods.Id
		c.connectionDelegate.shipIDUpdateForService(c.remoteService)
	} else {
		c.endHandshakeWithError(fmt.Errorf("access methods: invalid response: %s", dataString))
		return
	}

	c.smeState = smeApproved
	c.approveHandshake()
}

func (c *shipConnection) shipClose() {
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

func (c *shipConnection) resetHandshakeTimer(duration time.Duration) {
	c.handshakeTimer.Reset(cmiTimeout)
	c.handshakeTimerRunning = true
}
