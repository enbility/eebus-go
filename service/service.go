package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

const defaultPort int = 4711

type ServiceDetails struct {
	SKI                string
	ShipID             string
	RegisterAutoAccept bool
}

type ServiceDescription struct {
	// The brand of the device
	DeviceBrand string

	// The device model
	DeviceModel string

	// The EEBUS device type of the device model
	DeviceType model.DeviceTypeType

	// Serial number of the device
	DeviceSerialNumber string

	// The mDNS service identifier
	// Optional, if not set will be  generated using "DeviceBrand-DeviceModel-DeviceSerialNumber"
	DeviceIdentifier string

	// The EEBUS device type of supported remote devices
	RemoteDeviceTypes []model.DeviceTypeType

	// Network interface to use for the service
	// Optional, if not set all detected interfaces will be used
	Interfaces []string

	// The port address of the websocket server
	Port int

	// The certificate used for the service and its connections
	Certificate tls.Certificate

	// Wether remote devices should be automatically accepted
	// If enabled will automatically search for other services with
	// the same setting and automatically connect to them.
	// Has to be set on configuring the service!
	// TODO: if disabled, user verification needs to be implemented and supported
	RegisterAutoAccept bool
}

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSService struct {
	serviceDescription *ServiceDescription

	// The local service details
	localService *ServiceDetails

	// Connection Registrations
	connectionsHub *connectionsHub
}

func NewEEBUSService(ServiceDescription *ServiceDescription) *EEBUSService {
	return &EEBUSService{
		serviceDescription: ServiceDescription,
	}
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Start() {
	if s.serviceDescription.Port == 0 {
		s.serviceDescription.Port = defaultPort
	}

	leaf, err := x509.ParseCertificate(s.serviceDescription.Certificate.Certificate[0])
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
		ShipID:             s.serviceDescription.DeviceIdentifier,
		RegisterAutoAccept: s.serviceDescription.RegisterAutoAccept,
	}

	fmt.Println("Local SKI: ", ski)

	s.connectionsHub = newConnectionsHub(s.serviceDescription, s.localService)
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
