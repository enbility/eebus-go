package service

import (
	"crypto/x509"
	"errors"
	"fmt"
	"sync"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type RemoteService struct {
	Name       string `json:"name"`
	Ski        string `json:"ski"`
	Identifier string `json:"identifier"`
	Brand      string `json:"brand"`
	Type       string `json:"type"`
	Model      string `json:"model"`
}

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSService struct {
	Configuration *Configuration

	// The local service details
	LocalService *ServiceDetails

	// Connection Registrations
	connectionsHub ConnectionsHub

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

var _ ServiceProvider = (*EEBUSService)(nil)

func (s *EEBUSService) VisibleMDNSRecordsUpdated(entries []*MdnsEntry) {
	var remoteServices []RemoteService

	for _, entry := range entries {
		remoteService := RemoteService{
			Name:       entry.Name,
			Ski:        entry.Ski,
			Identifier: entry.Identifier,
			Brand:      entry.Brand,
			Type:       entry.Type,
			Model:      entry.Model,
		}

		remoteServices = append(remoteServices, remoteService)
	}
	s.serviceHandler.VisibleRemoteServicesUpdated(s, remoteServices)
}

// report a connection to a SKI
func (s *EEBUSService) RemoteSKIConnected(ski string) {
	s.serviceHandler.RemoteSKIConnected(s, ski)
}

// report a disconnection to a SKI
func (s *EEBUSService) RemoteSKIDisconnected(ski string) {
	s.serviceHandler.RemoteSKIDisconnected(s, ski)
}

// Provides the SHIP ID the remote service reported during the handshake process
func (s *EEBUSService) ServiceShipIDUpdate(ski string, shipdID string) {
	s.serviceHandler.ServiceShipIDUpdate(ski, shipdID)
}

// Provides the current pairing state for the remote service
// This is called whenever the state changes and can be used to
// provide user information for the pairing/connection process
func (s *EEBUSService) ServicePairingDetailUpdate(ski string, detail *ConnectionStateDetail) {
	s.serviceHandler.ServicePairingDetailUpdate(ski, detail)
}

// return if the user is still able to trust the connection
func (s *EEBUSService) AllowWaitingForTrust(ski string) bool {
	return s.serviceHandler.AllowWaitingForTrust(ski)
}

// Get the current pairing details for a given SKI
func (s *EEBUSService) PairingDetailForSki(ski string) *ConnectionStateDetail {
	return s.connectionsHub.PairingDetailForSki(ski)
}

// Starts browsing for any EEBUS mDNS entry
func (s *EEBUSService) StartBrowseMdnsEntries() {
	s.connectionsHub.StartBrowseMdnsSearch()
}

// Stop brwosing for any EEBUS mDNS entry
func (s *EEBUSService) StopBrowseMdnsEntries() {
	s.connectionsHub.StopBrowseMdnsSearch()
}

// Sets a custom logging implementation
// By default NoLogging is used, so no logs are printed
func (s *EEBUSService) SetLogging(logger logging.Logging) {
	if logger == nil {
		return
	}
	logging.SetLogging(logger)
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Setup() error {
	if s.Configuration.port == 0 {
		s.Configuration.port = defaultPort
	}

	sd := s.Configuration

	if len(sd.certificate.Certificate) == 0 {
		return errors.New("missing certificate")
	}

	leaf, err := x509.ParseCertificate(sd.certificate.Certificate[0])
	if err != nil {
		return err
	}

	ski, err := skiFromCertificate(leaf)
	if err != nil {
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
	s.LocalService.ShipID = sd.Identifier()
	s.LocalService.DeviceType = sd.deviceType
	s.LocalService.RegisterAutoAccept = sd.registerAutoAccept

	logging.Log().Info("Local SKI: ", ski)

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
		sd.heartbeatTimeout,
	)

	// Create the device entities and add it to the SPINE device
	for _, entityType := range sd.entityTypes {
		entityAddressId := model.AddressEntityType(len(s.spineLocalDevice.Entities()))
		entityAddress := []model.AddressEntityType{entityAddressId}
		entity := spine.NewEntityLocalImpl(s.spineLocalDevice, entityType, entityAddress)
		s.spineLocalDevice.AddEntity(entity)
	}

	// setup mDNS
	mdns := newMDNS(s.LocalService.SKI, s.Configuration)

	// Setup connections hub with mDNS and websocket connection handling
	s.connectionsHub = newConnectionsHub(s, mdns, s.spineLocalDevice, s.Configuration, s.LocalService)

	return nil
}

// Starts the service
func (s *EEBUSService) Start() {
	s.startOnce.Do(func() {
		s.connectionsHub.Start()
	})
}

// Shutdown all services and stop the server.
func (s *EEBUSService) Shutdown() {
	// Shut down all running connections
	s.connectionsHub.Shutdown()
}

func (s *EEBUSService) LocalDevice() *spine.DeviceLocalImpl {
	return s.spineLocalDevice
}

// Returns the Service detail of a given remote SKI
func (s *EEBUSService) RemoteServiceForSKI(ski string) *ServiceDetails {
	return s.connectionsHub.ServiceForSKI(ski)
}

// Sets the SKI as being paired or not
// and connect it if paired and not currently being connected
func (s *EEBUSService) RegisterRemoteSKI(ski string, enable bool) {
	s.connectionsHub.RegisterRemoteSKI(ski, enable)
}

// Triggers the pairing process for a SKI
func (s *EEBUSService) InitiatePairingWithSKI(ski string) {
	s.connectionsHub.InitiatePairingWithSKI(ski)
}

// Cancels the pairing process for a SKI
func (s *EEBUSService) CancelPairingWithSKI(ski string) {
	s.connectionsHub.CancelPairingWithSKI(ski)
}

// Close a connection to a remote SKI
func (s *EEBUSService) DisconnectSKI(ski string, reason string) {
	s.connectionsHub.DisconnectSKI(ski, reason)
}
