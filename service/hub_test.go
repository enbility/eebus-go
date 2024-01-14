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
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/cert"
	"github.com/enbility/eebus-go/mocks"
	shipapi "github.com/enbility/ship-go/api"
	shipmocks "github.com/enbility/ship-go/mocks"
	shipmodel "github.com/enbility/ship-go/model"
	"github.com/enbility/spine-go/model"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestHubSuite(t *testing.T) {
	suite.Run(t, new(HubSuite))
}

type testStruct struct {
	counter   int
	timeRange connectionInitiationDelayTimeRange
}

type HubSuite struct {
	suite.Suite

	serviceProvider *mocks.MockServiceProvider
	mdnsService     *mocks.MockMdnsService

	// serviceProvider  *mocks.ServiceProvider
	// mdnsService      *mocks.MdnsService
	shipConnection   *shipmocks.ShipConnection
	wsDataConnection *shipmocks.WebsocketDataConnection

	remoteSki string

	tests []testStruct

	sut *connectionsHubImpl
}

func (s *HubSuite) BeforeTest(suiteName, testName string) {
	s.remoteSki = "remotetestski"

	s.tests = []testStruct{
		{0, connectionInitiationDelayTimeRanges[0]},
		{1, connectionInitiationDelayTimeRanges[1]},
		{2, connectionInitiationDelayTimeRanges[2]},
		{3, connectionInitiationDelayTimeRanges[2]},
		{4, connectionInitiationDelayTimeRanges[2]},
		{5, connectionInitiationDelayTimeRanges[2]},
		{6, connectionInitiationDelayTimeRanges[2]},
		{7, connectionInitiationDelayTimeRanges[2]},
		{8, connectionInitiationDelayTimeRanges[2]},
		{9, connectionInitiationDelayTimeRanges[2]},
		{10, connectionInitiationDelayTimeRanges[2]},
	}

	ctrl := gomock.NewController(s.T())
	// use gomock mocks instead of mockery, as those will panic with a data race error in these tests

	s.serviceProvider = mocks.NewMockServiceProvider(ctrl)
	// s.serviceProvider = mocks.NewServiceProvider(s.T())
	s.serviceProvider.EXPECT().RemoteSKIConnected(gomock.Any()).Return().AnyTimes()
	s.serviceProvider.EXPECT().ServiceShipIDUpdate(gomock.Any(), gomock.Any()).Return().AnyTimes()
	s.serviceProvider.EXPECT().ServicePairingDetailUpdate(gomock.Any(), gomock.Any()).Return().AnyTimes()
	s.serviceProvider.EXPECT().RemoteSKIDisconnected(gomock.Any()).Return().AnyTimes()
	s.serviceProvider.EXPECT().AllowWaitingForTrust(gomock.Any()).Return(false).AnyTimes()

	s.mdnsService = mocks.NewMockMdnsService(ctrl)
	// s.mdnsService = mocks.NewMdnsService(s.T())
	s.mdnsService.EXPECT().SetupMdnsService().Return(nil).AnyTimes()
	s.mdnsService.EXPECT().AnnounceMdnsEntry().Return(nil).AnyTimes()
	s.mdnsService.EXPECT().UnannounceMdnsEntry().Return().AnyTimes()
	s.mdnsService.EXPECT().RegisterMdnsSearch(gomock.Any()).Return().AnyTimes()
	s.mdnsService.EXPECT().UnregisterMdnsSearch(gomock.Any()).Return().AnyTimes()

	s.wsDataConnection = shipmocks.NewWebsocketDataConnection(s.T())

	s.shipConnection = shipmocks.NewShipConnection(s.T())
	s.shipConnection.EXPECT().CloseConnection(mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	s.shipConnection.EXPECT().RemoteSKI().Return(s.remoteSki).Maybe()
	s.shipConnection.EXPECT().ApprovePendingHandshake().Return().Maybe()
	s.shipConnection.EXPECT().AbortPendingHandshake().Return().Maybe()
	s.shipConnection.EXPECT().DataHandler().Return(s.wsDataConnection).Maybe()
	s.shipConnection.EXPECT().ShipHandshakeState().Return(shipmodel.SmeStateComplete, nil).Maybe()

	localService := &api.ServiceDetails{
		SKI: "localSKI",
	}

	s.sut = &connectionsHubImpl{
		connections:              make(map[string]shipapi.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		connectionAttemptRunning: make(map[string]bool),
		remoteServices:           make(map[string]*api.ServiceDetails),
		serviceProvider:          s.serviceProvider,
		localService:             localService,
		mdns:                     s.mdnsService,
	}

	certificate, _ := cert.CreateCertificate("unit", "org", "DE", "CN")
	s.sut.configuration, _ = api.NewConfiguration("vendor", "brand", "model", "serial",
		model.DeviceTypeTypeGeneric, []model.EntityTypeType{model.EntityTypeTypeCEM},
		4567, certificate, 230, time.Second*4)
}

func (s *HubSuite) Test_NewConnectionsHub() {
	ski := "12af9e"
	localService := api.NewServiceDetails(ski)

	configuration := &api.Configuration{}
	configuration.SetInterfaces([]string{"en0"})

	hub := newConnectionsHub(s.serviceProvider, s.mdnsService, nil, configuration, localService)
	assert.NotNil(s.T(), hub)

	hub.Start()
}

func (s *HubSuite) Test_IsRemoteSKIPaired() {
	paired := s.sut.IsRemoteServiceForSKIPaired(s.remoteSki)
	assert.Equal(s.T(), false, paired)

	s.sut.registerConnection(s.shipConnection)
	s.sut.RegisterRemoteSKI(s.remoteSki, true)

	paired = s.sut.IsRemoteServiceForSKIPaired(s.remoteSki)
	assert.Equal(s.T(), true, paired)

	// remove the connection, so the test doesn't try to close it
	delete(s.sut.connections, s.remoteSki)
	s.sut.RegisterRemoteSKI(s.remoteSki, false)
	paired = s.sut.IsRemoteServiceForSKIPaired(s.remoteSki)
	assert.Equal(s.T(), false, paired)
}

func (s *HubSuite) Test_HandleConnecitonClosed() {
	s.sut.HandleConnectionClosed(s.shipConnection, false)

	s.sut.registerConnection(s.shipConnection)

	s.sut.HandleConnectionClosed(s.shipConnection, true)

	assert.Equal(s.T(), 0, len(s.sut.connections))
}

func (s *HubSuite) Test_Mdns() {
	s.sut.checkRestartMdnsSearch()

	pairedServices := s.sut.numberPairedServices()
	assert.Equal(s.T(), 0, len(s.sut.connections))
	assert.Equal(s.T(), 0, pairedServices)

	s.sut.RegisterRemoteSKI(s.remoteSki, true)
	pairedServices = s.sut.numberPairedServices()
	assert.Equal(s.T(), 0, len(s.sut.connections))
	assert.Equal(s.T(), 1, pairedServices)

	s.sut.StartBrowseMdnsSearch()

	s.sut.StopBrowseMdnsSearch()
}

func (s *HubSuite) Test_Ship() {
	s.sut.HandleShipHandshakeStateUpdate(s.remoteSki, shipmodel.ShipState{
		State: shipmodel.SmeStateError,
		Error: errors.New("test"),
	})

	s.sut.HandleShipHandshakeStateUpdate(s.remoteSki, shipmodel.ShipState{
		State: shipmodel.SmeHelloStateOk,
	})

	s.sut.ReportServiceShipID(s.remoteSki, "test")

	trust := s.sut.AllowWaitingForTrust(s.remoteSki)
	assert.Equal(s.T(), true, trust)

	trust = s.sut.AllowWaitingForTrust("test")
	assert.Equal(s.T(), false, trust)

	detail := s.sut.PairingDetailForSki(s.remoteSki)
	assert.NotNil(s.T(), detail)

	s.sut.registerConnection(s.shipConnection)

	detail = s.sut.PairingDetailForSki(s.remoteSki)
	assert.NotNil(s.T(), detail)
}

func (s *HubSuite) Test_MapShipMessageExchangeState() {
	state := s.sut.mapShipMessageExchangeState(shipmodel.CmiStateInitStart, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateQueued, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.CmiStateClientSend, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateInitiated, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeHelloStateReadyInit, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateInProgress, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeHelloStatePendingInit, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateReceivedPairingRequest, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeHelloStateOk, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateTrusted, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeHelloStateAbort, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateNone, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeHelloStateRemoteAbortDone, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateRemoteDeniedTrust, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmePinStateCheckInit, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStatePin, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeAccessMethodsRequest, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateInProgress, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeStateComplete, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateCompleted, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeStateError, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateError, state)

	state = s.sut.mapShipMessageExchangeState(shipmodel.SmeProtHStateTimeout, s.remoteSki)
	assert.Equal(s.T(), api.ConnectionStateInProgress, state)
}

func (s *HubSuite) Test_DisconnectSKI() {
	s.sut.DisconnectSKI(s.remoteSki, "none")
}

func (s *HubSuite) Test_RegisterConnection() {
	s.sut.registerConnection(s.shipConnection)
	assert.Equal(s.T(), 1, len(s.sut.connections))
	con := s.sut.connectionForSKI(s.remoteSki)
	assert.NotNil(s.T(), con)
}

func (s *HubSuite) Test_Shutdown() {
	s.mdnsService.EXPECT().ShutdownMdnsService()
	s.sut.Shutdown()
}

func (s *HubSuite) Test_VerifyPeerCertificate() {
	testCert, _ := cert.CreateCertificate("unit", "org", "DE", "CN")
	var rawCerts [][]byte
	rawCerts = append(rawCerts, testCert.Certificate...)
	err := s.sut.verifyPeerCertificate(rawCerts, nil)
	assert.Nil(s.T(), err)

	rawCerts = nil
	rawCerts = append(rawCerts, []byte{100})
	err = s.sut.verifyPeerCertificate(rawCerts, nil)
	assert.NotNil(s.T(), err)

	rawCerts = nil
	invalidCert, _ := createInvalidCertificate("unit", "org", "DE", "CN")
	rawCerts = append(rawCerts, invalidCert.Certificate...)

	err = s.sut.verifyPeerCertificate(rawCerts, nil)
	assert.NotNil(s.T(), err)
}

func (s *HubSuite) Test_ServeHTTP_01() {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	s.sut.ServeHTTP(w, req)

	server := httptest.NewServer(s.sut)
	wsURL := strings.Replace(server.URL, "http://", "ws://", -1)

	// Connect to the server
	con, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Nil(s.T(), err)
	con.Close()

	dialer := &websocket.Dialer{
		Subprotocols: []string{shipWebsocketSubProtocol},
	}
	con, _, err = dialer.Dial(wsURL, nil)
	assert.Nil(s.T(), err)

	con.Close()
	server.CloseClientConnections()
	server.Close()

	time.Sleep(time.Second)
}

func (s *HubSuite) Test_ServeHTTP_02() {
	server := httptest.NewUnstartedServer(s.sut)
	server.TLS = &tls.Config{
		Certificates:       []tls.Certificate{s.sut.configuration.Certificate()},
		ClientAuth:         tls.RequireAnyClientCert,
		CipherSuites:       cert.CiperSuites,
		InsecureSkipVerify: true,
	}
	server.StartTLS()
	wsURL := strings.Replace(server.URL, "https://", "wss://", -1)

	invalidCert, _ := createInvalidCertificate("unit", "org", "DE", "CN")
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{invalidCert},
			InsecureSkipVerify: true,
			CipherSuites:       cert.CiperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}
	con, _, err := dialer.Dial(wsURL, nil)
	assert.Nil(s.T(), err)

	con.Close()

	validCert, _ := cert.CreateCertificate("unit", "org", "DE", "CN")
	dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{validCert},
			InsecureSkipVerify: true,
			CipherSuites:       cert.CiperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}
	con, _, err = dialer.Dial(wsURL, nil)
	assert.Nil(s.T(), err)

	con.Close()
	server.CloseClientConnections()
	server.Close()

	time.Sleep(time.Second)
}

