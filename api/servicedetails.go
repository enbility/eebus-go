package api

import (
	"sync"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
)

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
	connectionStateDetail *ConnectionStateDetail

	mux sync.Mutex
}

// create a new ServiceDetails record with a SKI
func NewServiceDetails(ski string) *ServiceDetails {
	connState := NewConnectionStateDetail(ConnectionStateNone, nil)
	service := &ServiceDetails{
		SKI:                   util.NormalizeSKI(ski), // standardize the provided SKI strings
		connectionStateDetail: connState,
	}

	return service
}

func (s *ServiceDetails) ConnectionStateDetail() *ConnectionStateDetail {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.connectionStateDetail
}

func (s *ServiceDetails) SetConnectionStateDetail(detail *ConnectionStateDetail) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.connectionStateDetail = detail
}
