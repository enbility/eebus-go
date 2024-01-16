package api

import (
	"github.com/enbility/ship-go/logging"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

/* EEBUSService */

// interface for receiving data for specific events from EEBUSService
type EEBUSServiceHandler interface {
	// report all currently visible EEBUS services
	VisibleRemoteServicesUpdated(service EEBUSService, entries []shipapi.RemoteService)

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
	ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail)

	// return if the user is still able to trust the connection
	AllowWaitingForTrust(ski string) bool
}

type EEBUSService interface {
	Setup() error
	Start()
	Shutdown()
	SetLogging(logger logging.Logging)

	Configuration() *Configuration
	LocalService() *shipapi.ServiceDetails
	LocalDevice() spineapi.DeviceLocal
	RemoteServiceForSKI(ski string) *shipapi.ServiceDetails
	RegisterRemoteSKI(ski string, enable bool)
	InitiatePairingWithSKI(ski string)
	CancelPairingWithSKI(ski string)
	DisconnectSKI(ski string, reason string)

	// Passthough functions to ConnectionsHub
	PairingDetailForSki(ski string) *shipapi.ConnectionStateDetail
	StartBrowseMdnsEntries()
	StopBrowseMdnsEntries()
}
