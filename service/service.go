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
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/gorilla/websocket"
	"github.com/grandcat/zeroconf"
)

const defaultPort int = 4711
const shipWebsocketSubProtocol = "ship" // SHIP 10.2: sub protocol is required for websocket connections
const shipWebsocketPath = "/ship/"
const shipZeroConfServiceType = "_ship._tcp"
const shipZeroConfDomain = "local."

var ciperSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, // SHIP 9.1: required cipher suite
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, // SHIP 9.1: optional cipher suite
}

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSService struct {
	// The brand of the device
	DeviceBrand string

	// The device model
	DeviceModel string

	// The EEBUS device type of the device model
	DeviceType model.DeviceTypeType

	// Serial number of the device
	DeviceSerialNumber string

	// The mDNS service identifier
	// Optional, if not set will be  generated using "DeviceBrand-DeviceModel-DeviceSerialNumber"
	DeviceIdentifier string

	// Network interface to use for the service
	// Optional, if not set all detected interfaces will be used
	Interfaces []string

	// The port address of the websocket server
	Port int

	// The certificate used for the service and its connections
	Certificate tls.Certificate

	// The service SKI
	SKI string

	// Wether remote devices should be automatically accepted
	RemoteDeviceAutoAccept bool

	// Connection Registrations
	connectionsHub *ConnectionsHub

	// The web server for handling incoming websocket connections
	httpServer *http.Server

	// The zeroconf service for mDNS related tasks
	zc *zeroconf.Server

	// contains a websocket connection per connected SKI
	connectedServices map[string]*websocket.Conn
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Start() {
	if s.Port == 0 {
		s.Port = defaultPort
	}

	leaf, err := x509.ParseCertificate(s.Certificate.Certificate[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	ski, err := s.skiFromCertificate(leaf)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.SKI = ski

	fmt.Println("Local SKI: ", ski)

	s.connectionsHub = newConnectionsHub()
	go s.connectionsHub.run()

	go func() {
		if err := s.startWebsocketServer(); err != nil {
			fmt.Println("Error during websocket server starting: ", err)
		}
	}()
}

// Shutdown all services and stop the server.
func (s *EEBUSService) Shutdown() {
	s.MdnsShutdown()

	// Shut down all running connections
	s.connectionsHub.Shutdown()
}

// Connect to another EEBUS service
func (s *EEBUSService) connectToService(host, port string) error {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{s.Certificate},
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

	if _, err := s.skiFromCertificate(remoteCerts[0]); err != nil {
		// Close connection as the remote SKI can't be correct
		conn.Close()
		return err
	}

	remoteSKI := fmt.Sprintf("%0x", remoteCerts[0].SubjectKeyId)

	connectionHandler := &ConnectionHandler{
		Role:           ShipRoleClient,
		SKI:            remoteSKI,
		ShipID:         s.DeviceIdentifier,
		ConnectionsHub: s.connectionsHub,
	}

	connectionHandler.handleConnection(conn)

	return nil
}

// start the ship websocket server
func (s *EEBUSService) startWebsocketServer() error {
	addr := fmt.Sprintf(":%d", s.Port)
	fmt.Println("Starting websocket server on ", addr)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{s.Certificate},
			ClientAuth:   tls.RequireAnyClientCert, // SHIP 9: Client authentication is required
			CipherSuites: ciperSuites,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				skiFound := false
				for _, v := range rawCerts {
					cert, err := x509.ParseCertificate(v)
					if err != nil {
						return err
					}

					if _, err := s.skiFromCertificate(cert); err == nil {
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

	if err := s.httpServer.ListenAndServeTLS("", ""); err != nil {
		return err
	}

	return nil
}

// HTTP Server callback for handling incoming connection requests
func (s *EEBUSService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	ski, err := s.skiFromCertificate(r.TLS.PeerCertificates[0])
	if err != nil {
		fmt.Println(err)
		conn.Close()
	}

	connectionHandler := &ConnectionHandler{
		Role:           ShipRoleServer,
		ConnectionsHub: s.connectionsHub,
		SKI:            ski,
	}

	connectionHandler.handleConnection(conn)
}

func (s *EEBUSService) skiFromCertificate(cert *x509.Certificate) (string, error) {
	// check if the clients certificate provides a SKI
	subjectKeyId := cert.SubjectKeyId
	if len(subjectKeyId) != 20 {
		return "", errors.New("Client certificate does not provide a SKI")
	}

	return fmt.Sprintf("%0x", subjectKeyId), nil
}

// Announces the service to the network via mDNS
// A CEM service should always invoke this on startup
// Any other service should only invoke this when it is not connected to a CEM service and on startup
func (s *EEBUSService) MdnsAnnounce() error {
	var ifaces []net.Interface = nil
	if len(s.Interfaces) > 0 {
		ifaces = make([]net.Interface, len(s.Interfaces))
		for i, iface := range s.Interfaces {
			ifaceInt, err := net.InterfaceByName(iface)
			if err != nil {
				return err
			}
			ifaces[i] = *ifaceInt
		}
	}

	serviceIdentifier := fmt.Sprintf("%s-%s-%s", s.DeviceBrand, s.DeviceModel, s.DeviceSerialNumber)
	if len(s.DeviceIdentifier) > 0 {
		serviceIdentifier = s.DeviceIdentifier
	}

	mDNSServer, err := zeroconf.Register(
		serviceIdentifier,
		shipZeroConfServiceType,
		shipZeroConfDomain,
		s.Port,
		[]string{ // SHIP 7.3.2
			"txtvers=1",
			"path=" + shipWebsocketPath,
			"id=" + serviceIdentifier,
			"ski=" + s.SKI,
			"brand=" + s.DeviceBrand,
			"model=" + s.DeviceModel,
			"type=" + string(s.DeviceType),
			"register=" + fmt.Sprintf("%v", (s.RemoteDeviceAutoAccept == true)),
		},
		ifaces,
	)

	if err != nil {
		return fmt.Errorf("mDNS: registration failed: %w", err)
	}

	s.zc = mDNSServer

	return nil
}

// Stops the mDNS announcement on the network
// A CEM service should only invoke this on the service shutdown
// Any other service should invoke this always when it connected to a CEM and on shutdown
func (s *EEBUSService) MdnsShutdown() {
	if s.zc != nil {
		s.zc.Shutdown()
	}
}

func (s *EEBUSService) ConnectToSKI(ski string) error {
	fmt.Println("Searching for ski: ", ski)

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			for _, typ := range entry.Text {
				if typ == fmt.Sprintf("ski=%s", ski) {
					fmt.Println("Service discovered: ", entry.ServiceInstanceName())

					fmt.Printf("Connecting to %s:%d\n", entry.HostName, entry.Port)
					s.connectToService(entry.HostName, strconv.Itoa(int(entry.Port)))
					fmt.Printf("\n\n")
				}
			}
		}
		fmt.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	if err = resolver.Browse(ctx, shipZeroConfServiceType, shipZeroConfDomain, entries); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
