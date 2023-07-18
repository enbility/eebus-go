package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/ship"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/gorilla/websocket"
)

const shipWebsocketSubProtocol = "ship" // SHIP 10.2: sub protocol is required for websocket connections
const shipWebsocketPath = "/ship/"

// used for randomizing the connection initiation delay
// this limits the possibility of concurrent connection attempts from both sides
type connectionInitiationDelayTimeRange struct {
	// defines the minimum and maximum wait time for when to try to initate an connection
	min, max int
}

// defines the delay timeframes in seconds depening on the connection attempt counter
// the last item will be re-used for higher attempt counter values
var connectionInitiationDelayTimeRanges = []connectionInitiationDelayTimeRange{
	{min: 0, max: 3},
	{min: 3, max: 10},
	{min: 10, max: 20},
}

//go:generate mockgen -destination=mock_hub.go -package=service github.com/enbility/eebus-go/service ServiceProvider

// interface for reporting data from connectionsHub to the EEBUSService
type ServiceProvider interface {
	// report a newly discovered remote EEBUS service
	VisibleMDNSRecordsUpdated(entries []MdnsEntry)

	// report a connection to a SKI
	RemoteSKIConnected(ski string)

	// report a disconnection to a SKI
	RemoteSKIDisconnected(ski string)

	// provide the SHIP ID received during SHIP handshake process
	// the ID needs to be stored and then provided for remote services so it can be compared and verified
	ServiceShipIDUpdate(ski string, shipID string)

	// provides the current handshake state for a given SKI
	ServicePairingDetailUpdate(ski string, detail ConnectionStateDetail)

	// return if the user is still able to trust the connection
	AllowWaitingForTrust(ski string) bool
}

// handling all connections to remote services
type connectionsHub struct {
	connections map[string]*ship.ShipConnection

	// which attempt is it to initate an connection to the remote SKI
	connectionAttemptCounter map[string]int
	connectionAttemptRunning map[string]bool

	configuration *Configuration
	localService  *ServiceDetails

	serviceProvider ServiceProvider

	// The list of known remote services
	remoteServices map[string]*ServiceDetails

	// The web server for handling incoming websocket connections
	httpServer *http.Server

	// Handling mDNS related tasks
	mdns MdnsService

	// list of currently known/reported mDNS entries
	knownMdnsEntries []MdnsEntry

	// the SPINE local device
	spineLocalDevice *spine.DeviceLocalImpl

	muxCon        sync.Mutex
	muxConAttempt sync.Mutex
	muxReg        sync.Mutex
	muxMdns       sync.Mutex
}

func newConnectionsHub(serviceProvider ServiceProvider, mdns MdnsService, spineLocalDevice *spine.DeviceLocalImpl, configuration *Configuration, localService *ServiceDetails) *connectionsHub {
	hub := &connectionsHub{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		connectionAttemptRunning: make(map[string]bool),
		remoteServices:           make(map[string]*ServiceDetails),
		knownMdnsEntries:         make([]MdnsEntry, 0),
		serviceProvider:          serviceProvider,
		spineLocalDevice:         spineLocalDevice,
		configuration:            configuration,
		localService:             localService,
		mdns:                     mdns,
	}

	return hub
}

// start the ConnectionsHub with all its services
func (h *connectionsHub) start() {
	// start the websocket server
	if err := h.startWebsocketServer(); err != nil {
		logging.Log.Debug("error during websocket server starting:", err)
	}

	// start mDNS
	err := h.mdns.SetupMdnsService()
	if err != nil {
		logging.Log.Debug("error during mdns setup:", err)
	}

	h.checkRestartMdnsSearch()
}

var _ ship.ShipServiceDataProvider = (*connectionsHub)(nil)

// Returns if the provided SKI is from a registered service
func (h *connectionsHub) IsRemoteServiceForSKIPaired(ski string) bool {
	service := h.serviceForSKI(ski)

	return service.Trusted
}

