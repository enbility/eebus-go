package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

const defaultPort int = 4711

type ServiceDetails struct {
	// This is the SKI of the service
	// This needs to be persisted
	SKI string

	// This is the IPv4 address of the device running the service
	// This is optional only needed when this runs with
	// zeroconf as mDNS and the remote device is using the latest
	// avahi version and thus zeroconf can sometimes not detect
	// the IPv4 address and not initiate a connection
	IPv4 string

	// ShipID is the ship identifier of the service
	// This needs to be persisted
	ShipID string

	// The EEBUS device type of the device model
	deviceType model.DeviceTypeType

	// Flags if the service auto auto accepts other services
	registerAutoAccept bool

	// Flag if a user interaction marks this service as trusted
	userTrust bool
}

type ServiceDescription struct {
	// The vendors IANA PEN, optional but highly recommended.
	// If not set, brand will be used instead
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	vendorCode string

	// The deviceBrand of the device, required
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceBrand string

	// The device model, required
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceModel string

	// Serial number of the device, required
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	deviceSerialNumber string

	// An alternate mDNS service identifier
	// Optional, if not set will be generated using "Brand-Model-SerialNumber"
	// Used for mDNS service and SHIP identifier: SHIP - Specification 7.2
	alternateIdentifier string

	// An alternate SHIP identifier
	// Optional, if not set will be identical to alternateIdentifier or generated using "Brand-Model-SerialNumber"
	// Overwrites alternateIdentifier
	// Used for SHIP identifier: SHIP - Specification 7.2
	alternateShipIdentifier string

	// SPINE device type of the device model, required
	// Used for SPINE device type
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceType model.DeviceTypeType

	// SPINE device network feature set type, optional
	// SPINE Protocol Specification 6
	featureSet model.NetworkManagementFeatureSetType

	// Network interface to use for the service
	// Optional, if not set all detected interfaces will be used
	interfaces []string

	// The port address of the websocket server, required
	port int

	// The certificate used for the service and its connections, required
	certificate tls.Certificate

	// Wether remote devices should be automatically accepted
	// If enabled will automatically search for other services with
	// the same setting and automatically connect to them.
	// Has to be set on configuring the service!
	// TODO: if disabled, user verification needs to be implemented and supported
	// the spec defines that this should have a timeout and be activate
	// e.g via a physical button
	registerAutoAccept bool
}

// Setup a ServiceDescription with the required parameters
func NewServiceDescription(
	vendorCode,
	deviceBrand,
	deviceModel,
	serialNumber string,
	deviceType model.DeviceTypeType,
	port int,
	certificate tls.Certificate,
) (*ServiceDescription, error) {
	serviceDescription := &ServiceDescription{
		certificate: certificate,
		port:        port,
	}

	isRequired := "is required"

	if len(vendorCode) == 0 {
		return nil, fmt.Errorf("vendorCode %s", isRequired)
	} else {
		serviceDescription.vendorCode = vendorCode
	}
	if len(deviceBrand) == 0 {
		return nil, fmt.Errorf("brand %s", isRequired)
	} else {
		serviceDescription.deviceBrand = deviceBrand
	}
	if len(deviceModel) == 0 {
		return nil, fmt.Errorf("model %s", isRequired)
	} else {
		serviceDescription.deviceModel = deviceModel
	}
	if len(serialNumber) == 0 {
		return nil, fmt.Errorf("serialNumber %s", isRequired)
	} else {
		serviceDescription.deviceSerialNumber = serialNumber
	}
	if len(deviceType) == 0 {
		return nil, fmt.Errorf("deviceType %s", isRequired)
	} else {
		serviceDescription.deviceType = deviceType
	}

	// set default
	serviceDescription.featureSet = model.NetworkManagementFeatureSetTypeSmart

	return serviceDescription, nil
}

// define an alternative mDNS and SHIP identifier
// usually this is only used when no deviceCode is available or identical to the brand
// if this is not set, generated identifier is used
func (s *ServiceDescription) SetAlternateIdentifier(identifier string) {
	s.alternateIdentifier = identifier
}

// define an alternative identifier to be used for SHIP
// usually this is only used when no deviceCode is available or identical to the brand
//
// will overwrite the alternateIdentifier for the SHIP id, if that is set
// if this is not set, alternateIdentifier or generated identifier is used
func (s *ServiceDescription) SetAlternateShipIdentifier(identifier string) {
	s.alternateShipIdentifier = identifier
}

