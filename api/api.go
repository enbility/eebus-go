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

	// return if the service is running
	IsRunning() bool

	// add a use case to the service
	AddUseCase(useCase UseCaseInterface)

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

	// Defines wether incoming pairing requests should be automatically accepted or not
	//
	// Default: false
	SetAutoAccept(value bool)

	// Returns if the service has auto accept enabled or not
	IsAutoAcceptEnabled() bool

	// Returns the QR code text for the service
	// as defined in SHIP Requirements for Installation Process V1.0.0
	QRCodeText() string

	// Returns the Service detail of a remote SKI
	RemoteServiceForSKI(ski string) *shipapi.ServiceDetails

	// Sets the SKI as being paired
	RegisterRemoteSKI(ski string)

	// Sets the SKI as not being paired
	UnregisterRemoteSKI(ski string)

	// Disconnect from a connected remote SKI
	DisconnectSKI(ski string, reason string)

	// Cancels the pairing process for a SKI
	//
	// This should be called while the service is running and the end
	// user wants to cancel/disallow an incoming pairing request
	CancelPairingWithSKI(ski string)

	// Define wether the user is able to react to an incoming pairing request
	//
	// Call this with `true` e.g. if the user is currently using a web interface
	// where an incoming request can be accepted or denied
	//
	// Default is set to false, meaning every incoming pairing request will be
	// automatically denied
	UserIsAbleToApproveOrCancelPairingRequests(allow bool)
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
}