// The connection was closed, we need to clean up
func (h *connectionsHub) HandleConnectionClosed(connection *ship.ShipConnection, handshakeCompleted bool) {
	// only remove this connection if it is the registered one for the ski!
	// as we can have double connections but only one can be registered
	if existingC := h.connectionForSKI(connection.RemoteSKI); existingC != nil {
		if existingC.DataHandler == connection.DataHandler {
			h.muxCon.Lock()
			delete(h.connections, connection.RemoteSKI)
			h.muxCon.Unlock()
		}

		// connection close was after a completed handshake, so we can reset the attetmpt counter
		if handshakeCompleted {
			h.removeConnectionAttemptCounter(connection.RemoteSKI)
		}
	}

	h.serviceProvider.RemoteSKIDisconnected(connection.RemoteSKI)

	// Do not automatically reconnect if handshake failed and not already paired
	remoteService := h.serviceForSKI(connection.RemoteSKI)
	if !handshakeCompleted && !remoteService.Trusted {
		return
	}

	h.checkRestartMdnsSearch()
}

// return the number of paired services
func (h *connectionsHub) numberPairedServices() int {
	amount := 0

	h.muxReg.Lock()
	for _, service := range h.remoteServices {
		if service.Trusted {
			amount++
		}
	}
	h.muxReg.Unlock()

	return amount
}

// startup mDNS if a paired service is not connected
func (h *connectionsHub) checkRestartMdnsSearch() {
	countPairedServices := h.numberPairedServices()
	h.muxCon.Lock()
	countConnections := len(h.connections)
	h.muxCon.Unlock()

	if countPairedServices > countConnections {
		// if this is not a CEM also start the mDNS announcement
		if h.localService.DeviceType != model.DeviceTypeTypeEnergyManagementSystem {
			_ = h.mdns.AnnounceMdnsEntry()
		}

		h.mdns.RegisterMdnsSearch(h)
	}
}

func (h *connectionsHub) StartBrowseMdnsSearch() {
	// TODO: this currently collides with searching for a specific SKI
	h.mdns.RegisterMdnsSearch(h)
}

func (h *connectionsHub) StopBrowseMdnsSearch() {
	// TODO: this currently collides with searching for a specific SKI
	h.mdns.UnregisterMdnsSearch(h)
}

// Provides the SHIP ID the remote service reported during the handshake process
func (h *connectionsHub) ReportServiceShipID(ski string, shipdID string) {
	h.serviceProvider.RemoteSKIConnected(ski)

	h.serviceProvider.ServiceShipIDUpdate(ski, shipdID)
}

// return if the user is still able to trust the connection
func (h *connectionsHub) AllowWaitingForTrust(ski string) bool {
	if service := h.serviceForSKI(ski); service != nil {
		if service.Trusted {
			return true
		}
	}

	return h.serviceProvider.AllowWaitingForTrust(ski)
}

// Provides the current ship message exchange state for a given SKI and the corresponding error if state is error
func (h *connectionsHub) HandleShipHandshakeStateUpdate(ski string, state ship.ShipState) {
	service := h.serviceForSKI(ski)

	// overwrite service Paired value
	if state.State == ship.SmeHelloStateOk {
		h.RegisterRemoteSKI(ski, true)
	}

	pairingState := h.mapShipMessageExchangeState(state.State, ski)
	if state.Error != nil && state.Error != ErrConnectionNotFound {
		pairingState = ConnectionStateError
	}

	pairingDetail := ConnectionStateDetail{
		State: pairingState,
		Error: state.Error,
	}

	existingDetails := service.ConnectionStateDetail
	if existingDetails.State != pairingState || existingDetails.Error != state.Error {
		service.ConnectionStateDetail = pairingDetail

		h.serviceProvider.ServicePairingDetailUpdate(ski, pairingDetail)
	}
}

