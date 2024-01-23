package api

import (
	"github.com/enbility/ship-go/logging"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

/* Service */

type ServiceInterface interface {
	Setup() error
	Start()
	Shutdown()
	SetLogging(logger logging.LoggingInterface)

	Configuration() *Configuration
	LocalService() *shipapi.ServiceDetails
	LocalDevice() spineapi.DeviceLocalInterface
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

// interface for receiving data for specific events from Service
type ServiceReaderInterface interface {
	// report all currently visible EEBUS services
	VisibleRemoteServicesUpdated(service ServiceInterface, entries []shipapi.RemoteService)

	// report a connection to a SKI
	RemoteSKIConnected(service ServiceInterface, ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(service ServiceInterface, ski string)

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
