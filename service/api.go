package service

import "net"

/* EEBUSService */

//go:generate mockgen -destination=mock_service_test.go -package=service github.com/enbility/eebus-go/service EEBUSServiceHandler

// interface for receiving data for specific events
type EEBUSServiceHandler interface {
	// report all currently visible EEBUS services
	VisibleRemoteServicesUpdated(service *EEBUSService, entries []RemoteService)

	// report a connection to a SKI
	RemoteSKIConnected(service *EEBUSService, ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(service *EEBUSService, ski string)

	// Provides the SHIP ID the remote service reported during the handshake process
	// This needs to be persisted and passed on for future remote service connections
	// when using `PairRemoteService`
	ServiceShipIDUpdate(ski string, shipdID string)

	// Provides the current pairing state for the remote service
	// This is called whenever the state changes and can be used to
	// provide user information for the pairing/connection process
	ServicePairingDetailUpdate(ski string, detail *ConnectionStateDetail)

	// return if the user is still able to trust the connection
	AllowWaitingForTrust(ski string) bool
}

/* Hub */

//go:generate mockgen -destination=mock_hub_test.go -package=service github.com/enbility/eebus-go/service ServiceProvider,ConnectionsHub

// interface for reporting data from connectionsHub to the EEBUSService
type ServiceProvider interface {
	// report a newly discovered remote EEBUS service
	VisibleMDNSRecordsUpdated(entries []*MdnsEntry)

	// report a connection to a SKI
	RemoteSKIConnected(ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(ski string)

	// provide the SHIP ID received during SHIP handshake process
	// the ID needs to be stored and then provided for remote services so it can be compared and verified
	ServiceShipIDUpdate(ski string, shipID string)

	// provides the current handshake state for a given SKI
	ServicePairingDetailUpdate(ski string, detail *ConnectionStateDetail)

	// return if the user is still able to trust the connection
	AllowWaitingForTrust(ski string) bool
}

type ConnectionsHub interface {
	PairingDetailForSki(ski string) *ConnectionStateDetail
	StartBrowseMdnsSearch()
	StopBrowseMdnsSearch()
	Start()
	Shutdown()
	ServiceForSKI(ski string) *ServiceDetails
	RegisterRemoteSKI(ski string, enable bool)
	InitiatePairingWithSKI(ski string)
	CancelPairingWithSKI(ski string)
	DisconnectSKI(ski string, reason string)
}

/* Mdns */

//go:generate mockgen -destination=mock_mdns_test.go -package=service github.com/enbility/eebus-go/service MdnsSearch,MdnsService

// implemented by hubConnection, used by mdns
type MdnsSearch interface {
	ReportMdnsEntries(entries map[string]*MdnsEntry)
}

// implemented by mdns, used by hubConnection
type MdnsService interface {
	SetupMdnsService() error
	ShutdownMdnsService()
	AnnounceMdnsEntry() error
	UnannounceMdnsEntry()
	RegisterMdnsSearch(cb MdnsSearch)
	UnregisterMdnsSearch(cb MdnsSearch)
}

//go:generate mockery --name=MdnsProvider

type MdnsProvider interface {
	CheckAvailability() bool
	Shutdown()
	Announce(serviceName string, port int, txt []string) error
	Unannounce()
	ResolveEntries(cancelChan chan bool, callback func(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool))
}