// Provide the current pairing state for a given SKI
//
// returns:
//
//	ErrNotPaired if the SKI is not in the (to be) paired list
//	ErrNoConnectionFound if no connection for the SKI was found
func (h *connectionsHub) PairingDetailForSki(ski string) ConnectionStateDetail {
	service := h.serviceForSKI(ski)

	if conn := h.connectionForSKI(ski); conn != nil {
		shipState, shipError := conn.ShipHandshakeState()
		state := h.mapShipMessageExchangeState(shipState, ski)
		return ConnectionStateDetail{
			State: state,
			Error: shipError,
		}
	}

	return service.ConnectionStateDetail
}

// maps ShipMessageExchangeState to PairingState
func (h *connectionsHub) mapShipMessageExchangeState(state ship.ShipMessageExchangeState, ski string) ConnectionState {
	var connState ConnectionState

	// map the SHIP states to a public gState
	switch state {
	case ship.CmiStateInitStart:
		connState = ConnectionStateQueued
	case ship.CmiStateClientSend, ship.CmiStateClientWait, ship.CmiStateClientEvaluate,
		ship.CmiStateServerWait, ship.CmiStateServerEvaluate:
		connState = ConnectionStateInitiated
	case ship.SmeHelloStateReadyInit, ship.SmeHelloStateReadyListen, ship.SmeHelloStateReadyTimeout:
		connState = ConnectionStateInProgress
	case ship.SmeHelloStatePendingInit, ship.SmeHelloStatePendingListen, ship.SmeHelloStatePendingTimeout:
		connState = ConnectionStateReceivedPairingRequest
	case ship.SmeHelloStateOk:
		connState = ConnectionStateTrusted
	case ship.SmeHelloStateAbort, ship.SmeHelloStateAbortDone:
		connState = ConnectionStateNone
	case ship.SmeHelloStateRemoteAbortDone, ship.SmeHelloStateRejected:
		connState = ConnectionStateRemoteDeniedTrust
	case ship.SmePinStateCheckInit, ship.SmePinStateCheckListen, ship.SmePinStateCheckError,
		ship.SmePinStateCheckBusyInit, ship.SmePinStateCheckBusyWait, ship.SmePinStateCheckOk,
		ship.SmePinStateAskInit, ship.SmePinStateAskProcess, ship.SmePinStateAskRestricted,
		ship.SmePinStateAskOk:
		connState = ConnectionStatePin
	case ship.SmeAccessMethodsRequest, ship.SmeStateApproved:
		connState = ConnectionStateInProgress
	case ship.SmeStateComplete:
		connState = ConnectionStateCompleted
	case ship.SmeStateError:
		connState = ConnectionStateError
	default:
		connState = ConnectionStateInProgress
	}

	return connState
}

// Disconnect a connection to an SKI, used by a service implementation
// e.g. if heartbeats go wrong
func (h *connectionsHub) DisconnectSKI(ski string, reason string) {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	// The connection with the higher SKI should retain the connection
	con, ok := h.connections[ski]
	if !ok {
		return
	}

	con.CloseConnection(true, 0, reason)
}

// register a new ship Connection
func (h *connectionsHub) registerConnection(connection *ship.ShipConnection) {
	h.muxCon.Lock()
	h.connections[connection.RemoteSKI] = connection
	h.muxCon.Unlock()
}

// return the connection for a specific SKI
func (h *connectionsHub) connectionForSKI(ski string) *ship.ShipConnection {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	con, ok := h.connections[ski]
	if !ok {
		return nil
	}
	return con
}

// close all connections
func (h *connectionsHub) shutdown() {
	h.mdns.ShutdownMdnsService()
	for _, c := range h.connections {
		c.CloseConnection(false, 0, "")
	}
}

// return if there is a connection for a SKI
func (h *connectionsHub) isSkiConnected(ski string) bool {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	// The connection with the higher SKI should retain the connection
	_, ok := h.connections[ski]
	return ok
}

// Websocket connection handling