func (s *HubSuite) Test_ConnectFoundService_01() {
	service := s.sut.ServiceForSKI(s.remoteSki)

	err := s.sut.connectFoundService(service, "localhost", "80")
	assert.NotNil(s.T(), err)

	server := httptest.NewServer(s.sut)
	url, err := url.Parse(server.URL)
	assert.Nil(s.T(), err)

	err = s.sut.connectFoundService(service, url.Hostname(), url.Port())
	assert.NotNil(s.T(), err)

	server.CloseClientConnections()
	server.Close()

	time.Sleep(time.Second)
}

func (s *HubSuite) Test_ConnectFoundService_02() {
	service := s.sut.ServiceForSKI(s.remoteSki)

	server := httptest.NewUnstartedServer(s.sut)
	invalidCert, _ := createInvalidCertificate("unit", "org", "DE", "CN")
	server.TLS = &tls.Config{
		Certificates:       []tls.Certificate{invalidCert},
		ClientAuth:         tls.RequireAnyClientCert,
		CipherSuites:       cert.CiperSuites,
		InsecureSkipVerify: true,
	}
	server.StartTLS()

	url, err := url.Parse(server.URL)
	assert.Nil(s.T(), err)

	err = s.sut.connectFoundService(service, url.Hostname(), url.Port())
	assert.NotNil(s.T(), err)

	server.CloseClientConnections()
	server.Close()

	time.Sleep(time.Second)
}

