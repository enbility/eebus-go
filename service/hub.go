package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/gorilla/websocket"
	"github.com/grandcat/zeroconf"
)

const shipWebsocketSubProtocol = "ship" // SHIP 10.2: sub protocol is required for websocket connections
const shipWebsocketPath = "/ship/"
const shipZeroConfServiceType = "_ship._tcp"
const shipZeroConfDomain = "local."

type connectionsHub struct {
	connections map[string]*ConnectionHandler

	// Register reuqests from a new connection
	register chan *ConnectionHandler

	// Unregister requests from a closing connection
	unregister chan *ConnectionHandler

	serviceDescription *ServiceDescription
	localService       *ServiceDetails

	// The list of paired devices
	registeredServices []ServiceDetails

	// The web server for handling incoming websocket connections
	httpServer *http.Server

	// contains a websocket connection per connected SKI
	connectedServices map[string]*websocket.Conn

	// The zeroconf service for mDNS related tasks
	zc *zeroconf.Server

	zcEntries chan *zeroconf.ServiceEntry

	connectionDelegate ConnectionHandlerDelegate

	mux sync.Mutex
}

func newConnectionsHub(serviceDescription *ServiceDescription, localService *ServiceDetails, connectionDelegate ConnectionHandlerDelegate) *connectionsHub {
	return &connectionsHub{
		connections:        make(map[string]*ConnectionHandler),
		register:           make(chan *ConnectionHandler),
		unregister:         make(chan *ConnectionHandler),
		zcEntries:          make(chan *zeroconf.ServiceEntry),
		registeredServices: make([]ServiceDetails, 0),
		serviceDescription: serviceDescription,
		localService:       localService,
		connectionDelegate: connectionDelegate,
	}
}

// start the ConnectionsHub with all its services
func (h *connectionsHub) start() {
	go h.run()

	// start the websocket server
	go func() {
		if err := h.startWebsocketServer(); err != nil {
			fmt.Println("Error during websocket server starting: ", err)
		}
	}()

	// start mDNS announcement
	if err := h.mdnsRegister(); err != nil {
		fmt.Println(err)
	}

	// handle found mDNS entries
	go h.handleMdnsBrowseEntries(h.zcEntries)

	// Automatically search and connect to services with the same setting
	if h.serviceDescription.RegisterAutoAccept {
		go h.connectRemoteServices()
	}
}

// handle (dis-)connecting remote services
func (h *connectionsHub) run() {
	for {
		select {
		// connect to a paired service
		case c := <-h.register:
			// SHIP 12.2.2 recommends that the connection initiated with the higher SKI should retain the connection
			if existingC := h.connectionForSKI(c.remoteService.SKI); existingC != nil {
				// The connection initiated by the higher SKI should retain the connection
				// and the other one should be closed
				if (c.localService.SKI > c.remoteService.SKI && c.role == ShipRoleClient) ||
					(c.localService.SKI < c.remoteService.SKI && c.role == ShipRoleServer) {
					existingC.conn.Close()

					h.mux.Lock()
					delete(h.connections, c.remoteService.SKI)
					h.mux.Unlock()
				} else {
					c.conn.Close()
					continue
				}
			}

			h.mux.Lock()
			h.connections[c.remoteService.SKI] = c
			h.mux.Unlock()

			c.handleConnection()

			// TODO: shutdown mDNS if this is not a CEM, auto accept is disabled and all registered services are connected
		// disconnect from a no longer paired service
		case c := <-h.unregister:
			if chRegistered, ok := h.connections[c.remoteService.SKI]; ok {
				if chRegistered.conn == c.conn {
					h.mux.Lock()
					delete(h.connections, c.remoteService.SKI)
					h.mux.Unlock()
				}
			}
			// TODO: startup mDNS if this is not a CEM, auto accept is disabled and no registered service is connected
		}
	}
}

// return the connection for a specific SKI
func (h *connectionsHub) connectionForSKI(ski string) *ConnectionHandler {
	return h.connections[ski]
}

// close all connections
func (h *connectionsHub) shutdown() {
	h.mdnsShutdown()
	for _, c := range h.connections {
		c.shutdown(true)
	}
}

// return if there is a connection for a SKI
func (h *connectionsHub) isSkiConnected(ski string) bool {
	// The connection with the higher SKI should retain the connection
	_, ok := h.connections[ski]
	return ok
}

// mDNS handling