// start the ship websocket server
func (h *connectionsHub) startWebsocketServer() error {
	addr := fmt.Sprintf(":%d", h.configuration.port)
	logging.Log.Debug("starting websocket server on", addr)

	h.httpServer = &http.Server{
		Addr:    addr,
		Handler: h,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{h.configuration.certificate},
			ClientAuth:   tls.RequireAnyClientCert, // SHIP 9: Client authentication is required
			CipherSuites: ciperSuites,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				skiFound := false
				for _, v := range rawCerts {
					cert, err := x509.ParseCertificate(v)
					if err != nil {
						return err
					}

					if _, err := skiFromCertificate(cert); err == nil {
						skiFound = true
						break
					}
				}
				if !skiFound {
					return errors.New("no valid SKI provided in certificate")
				}

				return nil
			},
		},
	}

	go func() {
		if err := h.httpServer.ListenAndServeTLS("", ""); err != nil {
			logging.Log.Debug("websocket server error:", err)
			// TODO: decide how to handle this case
		}
	}()

	return nil
}

// Connection Handling

// HTTP Server callback for handling incoming connection requests
func (h *connectionsHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  ship.MaxMessageSize,
		WriteBufferSize: ship.MaxMessageSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
		Subprotocols:    []string{shipWebsocketSubProtocol}, // SHIP 10.2: Sub protocol "ship" is required
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Log.Debug("error during connection upgrading:", err)
		return
	}

	// check if the client supports the ship sub protocol
	if conn.Subprotocol() != shipWebsocketSubProtocol {
		logging.Log.Debug("client does not support the ship sub protocol")
		_ = conn.Close()
		return
	}

	// check if the clients certificate provides a SKI
	if len(r.TLS.PeerCertificates) == 0 {
		logging.Log.Debug("client does not provide a certificate")
		_ = conn.Close()
		return
	}

	ski, err := skiFromCertificate(r.TLS.PeerCertificates[0])
	if err != nil {
		logging.Log.Debug(err)
		_ = conn.Close()
		return
	}

	// normalize the incoming SKI
	remoteService := NewServiceDetails(ski)
	logging.Log.Debug("incoming connection request from", remoteService.SKI)

	// Check if the remote service is paired
	service := h.serviceForSKI(remoteService.SKI)
	if service.ConnectionStateDetail.State == ConnectionStateQueued {
		service.ConnectionStateDetail.State = ConnectionStateReceivedPairingRequest
		h.serviceProvider.ServicePairingDetailUpdate(ski, service.ConnectionStateDetail)
	}

	remoteService = service

	// don't allow a second connection
	if !h.keepThisConnection(conn, true, remoteService) {
		_ = conn.Close()
		return
	}

	dataHandler := ship.NewWebsocketConnection(conn, remoteService.SKI)
	shipConnection := ship.NewConnectionHandler(h, dataHandler, h.spineLocalDevice, ship.ShipRoleServer,
		h.localService.ShipID, remoteService.SKI, remoteService.ShipID)
	shipConnection.Run()

	h.registerConnection(shipConnection)
}

// Connect to another EEBUS service
//
// returns error contains a reason for failing the connection or nil if no further tries should be processed
func (h *connectionsHub) connectFoundService(remoteService *ServiceDetails, host, port string) error {
	if h.isSkiConnected(remoteService.SKI) {
		return nil
	}

	logging.Log.Debugf("initiating connection to %s at %s:%s", remoteService.SKI, host, port)

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{h.configuration.certificate},
			InsecureSkipVerify: true,
			CipherSuites:       ciperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}

	address := fmt.Sprintf("wss://%s:%s", host, port)
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		return err
	}

	tlsConn := conn.UnderlyingConn().(*tls.Conn)
	remoteCerts := tlsConn.ConnectionState().PeerCertificates

	if len(remoteCerts) == 0 || remoteCerts[0].SubjectKeyId == nil {
		// Close connection as we couldn't get the remote SKI
		errorString := fmt.Sprintf("closing connection to %s: could not get remote SKI from certificate", remoteService.SKI)
		conn.Close()
		return errors.New(errorString)
	}

	if _, err := skiFromCertificate(remoteCerts[0]); err != nil {
		// Close connection as the remote SKI can't be correct
		errorString := fmt.Sprintf("closing connection to %s: %s", remoteService.SKI, err)
		conn.Close()
		return errors.New(errorString)
	}

	remoteSKI := fmt.Sprintf("%0x", remoteCerts[0].SubjectKeyId)

	if remoteSKI != remoteService.SKI {
		errorString := fmt.Sprintf("closing connection to %s: SKI does not match %s", remoteService.SKI, remoteSKI)
		conn.Close()
		return errors.New(errorString)
	}

	if !h.keepThisConnection(conn, false, remoteService) {
		errorString := fmt.Sprintf("closing connection to %s: ignoring this connection", remoteService.SKI)
		return errors.New(errorString)
	}

	dataHandler := ship.NewWebsocketConnection(conn, remoteService.SKI)
	shipConnection := ship.NewConnectionHandler(h, dataHandler, h.spineLocalDevice, ship.ShipRoleClient,
		h.localService.ShipID, remoteService.SKI, remoteService.ShipID)
	shipConnection.Run()

	h.registerConnection(shipConnection)

	return nil
}