func (s *HubSuite) Test_ConnectFoundService_03() {
	service := s.sut.ServiceForSKI(s.remoteSki)

	server := httptest.NewUnstartedServer(s.sut)
	server.TLS = &tls.Config{
		Certificates:       []tls.Certificate{s.sut.configuration.Certificate()},
		ClientAuth:         tls.RequireAnyClientCert,
		CipherSuites:       cert.CiperSuites,
		InsecureSkipVerify: true,
	}
	server.StartTLS()

	url, err := url.Parse(server.URL)
	assert.Nil(s.T(), err)

	err = s.sut.connectFoundService(service, url.Hostname(), url.Port())
	assert.NotNil(s.T(), err)

	time.Sleep(time.Second)

	server.CloseClientConnections()
	server.Close()
}

func (s *HubSuite) Test_KeepThisConnection() {
	service := s.sut.ServiceForSKI(s.remoteSki)

	result := s.sut.keepThisConnection(nil, false, service)
	assert.Equal(s.T(), true, result)

	s.sut.registerConnection(s.shipConnection)

	result = s.sut.keepThisConnection(nil, false, service)
	assert.Equal(s.T(), false, result)

	result = s.sut.keepThisConnection(nil, true, service)
	assert.Equal(s.T(), true, result)
}

func (s *HubSuite) Test_prepareConnectionInitiation() {
	entry := &api.MdnsEntry{
		Ski:  s.remoteSki,
		Host: "somehost",
	}
	service := s.sut.ServiceForSKI(s.remoteSki)

	s.sut.prepareConnectionInitation(s.remoteSki, 0, entry)

	s.sut.setConnectionAttemptRunning(s.remoteSki, true)

	counter := s.sut.increaseConnectionAttemptCounter(s.remoteSki)
	assert.Equal(s.T(), 0, counter)
	s.sut.prepareConnectionInitation(s.remoteSki, 0, entry)

	s.sut.RegisterRemoteSKI(s.remoteSki, false)
	service.ConnectionStateDetail().SetState(api.ConnectionStateQueued)

	counter = s.sut.increaseConnectionAttemptCounter(s.remoteSki)
	assert.Equal(s.T(), 0, counter)

	s.sut.prepareConnectionInitation(s.remoteSki, 0, entry)
}

