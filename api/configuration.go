package api

import (
	"crypto/tls"
	"fmt"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/mdns"
	"github.com/enbility/spine-go/model"
)

const defaultPort int = 4711

// defines requires meta information about this service
type Configuration struct {
	// The vendors IANA PEN, optional but highly recommended.
	//
	// If not set, brand will be used instead
	//
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	vendorCode string

	// The deviceBrand of the device, required (may not be longer than 32 UTF8 characters)
	//
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	//
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceBrand string

	// The device model, required (may not be longer than 32 UTF8 characters)
	//
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	//
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceModel string

	// Serial number of the device, required (may not be longer than 32 UTF8 characters)
	//
	// Used for the Device Address: SPINE - Protocol Specification 7.1.1.2
	//
	// Used for mDNS txt record: SHIP - Requirements for Installation Process V1.0.0
	deviceSerialNumber string

	// Device categories of the device model, required
	//
	// Used for mDNS txt record: SHIP - Requirements for Installation Process V1.0.0
	deviceCategories []shipapi.DeviceCategoryType

	// An alternate mDNS service identifier
	//
	// Optional, if not set will be generated using "Brand-Model-SerialNumber"
	//
	// Used for mDNS service and SHIP identifier: SHIP - Specification 7.2
	alternateIdentifier string

	// An alternate mDNS service name
	//
	// Optional, if not set will be identical to alternateIdentifier or generated using "Brand-Model-SerialNumber"
	alternateMdnsServiceName string

	// SPINE device type of the device model, required
	//
	// Used for SPINE device type
	//
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceType model.DeviceTypeType

	// SPINE device network feature set type, optional
	//
	// SPINE Protocol Specification 6
	featureSet model.NetworkManagementFeatureSetType

	// SPINE entity types for each entity that should automatically be created, can be empty
	//
	// Each entity has to have a different type!
	entityTypes []model.EntityTypeType

	// Network interface to use for the service
	//
	// Optional, if not set all detected interfaces will be used
	interfaces []string

	// The port address of the websocket server, required
	port int

	// The certificate used for the service and its connections, required
	certificate tls.Certificate

	// The timeout to be used for sending heartbeats and applied to all
	// local entities created on setup of the service
	heartbeatTimeout time.Duration

	// Optional set which mDNS providers should be used
	mdnsProviderSelection mdns.MdnsProviderSelection
}

// Setup a Configuration with the required parameters
//
// Parameters:
//   - vendorCode: The vendors IANA PEN, optional but highly recommended.
//   - deviceBrand: The deviceBrand of the device, required (may not be longer than 32 UTF8 characters)
//   - deviceModel: The device model, required (may not be longer than 32 UTF8 characters)
//   - serialNumber: Serial number of the device, required (may not be longer than 32 UTF8 characters)
//   - deviceCategories: Device categories of the device model, required
//   - deviceType: SPINE device type of the device model, required
//   - entityTypes: SPINE entity types for each entity that should automatically be created, can be empty
//   - port: The port address of the websocket server, required
//   - certificate: The certificate used for the service and its connections, required
//   - heartbeatTimeout: The timeout to be used for sending heartbeats and applied to all local entities created on setup of the service
//   - mdnsProviderSelection: Optional set which mDNS providers should be used, default is `mdns.MdnsProviderSelectionAll`
//
// Returns:
//   - *Configuration: The created configuration
//   - error: An error if the configuration could not be created
func NewConfiguration(
	vendorCode,
	deviceBrand,
	deviceModel,
	serialNumber string,
	deviceCategories []shipapi.DeviceCategoryType,
	deviceType model.DeviceTypeType,
	entityTypes []model.EntityTypeType,
	port int,
	certificate tls.Certificate,
	heartbeatTimeout time.Duration,
) (*Configuration, error) {
	configuration := &Configuration{
		certificate:           certificate,
		port:                  port,
		heartbeatTimeout:      heartbeatTimeout,
		mdnsProviderSelection: mdns.MdnsProviderSelectionAll,
	}

	if port == 0 {
		configuration.port = defaultPort
	}

	isRequired := "is required"

	if len(vendorCode) == 0 {
		return nil, fmt.Errorf("vendorCode %s", isRequired)
	}
	configuration.vendorCode = vendorCode

	if len(deviceBrand) == 0 {
		return nil, fmt.Errorf("brand %s", isRequired)
	}
	configuration.deviceBrand = deviceBrand

	if len(deviceModel) == 0 {
		return nil, fmt.Errorf("model %s", isRequired)
	}
	configuration.deviceModel = deviceModel

	if len(serialNumber) == 0 {
		return nil, fmt.Errorf("serialNumber %s", isRequired)
	}
	configuration.deviceSerialNumber = serialNumber

	if len(deviceCategories) == 0 {
		return nil, fmt.Errorf("deviceCategories %s", isRequired)
	}
	configuration.deviceCategories = deviceCategories

	if len(deviceType) == 0 {
		return nil, fmt.Errorf("deviceType %s", isRequired)
	}
	configuration.deviceType = deviceType

	if len(entityTypes) == 0 {
		return nil, fmt.Errorf("entityTypes %s", isRequired)
	}
	configuration.entityTypes = entityTypes

	// set default
	configuration.featureSet = model.NetworkManagementFeatureSetTypeSmart

	return configuration, nil
}