// prevent double connections
// only keep the connection initiated by the higher SKI
//
// returns true if this connection is fine to be continue
// returns false if this connection should not be established or kept
func (h *connectionsHub) keepThisConnection(conn *websocket.Conn, incomingRequest bool, remoteService *ServiceDetails) bool {
	// SHIP 12.2.2 defines:
	// prevent double connections with SKI Comparison
	// the node with the hight SKI value kees the most recent connection and
	// and closes all other connections to the same SHIP node
	//
	// This is hard to implement without any flaws. Therefor I chose a
	// different approach: The connection initiated by the higher SKI will be kept

	remoteSKI := remoteService.SKI
	existingC := h.connectionForSKI(remoteSKI)
	if existingC == nil {
		return true
	}

	keep := false
	if incomingRequest {
		keep = remoteSKI > h.localService.SKI
	} else {
		keep = h.localService.SKI > remoteSKI
	}

	if keep {
		// we have an existing connection
		// so keep the new (most recent) and close the old one
		logging.Log.Debug("closing existing double connection")
		go existingC.CloseConnection(false, 0, "")
	} else {
		connType := "incoming"
		if !incomingRequest {
			connType = "outgoing"
		}
		logging.Log.Debugf("closing %s double connection, as the existing connection will be used", connType)
		go func() {
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "double connection"))
			time.Sleep(time.Millisecond * 100)
			conn.Close()
		}()
	}

	return keep
}

// return the service for a given SKI or an error if not found
func (h *connectionsHub) serviceForSKI(ski string) *ServiceDetails {
	h.muxReg.Lock()
	defer h.muxReg.Unlock()

	service, ok := h.remoteServices[ski]
	if !ok {
		service = NewServiceDetails(ski)
		service.ConnectionStateDetail.State = ConnectionStateNone
		h.remoteServices[ski] = service
	}

	return service
}

// Sets the SKI as being paired or not
// Should be used for services which completed the pairing process and
// which were stored as having the process completed
func (h *connectionsHub) RegisterRemoteSKI(ski string, enable bool) {
	service := h.serviceForSKI(ski)
	service.Trusted = enable

	if enable {
		h.checkRestartMdnsSearch()
		return
	}

	h.removeConnectionAttemptCounter(ski)

	service.ConnectionStateDetail.State = ConnectionStateNone

	h.serviceProvider.ServicePairingDetailUpdate(ski, service.ConnectionStateDetail)

	if existingC := h.connectionForSKI(ski); existingC != nil {
		existingC.CloseConnection(true, 4500, "User close")
	}
}

// Triggers the pairing process for a SKI
func (h *connectionsHub) InitiatePairingWithSKI(ski string) {
	conn := h.connectionForSKI(ski)

	// remotely initiated
	if conn != nil {
		conn.ApprovePendingHandshake()

		return
	}

	// locally initiated
	service := h.serviceForSKI(ski)
	service.ConnectionStateDetail.State = ConnectionStateQueued

	h.serviceProvider.ServicePairingDetailUpdate(ski, service.ConnectionStateDetail)

	// initiate a search and also a connection if it does not yet exist
	if !h.isSkiConnected(service.SKI) {
		h.mdns.RegisterMdnsSearch(h)
	}
}