func (s *HubSuite) Test_InitiateConnection() {
	entry := &api.MdnsEntry{
		Ski:  s.remoteSki,
		Host: "somehost",
	}
	service := s.sut.ServiceForSKI(s.remoteSki)

	result := s.sut.initateConnection(service, entry)
	assert.Equal(s.T(), false, result)

	entry.Addresses = []net.IP{[]byte("127.0.0.1")}

	result = s.sut.initateConnection(service, entry)
	assert.Equal(s.T(), false, result)

	s.sut.RegisterRemoteSKI(s.remoteSki, true)
	service.ConnectionStateDetail().SetState(api.ConnectionStateQueued)

	result = s.sut.initateConnection(service, entry)
	assert.Equal(s.T(), false, result)
}

func (s *HubSuite) Test_IncreaseConnectionAttemptCounter() {
	for _, test := range s.tests {
		s.sut.increaseConnectionAttemptCounter(s.remoteSki)

		s.sut.muxConAttempt.Lock()
		counter, exists := s.sut.connectionAttemptCounter[s.remoteSki]
		timeRange := connectionInitiationDelayTimeRanges[counter]
		s.sut.muxConAttempt.Unlock()

		assert.Equal(s.T(), true, exists)
		assert.Equal(s.T(), test.timeRange.min, timeRange.min)
		assert.Equal(s.T(), test.timeRange.max, timeRange.max)
	}
}