// Announces the service to the network via mDNS
// A CEM service should always invoke this on startup
// Any other service should only invoke this when it is not connected to a CEM service and on startup
func (h *connectionsHub) mdnsRegister() error {
	var ifaces []net.Interface = nil
	if len(h.serviceDescription.Interfaces) > 0 {
		ifaces = make([]net.Interface, len(h.serviceDescription.Interfaces))
		for i, iface := range h.serviceDescription.Interfaces {
			ifaceInt, err := net.InterfaceByName(iface)
			if err != nil {
				return err
			}
			ifaces[i] = *ifaceInt
		}
	}

	serviceIdentifier := fmt.Sprintf("%s-%s-%s", h.serviceDescription.Brand, h.serviceDescription.Model, h.serviceDescription.SerialNumber)
	if len(h.serviceDescription.Identifier) > 0 {
		serviceIdentifier = h.serviceDescription.Identifier
	}

	mDNSServer, err := zeroconf.Register(
		serviceIdentifier,
		shipZeroConfServiceType,
		shipZeroConfDomain,
		h.serviceDescription.Port,
		[]string{ // SHIP 7.3.2
			"txtvers=1",
			"path=" + shipWebsocketPath,
			"id=" + serviceIdentifier,
			"ski=" + h.localService.SKI,
			"brand=" + h.serviceDescription.Brand,
			"model=" + h.serviceDescription.Model,
			"type=" + string(h.serviceDescription.DeviceType),
			"register=" + fmt.Sprintf("%v", h.serviceDescription.RegisterAutoAccept),
		},
		ifaces,
	)

	if err != nil {
		return fmt.Errorf("mDNS: registration failed: %w", err)
	}

	h.zc = mDNSServer

	return nil
}

// Stops the mDNS announcement on the network
// A CEM service should only invoke this on the service shutdown
// Any other service should invoke this always when it connected to a CEM and on shutdown
func (h *connectionsHub) mdnsShutdown() {
	if h.zc != nil {
		h.zc.Shutdown()
	}
}

// Websocket connection handling

// start the ship websocket server
func (h *connectionsHub) startWebsocketServer() error {
	addr := fmt.Sprintf(":%d", h.serviceDescription.Port)
	fmt.Println("Starting websocket server on ", addr)

	h.httpServer = &http.Server{
		Addr:    addr,
		Handler: h,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{h.serviceDescription.Certificate},
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
		fmt.Println("Error during connection upgrading: ", err)
		return
	}

	// check if the client supports the ship sub protocol
	if conn.Subprotocol() != shipWebsocketSubProtocol {
		fmt.Println("Client does not support the ship sub protocol")
		conn.Close()
		return
	}

	// check if the clients certificate provides a SKI
	if len(r.TLS.PeerCertificates) == 0 {
		fmt.Println("Client does not provide a certificate")
		conn.Close()
		return
	}

	ski, err := skiFromCertificate(r.TLS.PeerCertificates[0])
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	remoteService := ServiceDetails{
		SKI: ski,
	}
	// check if we already know this remote service
	if remoteS, err := h.registeredServiceForSKI(ski); err == nil {
		remoteService = remoteS
	}

	connectionHandler := newConnectionHandler(h.unregister, h.connectionDelegate, ShipRoleClient, h.localService, &remoteService, conn)

	h.register <- connectionHandler
}

var connectedServicesRunning bool

// handle resolved mDNS entries
func (h *connectionsHub) handleMdnsBrowseEntries(results <-chan *zeroconf.ServiceEntry) {
	for entry := range results {
		fmt.Println("Found service: ", entry.ServiceInstanceName())

		var deviceType, ski string
		var registerAuto bool

		for _, element := range entry.Text {
			if strings.HasPrefix(element, "type=") {
				deviceType = strings.Split(element, "=")[1]
			}

			if strings.HasPrefix(element, "ski=") {
				ski = strings.Split(element, "=")[1]
			}

			if strings.HasPrefix(element, "register=") {
				registerAuto = (strings.Split(element, "=")[1] == "true")
			}
		}

		fmt.Println("SKI: ", ski, " DeviceType: ", deviceType, " RegisterAuto: ", registerAuto)

		if len(ski) == 0 || len(deviceType) == 0 {
			continue
		}

		// Only try to connect to compatible services
		compatibleDeviceType := false
		for _, element := range h.serviceDescription.RemoteDeviceTypes {
			if string(element) == deviceType {
				compatibleDeviceType = true
				break
			}
		}

		if !compatibleDeviceType {
			continue
		}

		// If local and remote registration are set to auto acceppt, we can connect to the remote service
		if h.serviceDescription.RegisterAutoAccept && registerAuto {
			remoteService := ServiceDetails{
				SKI:                ski,
				registerAutoAccept: true,
				deviceType:         model.DeviceTypeType(deviceType),
			}
			if !h.isSkiConnected(ski) {
				h.connectFoundService(remoteService, entry.HostName, strconv.Itoa(int(entry.Port)))
			}
		} else {
			// Check if the remote service is paired
			registeredService, err := h.registeredServiceForSKI(ski)
			if err == nil && !h.isSkiConnected(ski) {
				h.connectFoundService(registeredService, entry.HostName, strconv.Itoa(int(entry.Port)))
			}
		}

		registeredServiceMissing := false
		for _, service := range h.registeredServices {
			if !h.isSkiConnected(service.SKI) {
				registeredServiceMissing = true
				break
			}
		}
		if !registeredServiceMissing && !h.serviceDescription.RegisterAutoAccept {
			fmt.Println("Exit discovery")
			return
		}

	}
}

