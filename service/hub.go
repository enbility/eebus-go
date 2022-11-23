package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/gorilla/websocket"
)

const shipWebsocketSubProtocol = "ship" // SHIP 10.2: sub protocol is required for websocket connections
const shipWebsocketPath = "/ship/"
const shipZeroConfServiceType = "_ship._tcp"
const shipZeroConfDomain = "local."

// implemented by connectionsHub and used by shipConnection
type connectionHubHandler interface {
	HandleConnectionClosing(connection *shipConnection)
}

// interface for interactions back into a connection
// implemented by the connection
type connectionInteraction interface {
	ReportUserTrust(bool)
}

type connectionsHub struct {
	connections map[string]*shipConnection

	serviceDescription *ServiceDescription
	localService       *ServiceDetails

	// The list of paired devices
	registeredServices []ServiceDetails

	// The web server for handling incoming websocket connections
	httpServer *http.Server

	// Handling mDNS related tasks
	mdns *mdns

	connectionDelegate interactionShipSpine

	muxCon  sync.Mutex
	muxReg  sync.Mutex
	muxMdns sync.Mutex
}

func newConnectionsHub(serviceDescription *ServiceDescription, localService *ServiceDetails, connectionDelegate interactionShipSpine) (*connectionsHub, error) {
	hub := &connectionsHub{
		connections:        make(map[string]*shipConnection),
		registeredServices: make([]ServiceDetails, 0),
		serviceDescription: serviceDescription,
		localService:       localService,
		connectionDelegate: connectionDelegate,
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

var _ connectionHubHandler = (*connectionsHub)(nil)

// The connection was closed, we need to clean up
func (h *connectionsHub) HandleConnectionClosing(connection *shipConnection) {
	// only remove this connection if it is the registered one for the ski!
	// as we can have double connections but only one can be registered
	if existingC := h.connectionForSKI(connection.remoteService.SKI); existingC != nil {
		if connection.dataHandler == connection.dataHandler {
			h.muxCon.Lock()
			delete(h.connections, connection.remoteService.SKI)
			h.muxCon.Unlock()
		}
	}

	// startup mDNS if a paired service is not connected
	if len(h.connections) == 0 && len(h.registeredServices) > 0 {
		logging.Log.Debug("Starting mDNS")
		// if this is not a CEM also start the mDNS announcement
		if connection.localService.deviceType != model.DeviceTypeTypeEnergyManagementSystem {
			_ = h.mdns.Announce()
		}
		h.mdns.RegisterMdnsSearch(h)
	}
}

// register a new ship Connection
func (h *connectionsHub) registerConnection(connection *shipConnection) {
	// check if we already have a connection for the SKI
	if existingC := h.connectionForSKI(connection.remoteService.SKI); existingC != nil {
		// TODO: provide a reason
		go connection.CloseConnection(true)
		return
	}

	h.muxCon.Lock()
	h.connections[connection.remoteService.SKI] = connection
	h.muxCon.Unlock()

	// shutdown mDNS if this is not a CEM
	if connection.localService.deviceType != model.DeviceTypeTypeEnergyManagementSystem {
		h.mdns.Unannounce()
		h.mdns.UnregisterMdnsSearch(h)
	}
}

// return the connection for a specific SKI
func (h *connectionsHub) connectionForSKI(ski string) *shipConnection {
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

func (h *connectionsHub) disconnectSKI(ski string) {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	// The connection with the higher SKI should retain the connection
	con, ok := h.connections[ski]
	if !ok {
		return
	}

	con.CloseConnection(true)
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
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
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
	_, err = h.registeredServiceForSKI(ski)
	if err != nil {
		logging.Log.Debug("ski", ski, "is not registered, closing the connection")
		return
	}

	remoteService := &ServiceDetails{
		SKI: ski,
	}
	// check if we already know this remote service
	if remoteS, err := h.registeredServiceForSKI(ski); err == nil {
		remoteService = remoteS
	}

	// don't allow a second connection
	if existingC := h.connectionForSKI(ski); existingC != nil {
		logging.Log.Debug("incoming double connection closed")
		conn.Close()
		return
	}

	dataConnection := newWebsocketConnection(conn, ski)
	shipConnection := newConnectionHandler(h, h.connectionDelegate, ShipRoleClient, h.localService, remoteService, dataConnection)

	h.registerConnection(shipConnection)
}

// Connect to another EEBUS service
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

	// don't allow a second connection
	if existingC := h.connectionForSKI(remoteService.SKI); existingC != nil {
		errorString := "outgoing double connection closed"
		logging.Log.Debug(errorString)
		conn.Close()
		return errors.New(errorString)
	}

	dataConnection := newWebsocketConnection(conn, remoteService.SKI)
	shipConnection := newConnectionHandler(h, h.connectionDelegate, ShipRoleServer, h.localService, remoteService, dataConnection)

	h.registerConnection(shipConnection)

	return nil
}

func (h *connectionsHub) registeredServiceForSKI(ski string) (*ServiceDetails, error) {
	h.muxReg.Lock()
	defer h.muxReg.Unlock()
	for _, service := range h.registeredServices {
		if service.SKI == ski {
			return &service, nil
		}
	}
	return &ServiceDetails{}, fmt.Errorf("No registered service found for SKI %s", ski)
}

// Adds a new device to the list of known devices which can be connected to
// and connect it if it is currently not connected
func (h *connectionsHub) registerRemoteService(service ServiceDetails) {

	// standardize the provided SKI strings
	service.SKI = util.NormalizeSKI(service.SKI)

	// check if it is already registered
	if _, err := h.registeredServiceForSKI(service.SKI); err != nil {
		h.muxReg.Lock()
		h.registeredServices = append(h.registeredServices, service)
		h.muxReg.Unlock()
	}

	if !h.isSkiConnected(service.SKI) {
		h.mdns.RegisterMdnsSearch(h)
	}
}

// Update known device in the list of known devices which can be connected to
func (h *connectionsHub) updateRemoteServiceTrust(ski string, trusted bool) {
	h.muxReg.Lock()

	var conn *shipConnection

	for i, device := range h.registeredServices {
		if device.SKI == ski {
			h.registeredServices[i].userTrust = trusted

			conn = h.connectionForSKI(ski)
			break
		}
	}

	h.muxReg.Unlock()

	if conn != nil {
		conn.ReportUserTrust(trusted)
	}
}

// Remove a device from the list of known devices which can be connected to
// and disconnect it if it is currently connected
func (h *connectionsHub) unregisterRemoteService(ski string) error {
	newRegisteredDevice := make([]ServiceDetails, 0)

	h.muxReg.Lock()
	for _, device := range h.registeredServices {
		if device.SKI != ski {
			newRegisteredDevice = append(newRegisteredDevice, device)
		}
	}

	h.registeredServices = newRegisteredDevice
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
		remoteService, err := h.registeredServiceForSKI(ski)
		if err != nil {
			continue
		}

		// patch the addresses list if an IPv4 address was provided
		if remoteService.IPv4 != "" {
			if ip := net.ParseIP(remoteService.IPv4); ip != nil {
				entry.Addresses = []net.IP{ip}
			}
		}

		logging.Log.Debug("Trying to connect to", ski, "at", entry.Host)
		if err = h.connectFoundService(remoteService, entry.Host, strconv.Itoa(entry.Port)); err != nil {
			// connecting via the host failed, so try all of the provided addresses
			for _, address := range entry.Addresses {
				logging.Log.Debug("Trying to connect to", ski, "at", address)
				if err = h.connectFoundService(remoteService, address.String(), strconv.Itoa(entry.Port)); err == nil {
					break
				}
			}
			if err != nil {
				continue
			}
		}
	}
}
