package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

const defaultPort int = 4711

type ServiceDetails struct {
	// This is the SKI of the service
	// This needs to be peristed
	SKI string

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
	// The brand of the device, required
	Brand string

	// The device model, required
	Model string

	// The EEBUS device type of the device model, required
	DeviceType model.DeviceTypeType

	// Serial number of the device, required
	SerialNumber string

	// The mDNS service identifier, will also be used as SHIP ID
	// Optional, if not set will be  generated using "Brand-Model-SerialNumber"
	Identifier string

	// The vendors IANA PEN, optional
	VendorCode string

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
	RegisterAutoAccept bool
}

type EEBUSServiceDelegate interface {
	// RemoteServicesListUpdated(services []ServiceDetails)

	// handle a request to trust a remote service
	RemoteServiceTrustRequested(ski string)

	// report the Ship ID of a newly trusted connection
	RemoteServiceShipIDReported(ski string, shipID string)
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

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Setup() error {
	if s.ServiceDescription.Port == 0 {
		s.ServiceDescription.Port = defaultPort
	}

	sd := s.ServiceDescription

	leaf, err := x509.ParseCertificate(sd.Certificate.Certificate[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	ski, err := skiFromCertificate(leaf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.LocalService = &ServiceDetails{
		SKI:                ski,
		ShipID:             sd.Identifier,
		deviceType:         sd.DeviceType,
		registerAutoAccept: sd.RegisterAutoAccept,
	}

	fmt.Println("Local SKI: ", ski)

	vendor := sd.VendorCode
	if vendor == "" {
		vendor = sd.Brand
	}

	// Create the SPINE device address, according to Protocol Specification 7.1.1.2
	deviceAdress := fmt.Sprintf("d:_i:%s_%s%s-%s", vendor, sd.Brand, sd.Model, sd.SerialNumber)

	// Create the local SPINE device
	s.spineLocalDevice = spine.NewDeviceLocalImpl(
		sd.Brand,
		sd.Model,
		sd.Identifier,
		sd.SerialNumber,
		deviceAdress,
		sd.DeviceType,
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
		return errors.New(fmt.Sprintf("Unknown device type: %s", sd.DeviceType))
	}
	entity := spine.NewEntityLocalImpl(s.spineLocalDevice, entityType, entityAddress)
	s.spineLocalDevice.AddEntity(entity)

	// Setup connections hub with mDNS and websocket connection handling
	s.connectionsHub = newConnectionsHub(s.ServiceDescription, s.LocalService, s)

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
func (s *EEBUSService) AddEntity(entity *spine.EntityLocalImpl) {
	s.spineLocalDevice.AddEntity(entity)
}

// Remove an entity, used for disconnected EVs
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
func (s *EEBUSService) RegisterRemoteService(service ServiceDetails) error {
	return s.connectionsHub.registerRemoteService(service)
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
}

func (s *EEBUSService) removeRemoteDeviceConnection(ski string) {
	s.spineLocalDevice.RemoveRemoteDevice(ski)
}
