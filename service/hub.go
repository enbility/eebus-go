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
	"strconv"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/gorilla/websocket"
)

const shipWebsocketSubProtocol = "ship" // SHIP 10.2: sub protocol is required for websocket connections
const shipWebsocketPath = "/ship/"
const shipZeroConfServiceType = "_ship._tcp"
const shipZeroConfDomain = "local."

// used for randomizing the connection initation delay
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
	{min: 10, max: 30},
	{min: 30, max: 60},
	{min: 60, max: 180},
	{min: 180, max: 360},
}

type connectionsHub struct {
	connections map[string]*ship.ShipConnection

	// which attempt is it to initate an connection to the remote SKI
	connectionAttemptCounter map[string]int

	serviceDescription *ServiceDescription
	localService       *ServiceDetails

	// The list of paired devices
	pairedServices []ServiceDetails

	// The web server for handling incoming websocket connections
	httpServer *http.Server

	// Handling mDNS related tasks
	mdns *mdns

	// the SPINE local device
	spineLocalDevice *spine.DeviceLocalImpl

	muxCon        sync.Mutex
	muxConAttempt sync.Mutex
	muxReg        sync.Mutex
	muxMdns       sync.Mutex
}

func newConnectionsHub(spineLocalDevice *spine.DeviceLocalImpl, serviceDescription *ServiceDescription, localService *ServiceDetails) (*connectionsHub, error) {
	hub := &connectionsHub{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		pairedServices:           make([]ServiceDetails, 0),
		spineLocalDevice:         spineLocalDevice,
		serviceDescription:       serviceDescription,
		localService:             localService,
	}

	localService.SKI = util.NormalizeSKI(localService.SKI)

	mdns, err := newMDNS(localService.SKI, serviceDescription)
	if err != nil {
		return nil, err
	}

	hub.mdns = mdns

	return hub, nil
}

// start the ConnectionsHub with all its services
func (h *connectionsHub) start() {
	// start the websocket server
	go func() {
		if err := h.startWebsocketServer(); err != nil {
			logging.Log.Error("Error during websocket server starting: ", err)
		}
	}()

	if err := h.mdns.Announce(); err != nil {
		logging.Log.Error("Error registering mDNS Service:", err)
	}
}

var _ ship.ShipServiceDataProvider = (*connectionsHub)(nil)

// Returns if the provided SKI is from a registered service
func (h *connectionsHub) IsRemoteServiceForSKIPaired(ski string) bool {
	if _, err := h.PairedServiceForSKI(ski); err != nil {
		return false
	}

	return true
}

// The connection was closed, we need to clean up
func (h *connectionsHub) HandleConnectionClosing(connection *ship.ShipConnection) {
	// only remove this connection if it is the registered one for the ski!
	// as we can have double connections but only one can be registered
	if existingC := h.connectionForSKI(connection.RemoteSKI); existingC != nil {
		if connection.DataHandler == connection.DataHandler {
			h.muxCon.Lock()
			delete(h.connections, connection.RemoteSKI)
			h.muxCon.Unlock()
		}
	}

	// startup mDNS if a paired service is not connected
	if len(h.connections) == 0 && len(h.pairedServices) > 0 {
		logging.Log.Debug("Starting mDNS")
		// if this is not a CEM also start the mDNS announcement
		if h.localService.deviceType != model.DeviceTypeTypeEnergyManagementSystem {
			_ = h.mdns.Announce()
		}
		h.mdns.RegisterMdnsSearch(h)
	}
}

// Disconnect a connection to an SKI, used by a service implementation
// e.g. if heartbeats go wrong
func (h *connectionsHub) DisconnectSKI(ski string) {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	// The connection with the higher SKI should retain the connection
	con, ok := h.connections[ski]
	if !ok {
		return
	}

	con.CloseConnection(true)
}

