package spine

import "github.com/enbility/eebus-go/spine/model"

//go:generate mockery --name=SpineDataProcessing

// Used to pass an incoming SPINE message from a SHIP connection to the proper DeviceRemoteImpl
//
// Implemented by DeviceRemoteImpl, used by ShipConnection
type SpineDataProcessing interface {
	HandleIncomingSpineMesssage(message []byte) (*model.MsgCounterType, error)
}

//go:generate mockery --name=SpineDataConnection

// Used to pass an outgoing SPINE message from a DeviceLocalImpl to the SHIP connection
//
// Implemented by ShipConnection, used by DeviceLocalImpl
type SpineDataConnection interface {
	WriteSpineMessage(message []byte)
}
