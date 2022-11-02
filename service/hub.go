package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/gorilla/websocket"
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

	// Handling mDNS related tasks
	mdns *mdns

	connectionDelegate ConnectionHandlerDelegate

	muxCon sync.Mutex
	muxReg sync.Mutex
}

func newConnectionsHub(serviceDescription *ServiceDescription, localService *ServiceDetails, connectionDelegate ConnectionHandlerDelegate) (*connectionsHub, error) {
	hub := &connectionsHub{
		connections:        make(map[string]*ConnectionHandler),
		register:           make(chan *ConnectionHandler),
		unregister:         make(chan *ConnectionHandler),
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
	go h.run()

	// start the websocket server
	go func() {
		if err := h.startWebsocketServer(); err != nil {
			fmt.Println("Error during websocket server starting: ", err)
		}
	}()

	if err := h.mdns.Announce(); err != nil {
		fmt.Println("Error registering mDNS Service:", err)
	}

	// Automatically search and connect to services with the same setting
	if h.serviceDescription.RegisterAutoAccept {
		h.mdns.RegisterMdnsSearch(h)
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

					h.muxCon.Lock()
					delete(h.connections, c.remoteService.SKI)
					h.muxCon.Unlock()
				} else {
					c.conn.Close()
					continue
				}
			}

			h.muxCon.Lock()
			h.connections[c.remoteService.SKI] = c
			h.muxCon.Unlock()

			c.startup()

			// shutdown mDNS if this is not a CEM
			if c.localService.deviceType != model.DeviceTypeTypeEnergyManagementSystem {
				h.mdns.Unannounce()
				h.mdns.UnregisterMdnsSearch(h)
			}

		// disconnect from a no longer connected or paired service
		case c := <-h.unregister:
			h.muxCon.Lock()
			chRegistered, ok := h.connections[c.remoteService.SKI]
			h.muxCon.Unlock()

			if ok {
				if chRegistered.conn == c.conn {
					h.muxCon.Lock()
					delete(h.connections, c.remoteService.SKI)
					h.muxCon.Unlock()
				}
			}
			// startup mDNS if a paired service is not connected
			if len(h.connections) == 0 && len(h.registeredServices) > 0 {
				fmt.Println("Starting mDNS")
				// if this is not a CEM also start the mDNS announcement
				if c.localService.deviceType != model.DeviceTypeTypeEnergyManagementSystem {
					_ = h.mdns.Announce()
				}
				h.mdns.RegisterMdnsSearch(h)
			}
		}
	}
}

// return the connection for a specific SKI
func (h *connectionsHub) connectionForSKI(ski string) *ConnectionHandler {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	return h.connections[ski]
}

// close all connections
func (h *connectionsHub) shutdown() {
	h.muxCon.Lock()
	defer h.muxCon.Unlock()

	h.mdns.shutdown()
	for _, c := range h.connections {
		c.shutdown(true)
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

	ski = util.NormalizeSKI(ski)
	fmt.Println("Incoming connection request from ", ski)

	// Check if the remote service is paired
	_, err = h.registeredServiceForSKI(ski)
	if err != nil {
		fmt.Println("SKI is not registered!")
		return
	}

	remoteService := &ServiceDetails{
		SKI: ski,
	}
	// check if we already know this remote service
	if remoteS, err := h.registeredServiceForSKI(ski); err == nil {
		remoteService = remoteS
	}

	connectionHandler := newConnectionHandler(h.unregister, h.connectionDelegate, ShipRoleServer, h.localService, remoteService, conn)

	h.register <- connectionHandler
}

// Connect to another EEBUS service
func (h *connectionsHub) connectFoundService(remoteService *ServiceDetails, host, port string) error {
	if h.isSkiConnected(remoteService.SKI) {
		return nil
	}

	fmt.Println("Initiating connection to", remoteService.SKI)

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

	if len(remoteCerts) == 0 || remoteCerts[0].SubjectKeyId == nil {
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

	connectionHandler := newConnectionHandler(h.unregister, h.connectionDelegate, ShipRoleClient, h.localService, remoteService, conn)

	h.register <- connectionHandler

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
	h.muxReg.Lock()

	// standardize the provided SKI strings
	service.SKI = util.NormalizeSKI(service.SKI)
	h.registeredServices = append(h.registeredServices, service)

	h.muxReg.Unlock()

	if !h.isSkiConnected(service.SKI) {
		h.mdns.RegisterMdnsSearch(h)
	}
}

// Update known device in the list of known devices which can be connected to
func (h *connectionsHub) updateRemoteServiceTrust(ski string, trusted bool) {
	h.muxReg.Lock()
	defer h.muxReg.Unlock()

	for i, device := range h.registeredServices {
		if device.SKI == ski {
			h.registeredServices[i].userTrust = true

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
		existingC.shutdown(true)
	}

	return nil
}

// Process reported mDNS services
func (h *connectionsHub) ReportMdnsEntries(entries map[string]MdnsEntry) {
	for ski, entry := range entries {
		// check if this ski is already connected
		if h.isSkiConnected(ski) {
			continue
		}

		var remoteService *ServiceDetails
		var err error

		// If local and remote registration are set to auto acceppt, we can connect to the remote service
		if h.serviceDescription.RegisterAutoAccept && entry.Register {
			remoteService = &ServiceDetails{
				SKI:                ski,
				registerAutoAccept: true,
				deviceType:         model.DeviceTypeType(entry.Type),
			}
		} else {
			// Check if the remote service is paired
			remoteService, err = h.registeredServiceForSKI(ski)
			if err != nil {
				continue
			}
		}

		fmt.Println("Trying to connect to", ski, "at", entry.Host)
		if err = h.connectFoundService(remoteService, entry.Host, strconv.Itoa(entry.Port)); err != nil {
			// connecting via the host failed, so try all of the provided addresses
			for _, address := range entry.Addresses {
				fmt.Println("Trying to connect to", ski, "at", address)
				if err = h.connectFoundService(remoteService, address.String(), strconv.Itoa(entry.Port)); err == nil {
					break
				}
			}
			if err != nil {
				continue
			}
		}

		h.muxReg.Lock()
		registeredServiceMissing := false
		for _, service := range h.registeredServices {
			if !h.isSkiConnected(service.SKI) {
				registeredServiceMissing = true
				break
			}
		}
		h.muxReg.Unlock()

		if !registeredServiceMissing && !h.serviceDescription.RegisterAutoAccept {
			h.mdns.UnregisterMdnsSearch(h)
			break
		}
	}
}