// Searches via mDNS for registered SHIP services that are not yet connected
// TODO: This should be done in a seperate struct being triggered by a channel
//   to make sure it is not running multiple times at the same time and gets
//   new remote services information while running safely
func (h *connectionsHub) connectRemoteServices() error {
	h.mux.Lock()
	if connectedServicesRunning {
		h.mux.Unlock()
		return nil
	}
	connectedServicesRunning = true
	h.mux.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	if err = resolver.Browse(ctx, shipZeroConfServiceType, shipZeroConfDomain, h.zcEntries); err != nil {
		return err
	}

	<-ctx.Done()

	h.mux.Lock()
	connectedServicesRunning = false
	h.mux.Unlock()

	return nil
}

// Connect to another EEBUS service
func (h *connectionsHub) connectFoundService(remoteService ServiceDetails, host, port string) error {
	fmt.Println("Connecting to ", remoteService.SKI)

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{h.serviceDescription.Certificate},
			InsecureSkipVerify: true,
			CipherSuites:       ciperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}

	address := fmt.Sprintf("wss://%s:%s", host, port)
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	tlsConn := conn.UnderlyingConn().(*tls.Conn)
	remoteCerts := tlsConn.ConnectionState().PeerCertificates

	if remoteCerts == nil || len(remoteCerts) == 0 || remoteCerts[0].SubjectKeyId == nil {
		// Close connection as we couldn't get the remote SKI
		conn.Close()
		return errors.New("Could not get remote SKI")
	}

	if _, err := skiFromCertificate(remoteCerts[0]); err != nil {
		// Close connection as the remote SKI can't be correct
		conn.Close()
		return err
	}

	remoteSKI := fmt.Sprintf("%0x", remoteCerts[0].SubjectKeyId)

	if remoteSKI != remoteService.SKI {
		conn.Close()
		return errors.New("Remote SKI does not match")
	}

	connectionHandler := newConnectionHandler(h.unregister, h.connectionDelegate, ShipRoleServer, h.localService, &remoteService, conn)

	h.register <- connectionHandler

	return nil
}

func (h *connectionsHub) registeredServiceForSKI(ski string) (ServiceDetails, error) {
	for _, service := range h.registeredServices {
		if service.SKI == ski {
			return service, nil
		}
	}
	return ServiceDetails{}, fmt.Errorf("No registered service found for SKI %s", ski)
}

// Adds a new device to the list of known devices which can be connected to
// and connect it if it is currently not connected
func (h *connectionsHub) registerRemoteService(service ServiceDetails) error {
	h.mux.Lock()
	h.registeredServices = append(h.registeredServices, service)
	h.mux.Unlock()

	if !h.isSkiConnected(service.SKI) {
		go h.connectRemoteServices()
	}

	return nil
}

// Update known device in the list of known devices which can be connected to
func (h *connectionsHub) updateRemoteServiceTrust(ski string, trusted bool) {
	for i, device := range h.registeredServices {
		if device.SKI == ski {
			h.mux.Lock()
			h.registeredServices[i].userTrust = true
			h.mux.Unlock()

			conn := h.connectionForSKI(ski)
			if conn != nil {
				if conn.smeState >= smeHelloState {
					conn.shipTrustChannel <- trusted
				} else {
					conn.remoteService.userTrust = trusted
				}
			} else {
				continue
			}
			break
		}
	}
}

// Remove a device from the list of known devices which can be connected to
// and disconnect it if it is currently connected
func (h *connectionsHub) unregisterRemoteService(ski string) error {
	h.mux.Lock()

	newRegisteredDevice := make([]ServiceDetails, 0)

	for _, device := range h.registeredServices {
		if device.SKI != ski {
			newRegisteredDevice = append(newRegisteredDevice, device)
		}
	}

	h.registeredServices = newRegisteredDevice
	h.mux.Unlock()

	if existingC := h.connectionForSKI(ski); existingC != nil {
		existingC.shutdown(true)
	}

	return nil
}