// define which network interfaces should be considered instead of all existing
// expects a list of network interface names
func (s *ServiceDescription) SetInterfaces(ifaces []string) {
	s.interfaces = ifaces
}

// define wether this service should announce auto accept
// TODO: this needs to be redesigned!
func (s *ServiceDescription) SetRegisterAutoAccept(auto bool) {
	s.registerAutoAccept = auto
}

// generates a standard identifier used for mDNS ID and SHIP ID
// Brand-Model-SerialNumber
func (s *ServiceDescription) generateIdentifier() string {
	return fmt.Sprintf("%s-%s-%s", s.deviceBrand, s.deviceModel, s.deviceSerialNumber)
}

// return the identifier to be used for mDNS
// returns in this order:
// - alternateIdentifier
// - generateIdentifier
func (s *ServiceDescription) mDNSIdentifier() string {
	// SHIP identifier is identical to the mDNS ID
	if len(s.alternateIdentifier) > 0 {
		return s.alternateIdentifier
	}

	return s.generateIdentifier()
}

// return the identifier to be used for mDNS
// returns in this order:
// - alternateShipIdentifier
// - alternateIdentifier
// - generateIdentifier
func (s *ServiceDescription) shipIdentifier() string {
	if len(s.alternateShipIdentifier) > 0 {
		return s.alternateShipIdentifier
	}

	return s.mDNSIdentifier()
}

