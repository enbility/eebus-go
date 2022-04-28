package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/entity"
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
	DeviceBrand string

	// The device model, required
	DeviceModel string

	// The EEBUS device type of the device model, required
	DeviceType model.DeviceTypeType

	// Serial number of the device, required
	DeviceSerialNumber string

	// The mDNS service identifier
	// Optional, if not set will be  generated using "DeviceBrand-DeviceModel-DeviceSerialNumber"
	DeviceIdentifier string

	// The vendors IANA PEN, optional
	IANAPEN string

	// The EEBUS device type of supported remote devices, required
	RemoteDeviceTypes []model.DeviceTypeType

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
	serviceDescription *ServiceDescription

	// The local service details
	localService *ServiceDetails

	// Connection Registrations
	connectionsHub *connectionsHub

	// The SPINE specific device definition
	spineLocalDevice *spine.DeviceLocalImpl

	serviceDelegate EEBUSServiceDelegate
}

func NewEEBUSService(ServiceDescription *ServiceDescription, serviceDelegate EEBUSServiceDelegate) *EEBUSService {
	return &EEBUSService{
		serviceDescription: ServiceDescription,
		serviceDelegate:    serviceDelegate,
	}
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Start() {
	if s.serviceDescription.Port == 0 {
		s.serviceDescription.Port = defaultPort
	}

	sd := s.serviceDescription

	leaf, err := x509.ParseCertificate(sd.Certificate.Certificate[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	ski, err := skiFromCertificate(leaf)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.localService = &ServiceDetails{
		SKI:                ski,
		ShipID:             sd.DeviceIdentifier,
		deviceType:         sd.DeviceType,
		registerAutoAccept: sd.RegisterAutoAccept,
	}

	fmt.Println("Local SKI: ", ski)

	vendor := sd.IANAPEN
	if vendor == "" {
		vendor = sd.DeviceBrand
	}

	// Create the SPINE device address, according to Protocol Specification 7.1.1.2
	deviceAdress := fmt.Sprintf("d:_i:%s_%s%s-%s", vendor, sd.DeviceBrand, sd.DeviceModel, sd.DeviceSerialNumber)

	s.spineLocalDevice = spine.NewDeviceLocalImpl(
		sd.DeviceBrand,
		sd.DeviceModel,
		deviceAdress,
		sd.DeviceSerialNumber,
		sd.DeviceType,
	)

	if s.localService.deviceType == model.DeviceTypeTypeEnergyManagementSystem {
		e1 := entity.NewCEM(s.spineLocalDevice, []model.AddressEntityType{1})
		s.spineLocalDevice.AddEntity(e1)
	} else {
		e1 := entity.NewEVSE(s.spineLocalDevice, []model.AddressEntityType{1})
		s.spineLocalDevice.AddEntity(e1)
	}

	s.connectionsHub = newConnectionsHub(s.serviceDescription, s.localService, s)
	s.connectionsHub.start()
}

// Shutdown all services and stop the server.
func (s *EEBUSService) Shutdown() {
	// Shut down all running connections
	s.connectionsHub.shutdown()
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

var _ ConnectionHandlerDelegate = (*EEBUSService)(nil)

func (s *EEBUSService) requestUserTrustForService(service *ServiceDetails) {
	s.serviceDelegate.RemoteServiceTrustRequested(service.SKI)
}

func (s *EEBUSService) shipIDUpdateForService(service *ServiceDetails) {
	s.serviceDelegate.RemoteServiceShipIDReported(service.SKI, service.ShipID)
}

func (s *EEBUSService) addRemoteDeviceConnection(ski, deviceCode string, deviceType model.DeviceTypeType, readC <-chan []byte, writeC chan<- []byte) {
	s.spineLocalDevice.AddRemoteDevice(ski, deviceCode, deviceType, readC, writeC)
}

func (s *EEBUSService) removeRemoteDeviceConnection(ski string) {
	s.spineLocalDevice.RemoveRemoteDevice(ski)
}
