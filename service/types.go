package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

const defaultPort int = 4711

// connection state for global usage, e.g. UI
type ConnectionState uint

const (
	ConnectionStateNone                   ConnectionState = iota // The initial state, when no pairing exists
	ConnectionStateQueued                                        // The connection request has been started and is pending connection initialization
	ConnectionStateInitiated                                     // This service initiated the connection process
	ConnectionStateReceivedPairingRequest                        // A remote service initiated the connection process
	ConnectionStateInProgress                                    // The connection handshake is in progress
	ConnectionStateTrusted                                       // The connection is trusted on both ends
	ConnectionStatePin                                           // PIN processing, not supported right now!
	ConnectionStateCompleted                                     // The connection handshake is completed from both ends
	ConnectionStateRemoteDeniedTrust                             // The remote service denied trust
	ConnectionStateError                                         // The connection handshake resulted in an error
)

// the connection state of a service and error if applicable
type ConnectionStateDetail struct {
	State ConnectionState
	Error error
}

// generic service details about the local or any remote service
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

	// shipID is the SHIP identifier of the service
	// This needs to be persisted
	ShipID string

	// The EEBUS device type of the device model
	DeviceType model.DeviceTypeType

	// Flags if the service auto auto accepts other services
	RegisterAutoAccept bool

	// Flags if the service is trusted and should be reconnected to
	// Should be enabled after the connection process resulted
	// ConnectionStateDetail == ConnectionStateTrusted the first time
	Trusted bool

	// the current connection state details
	ConnectionStateDetail ConnectionStateDetail
}

// create a new ServiceDetails record with a SKI
func NewServiceDetails(ski string) *ServiceDetails {
	service := &ServiceDetails{
		SKI: util.NormalizeSKI(ski), // standardize the provided SKI strings
	}

	return service
}

// defines requires meta information about this service
type Configuration struct {
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

	// An alternate mDNS service name
	// Optional, if not set will be identical to alternateIdentifier or generated using "Brand-Model-SerialNumber"
	alternateMdnsServiceName string

	// SPINE device type of the device model, required
	// Used for SPINE device type
	// Used for mDNS txt record: SHIP - Specification 7.3.2
	deviceType model.DeviceTypeType

	// SPINE device network feature set type, optional
	// SPINE Protocol Specification 6
	featureSet model.NetworkManagementFeatureSetType

	// SPINE entity types for each entity that should be created
	// Each entity has to have a different type!
	entityTypes []model.EntityTypeType

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

	// The sites grid voltage
	// This is useful when e.g. power values are not available and therefor
	// need to be calculated using the current values
	voltage float64

	// The timeout to be used for sending heartbeats
	heartbeatTimeout time.Duration
}

// Setup a Configuration with the required parameters
func NewConfiguration(
	vendorCode,
	deviceBrand,
	deviceModel,
	serialNumber string,
	deviceType model.DeviceTypeType,
	entityTypes []model.EntityTypeType,
	port int,
	certificate tls.Certificate,
	voltage float64,
	heartbeatTimeout time.Duration,
) (*Configuration, error) {
	configuration := &Configuration{
		certificate:      certificate,
		port:             port,
		voltage:          voltage,
		heartbeatTimeout: heartbeatTimeout,
	}

	isRequired := "is required"

	if len(vendorCode) == 0 {
		return nil, fmt.Errorf("vendorCode %s", isRequired)
	} else {
		configuration.vendorCode = vendorCode
	}
	if len(deviceBrand) == 0 {
		return nil, fmt.Errorf("brand %s", isRequired)
	} else {
		configuration.deviceBrand = deviceBrand
	}
	if len(deviceModel) == 0 {
		return nil, fmt.Errorf("model %s", isRequired)
	} else {
		configuration.deviceModel = deviceModel
	}
	if len(serialNumber) == 0 {
		return nil, fmt.Errorf("serialNumber %s", isRequired)
	} else {
		configuration.deviceSerialNumber = serialNumber
	}
	if len(deviceType) == 0 {
		return nil, fmt.Errorf("deviceType %s", isRequired)
	} else {
		configuration.deviceType = deviceType
	}
	if len(entityTypes) == 0 {
		return nil, fmt.Errorf("entityTypes %s", isRequired)
	} else {
		configuration.entityTypes = entityTypes
	}
	// set default
	configuration.featureSet = model.NetworkManagementFeatureSetTypeSmart

	return configuration, nil
}

// define an alternative mDNS and SHIP identifier
// usually this is only used when no deviceCode is available or identical to the brand
// if this is not set, generated identifier is used
func (s *Configuration) SetAlternateIdentifier(identifier string) {
	s.alternateIdentifier = identifier
}

// define an alternative mDNS service name
// this is normally not needed or used
func (s *Configuration) SetAlternateMdnsServiceName(name string) {
	s.alternateMdnsServiceName = name
}

// define which network interfaces should be considered instead of all existing
// expects a list of network interface names
func (s *Configuration) SetInterfaces(ifaces []string) {
	s.interfaces = ifaces
}

// define wether this service should announce auto accept
// TODO: this needs to be redesigned!
func (s *Configuration) SetRegisterAutoAccept(auto bool) {
	s.registerAutoAccept = auto
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

// return the sites predefined grid voltage
func (s *Configuration) Voltage() float64 {
	return s.voltage
}

// ErrServiceNotPaired if the given SKI is not paired yet
var ErrServiceNotPaired = errors.New("the provided SKI is not paired")

// ErrConnectionNotFound that there was no active connection for a given SKI found
var ErrConnectionNotFound = errors.New("no connection for provided SKI found")