// Cancels the pairing process for a SKI
func (h *connectionsHub) CancelPairingWithSKI(ski string) {
	h.removeConnectionAttemptCounter(ski)

	if existingC := h.connectionForSKI(ski); existingC != nil {
		existingC.AbortPendingHandshake()
	}

	service := h.serviceForSKI(ski)
	service.ConnectionStateDetail.State = ConnectionStateNone
	service.Trusted = false

	h.serviceProvider.ServicePairingDetailUpdate(ski, service.ConnectionStateDetail)
}

// Process reported mDNS services
func (h *connectionsHub) ReportMdnsEntries(entries map[string]MdnsEntry) {
	h.muxMdns.Lock()
	defer h.muxMdns.Unlock()

	var mdnsEntries []MdnsEntry

	for ski, entry := range entries {
		mdnsEntries = append(mdnsEntries, entry)

		// check if this ski is already connected
		if h.isSkiConnected(ski) {
			continue
		}

		// Check if the remote service is paired or queued for connection
		service := h.serviceForSKI(ski)
		pairingState := service.ConnectionStateDetail.State
		if !h.IsRemoteServiceForSKIPaired(ski) && pairingState != ConnectionStateQueued {
			continue
		}

		// patch the addresses list if an IPv4 address was provided
		if service.IPv4 != "" {
			if ip := net.ParseIP(service.IPv4); ip != nil {
				entry.Addresses = []net.IP{ip}
			}
		}

		h.coordinateConnectionInitations(ski, entry)
	}

	sort.Slice(mdnsEntries, func(i, j int) bool {
		item1 := mdnsEntries[i]
		item2 := mdnsEntries[j]
		a := strings.ToLower(item1.Brand + item1.Model + item1.Ski)
		b := strings.ToLower(item2.Brand + item2.Model + item2.Ski)
		return a < b
	})

	h.knownMdnsEntries = mdnsEntries

	h.serviceProvider.VisibleMDNSRecordsUpdated(mdnsEntries)
}

// coordinate connection initiation attempts to a remove service
func (h *connectionsHub) coordinateConnectionInitations(ski string, entry MdnsEntry) {
	if h.isConnectionAttemptRunning(ski) {
		return
	}

	h.setConnectionAttemptRunning(ski, true)

	counter, duration := h.getConnectionInitiationDelayTime(ski)

	service := h.serviceForSKI(ski)
	if service.ConnectionStateDetail.State == ConnectionStateQueued {
		go h.prepareConnectionInitation(ski, counter, entry)
		return
	}

	logging.Log.Debugf("delaying connection to %s by %s to minimize double connection probability", ski, duration)

	// we do not stop this thread and just let the timer run out
	// otherwise we would need a stop channel for each ski
	go func() {
		// wait
		<-time.After(duration)

		h.prepareConnectionInitation(ski, counter, entry)
	}()
}

// invoked by coordinateConnectionInitations either with a delay or directly
// when initating a pairing process
func (h *connectionsHub) prepareConnectionInitation(ski string, counter int, entry MdnsEntry) {
	h.setConnectionAttemptRunning(ski, false)

	// check if the remoteService still exists
	service := h.serviceForSKI(ski)

	// check if the current counter is still the same, otherwise this counter is irrelevant
	currentCounter, exists := h.getCurrentConnectionAttemptCounter(ski)
	if !exists || currentCounter != counter {
		return
	}

	// connection attempt is not relevant if the device is no longer paired
	// or it is not queued for pairing
	pairingState := h.serviceForSKI(ski).ConnectionStateDetail.State
	if !h.IsRemoteServiceForSKIPaired(ski) && pairingState != ConnectionStateQueued {
		return
	}

	// connection attempt is not relevant if the device is already connected
	if h.isSkiConnected(ski) {
		return
	}

	// now initiate the connection
	if success := h.initateConnection(service, entry); !success {
		h.checkRestartMdnsSearch()
	}
}

