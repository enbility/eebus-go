package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

// Used to pass an incoming SPINE message from a SHIP connection to the proper DeviceRemoteImpl
type ReadMessageI interface {
	ReadMessage(message []byte) (*model.MsgCounterType, error)
}

// Used to pass an outgoing SPINE message from a DeviceLocalImpl to the SHIP connection
type WriteMessageI interface {
	WriteMessage(message []byte)
}
