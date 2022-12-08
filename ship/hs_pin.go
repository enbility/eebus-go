package ship

import (
	"encoding/json"
	"errors"

	"github.com/enbility/eebus-go/ship/model"
)

// Handshake Pin covers the states smePin...

func (c *ShipConnection) handshakePin_Init() {
	c.setState(smePinStateCheckInit)

	pinState := model.ConnectionPinState{
		ConnectionPinState: model.ConnectionPinStateType{
			PinState: model.PinStateTypeNone,
		},
	}

	if err := c.sendShipModel(model.MsgTypeControl, pinState); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setState(smePinStateCheckListen)
}

func (c *ShipConnection) handshakePin_smePinStateCheckListen(message []byte) {
	_, data := c.parseMessage(message, true)

	var connectionPinState model.ConnectionPinState
	if err := json.Unmarshal([]byte(data), &connectionPinState); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	switch connectionPinState.ConnectionPinState.PinState {
	case model.PinStateTypeNone:
		c.setState(smePinStateCheckOk)
		c.handleState(false, nil)
	case model.PinStateTypeRequired:
		c.endHandshakeWithError(errors.New("Got pin state: required (unsupported)"))
	case model.PinStateTypeOptional:
		c.endHandshakeWithError(errors.New("Got pin state: optional (unsupported)"))
	case model.PinStateTypePinOk:
		c.endHandshakeWithError(errors.New("Got pin state: ok (unsupported)"))
	default:
		c.endHandshakeWithError(errors.New("Got invalid pin state"))
	}
}