// attempt to establish a connection to a remote service
// returns true if successful
func (h *connectionsHub) initateConnection(remoteService *ServiceDetails, entry MdnsEntry) bool {
	var err error

	// try connecting via an IP address first
	for _, address := range entry.Addresses {
		// connection attempt is not relevant if the device is no longer paired
		// or it is not queued for pairing
		pairingState := h.serviceForSKI(remoteService.SKI).ConnectionStateDetail.State
		if !h.IsRemoteServiceForSKIPaired(remoteService.SKI) && pairingState != ConnectionStateQueued {
			return false
		}

		logging.Log.Debug("trying to connect to", remoteService.SKI, "at", address)
		if err = h.connectFoundService(remoteService, address.String(), strconv.Itoa(entry.Port)); err != nil {
			logging.Log.Debug("connection to", remoteService.SKI, "failed: ", err)
		} else {
			return true
		}
	}

	// connectdion via IP address failed, try hostname
	if len(entry.Host) > 0 {
		logging.Log.Debug("trying to connect to", remoteService.SKI, "at", entry.Host)
		if err = h.connectFoundService(remoteService, entry.Host, strconv.Itoa(entry.Port)); err != nil {
			logging.Log.Debugf("connection to %s failed: %s", remoteService.SKI, err)
		} else {
			return true
		}
	}

	// no connection could be estabished via any of the provided addresses
	// because no service was reachable at any of the addresses
	return false
}

// increase the connection attempt counter for the given ski
func (h *connectionsHub) increaseConnectionAttemptCounter(ski string) int {
	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	currentCounter := 0
	if counter, exists := h.connectionAttemptCounter[ski]; exists {
		currentCounter = counter + 1

		if currentCounter >= len(connectionInitiationDelayTimeRanges)-1 {
			currentCounter = len(connectionInitiationDelayTimeRanges) - 1
		}
	}

	h.connectionAttemptCounter[ski] = currentCounter

	return currentCounter
}

// remove the connection attempt counter for the given ski
func (h *connectionsHub) removeConnectionAttemptCounter(ski string) {
	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	delete(h.connectionAttemptCounter, ski)
}

// get the current attempt counter
func (h *connectionsHub) getCurrentConnectionAttemptCounter(ski string) (int, bool) {
	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	counter, exists := h.connectionAttemptCounter[ski]

	return counter, exists
}

// get the connection initiation delay time range for a given ski
// returns the current counter and the duration
func (h *connectionsHub) getConnectionInitiationDelayTime(ski string) (int, time.Duration) {
	counter := h.increaseConnectionAttemptCounter(ski)

	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	timeRange := connectionInitiationDelayTimeRanges[counter]

	// get range in Milliseconds
	min := timeRange.min * 1000
	max := timeRange.max * 1000

	// seed with the local SKI for initializing rand
	// TODO: remove when upping minimum go version to 1.20
	i := new(big.Int)
	hex := fmt.Sprintf("0x%s", h.localService.SKI)
	randSource := rand.NewSource(time.Now().UnixNano())
	if _, err := fmt.Sscan(hex, i); err == nil {
		randSource = rand.NewSource(i.Int64() + time.Now().UnixNano())
	}

	r := rand.New(randSource)
	duration := r.Intn(max-min) + min

	return counter, time.Duration(duration) * time.Millisecond
}

// set if a connection attempt is running/in progress
func (h *connectionsHub) setConnectionAttemptRunning(ski string, active bool) {
	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	h.connectionAttemptRunning[ski] = active
}

// return if a connection attempt is runnning/in progress
func (h *connectionsHub) isConnectionAttemptRunning(ski string) bool {
	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	running, exists := h.connectionAttemptRunning[ski]
	if !exists {
		return false
	}

	return running
}
