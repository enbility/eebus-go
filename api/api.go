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
	// Setup the service
	//
	// Returns error with description of the error that cannot be recovered from
	Setup() error

	// start the service
	//
	// Returns error with description of the error that cannot be recovered from
	Start() error

	// shutdown the service
	Shutdown()

	// return if the service is running
	IsRunning() bool

	// add a use case to the service
	AddUseCase(useCase UseCaseInterface) error

	// set logging interface
	SetLogging(logger logging.LoggingInterface)

	// return the configuration
	Configuration() *Configuration

	// return the local service details
	LocalService() *shipapi.ServiceDetails

	// return the local device
	LocalDevice() spineapi.DeviceLocalInterface

	// Passthough functions to HubInterface

	// Provide the current pairing state for a ServiceIdentity
	PairingDetailFor(identity shipapi.ServiceIdentity) *shipapi.ConnectionStateDetail

	// Return the remote service details for a given ServiceIdentity
	RemoteServiceFor(identity shipapi.ServiceIdentity) *shipapi.ServiceDetails

	// Defines wether incoming pairing requests should be automatically accepted or not
	//
	// Default: false
	SetAutoAccept(value bool)

	// Returns if the service has auto accept enabled or not
	IsAutoAcceptEnabled() bool

	// Generate a QR code string
	//
	// If a pairing config with a secret is set: generates SHIP Pairing Service QR format
	// Otherwise: generates standard SHIP QR format
	//
	// Must be called after Setup() as it requires the hub to be initialized.
	QRCodeText() (string, error)

	// Pair a remote service using ServiceIdentity
	//
	// Parameters:
	// - identity: ServiceIdentity containing SKI, fingerprint, and/or SHIP ID
	//
	// Note: The SHIP ID is optional, but should be provided if available.
	// if provided, it will be used to validate the remote service is
	// providing this SHIP ID during the handshake process and will reject
	// the connection if it does not match.
	RegisterRemoteService(identity shipapi.ServiceIdentity)

	// Unpair a remote service using ServiceIdentity
	UnregisterRemoteService(identity shipapi.ServiceIdentity)

	// Disconnect a connection using ServiceIdentity
	DisconnectService(identity shipapi.ServiceIdentity, reason string)

	// Cancels the pairing process for a ServiceIdentity
	//
	// This should be called while the service is running and the end
	// user wants to cancel/disallow an incoming pairing request
	CancelPairing(identity shipapi.ServiceIdentity)

	// Define wether the user is able to react to an incoming pairing request
	//
	// Call this with `true` e.g. if the user is currently using a web interface
	// where an incoming request can be accepted or denied
	//
	// Default is set to false, meaning every incoming pairing request will be
	// automatically denied
	UserIsAbleToApproveOrCancelPairingRequests(allow bool)

	// Calculate SHA-256 fingerprint of local certificate
	GetLocalCertificateFingerprint() (string, error)

	// **************************
	// SHIP Pairing Service APIs
	// **************************

	// Start announcing pairing to a specific target device
	// Used by devZ only.
	//
	// Parameters:
	//  - target: Pairing target
	StartAnnouncementTo(target shipapi.PairingTarget) error

	// Stop announcing pairing to a specific target device
	// Used by devZ only.
	//
	// Parameters:
	//  - shipID: Target SHIP ID
	StopAnnouncementTo(shipID string) error

	// Return true if currently announcing to a specific target device
	// Used by devZ only.
	//
	// Parameters:
	//  - shipID: Target SHIP ID
	IsAnnouncingTo(shipID string) bool

	// SHIP Pairing: Get Active Announcements.
	// Used by devZ only.
	//
	// Returns: List of SHIP IDs currently being announced to
	GetActiveAnnouncements() []string

	// SHIP Pairing: Get the SHIP ID and Fingerprint of controlbox paired via SHIP Pairing
	// Used by devA only.
	//
	// Returns: the ServiceDetails of any trusted AddCu device. Or nil if none
	GetTrustedAddCuDevice() *shipapi.ServiceDetails
}

// interface for receiving data for specific events from Service
//
// some are passthrough readers, because service needs to coordinate
// everything with SPINE
//
// implemented by the eebus service implementation, used by service
type ServiceReaderInterface interface {
	// report a connection to a remote service
	RemoteServiceConnected(service ServiceInterface, identity shipapi.ServiceIdentity)

	// report a disconnection from a remote service
	RemoteServiceDisconnected(service ServiceInterface, identity shipapi.ServiceIdentity)

	// report all currently visible EEBUS services
	VisibleRemoteMdnsServicesUpdated(service ServiceInterface, entries []shipapi.RemoteMdnsService)

	// report that service information has been updated
	// This includes updates to ShipID, fingerprint, or other service details
	// discovered during handshake
	ServiceUpdated(identity shipapi.ServiceIdentity)

	// Provides the current pairing state for the remote service
	// This is called whenever the state changes and can be used to
	// provide user information for the pairing/connection process
	ServicePairingDetailUpdate(identity shipapi.ServiceIdentity, detail *shipapi.ConnectionStateDetail)

	// ****************************
	// SHIP Pairing Service Events
	// ****************************

	// Called when a device is automatically trusted via SHIP pairing
	ServiceAutoTrusted(service ServiceInterface, identity shipapi.ServiceIdentity)

	// Called when SHIP pairing fails for a device
	ServiceAutoTrustFailed(service ServiceInterface, identity shipapi.ServiceIdentity, reason error)

	// Called when device trust is automatically removed
	// This can happen due to device replacement timeout or new device pairing
	ServiceAutoTrustRemoved(service ServiceInterface, identity shipapi.ServiceIdentity, reason string)
}
