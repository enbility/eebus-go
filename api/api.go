package api

import (
	"github.com/enbility/ship-go/logging"

	spineapi "github.com/enbility/spine-go/api"
)

// //go:generate mockery
//go:generate mockgen -destination=../mocks/mockgen_api.go -package=mocks github.com/enbility/eebus-go/api ServiceProvider,MdnsService

/* EEBUSService */

// interface for receiving data for specific events from EEBUSService
type EEBUSServiceHandler interface {
	// report all currently visible EEBUS services
	VisibleRemoteServicesUpdated(service EEBUSService, entries []RemoteService)

	// report a connection to a SKI
	RemoteSKIConnected(service EEBUSService, ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(service EEBUSService, ski string)

	// Provides the SHIP ID the remote service reported during the handshake process
	// This needs to be persisted and passed on for future remote service connections
	// when using `PairRemoteService`
	ServiceShipIDUpdate(ski string, shipdID string)

	// Provides the current pairing state for the remote service
	// This is called whenever the state changes and can be used to
	// provide user information for the pairing/connection process
	ServicePairingDetailUpdate(ski string, detail *ConnectionStateDetail)

	// return if the user is still able to trust the connection
	AllowWaitingForTrust(ski string) bool
}

type EEBUSService interface {
	Setup() error
	Start()
	Shutdown()
	SetLogging(logger logging.Logging)

	LocalDevice() spineapi.DeviceLocal
	RemoteServiceForSKI(ski string) *ServiceDetails
	RegisterRemoteSKI(ski string, enable bool)
	InitiatePairingWithSKI(ski string)
	CancelPairingWithSKI(ski string)
	DisconnectSKI(ski string, reason string)

	// Passthough functions to ConnectionsHub
	PairingDetailForSki(ski string) *ConnectionStateDetail
	StartBrowseMdnsEntries()
	StopBrowseMdnsEntries()
}

/* Hub */

// interface for reporting data from connectionsHub to the Service
type ServiceProvider interface {
	// report a newly discovered remote EEBUS service
	VisibleMDNSRecordsUpdated(entries []*MdnsEntry)

	// report a connection to a SKI
	RemoteSKIConnected(ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(ski string)

	// provide the SHIP ID received during SHIP handshake process
	// the ID needs to be stored and then provided for remote services so it can be compared and verified
	ServiceShipIDUpdate(ski string, shipID string)

	// provides the current handshake state for a given SKI
	ServicePairingDetailUpdate(ski string, detail *ConnectionStateDetail)

	// return if the user is still able to trust the connection
	AllowWaitingForTrust(ski string) bool
}

type ConnectionsHub interface {
	PairingDetailForSki(ski string) *ConnectionStateDetail
	StartBrowseMdnsSearch()
	StopBrowseMdnsSearch()
	Start()
	Shutdown()
	ServiceForSKI(ski string) *ServiceDetails
	RegisterRemoteSKI(ski string, enable bool)
	InitiatePairingWithSKI(ski string)
	CancelPairingWithSKI(ski string)
	DisconnectSKI(ski string, reason string)
}

/* Mdns */

// implemented by hubConnection, used by mdns
type MdnsSearch interface {
	ReportMdnsEntries(entries map[string]*MdnsEntry)
}

// implemented by mdns, used by hubConnection
type MdnsService interface {
	SetupMdnsService() error
	ShutdownMdnsService()
	AnnounceMdnsEntry() error
	UnannounceMdnsEntry()
	RegisterMdnsSearch(cb MdnsSearch)
	UnregisterMdnsSearch(cb MdnsSearch)
}
