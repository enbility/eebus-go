package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/gorilla/websocket"
	"github.com/grandcat/zeroconf"
)

const defaultPort int = 4711
const shipWebsocketSubProtocol = "ship" // SHIP 10.2: sub protocol is required for websocket connections

var ciperSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, // SHIP 9.1: required cipher suite
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, // SHIP 9.1: optional cipher suite
}

// A service is the central element of an EEBUS service
// including its websocket server and a zeroconf service.
type EEBUSService struct {
	ManufacturerData model.DeviceClassificationManufacturerDataType
	DeviceType       model.DeviceTypeEnumType
	Port             int
	Certificate      tls.Certificate

	connectionsHub *ConnectionsHub

	httpServer *http.Server
	zc         *zeroconf.Server

	// contains a websocket connection per connected SKI
	connectedServices map[string]*websocket.Conn
}

// Starts the service by initializeing mDNS and the server.
func (s *EEBUSService) Start() {
	if s.Port == 0 {
		s.Port = defaultPort
	}

	s.connectionsHub = newConnectionsHub()
	go s.connectionsHub.run()

	go func() {
		if err := s.startServer(); err != nil {
			fmt.Println("Error during websocket server starting: ", err)
		}
	}()
}

// Shutdown all services and stop the server.
func (s *EEBUSService) Shutdown() {
	s.zc.Shutdown()
}

// Connect to another EEBUS service
func (s *EEBUSService) ConnectToService(host, port string) error {
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
		ConnectionsHub: s.connectionsHub,
	}

	connectionHandler.handleConnection(conn)

	return nil
}

// start the ship websocket server
func (s *EEBUSService) startServer() error {
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

// Handling incoming connection requests
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
		SKI:            fmt.Sprintf("%0x", ski),
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