type EEBUSServiceDelegate interface {
	// RemoteServicesListUpdated(services []ServiceDetails)

	// handle a request to trust a remote service
	RemoteServiceTrustRequested(service *EEBUSService, ski string)

	// report the Ship ID of a newly trusted connection
	RemoteServiceShipIDReported(service *EEBUSService, ski string, shipID string)

	// report a connection to a SKI
	RemoteSKIConnected(service *EEBUSService, ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(service *EEBUSService, ski string)
}

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSService struct {
	ServiceDescription *ServiceDescription

	// The local service details
	LocalService *ServiceDetails

	// Connection Registrations
	connectionsHub *connectionsHub

	// The SPINE specific device definition
	spineLocalDevice *spine.DeviceLocalImpl

	serviceDelegate EEBUSServiceDelegate

	startOnce sync.Once
}

// creates a new EEBUS service
func NewEEBUSService(ServiceDescription *ServiceDescription, serviceDelegate EEBUSServiceDelegate) *EEBUSService {
	return &EEBUSService{
		ServiceDescription: ServiceDescription,
		serviceDelegate:    serviceDelegate,
	}
}

// Sets a custom logging implementation
// By default NoLogging is used, so no logs are printed
func (s *EEBUSService) SetLogging(logger logging.Logging) {
	if logger == nil {
		return
	}
	logging.Log = logger
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Setup() error {
	if s.ServiceDescription.port == 0 {
		s.ServiceDescription.port = defaultPort
	}

	sd := s.ServiceDescription

	leaf, err := x509.ParseCertificate(sd.certificate.Certificate[0])
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	ski, err := skiFromCertificate(leaf)
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	s.LocalService = &ServiceDetails{
		SKI:                ski,
		ShipID:             sd.shipIdentifier(),
		deviceType:         sd.deviceType,
		registerAutoAccept: sd.registerAutoAccept,
	}

	logging.Log.Info("Local SKI: ", ski)

	vendor := sd.vendorCode
	if vendor == "" {
		vendor = sd.deviceBrand
	}

	serial := sd.deviceSerialNumber
	if serial != "" {
		serial = fmt.Sprintf("-%s", serial)
	}

	// Create the SPINE device address, according to Protocol Specification 7.1.1.2
	deviceAdress := fmt.Sprintf("d:_i:%s_%s%s", vendor, sd.deviceModel, serial)

	// Create the local SPINE device
	s.spineLocalDevice = spine.NewDeviceLocalImpl(
		sd.deviceBrand,
		sd.deviceModel,
		sd.deviceSerialNumber,
		sd.shipIdentifier(),
		deviceAdress,
		sd.deviceType,
		sd.featureSet,
	)

	// Create the device entity and add it to the SPINE device
	entityAddress := []model.AddressEntityType{1}
	var entityType model.EntityTypeType
	switch sd.deviceType {
	case model.DeviceTypeTypeEnergyManagementSystem:
		entityType = model.EntityTypeTypeCEM
	case model.DeviceTypeTypeChargingStation:
		entityType = model.EntityTypeTypeEVSE
	default:
		logging.Log.Errorf("Unknown device type: %s", sd.deviceType)
	}
	entity := spine.NewEntityLocalImpl(s.spineLocalDevice, entityType, entityAddress)
	s.spineLocalDevice.AddEntity(entity)

	// Setup connections hub with mDNS and websocket connection handling
	hub, err := newConnectionsHub(s.ServiceDescription, s.LocalService, s)
	if err != nil {
		return err
	}

	s.connectionsHub = hub

	return nil
}

// Starts the service
func (s *EEBUSService) Start() {
	s.startOnce.Do(func() {
		s.connectionsHub.start()
	})
}

// Shutdown all services and stop the server.
func (s *EEBUSService) Shutdown() {
	// Shut down all running connections
	s.connectionsHub.shutdown()
}

func (s *EEBUSService) LocalDevice() *spine.DeviceLocalImpl {
	return s.spineLocalDevice
}

// return the local entity 1
func (s *EEBUSService) LocalEntity() *spine.EntityLocalImpl {
	return s.spineLocalDevice.Entity([]model.AddressEntityType{1})
}

// Add a new entity, used for connected EVs
// Only for EVSE implementations
func (s *EEBUSService) AddEntity(entity *spine.EntityLocalImpl) {
	s.spineLocalDevice.AddEntity(entity)
}

// Remove an entity, used for disconnected EVs
// Only for EVSE implementations
func (s *EEBUSService) RemoveEntity(entity *spine.EntityLocalImpl) {
	s.spineLocalDevice.RemoveEntity(entity)
}

// return all remote devices
func (s *EEBUSService) RemoteDevices() []*spine.DeviceRemoteImpl {
	return s.spineLocalDevice.RemoteDevices()
}

func (s *EEBUSService) RemoteDeviceForSki(ski string) *spine.DeviceRemoteImpl {
	return s.spineLocalDevice.RemoteDeviceForSki(ski)
}

// return a specific remote device of a given DeviceType
func (s *EEBUSService) RemoteDeviceOfType(deviceType model.DeviceTypeType) *spine.DeviceRemoteImpl {
	for _, device := range s.spineLocalDevice.RemoteDevices() {
		if *device.DeviceType() == deviceType {
			return device
		}
	}
	return nil
}

// Adds a new device to the list of known devices which can be connected to
// and connect it if it is currently not connected
func (s *EEBUSService) RegisterRemoteService(service ServiceDetails) {
	s.connectionsHub.registerRemoteService(service)
}

// Returns if the provided SKI is from a registered service
func (s *EEBUSService) IsRemoteServiceRegisteredForSKI(ski string) bool {
	if _, err := s.connectionsHub.registeredServiceForSKI(ski); err != nil {
		return false
	}

	return true
}

// Remove a device from the list of known devices which can be connected to
// and disconnect it if it is currently connected
func (s *EEBUSService) UnregisterRemoteService(ski string) error {
	return s.connectionsHub.unregisterRemoteService(ski)
}

// Mark a remote service to be trusted or not
// Should also be called if the user can't somehow choose to trust,
// e.g. if the UI disappeared
func (s *EEBUSService) UpdateRemoteServiceTrust(ski string, trusted bool) {
	s.connectionsHub.updateRemoteServiceTrust(ski, trusted)
}

// ConnectionHandlerDelegate
func (s *EEBUSService) requestUserTrustForService(details *ServiceDetails) {
	s.serviceDelegate.RemoteServiceTrustRequested(s, details.SKI)
}

func (s *EEBUSService) shipIDUpdateForService(details *ServiceDetails) {
	s.serviceDelegate.RemoteServiceShipIDReported(s, details.SKI, details.ShipID)
}

func (s *EEBUSService) addRemoteDeviceConnection(ski string, readC <-chan []byte, writeC chan<- []byte) {
	s.spineLocalDevice.AddRemoteDevice(ski, readC, writeC)
	s.serviceDelegate.RemoteSKIConnected(s, ski)
}

func (s *EEBUSService) removeRemoteDeviceConnection(ski string) {
	remoteDevice := s.spineLocalDevice.RemoteDeviceForSki(ski)

	s.spineLocalDevice.RemoveRemoteDevice(ski)
	s.serviceDelegate.RemoteSKIDisconnected(s, ski)

	// inform about the disconnection
	payload := spine.EventPayload{
		Ski:        ski,
		EventType:  spine.EventTypeDeviceChange,
		ChangeType: spine.ElementChangeRemove,
		Device:     remoteDevice,
	}
	spine.Events.Publish(payload)
}

// Close a connection to a remote SKI
func (s *EEBUSService) DisconnectSKI(ski string) {
	s.connectionsHub.disconnectSKI(ski)
}
