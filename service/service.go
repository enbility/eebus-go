package service

import (
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// interface for receiving data for specific events
type EEBUSServiceHandler interface {
	// RemoteServicesListUpdated(services []ServiceDetails)

	// report a connection to a SKI
	RemoteSKIConnected(service *EEBUSService, ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(service *EEBUSService, ski string)

	// Provides the SHIP ID the remote service reported during the handshake process
	// This needs to be persisted and passed on for future remote service connections
	// when using `PairRemoteService`
	ReportServiceShipID(ski string, shipdID string)
}

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSService struct {
	Configuration *Configuration

	// The local service details
	LocalService *ServiceDetails

	// Connection Registrations
	connectionsHub *connectionsHub

	// The SPINE specific device definition
	spineLocalDevice *spine.DeviceLocalImpl

	serviceHandler EEBUSServiceHandler

	startOnce sync.Once
}

// creates a new EEBUS service
func NewEEBUSService(configuration *Configuration, serviceHandler EEBUSServiceHandler) *EEBUSService {
	return &EEBUSService{
		Configuration:  configuration,
		serviceHandler: serviceHandler,
	}
}

var _ serviceProvider = (*EEBUSService)(nil)

// report a connection to a SKI
func (s *EEBUSService) RemoteSKIConnected(ski string) {
	s.serviceHandler.RemoteSKIConnected(s, ski)
}

// report a disconnection to a SKI
func (s *EEBUSService) RemoteSKIDisconnected(ski string) {
	s.serviceHandler.RemoteSKIDisconnected(s, ski)
}

// Provides the SHIP ID the remote service reported during the handshake process
func (s *EEBUSService) ReportServiceShipID(ski string, shipdID string) {
	s.serviceHandler.ReportServiceShipID(ski, shipdID)
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
	if s.Configuration.port == 0 {
		s.Configuration.port = defaultPort
	}

	sd := s.Configuration

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

	// Initialize the local service
	// The ShipID is defined in SHIP Spec 3. as
	//   Each SHIP node has a globally unique SHIP ID. The SHIP ID is used to uniquely identify a SHIP node,
	//   e.g. in its service discovery. This ID is present in the mDNS/DNS-SD local service discovery;
	// In SHIP 13.4.6.2 the accessMethods.id is defined as
	//   The originator's unique ID
	// I assume those two to mean the same.
	// TODO: clarify
	s.LocalService = NewServiceDetails(ski)
	s.LocalService.SetShipID(sd.Identifier())
	s.LocalService.SetDeviceType(sd.deviceType)
	s.LocalService.SetRegisterAutoAccept(sd.registerAutoAccept)

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
		sd.Identifier(),
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
	s.connectionsHub = newConnectionsHub(s, s.spineLocalDevice, s.Configuration, s.LocalService)

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
func (s *EEBUSService) PairRemoteService(service *ServiceDetails) {
	s.connectionsHub.PairRemoteService(service)
}

// Returns if the provided SKI is from a registered service
func (s *EEBUSService) IsRemoteServiceForSKIPaired(ski string) bool {
	return s.connectionsHub.IsRemoteServiceForSKIPaired(ski)
}

// Remove a device from the list of known devices which can be connected to
// and disconnect it if it is currently connected
func (s *EEBUSService) UnpairRemoteService(ski string) error {
	return s.connectionsHub.UnpairRemoteService(ski)
}

// Close a connection to a remote SKI
func (s *EEBUSService) DisconnectSKI(ski string, reason string) {
	s.connectionsHub.DisconnectSKI(ski, reason)
}
