package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

// Used to pass an incoming SPINE message from a SHIP connection to the proper DeviceRemoteImpl
//
// Implemented by DeviceRemoteImpl, used by ShipConnection
type SpineDataProcessing interface {
	HandleIncomingSpineMesssage(message []byte) (*model.MsgCounterType, error)
}

// Used to pass an outgoing SPINE message from a DeviceLocalImpl to the SHIP connection
//
// Implemented by ShipConnection, used by DeviceLocalImpl
type SpineDataConnection interface {
	WriteSpineMessage(message []byte)
}