// register a new ship Connection
func (h *connectionsHub) registerConnection(connection *ship.ShipConnection) {
	h.muxCon.Lock()
	h.connections[connection.RemoteSKI] = connection
	h.muxCon.Unlock()

	// shutdown mDNS if this is not a CEM
	if h.localService.deviceType != model.DeviceTypeTypeEnergyManagementSystem {
		h.mdns.Unannounce()
		h.mdns.UnregisterMdnsSearch(h)
	}
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
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	h.mdns.shutdown()
	for _, c := range h.connections {
		c.CloseConnection(true)
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
	addr := fmt.Sprintf(":%d", h.serviceDescription.port)
	logging.Log.Debug("Starting websocket server on ", addr)

	h.httpServer = &http.Server{
		Addr:    addr,
		Handler: h,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{h.serviceDescription.certificate},
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
					return errors.New("No valid SKI provided in certificate")
				}

				return nil
			},
		},
	}

	if err := h.httpServer.ListenAndServeTLS("", ""); err != nil {
		return err
	}

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
		logging.Log.Error("Error during connection upgrading: ", err)
		return
	}

	// check if the client supports the ship sub protocol
	if conn.Subprotocol() != shipWebsocketSubProtocol {
		logging.Log.Error("Client does not support the ship sub protocol")
		conn.Close()
		return
	}

	// check if the clients certificate provides a SKI
	if len(r.TLS.PeerCertificates) == 0 {
		logging.Log.Error("Client does not provide a certificate")
		conn.Close()
		return
	}

	ski, err := skiFromCertificate(r.TLS.PeerCertificates[0])
	if err != nil {
		logging.Log.Error(err)
		conn.Close()
		return
	}

	ski = util.NormalizeSKI(ski)
	logging.Log.Debug("Incoming connection request from ", ski)

	// Check if the remote service is paired
	_, err = h.PairedServiceForSKI(ski)
	if err != nil {
		logging.Log.Debug("ski", ski, "is not registered, closing the connection")
		return
	}

	remoteService := &ServiceDetails{
		SKI: ski,
	}
	// check if we already know this remote service
	if remoteS, err := h.PairedServiceForSKI(ski); err == nil {
		remoteService = remoteS
	}

	// don't allow a second connection
	if !h.keepThisConnection(conn, true, remoteService.SKI) {
		return
	}

	dataHandler := ship.NewWebsocketConnection(conn, remoteService.SKI)
	shipConnection := ship.NewConnectionHandler(h, dataHandler, h.spineLocalDevice, ship.ShipRoleServer, h.localService.ShipID, remoteService.SKI, remoteService.ShipID)
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

	logging.Log.Debug("Initiating connection to", remoteService.SKI)

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{h.serviceDescription.certificate},
			InsecureSkipVerify: true,
			CipherSuites:       ciperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}

	address := fmt.Sprintf("wss://%s:%s", host, port)
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	tlsConn := conn.UnderlyingConn().(*tls.Conn)
	remoteCerts := tlsConn.ConnectionState().PeerCertificates

	if len(remoteCerts) == 0 || remoteCerts[0].SubjectKeyId == nil {
		// Close connection as we couldn't get the remote SKI
		errorString := "closing, could not get remote SKI"
		logging.Log.Error(errorString)
		conn.Close()
		return errors.New(errorString)
	}

	if _, err := skiFromCertificate(remoteCerts[0]); err != nil {
		// Close connection as the remote SKI can't be correct
		logging.Log.Errorf("closing", err)
		conn.Close()
		return err
	}

	remoteSKI := fmt.Sprintf("%0x", remoteCerts[0].SubjectKeyId)

	if remoteSKI != remoteService.SKI {
		errorString := "remote SKI does not match"
		logging.Log.Error(errorString)
		conn.Close()
		return errors.New(errorString)
	}

	if !h.keepThisConnection(conn, false, remoteService.SKI) {
		return nil
	}

	dataHandler := ship.NewWebsocketConnection(conn, remoteService.SKI)
	shipConnection := ship.NewConnectionHandler(h, dataHandler, h.spineLocalDevice, ship.ShipRoleClient, h.localService.ShipID, remoteService.SKI, remoteService.ShipID)
	shipConnection.Run()

	h.registerConnection(shipConnection)

	return nil
}