// Returns the configuration vendor code
func (s *Configuration) VendorCode() string {
	return s.vendorCode
}

// Returns the configuration device brand
func (s *Configuration) DeviceBrand() string {
	return s.deviceBrand
}

// Returns the configuration device model
func (s *Configuration) DeviceModel() string {
	return s.deviceModel
}

// Returns the configuration device serial number
func (s *Configuration) DeviceSerialNumber() string {
	return s.deviceSerialNumber
}

// Returns the configuration device categories
func (c *Configuration) DeviceCategories() []shipapi.DeviceCategoryType {
	return c.deviceCategories
}

// set an alternative mDNS and SHIP identifier
// usually this is only used when no deviceCode is available or identical to the brand
// if this is not set, generated identifier is used
func (s *Configuration) SetAlternateIdentifier(identifier string) {
	s.alternateIdentifier = identifier
}

// set an alternative mDNS service name
// this is normally not needed or used
func (s *Configuration) SetAlternateMdnsServiceName(name string) {
	s.alternateMdnsServiceName = name
}

// set the mDNS provider selection
func (s *Configuration) SetMdnsProviderSelection(providerSelection mdns.MdnsProviderSelection) {
	s.mdnsProviderSelection = providerSelection
}

// Returns the mDNS provider selection
func (s *Configuration) MdnsProviderSelection() mdns.MdnsProviderSelection {
	return s.mdnsProviderSelection
}

// Returns the configuration device type
func (s *Configuration) DeviceType() model.DeviceTypeType {
	return s.deviceType
}

// Returns the configuration network management feature set type
func (s *Configuration) FeatureSet() model.NetworkManagementFeatureSetType {
	return s.featureSet
}

// Returns the configuration entity types
func (s *Configuration) EntityTypes() []model.EntityTypeType {
	return s.entityTypes
}

// Returns the configuration network interfaces
func (s *Configuration) Interfaces() []string {
	return s.interfaces
}

// define which network interfaces should be considered instead of all existing
// expects a list of network interface names
func (s *Configuration) SetInterfaces(ifaces []string) {
	s.interfaces = ifaces
}

// generates a standard identifier used for mDNS ID and SHIP ID
// Brand-Model-SerialNumber
func (s *Configuration) generateIdentifier() string {
	return fmt.Sprintf("%s-%s-%s", s.deviceBrand, s.deviceModel, s.deviceSerialNumber)
}

// return the identifier to be used for mDNS and SHIP ID
// returns in this order:
// - alternateIdentifier
// - generateIdentifier
func (s *Configuration) Identifier() string {
	// SHIP identifier is identical to the mDNS ID
	if len(s.alternateIdentifier) > 0 {
		return s.alternateIdentifier
	}

	return s.generateIdentifier()
}

// return the name to be used as the mDNS service name
// returns in this order:
// - alternateMdnsServiceName
// - generateIdentifier
func (s *Configuration) MdnsServiceName() string {
	// SHIP identifier is identical to the mDNS ID
	if len(s.alternateMdnsServiceName) > 0 {
		return s.alternateMdnsServiceName
	}

	return s.generateIdentifier()
}

// Returns the certificate
func (s *Configuration) Certificate() tls.Certificate {
	return s.certificate
}

// Returns the port
func (s *Configuration) Port() int {
	return s.port
}

// Set the certificate
func (s *Configuration) SetCertificate(cert tls.Certificate) {
	s.certificate = cert
}

// Returns the heartbeat timeout
func (s *Configuration) HeartbeatTimeout() time.Duration {
	return s.heartbeatTimeout
}
