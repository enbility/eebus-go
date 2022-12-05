package ship

import (
	"encoding/json"
	"errors"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/ship/model"
)

// Handshake Prot covers the states smeProt...

func (c *ShipConnection) handshakeProtocol_Init() {
	switch c.role {
	case ShipRoleServer:
		c.setState(smeProtHStateServerInit)
		c.setHandshakeTimer(timeoutTimerTypeWaitForReady, cmiTimeout)
		c.setState(smeProtHStateServerListenProposal)
	case ShipRoleClient:
		c.setState(smeProtHStateClientInit)
		c.handshakeProtocol_smeProtHStateClientInit()
	}
}

// provide a ship.MessageProtocolHandshake struct
func (c *ShipConnection) protocolHandshake() model.MessageProtocolHandshake {
	protocolHandshake := model.MessageProtocolHandshake{
		MessageProtocolHandshake: model.MessageProtocolHandshakeType{
			Version: model.Version{Major: 1, Minor: 0},
			Formats: model.MessageProtocolFormatsType{
				Format: []model.MessageProtocolFormatType{model.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	return protocolHandshake
}

func (c *ShipConnection) handshakeProtocol_smeProtHStateServerListenProposal(message []byte) {
	_, data := c.parseMessage(message, true)

	messageProtocolHandshake := model.MessageProtocolHandshake{}
	if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != model.ProtocolHandshakeTypeTypeAnnounceMax {
		c.endHandshakeWithError(errors.New("Invalid protocol handshake request"))
		return
	}

	c.stopHandshakeTimer()

	protocolHandshake := c.protocolHandshake()
	protocolHandshake.MessageProtocolHandshake.HandshakeType = model.ProtocolHandshakeTypeTypeSelect

	if err := c.sendShipModel(model.MsgTypeControl, protocolHandshake); err != nil {
		c.endHandshakeWithError(err)
	}

	c.setHandshakeTimer(timeoutTimerTypeWaitForReady, cmiTimeout)

	c.setState(smeProtHStateServerListenConfirm)
}

func (c *ShipConnection) handshakeProtocol_smeProtHStateServerListenConfirm(message []byte) {
	_, data := c.parseMessage(message, true)

	var messageProtocolHandshake model.MessageProtocolHandshake
	if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
		logging.Log.Debug(err)
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeUnexpectedMessage)
		return
	}

	if messageProtocolHandshake.MessageProtocolHandshake.HandshakeType != model.ProtocolHandshakeTypeTypeSelect {
		logging.Log.Debug("invalid protocol handshake response")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	c.stopHandshakeTimer()

	c.setState(smeProtHStateServerOk)
	c.handleState(false, nil)
}

func (c *ShipConnection) handshakeProtocol_smeProtHStateClientInit() {
	c.setState(smeProtHStateClientInit)

	protocolHandshake := c.protocolHandshake()
	protocolHandshake.MessageProtocolHandshake.HandshakeType = model.ProtocolHandshakeTypeTypeAnnounceMax

	if err := c.sendShipModel(model.MsgTypeControl, protocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setState(smeProtHStateClientListenChoice)
}

func (c *ShipConnection) handshakeProtocol_smeProtHStateClientListenChoice(message []byte) {
	_, data := c.parseMessage(message, true)

	messageProtocolHandshake := model.MessageProtocolHandshake{}
	if err := json.Unmarshal([]byte(data), &messageProtocolHandshake); err != nil {
		logging.Log.Debug(err)
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeUnexpectedMessage)
		return
	}

	msgHandshake := messageProtocolHandshake.MessageProtocolHandshake

	if msgHandshake.HandshakeType != model.ProtocolHandshakeTypeTypeSelect {
		logging.Log.Debug("invalid protocol handshake response")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	if msgHandshake.Version.Major != 1 {
		logging.Log.Debug("unsupported protocol major version")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	if msgHandshake.Version.Minor != 0 {
		logging.Log.Debug("unsupported protocol minor version")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	if msgHandshake.Formats.Format == nil || len(msgHandshake.Formats.Format) == 0 {
		logging.Log.Debug("format is missing")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	if len(msgHandshake.Formats.Format) != 1 {
		logging.Log.Debug("unsupported format response")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	if msgHandshake.Formats.Format[0] != model.MessageProtocolFormatTypeUTF8 {
		logging.Log.Debug("unsupported format")
		c.abortProtocolHandshake(model.MessageProtocolHandshakeErrorErrorTypeSelectionMismatch)
		return
	}

	c.stopHandshakeTimer()

	protocolHandshake := c.protocolHandshake()
	protocolHandshake.MessageProtocolHandshake.HandshakeType = model.ProtocolHandshakeTypeTypeSelect

	if err := c.sendShipModel(model.MsgTypeControl, protocolHandshake); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setState(smeProtHStateClientOk)
	c.handleState(false, nil)
}

func (c *ShipConnection) abortProtocolHandshake(err model.MessageProtocolHandshakeErrorErrorType) {
	c.stopHandshakeTimer()

	msg := model.MessageProtocolHandshakeError{
		Error: err,
	}

	_ = c.sendShipModel(model.MsgTypeControl, msg)

	c.CloseConnection(false, "")
}
