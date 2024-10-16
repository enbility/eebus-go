package service

import (
	"crypto/x509"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/hub"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/mdns"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type Service struct {
	configuration *api.Configuration

	// The local service details
	localService *shipapi.ServiceDetails

	// mDNS manager
	mdns shipapi.MdnsInterface

	// Connection Registrations
	connectionsHub shipapi.HubInterface

	// The SPINE specific device definition
	spineLocalDevice spineapi.DeviceLocalInterface

	serviceHandler api.ServiceReaderInterface

	usecases []api.UseCaseInterface

	// defines wether a user interaction to accept pairing is possible
	isPairingPossible bool

	// return if the service is running
	isRunning bool

	mux        sync.Mutex
	muxRunning sync.Mutex
}

// creates a new EEBUS service
func NewService(configuration *api.Configuration, serviceHandler api.ServiceReaderInterface) *Service {
	return &Service{
		configuration:  configuration,
		serviceHandler: serviceHandler,
	}
}

var _ api.ServiceInterface = (*Service)(nil)

// Starts the service by initializeing mDNS and the server.
func (s *Service) Setup() error {
	sd := s.configuration

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
	s.localService = shipapi.NewServiceDetails(ski)
	s.localService.SetShipID(sd.Identifier())
	s.localService.SetDeviceType(string(sd.DeviceType()))

	logging.Log().Info("Local SKI:", ski)

	vendor := sd.VendorCode()
	if vendor == "" {
		vendor = sd.DeviceBrand()
	}

	serial := sd.DeviceSerialNumber()
	if serial != "" {
		serial = fmt.Sprintf("-%s", serial)
	}

	// Create the SPINE device address, according to Protocol Specification 7.1.1.2
	var deviceAddress string
	vendorType := "i"
	if _, err := strconv.Atoi(vendor); err != nil {
		vendorType = "n"
	}
	deviceAddress = fmt.Sprintf("d:_%s:%s_%s%s", vendorType, vendor, sd.DeviceModel(), serial)

	if len(deviceAddress) > 256 {
		return fmt.Errorf("generated device address may not be longer than 256 characters: %s", deviceAddress)
	}

	// Create the local SPINE device
	s.spineLocalDevice = spine.NewDeviceLocal(
		sd.DeviceBrand(),
		sd.DeviceModel(),
		sd.DeviceSerialNumber(),
		sd.Identifier(),
		deviceAddress,
		sd.DeviceType(),
		sd.FeatureSet(),
	)

	// Create the device entities and add it to the SPINE device
	for _, entityType := range sd.EntityTypes() {
		entityAddressId := model.AddressEntityType(len(s.spineLocalDevice.Entities()))
		entityAddress := []model.AddressEntityType{entityAddressId}
		entity := spine.NewEntityLocal(s.spineLocalDevice, entityType, entityAddress, sd.HeartbeatTimeout())
		s.spineLocalDevice.AddEntity(entity)
	}

	// setup mDNS
	s.mdns = mdns.NewMDNS(
		s.localService.SKI(),
		sd.DeviceBrand(),
		sd.DeviceModel(),
		string(sd.DeviceType()),
		sd.DeviceSerialNumber(),
		sd.DeviceCategories(),
		sd.Identifier(),
		sd.MdnsServiceName(),
		sd.Port(),
		sd.Interfaces(),
		sd.MdnsProviderSelection(),
	)

	// Setup connections hub with mDNS and websocket connection handling
	s.connectionsHub = hub.NewHub(s, s.mdns, s.configuration.Port(), s.configuration.Certificate(), s.localService)

	return nil
}

// Starts the service
func (s *Service) Start() {
	s.muxRunning.Lock()
	defer s.muxRunning.Unlock()

	// make sure we do not start twice while the service is already running
	if s.isRunning {
		return
	}

	s.connectionsHub.Start()

	s.isRunning = true
}

// Shutdown all services and stop the server.
func (s *Service) Shutdown() {
	s.muxRunning.Lock()
	defer s.muxRunning.Unlock()

	// if the service is not running, we do not need to shut it down
	if !s.isRunning {
		return
	}

	// Shut down all running connections
	s.connectionsHub.Shutdown()

	s.isRunning = false
}

// return if the service is running
func (s *Service) IsRunning() bool {
	s.muxRunning.Lock()
	defer s.muxRunning.Unlock()

	return s.isRunning
}

// add a use case to the service
func (s *Service) AddUseCase(useCase api.UseCaseInterface) {
	s.usecases = append(s.usecases, useCase)

	useCase.AddFeatures()
	useCase.AddUseCase()
}

func (s *Service) Configuration() *api.Configuration {
	return s.configuration
}

func (s *Service) LocalService() *shipapi.ServiceDetails {
	return s.localService
}

func (s *Service) LocalDevice() spineapi.DeviceLocalInterface {
	return s.spineLocalDevice
}

// Sets a custom logging implementation
// By default NoLogging is used, so no logs are printed
func (s *Service) SetLogging(logger logging.LoggingInterface) {
	if logger == nil {
		return
	}
	logging.SetLogging(logger)
}

// Get the current pairing details for a given SKI
func (s *Service) PairingDetailForSki(ski string) *shipapi.ConnectionStateDetail {
	return s.connectionsHub.PairingDetailForSki(ski)
}

// Returns the Service detail of a given remote SKI
func (s *Service) RemoteServiceForSKI(ski string) *shipapi.ServiceDetails {
	return s.connectionsHub.ServiceForSKI(ski)
}

func (s *Service) SetAutoAccept(value bool) {
	s.localService.SetAutoAccept(value)
	s.connectionsHub.SetAutoAccept(value)
}

func (s *Service) IsAutoAcceptEnabled() bool {
	return s.localService.AutoAccept()
}

// Returns the QR code text for the service
// as defined in SHIP Requirements for Installation Process V1.0.0
func (s *Service) QRCodeText() string {
	return s.mdns.QRCodeText()
}

// Sets the SKI as being paired
// and connect it if paired and not currently being connected
func (s *Service) RegisterRemoteSKI(ski string) {
	s.connectionsHub.RegisterRemoteSKI(ski)
}

// Sets the SKI as not being paired
// and disconnects it if connected
func (s *Service) UnregisterRemoteSKI(ski string) {
	s.connectionsHub.UnregisterRemoteSKI(ski)
}

// Close a connection to a remote SKI
func (s *Service) DisconnectSKI(ski string, reason string) {
	s.connectionsHub.DisconnectSKI(ski, reason)
}

// Cancels the pairing process for a SKI
func (s *Service) CancelPairingWithSKI(ski string) {
	s.connectionsHub.CancelPairingWithSKI(ski)
}

// Define wether the user is able to react to an incoming pairing request
//
// Call this with `true` e.g. if the user is currently using a web interface
// where an incoming request can be accepted or denied
//
// Default is set to false, meaning every incoming pairing request will be
// automatically denied
func (s *Service) UserIsAbleToApproveOrCancelPairingRequests(allow bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.isPairingPossible = allow
}
