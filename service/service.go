package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

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
	VendorCode string

	// The brand of the device, required
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	Brand string

	// The device model, required
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	Model string

	// Serial number of the device, required
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	SerialNumber string

	// An alternate mDNS service identifier
	// Optional, if not set will be  generated using "Brand-Model-SerialNumber"
	// Used for mDNS service identifier: SHIP - Specification 7.2
	AlternateIdentifier string

	// SPINE device type of the device model, required
	// Used for SPINE device type
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	DeviceType model.DeviceTypeType

	// SPINE device network feature set type, optional
	// SPINE Protocol Specification 6
	FeatureSet model.NetworkManagementFeatureSetType

	// Network interface to use for the service
	// Optional, if not set all detected interfaces will be used
	Interfaces []string

	// The port address of the websocket server, required
	Port int

	// The certificate used for the service and its connections, required
	Certificate tls.Certificate

	// Wether remote devices should be automatically accepted
	// If enabled will automatically search for other services with
	// the same setting and automatically connect to them.
	// Has to be set on configuring the service!
	// TODO: if disabled, user verification needs to be implemented and supported
	// the spec defines that this should have a timeout and be activate
	// e.g via a physical button
	RegisterAutoAccept bool
}

// Setup a ServiceDescription with the required parameters
func NewServiceDescription(
	vendorCode,
	brand,
	model,
	serialNumber,
	alternateIdentifier string,
	deviceType model.DeviceTypeType,
	port int,
	certificate tls.Certificate,
) (*ServiceDescription, error) {
	serviceDescription := &ServiceDescription{
		Certificate: certificate,
		Port:        port,
	}

	isRequired := "is required"

	if len(vendorCode) == 0 {
		return nil, fmt.Errorf("vendorCode %s", isRequired)
	} else {
		serviceDescription.VendorCode = vendorCode
	}
	if len(brand) == 0 {
		return nil, fmt.Errorf("brand %s", isRequired)
	} else {
		serviceDescription.Brand = brand
	}
	if len(model) == 0 {
		return nil, fmt.Errorf("model %s", isRequired)
	} else {
		serviceDescription.Model = model
	}
	if len(serialNumber) == 0 {
		return nil, fmt.Errorf("serialNumber %s", isRequired)
	} else {
		serviceDescription.SerialNumber = serialNumber
	}
	if len(deviceType) == 0 {
		return nil, fmt.Errorf("deviceType %s", isRequired)
	} else {
		serviceDescription.DeviceType = deviceType
	}

	return serviceDescription, nil
}

type EEBUSServiceDelegate interface {
	// RemoteServicesListUpdated(services []ServiceDetails)

	// handle a request to trust a remote service
	RemoteServiceTrustRequested(ski string)

	// report the Ship ID of a newly trusted connection
	RemoteServiceShipIDReported(ski string, shipID string)

	// report a connection to a SKI
	RemoteSKIConnected(ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(ski string)
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
}

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
	if s.ServiceDescription.Port == 0 {
		s.ServiceDescription.Port = defaultPort
	}

	sd := s.ServiceDescription

	leaf, err := x509.ParseCertificate(sd.Certificate.Certificate[0])
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	ski, err := skiFromCertificate(leaf)
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	// SHIP identifier is identical to the mDNS ID
	// Brand-Model-SerialNumber or AlternateIdentifier
	shipID := fmt.Sprintf("%s-%s-%s", sd.Brand, sd.Model, sd.SerialNumber)
	if len(sd.AlternateIdentifier) > 0 {
		shipID = sd.AlternateIdentifier
	}

	s.LocalService = &ServiceDetails{
		SKI:                ski,
		ShipID:             shipID,
		deviceType:         sd.DeviceType,
		registerAutoAccept: sd.RegisterAutoAccept,
	}

	logging.Log.Info("Local SKI: ", ski)

	vendor := sd.VendorCode
	if vendor == "" {
		vendor = sd.Brand
	}

	serial := sd.SerialNumber
	if serial != "" {
		serial = fmt.Sprintf("-%s", serial)
	}

	// Create the SPINE device address, according to Protocol Specification 7.1.1.2
	deviceAdress := fmt.Sprintf("d:_i:%s_%s%s", vendor, sd.Model, serial)

	// Create the local SPINE device
	s.spineLocalDevice = spine.NewDeviceLocalImpl(
		sd.Brand,
		sd.Model,
		sd.SerialNumber,
		shipID,
		deviceAdress,
		sd.DeviceType,
		sd.FeatureSet,
	)

	// Create the device entity and add it to the SPINE device
	entityAddress := []model.AddressEntityType{1}
	var entityType model.EntityTypeType
	switch sd.DeviceType {
	case model.DeviceTypeTypeEnergyManagementSystem:
		entityType = model.EntityTypeTypeCEM
	case model.DeviceTypeTypeChargingStation:
		entityType = model.EntityTypeTypeEVSE
	default:
		logging.Log.Errorf("Unknown device type: %s", sd.DeviceType)
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
	s.connectionsHub.start()
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
func (s *EEBUSService) requestUserTrustForService(service *ServiceDetails) {
	s.serviceDelegate.RemoteServiceTrustRequested(service.SKI)
}

func (s *EEBUSService) shipIDUpdateForService(service *ServiceDetails) {
	s.serviceDelegate.RemoteServiceShipIDReported(service.SKI, service.ShipID)
}

func (s *EEBUSService) addRemoteDeviceConnection(ski string, readC <-chan []byte, writeC chan<- []byte) {
	s.spineLocalDevice.AddRemoteDevice(ski, readC, writeC)
	s.serviceDelegate.RemoteSKIConnected(ski)
}

func (s *EEBUSService) removeRemoteDeviceConnection(ski string) {
	remoteDevice := s.spineLocalDevice.RemoteDeviceForSki(ski)

	s.spineLocalDevice.RemoveRemoteDevice(ski)
	s.serviceDelegate.RemoteSKIDisconnected(ski)

	// inform about the disconnection
	payload := spine.EventPayload{
		Ski:        ski,
		EventType:  spine.EventTypeDeviceChange,
		ChangeType: spine.ElementChangeRemove,
		Device:     remoteDevice,
	}
	spine.Events.Publish(payload)
}
