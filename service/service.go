package service

import (
	"crypto/x509"
	"errors"
	"fmt"
	"sync"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/cert"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSServiceImpl struct {
	Configuration *api.Configuration

	// The local service details
	LocalService *api.ServiceDetails

	// Connection Registrations
	connectionsHub api.ConnectionsHub

	// The SPINE specific device definition
	spineLocalDevice spineapi.DeviceLocal

	serviceHandler api.EEBUSServiceHandler

	startOnce sync.Once
}

// creates a new EEBUS service
func NewEEBUSService(configuration *api.Configuration, serviceHandler api.EEBUSServiceHandler) *EEBUSServiceImpl {
	return &EEBUSServiceImpl{
		Configuration:  configuration,
		serviceHandler: serviceHandler,
	}
}

var _ api.ServiceProvider = (*EEBUSServiceImpl)(nil)

func (s *EEBUSServiceImpl) VisibleMDNSRecordsUpdated(entries []*api.MdnsEntry) {
	var remoteServices []api.RemoteService

	for _, entry := range entries {
		remoteService := api.RemoteService{
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
func (s *EEBUSServiceImpl) RemoteSKIConnected(ski string) {
	s.serviceHandler.RemoteSKIConnected(s, ski)
}

// report a disconnection to a SKI
func (s *EEBUSServiceImpl) RemoteSKIDisconnected(ski string) {
	s.serviceHandler.RemoteSKIDisconnected(s, ski)
}

// Provides the SHIP ID the remote service reported during the handshake process
func (s *EEBUSServiceImpl) ServiceShipIDUpdate(ski string, shipdID string) {
	s.serviceHandler.ServiceShipIDUpdate(ski, shipdID)
}

// Provides the current pairing state for the remote service
// This is called whenever the state changes and can be used to
// provide user information for the pairing/connection process
func (s *EEBUSServiceImpl) ServicePairingDetailUpdate(ski string, detail *api.ConnectionStateDetail) {
	s.serviceHandler.ServicePairingDetailUpdate(ski, detail)
}

// return if the user is still able to trust the connection
func (s *EEBUSServiceImpl) AllowWaitingForTrust(ski string) bool {
	return s.serviceHandler.AllowWaitingForTrust(ski)
}

var _ api.EEBUSService = (*EEBUSServiceImpl)(nil)

// Get the current pairing details for a given SKI
func (s *EEBUSServiceImpl) PairingDetailForSki(ski string) *api.ConnectionStateDetail {
	return s.connectionsHub.PairingDetailForSki(ski)
}

// Starts browsing for any EEBUS mDNS entry
func (s *EEBUSServiceImpl) StartBrowseMdnsEntries() {
	s.connectionsHub.StartBrowseMdnsSearch()
}

// Stop brwosing for any EEBUS mDNS entry
func (s *EEBUSServiceImpl) StopBrowseMdnsEntries() {
	s.connectionsHub.StopBrowseMdnsSearch()
}

// Sets a custom logging implementation
// By default NoLogging is used, so no logs are printed
func (s *EEBUSServiceImpl) SetLogging(logger logging.Logging) {
	if logger == nil {
		return
	}
	logging.SetLogging(logger)
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSServiceImpl) Setup() error {
	sd := s.Configuration

	if len(sd.Certificate().Certificate) == 0 {
		return errors.New("missing certificate")
	}

	leaf, err := x509.ParseCertificate(sd.Certificate().Certificate[0])
	if err != nil {
		return err
	}

	ski, err := cert.SkiFromCertificate(leaf)
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
	s.LocalService = api.NewServiceDetails(ski)
	s.LocalService.ShipID = sd.Identifier()
	s.LocalService.DeviceType = sd.DeviceType()
	s.LocalService.RegisterAutoAccept = sd.RegisterAutoAccept()

	logging.Log().Info("Local SKI: ", ski)

	vendor := sd.VendorCode()
	if vendor == "" {
		vendor = sd.DeviceBrand()
	}

	serial := sd.DeviceSerialNumber()
	if serial != "" {
		serial = fmt.Sprintf("-%s", serial)
	}

	// Create the SPINE device address, according to Protocol Specification 7.1.1.2
	deviceAdress := fmt.Sprintf("d:_i:%s_%s%s", vendor, sd.DeviceModel(), serial)

	// Create the local SPINE device
	s.spineLocalDevice = spine.NewDeviceLocalImpl(
		sd.DeviceBrand(),
		sd.DeviceModel(),
		sd.DeviceSerialNumber(),
		sd.Identifier(),
		deviceAdress,
		sd.DeviceType(),
		sd.FeatureSet(),
		sd.HeartbeatTimeout(),
	)

	// Create the device entities and add it to the SPINE device
	for _, entityType := range sd.EntityTypes() {
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
func (s *EEBUSServiceImpl) Start() {
	s.startOnce.Do(func() {
		s.connectionsHub.Start()
	})
}

// Shutdown all services and stop the server.
func (s *EEBUSServiceImpl) Shutdown() {
	// Shut down all running connections
	s.connectionsHub.Shutdown()
}

func (s *EEBUSServiceImpl) LocalDevice() spineapi.DeviceLocal {
	return s.spineLocalDevice
}

// Returns the Service detail of a given remote SKI
func (s *EEBUSServiceImpl) RemoteServiceForSKI(ski string) *api.ServiceDetails {
	return s.connectionsHub.ServiceForSKI(ski)
}

// Sets the SKI as being paired or not
// and connect it if paired and not currently being connected
func (s *EEBUSServiceImpl) RegisterRemoteSKI(ski string, enable bool) {
	s.connectionsHub.RegisterRemoteSKI(ski, enable)
}

// Triggers the pairing process for a SKI
func (s *EEBUSServiceImpl) InitiatePairingWithSKI(ski string) {
	s.connectionsHub.InitiatePairingWithSKI(ski)
}

// Cancels the pairing process for a SKI
func (s *EEBUSServiceImpl) CancelPairingWithSKI(ski string) {
	s.connectionsHub.CancelPairingWithSKI(ski)
}

// Close a connection to a remote SKI
func (s *EEBUSServiceImpl) DisconnectSKI(ski string, reason string) {
	s.connectionsHub.DisconnectSKI(ski, reason)
}
