package service

import (
	shipapi "github.com/enbility/ship-go/api"
)

var _ shipapi.HubReaderInterface = (*Service)(nil)

// report a connection to a SKI
//
// is triggered whenever a SHIP connected was successful completed
func (s *Service) RemoteSKIConnected(ski string) {
	s.serviceHandler.RemoteSKIConnected(s, ski)
}

// report a disconnection to a SKI
//
// is triggered whenever a SHIP connect was closed, is also triggered when the SHIP
// process wasn't successfully completed
//
// NOTE: The connection may not have been reported as connected before!
func (s *Service) RemoteSKIDisconnected(ski string) {
	if s.spineLocalDevice != nil {
		s.spineLocalDevice.RemoveRemoteDeviceConnection(ski)
	}

	s.serviceHandler.RemoteSKIDisconnected(s, ski)
}

// report an approved handshake by a remote device
func (s *Service) SetupRemoteDevice(ski string, writeI shipapi.ShipConnectionDataWriterInterface) shipapi.ShipConnectionDataReaderInterface {
	return s.LocalDevice().SetupRemoteDevice(ski, writeI)
}

// report all currently visible EEBUS services
func (s *Service) VisibleRemoteServicesUpdated(entries []shipapi.RemoteService) {
	s.serviceHandler.VisibleRemoteServicesUpdated(s, entries)
}

// Provides the SHIP ID the remote service reported during the handshake process
// This needs to be persisted and passed on for future remote service connections
// when using `PairRemoteService`
func (s *Service) ServiceShipIDUpdate(ski string, shipdID string) {
	s.serviceHandler.ServiceShipIDUpdate(ski, shipdID)
}

// Provides the current pairing state for the remote service
// This is called whenever the state changes and can be used to
// provide user information for the pairing/connection process
func (s *Service) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
	s.serviceHandler.ServicePairingDetailUpdate(ski, detail)
}

// return if the user is still able to trust the connection
func (s *Service) AllowWaitingForTrust(ski string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.isPairingPossible
}