// prevent double connections
// only keep the connection initiated by the higher SKI
//
// returns true if this connection is fine to be continue
// returns false if this connection should not be established or kept
func (h *connectionsHub) keepThisConnection(conn *websocket.Conn, incomingRequest bool, remoteSKI string) bool {
	// SHIP 12.2.2 defines:
	// prevent double connections with SKI Comparison
	// the node with the hight SKI value kees the most recent connection and
	// and closes all other connections to the same SHIP node
	//
	// This is hard to implement without any flaws. Therefor I chose a
	// different approach: The connection initiated by the higher SKI will be kept

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
		go existingC.CloseConnection(true)
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

func (h *connectionsHub) PairedServiceForSKI(ski string) (*ServiceDetails, error) {
	h.muxReg.Lock()
	defer h.muxReg.Unlock()
	for _, service := range h.pairedServices {
		if service.SKI == ski {
			return &service, nil
		}
	}
	return &ServiceDetails{}, fmt.Errorf("No registered service found for SKI %s", ski)
}

// Adds a new device to the list of known devices which can be connected to
// and connect it if it is currently not connected
func (h *connectionsHub) PairRemoteService(service ServiceDetails) {

	// standardize the provided SKI strings
	service.SKI = util.NormalizeSKI(service.SKI)

	// check if it is already registered
	if _, err := h.PairedServiceForSKI(service.SKI); err != nil {
		h.muxReg.Lock()
		h.pairedServices = append(h.pairedServices, service)
		h.muxReg.Unlock()
	}

	if !h.isSkiConnected(service.SKI) {
		h.mdns.RegisterMdnsSearch(h)
	}
}

// Remove a device from the list of known devices which can be connected to
// and disconnect it if it is currently connected
func (h *connectionsHub) UnpairRemoteService(ski string) error {
	h.removeConnectionAttemptCounter(ski)

	newRegisteredDevice := make([]ServiceDetails, 0)

	h.muxReg.Lock()
	for _, device := range h.pairedServices {
		if device.SKI != ski {
			newRegisteredDevice = append(newRegisteredDevice, device)
		}
	}

	h.pairedServices = newRegisteredDevice
	h.muxReg.Unlock()

	if existingC := h.connectionForSKI(ski); existingC != nil {
		existingC.CloseConnection(true)
	}

	return nil
}

// Process reported mDNS services
func (h *connectionsHub) ReportMdnsEntries(entries map[string]MdnsEntry) {
	h.muxMdns.Lock()
	defer h.muxMdns.Unlock()
	for ski, entry := range entries {
		// check if this ski is already connected
		if h.isSkiConnected(ski) {
			continue
		}

		// Check if the remote service is paired
		remoteService, err := h.PairedServiceForSKI(ski)
		if err != nil {
			continue
		}

		// patch the addresses list if an IPv4 address was provided
		if remoteService.IPv4 != "" {
			if ip := net.ParseIP(remoteService.IPv4); ip != nil {
				entry.Addresses = []net.IP{ip}
			}
		}

		// check if an connection initiation attempt is already ongoing
		if h.checkConnectionInitiationAttemptExists(ski) {
			continue
		}

		h.coordinateConnectionInitations(remoteService, entry)
	}
}

// coordinate connection initiation attempts to a remove service
func (h *connectionsHub) coordinateConnectionInitations(remoteService *ServiceDetails, entry MdnsEntry) {
	counter, duration := h.getConnectionInitiationDelayTime(remoteService.SKI)

	// we do not stop this thread and just let the timer run out
	// otherwise we would need a stop channel for each ski
	go func() {
		// wait
		<-time.After(duration)

		// check if the remoteService still exists
		if remoteService == nil {
			return
		}

		// check if the current counter is still the same, otherwise this counter is irrelevant
		currentCounter, exists := h.getCurrentConnectionAttemptCounter(remoteService.SKI)
		if !exists || currentCounter != counter {
			return
		}

		// connection attempt is not relevant if the device is no longer paired
		if !h.IsRemoteServiceForSKIPaired(remoteService.SKI) {
			return
		}

		// connection attempt is not relevant if the device is already connected
		if h.isSkiConnected(remoteService.SKI) {
			return
		}

		// now initiate the connection
		if h.initatingConnection(remoteService, entry) {
			return
		} else {
			// attempt failed, initate a new attempt
			h.coordinateConnectionInitations(remoteService, entry)
		}
	}()
}

// attempt to establish a connection to a remote service
// returns true if successful
func (h *connectionsHub) initatingConnection(remoteService *ServiceDetails, entry MdnsEntry) bool {
	var err error

	logging.Log.Debug("Trying to connect to", remoteService.SKI, "at", entry.Host)
	if err = h.connectFoundService(remoteService, entry.Host, strconv.Itoa(entry.Port)); err != nil {
		logging.Log.Debug("Connection failed: ", err)
		// connecting via the host failed, so try all of the provided addresses
		for _, address := range entry.Addresses {
			logging.Log.Debug("Trying to connect to", remoteService.SKI, "at", address)
			if err = h.connectFoundService(remoteService, address.String(), strconv.Itoa(entry.Port)); err != nil {
				logging.Log.Debug("Connection failed: ", err)
			} else {
				break
			}
		}
	}
	// no connection could be estabiled via any of the provided addresses
	if err != nil {
		h.increaseConnectionAttemptCounter(remoteService.SKI)
		return false
	}

	h.removeConnectionAttemptCounter(remoteService.SKI)

	return true
}

// returns if an connection initiation attempt already exists
func (h *connectionsHub) checkConnectionInitiationAttemptExists(ski string) bool {
	h.muxConAttempt.Lock()
	defer h.muxConAttempt.Unlock()

	if _, exists := h.connectionAttemptCounter[ski]; exists {
		return true
	}

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
	i := new(big.Int)
	hex := fmt.Sprintf("0x%s", h.localService.SKI)
	if _, err := fmt.Sscan(hex, i); err == nil {
		rand.Seed(i.Int64())
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	duration := rand.Intn(max-min) + min

	return counter, time.Duration(duration) * time.Millisecond
}