func (s *HubSuite) Test_RemoveConnectionAttemptCounter() {
	s.sut.increaseConnectionAttemptCounter(s.remoteSki)
	_, exists := s.sut.connectionAttemptCounter[s.remoteSki]
	assert.Equal(s.T(), true, exists)

	s.sut.removeConnectionAttemptCounter(s.remoteSki)
	_, exists = s.sut.connectionAttemptCounter[s.remoteSki]
	assert.Equal(s.T(), false, exists)
}

func (s *HubSuite) Test_GetCurrentConnectionAttemptCounter() {
	s.sut.increaseConnectionAttemptCounter(s.remoteSki)
	_, exists := s.sut.connectionAttemptCounter[s.remoteSki]
	assert.Equal(s.T(), exists, true)
	s.sut.increaseConnectionAttemptCounter(s.remoteSki)

	value, exists := s.sut.getCurrentConnectionAttemptCounter(s.remoteSki)
	assert.Equal(s.T(), 1, value)
	assert.Equal(s.T(), true, exists)
}

func (s *HubSuite) Test_GetConnectionInitiationDelayTime() {
	counter, duration := s.sut.getConnectionInitiationDelayTime(s.remoteSki)
	assert.Equal(s.T(), 0, counter)
	assert.LessOrEqual(s.T(), float64(s.tests[counter].timeRange.min), float64(duration/time.Second))
	assert.GreaterOrEqual(s.T(), float64(s.tests[counter].timeRange.max), float64(duration/time.Second))
}

func (s *HubSuite) Test_ConnectionAttemptRunning() {
	s.sut.setConnectionAttemptRunning(s.remoteSki, true)
	status := s.sut.isConnectionAttemptRunning(s.remoteSki)
	assert.Equal(s.T(), true, status)
	s.sut.setConnectionAttemptRunning(s.remoteSki, false)
	status = s.sut.isConnectionAttemptRunning(s.remoteSki)
	assert.Equal(s.T(), false, status)
}

func (s *HubSuite) Test_InitiatePairingWithSKI() {
	s.sut.InitiatePairingWithSKI(s.remoteSki)
	assert.Equal(s.T(), 0, len(s.sut.connections))

	s.sut.registerConnection(s.shipConnection)
	s.sut.InitiatePairingWithSKI(s.remoteSki)
	assert.Equal(s.T(), 1, len(s.sut.connections))
}

func (s *HubSuite) Test_CancelPairingWithSKI() {
	s.sut.CancelPairingWithSKI(s.remoteSki)
	assert.Equal(s.T(), 0, len(s.sut.connections))
	assert.Equal(s.T(), 0, len(s.sut.connectionAttemptRunning))

	s.sut.registerConnection(s.shipConnection)
	assert.Equal(s.T(), 1, len(s.sut.connections))

	s.sut.CancelPairingWithSKI(s.remoteSki)
	assert.Equal(s.T(), 0, len(s.sut.connectionAttemptRunning))
}

func (s *HubSuite) Test_ReportMdnsEntries() {
	testski1 := "test1"
	testski2 := "test2"

	entries := make(map[string]*api.MdnsEntry)

	s.serviceProvider.EXPECT().VisibleMDNSRecordsUpdated(gomock.Any()).AnyTimes()
	s.sut.ReportMdnsEntries(entries)

	entries[testski1] = &api.MdnsEntry{
		Ski: testski1,
	}
	service1 := s.sut.ServiceForSKI(testski1)
	service1.Trusted = true
	service1.IPv4 = "127.0.0.1"

	entries[testski2] = &api.MdnsEntry{
		Ski: testski2,
	}
	service2 := s.sut.ServiceForSKI(testski2)
	service2.Trusted = true
	service2.IPv4 = "127.0.0.1"

	s.sut.ReportMdnsEntries(entries)
}

func createInvalidCertificate(organizationalUnit, organization, country, commonName string) (tls.Certificate, error) {
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

	subject := pkix.Name{
		OrganizationalUnit: []string{organizationalUnit},
		Organization:       []string{organization},
		Country:            []string{country},
		CommonName:         commonName,
	}

	// Create a random serial big int value
	maxValue := new(big.Int)
	maxValue.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(maxValue, big.NewInt(1))
	serialNumber, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),                                // Valid starting now
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10), // Valid for 10 years
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SubjectKeyId:          ski[:19],
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
