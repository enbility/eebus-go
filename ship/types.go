package ship

import (
	"time"

	"github.com/enbility/eebus-go/ship/model"
)

type shipRole string

const (
	ShipRoleServer shipRole = "server"
	ShipRoleClient shipRole = "client"
)

const (
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second // SHIP 4.2: ping interval + pong timeout
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 50 * time.Second // SHIP 4.2: ping interval

	// SHIP 9.2: Set maximum fragment length to 1024 bytes
	MaxMessageSize = 1024
)

const (
	cmiTimeout              = 10 * time.Second // SHIP 4.2
	cmiCloseTimeout         = 100 * time.Millisecond
	tHelloInit              = 60 * time.Second // SHIP 13.4.4.1.3
	tHelloInc               = 60 * time.Second
	tHelloProlongThrInc     = 30 * time.Second
	tHelloProlongWaitingGap = 15 * time.Second
	tHelloProlongMin        = 1 * time.Second
)

type timeoutTimerType uint

const (
	// SHIP 13.4.4.1.3: The communication partner must send its "READY" state (or request for prolongation") before the timer expires.
	timeoutTimerTypeWaitForReady timeoutTimerType = iota
	// SHIP 13.4.4.1.3: Local timer to request for prolongation at the communication partner in time (i.e. before the communication partner's Wait-For-Ready-Timer expires).
	timeoutTimerTypeSendProlongationRequest
	// SHIP 13.4.4.1.3: Detection of response timeout on prolongation request.
	timeoutTimerTypeProlongRequestReply
)

type ShipState struct {
	State ShipMessageExchangeState
	Error error
}

type ShipMessageExchangeState uint

const (
	// Connection Mode Initialisation (CMI) SHIP 13.4.3
	CmiStateInitStart ShipMessageExchangeState = iota
	CmiStateClientSend
	CmiStateClientWait
	CmiStateClientEvaluate
	CmiStateServerWait
	CmiStateServerEvaluate
	// Connection Data Preparation SHIP 13.4.4
	SmeHelloState
	SmeHelloStateReadyInit
	SmeHelloStateReadyListen
	SmeHelloStateReadyTimeout
	SmeHelloStatePendingInit
	SmeHelloStatePendingListen
	SmeHelloStatePendingTimeout
	SmeHelloStateOk
	SmeHelloStateAbort
	SmeHelloStateAbortDone

	// Connection State Protocol Handhsake SHIP 13.4.4.2
	SmeProtHStateServerInit
	SmeProtHStateClientInit
	SmeProtHStateServerListenProposal
	SmeProtHStateServerListenConfirm
	SmeProtHStateClientListenChoice
	SmeProtHStateTimeout
	SmeProtHStateClientOk
	SmeProtHStateServerOk
	// Connection PIN State 13.4.5
	SmePinStateCheckInit
	SmePinStateCheckListen
	SmePinStateCheckError
	SmePinStateCheckBusyInit
	SmePinStateCheckBusyWait
	SmePinStateCheckOk
	SmePinStateAskInit
	SmePinStateAskProcess
	SmePinStateAskRestricted
	SmePinStateAskOk
	// ConnectionAccess Methods Identification 13.4.6
	SmeAccessMethodsRequest

	// Handshake approved on both ends
	SmeApproved

	// Handshake process is successfully completed
	SmeComplete

	// Handshake ended with an error
	SmeError
)

var shipInit []byte = []byte{model.MsgTypeInit, 0x00}

//go:generate mockgen -destination=mock_types.go -package=ship github.com/enbility/eebus-go/ship ShipDataConnection,ShipDataProcessing,ShipServiceDataProvider

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
	HandleConnectionClosed(*ShipConnection, bool)

	// report the ship ID provided during the handshake
	ReportServiceShipID(string, string)

	// check if the user is still able to trust the connection
	AllowWaitingForTrust(string) bool

	// report the updated SHIP handshake state and optional error message for a SKI
	HandleShipHandshakeStateUpdate(string, ShipState)
}
