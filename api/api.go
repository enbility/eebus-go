package api

import (
	"github.com/enbility/ship-go/logging"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

/* Service */

// central service interface
//
// implemented by service, used by the eebus service implementation
type ServiceInterface interface {
	// setup the service
	Setup() error

	// start the service
	Start()

	// shutdown the service
	Shutdown()

	// set logging interface
	SetLogging(logger logging.LoggingInterface)

	// return the configuration
	Configuration() *Configuration

	// return the local service details
	LocalService() *shipapi.ServiceDetails

	// return the local device
	LocalDevice() spineapi.DeviceLocalInterface

	// Passthough functions to HubInterface

	// Provide the current pairing state for a SKI
	PairingDetailForSki(ski string) *shipapi.ConnectionStateDetail

	// Returns the Service detail of a given remote SKI
	RemoteServiceForSKI(ski string) *shipapi.ServiceDetails

	// Sets the SKI as being paired or not
	RegisterRemoteSKI(ski string, enable bool)

	// Disconnect a connection to an SKI
	DisconnectSKI(ski string, reason string)

	// Triggers the pairing process for a SKI
	InitiateOrApprovePairingWithSKI(ski string)

	// Cancels the pairing process for a SKI
	CancelPairingWithSKI(ski string)
}

// interface for receiving data for specific events from Service
//
// some are passthrough readers, because service needs to coordinate
// everything with SPINE
//
// implemented by the eebus service implementation, used by service
type ServiceReaderInterface interface {
	// report a connection to a SKI
	RemoteSKIConnected(service ServiceInterface, ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(service ServiceInterface, ski string)

	// report all currently visible EEBUS services
	VisibleRemoteServicesUpdated(service ServiceInterface, entries []shipapi.RemoteService)

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
