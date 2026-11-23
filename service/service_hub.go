package service

import (
	shipapi "github.com/enbility/ship-go/api"
)

var _ shipapi.HubReaderInterface = (*Service)(nil)

// report a connection to a remote service
//
// is triggered whenever a SHIP connected was successful completed
func (s *Service) RemoteServiceConnected(identity shipapi.ServiceIdentity) {
	s.serviceHandler.RemoteServiceConnected(s, identity)
}

// report a disconnection from a remote service
//
// is triggered whenever a SHIP connect was closed, is also triggered when the SHIP
// process wasn't successfully completed
//
// NOTE: The connection may not have been reported as connected before!
func (s *Service) RemoteServiceDisconnected(identity shipapi.ServiceIdentity) {
	if s.spineLocalDevice != nil {
		s.spineLocalDevice.RemoveRemoteDeviceConnection(identity.SKI)
	}

	s.serviceHandler.RemoteServiceDisconnected(s, identity)
}

// report an approved handshake by a remote device
func (s *Service) SetupRemoteService(identity shipapi.ServiceIdentity, writeI shipapi.ShipConnectionDataWriterInterface) shipapi.ShipConnectionDataReaderInterface {
	return s.LocalDevice().SetupRemoteDevice(identity.SKI, writeI)
}

// report all currently visible EEBUS services
func (s *Service) VisibleRemoteMdnsServicesUpdated(entries []shipapi.RemoteMdnsService) {
	s.serviceHandler.VisibleRemoteMdnsServicesUpdated(s, entries)
}

// report that service information has been updated
// This includes updates to ShipID, fingerprint, or other service details
// discovered during handshake
func (s *Service) ServiceUpdated(identity shipapi.ServiceIdentity) {
	s.serviceHandler.ServiceUpdated(identity)
}

// Provides the current pairing state for the remote service
// This is called whenever the state changes and can be used to
// provide user information for the pairing/connection process
func (s *Service) ServicePairingDetailUpdate(identity shipapi.ServiceIdentity, detail *shipapi.ConnectionStateDetail) {
	s.serviceHandler.ServicePairingDetailUpdate(identity, detail)
}

// return if the user is still able to trust the connection
func (s *Service) AllowWaitingForTrust(identity shipapi.ServiceIdentity) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.isPairingPossible
}

// Called when a device is automatically trusted via SHIP pairing
func (s *Service) ServiceAutoTrusted(identity shipapi.ServiceIdentity) {
	s.serviceHandler.ServiceAutoTrusted(s, identity)
}

// Called when SHIP pairing fails for a device
func (s *Service) ServiceAutoTrustFailed(identity shipapi.ServiceIdentity, reason error) {
	s.serviceHandler.ServiceAutoTrustFailed(s, identity, reason)
}

// Called when device trust is automatically removed This can happen due to device replacement timeout or new device pairing
func (s *Service) ServiceAutoTrustRemoved(identity shipapi.ServiceIdentity, reason string) {
	s.serviceHandler.ServiceAutoTrustRemoved(s, identity, reason)
}
