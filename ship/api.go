package ship

//go:generate mockgen -destination=mock_types_test.go -package=ship github.com/enbility/eebus-go/ship ShipDataConnection,ShipDataProcessing,ShipServiceDataProvider
//go:generate mockery --name=ShipDataConnection
//go:generate mockery --name=ShipConnection

type ShipConnection interface {
	DataHandler() ShipDataConnection
	CloseConnection(safe bool, code int, reason string)
	RemoteSKI() string
	ApprovePendingHandshake()
	AbortPendingHandshake()
	ShipHandshakeState() (ShipMessageExchangeState, error)
}

// interface for handling the actual remote device data connection
//
// implemented by websocketConnection, used by ShipConnection
type ShipDataConnection interface {
	// initialize data processing
	InitDataProcessing(ShipDataProcessing)

	// send data via the connection to the remote device
	WriteMessageToDataConnection([]byte) error

	// close the data connection
	CloseDataConnection(closeCode int, reason string)

	// report if the data connection is closed and the error if availab le
	IsDataConnectionClosed() (bool, error)
}

// interface for handling incoming data
//
// implemented by shipConnection, used by websocketConnection
type ShipDataProcessing interface {
	// called for each incoming message
	HandleIncomingShipMessage([]byte)

	// called if the data connection is closed unsafe
	// e.g. due to connection issues
	ReportConnectionError(error)
}

// interface for getting service wide information
//
// implemented by connectionsHub, used by shipConnection
type ShipServiceDataProvider interface {
	// check if the SKI is paired
	IsRemoteServiceForSKIPaired(string) bool

	// report closing of a connection and if handshake did complete
	HandleConnectionClosed(ShipConnection, bool)

	// report the ship ID provided during the handshake
	ReportServiceShipID(string, string)

	// check if the user is still able to trust the connection
	AllowWaitingForTrust(string) bool

	// report the updated SHIP handshake state and optional error message for a SKI
	HandleShipHandshakeStateUpdate(string, ShipState)
}
