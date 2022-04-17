package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
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

	if len(remoteCerts[0].SubjectKeyId) != 20 {
		// Close connection as the remote SKI can't be correct
		conn.Close()
		return errors.New("Remote SKI does not have the proper length")
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

// Create a ship compatible self signed certificate
func CreateCertificate() (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Create the EEBUS service SKI using the private key
	asn1, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	// SHIP 12.2: Required to be created according to RFC 3280 4.2.1.2
	ski := sha1.Sum(asn1)

	skiString := fmt.Sprintf("%0x", ski)
	fmt.Println("Local SKI: ", skiString)

	subject := pkix.Name{
		OrganizationalUnit: []string{"Demo"},
		Organization:       []string{"Demo"},
		Country:            []string{"DE"},
	}

	template := x509.Certificate{
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		SerialNumber:          big.NewInt(1),
		Subject:               subject,
		NotBefore:             time.Now(),                                // Valid starting now
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10), // Valid for 10 years
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SubjectKeyId:          ski[:],
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	tlsCertificate := tls.Certificate{
		Certificate:                  [][]byte{certBytes},
		PrivateKey:                   privateKey,
		SupportedSignatureAlgorithms: []tls.SignatureScheme{tls.ECDSAWithP256AndSHA256},
	}

	return tlsCertificate, nil
}

// start the ship websocket server
func (s *EEBUSService) startServer() error {
	addr := fmt.Sprintf(":%d", s.Port)
	fmt.Println("Starting websocket server on ", addr)
	connectionHandler := &ConnectionHandler{
		Role:           ShipRoleServer,
		ConnectionsHub: s.connectionsHub,
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: connectionHandler,
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

					if len(cert.SubjectKeyId) == 20 {
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
