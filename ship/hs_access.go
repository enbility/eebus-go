package ship

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/DerAndereAndi/eebus-go/ship/model"
)

// Handshake Access covers the states smeAccess...

func (c *ShipConnection) handshakeAccessMethods_Init() {
	// Access Methods
	accessMethodsRequest := model.AccessMethodsRequest{
		AccessMethodsRequest: model.AccessMethodsRequestType{},
	}

	if err := c.sendShipModel(model.MsgTypeControl, accessMethodsRequest); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setHandshakeTimer(timeoutTimerTypeWaitForReady, cmiTimeout)
	c.setState(smeAccessMethodsRequest)
}

func (c *ShipConnection) handshakeAccessMethods_Request(message []byte) {
	_, data := c.parseMessage(message, true)

	dataString := string(data)

	if strings.Contains(dataString, "\"accessMethodsRequest\":{") {
		methodsId := c.localShipID

		accessMethods := model.AccessMethods{
			AccessMethods: model.AccessMethodsType{
				Id: &methodsId,
			},
		}
		if err := c.sendShipModel(model.MsgTypeControl, accessMethods); err != nil {
			c.endHandshakeWithError(err)
		}
		return
	} else if strings.Contains(dataString, "\"accessMethods\":{") {
		// compare SHIP ID to stored value on pairing. SKI + SHIP ID should be verified on connection
		// otherwise close connection with error "close 4450: SHIP id mismatch"

		var accessMethods model.AccessMethods
		if err := json.Unmarshal([]byte(data), &accessMethods); err != nil {
			c.endHandshakeWithError(err)
			return
		}

		if accessMethods.AccessMethods.Id == nil {
			c.endHandshakeWithError(errors.New("Access methods response does not contain SHIP ID"))
			return
		}

		if len(c.remoteShipID) > 0 && c.remoteShipID != *accessMethods.AccessMethods.Id {
			c.endHandshakeWithError(errors.New("SHIP id mismatch"))
			return
		}

		c.remoteShipID = *accessMethods.AccessMethods.Id
		// c.interactionHandler.shipIDUpdateForService(c.remoteService)
	} else {
		c.endHandshakeWithError(fmt.Errorf("access methods: invalid response: %s", dataString))
		return
	}

	c.setState(smeApproved)
	c.approveHandshake()
}
