package service

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/enbility/eebus-go/ship"
	"github.com/enbility/eebus-go/spine/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

	serviceProvider *MockServiceProvider
	mdnsService     *MockMdnsService

	tests []testStruct

	sut *connectionsHubImpl
}

func (s *HubSuite) SetupSuite() {
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

	s.serviceProvider = NewMockServiceProvider(ctrl)
	s.serviceProvider.EXPECT().RemoteSKIConnected(gomock.Any()).AnyTimes()
	s.serviceProvider.EXPECT().ServiceShipIDUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	s.serviceProvider.EXPECT().ServicePairingDetailUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	s.serviceProvider.EXPECT().RemoteSKIDisconnected(gomock.Any()).AnyTimes()
	s.serviceProvider.EXPECT().AllowWaitingForTrust(gomock.Any()).AnyTimes()

	s.mdnsService = NewMockMdnsService(ctrl)
	s.mdnsService.EXPECT().SetupMdnsService().AnyTimes()
	s.mdnsService.EXPECT().AnnounceMdnsEntry().AnyTimes()
	s.mdnsService.EXPECT().UnannounceMdnsEntry().AnyTimes()
	s.mdnsService.EXPECT().RegisterMdnsSearch(gomock.Any()).AnyTimes()
	s.mdnsService.EXPECT().UnregisterMdnsSearch(gomock.Any()).AnyTimes()
}

func (s *HubSuite) BeforeTest(suiteName, testName string) {
	localService := &ServiceDetails{
		SKI: "localSKI",
	}

	s.sut = &connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		connectionAttemptRunning: make(map[string]bool),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
		localService:             localService,
		mdns:                     s.mdnsService,
	}

	certificate, _ := CreateCertificate("unit", "org", "DE", "CN")
	s.sut.configuration, _ = NewConfiguration("vendor", "brand", "model", "serial",
		model.DeviceTypeTypeGeneric, []model.EntityTypeType{model.EntityTypeTypeCEM},
		4567, certificate, 230, time.Second*4)
}

func (s *HubSuite) Test_NewConnectionsHub() {
	ski := "12af9e"
	localService := NewServiceDetails(ski)
	configuration := &Configuration{
		interfaces: []string{"en0"},
	}

	hub := newConnectionsHub(s.serviceProvider, s.mdnsService, nil, configuration, localService)
	assert.NotNil(s.T(), hub)

	hub.Start()
}

func (s *HubSuite) Test_IsRemoteSKIPaired() {
	ski := "test"

	paired := s.sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)

	// mark it as connected, so mDNS is not triggered
	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	s.sut.registerConnection(con)
	s.sut.RegisterRemoteSKI(ski, true)

	paired = s.sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), true, paired)

	// remove the connection, so the test doesn't try to close it
	delete(s.sut.connections, ski)
	s.sut.RegisterRemoteSKI(ski, false)
	paired = s.sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)
}

func (s *HubSuite) Test_HandleConnecitonClosed() {
	ski := "test"

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}

	s.sut.HandleConnectionClosed(con, false)

	s.sut.registerConnection(con)

	s.sut.HandleConnectionClosed(con, true)

	assert.Equal(s.T(), 0, len(s.sut.connections))
}

func (s *HubSuite) Test_Mdns() {
	s.sut.checkRestartMdnsSearch()

	pairedServices := s.sut.numberPairedServices()
	assert.Equal(s.T(), 0, len(s.sut.connections))
	assert.Equal(s.T(), 0, pairedServices)

	ski := "testski"

	s.sut.RegisterRemoteSKI(ski, true)
	pairedServices = s.sut.numberPairedServices()
	assert.Equal(s.T(), 0, len(s.sut.connections))
	assert.Equal(s.T(), 1, pairedServices)

	s.sut.StartBrowseMdnsSearch()

	s.sut.StopBrowseMdnsSearch()
}

func (s *HubSuite) Test_Ship() {
	ski := "testski"

	s.sut.HandleShipHandshakeStateUpdate(ski, ship.ShipState{
		State: ship.SmeStateError,
		Error: errors.New("test"),
	})

	s.sut.HandleShipHandshakeStateUpdate(ski, ship.ShipState{
		State: ship.SmeHelloStateOk,
	})

	s.sut.ReportServiceShipID(ski, "test")

	trust := s.sut.AllowWaitingForTrust(ski)
	assert.Equal(s.T(), true, trust)

	trust = s.sut.AllowWaitingForTrust("test")
	assert.Equal(s.T(), false, trust)

	detail := s.sut.PairingDetailForSki(ski)
	assert.NotNil(s.T(), detail)

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	s.sut.registerConnection(con)

	detail = s.sut.PairingDetailForSki(ski)
	assert.NotNil(s.T(), detail)
}

func (s *HubSuite) Test_MapShipMessageExchangeState() {
	ski := "test"

	state := s.sut.mapShipMessageExchangeState(ship.CmiStateInitStart, ski)
	assert.Equal(s.T(), ConnectionStateQueued, state)

	state = s.sut.mapShipMessageExchangeState(ship.CmiStateClientSend, ski)
	assert.Equal(s.T(), ConnectionStateInitiated, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeHelloStateReadyInit, ski)
	assert.Equal(s.T(), ConnectionStateInProgress, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeHelloStatePendingInit, ski)
	assert.Equal(s.T(), ConnectionStateReceivedPairingRequest, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeHelloStateOk, ski)
	assert.Equal(s.T(), ConnectionStateTrusted, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeHelloStateAbort, ski)
	assert.Equal(s.T(), ConnectionStateNone, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeHelloStateRemoteAbortDone, ski)
	assert.Equal(s.T(), ConnectionStateRemoteDeniedTrust, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmePinStateCheckInit, ski)
	assert.Equal(s.T(), ConnectionStatePin, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeAccessMethodsRequest, ski)
	assert.Equal(s.T(), ConnectionStateInProgress, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeStateComplete, ski)
	assert.Equal(s.T(), ConnectionStateCompleted, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeStateError, ski)
	assert.Equal(s.T(), ConnectionStateError, state)

	state = s.sut.mapShipMessageExchangeState(ship.SmeProtHStateTimeout, ski)
	assert.Equal(s.T(), ConnectionStateInProgress, state)
}

func (s *HubSuite) Test_DisconnectSKI() {
	ski := "test"
	s.sut.DisconnectSKI(ski, "none")
}

func (s *HubSuite) Test_RegisterConnection() {
	ski := "test"
	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	s.sut.registerConnection(con)
	assert.Equal(s.T(), 1, len(s.sut.connections))
	con = s.sut.connectionForSKI(ski)
	assert.NotNil(s.T(), con)
}

func (s *HubSuite) Test_Shutdown() {
	s.mdnsService.EXPECT().ShutdownMdnsService()
	s.sut.Shutdown()
}

func (s *HubSuite) Test_VerifyPeerCertificate() {
	testCert, _ := CreateCertificate("unit", "org", "DE", "CN")
	var rawCerts [][]byte
	rawCerts = append(rawCerts, testCert.Certificate...)
	err := s.sut.verifyPeerCertificate(rawCerts, nil)
	assert.Nil(s.T(), err)

	rawCerts = nil
	rawCerts = append(rawCerts, []byte{100})
	err = s.sut.verifyPeerCertificate(rawCerts, nil)
	assert.NotNil(s.T(), err)

	rawCerts = nil
	invalidCert, _ := CreateInvalidCertificate("unit", "org", "DE", "CN")
	rawCerts = append(rawCerts, invalidCert.Certificate...)

	err = s.sut.verifyPeerCertificate(rawCerts, nil)
	assert.NotNil(s.T(), err)
}

func (s *HubSuite) Test_ServeHTTP() {
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
	server.Close()

	server = httptest.NewUnstartedServer(s.sut)
	server.TLS = &tls.Config{
		Certificates:       []tls.Certificate{s.sut.configuration.certificate},
		ClientAuth:         tls.RequireAnyClientCert,
		CipherSuites:       ciperSuites,
		InsecureSkipVerify: true,
	}
	server.StartTLS()
	wsURL = strings.Replace(server.URL, "https://", "wss://", -1)

	invalidCert, _ := CreateInvalidCertificate("unit", "org", "DE", "CN")
	dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{invalidCert},
			InsecureSkipVerify: true,
			CipherSuites:       ciperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}
	con, _, err = dialer.Dial(wsURL, nil)
	assert.Nil(s.T(), err)

	con.Close()

	validCert, _ := CreateCertificate("unit", "org", "DE", "CN")
	dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{validCert},
			InsecureSkipVerify: true,
			CipherSuites:       ciperSuites,
		},
		Subprotocols: []string{shipWebsocketSubProtocol},
	}
	con, _, err = dialer.Dial(wsURL, nil)
	assert.Nil(s.T(), err)

	con.Close()
	server.Close()
}

func (s *HubSuite) Test_IncreaseConnectionAttemptCounter() {
	ski := "test"

	for _, test := range s.tests {
		s.sut.increaseConnectionAttemptCounter(ski)

		s.sut.muxConAttempt.Lock()
		counter, exists := s.sut.connectionAttemptCounter[ski]
		timeRange := connectionInitiationDelayTimeRanges[counter]
		s.sut.muxConAttempt.Unlock()

		assert.Equal(s.T(), true, exists)
		assert.Equal(s.T(), test.timeRange.min, timeRange.min)
		assert.Equal(s.T(), test.timeRange.max, timeRange.max)
	}
}

func (s *HubSuite) Test_RemoveConnectionAttemptCounter() {
	ski := "test"

	s.sut.increaseConnectionAttemptCounter(ski)
	_, exists := s.sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), true, exists)

	s.sut.removeConnectionAttemptCounter(ski)
	_, exists = s.sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), false, exists)
}

func (s *HubSuite) Test_GetCurrentConnectionAttemptCounter() {
	ski := "test"

	s.sut.increaseConnectionAttemptCounter(ski)
	_, exists := s.sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), exists, true)
	s.sut.increaseConnectionAttemptCounter(ski)

	value, exists := s.sut.getCurrentConnectionAttemptCounter(ski)
	assert.Equal(s.T(), 1, value)
	assert.Equal(s.T(), true, exists)
}

func (s *HubSuite) Test_GetConnectionInitiationDelayTime() {
	ski := "test"

	counter, duration := s.sut.getConnectionInitiationDelayTime(ski)
	assert.Equal(s.T(), 0, counter)
	assert.LessOrEqual(s.T(), float64(s.tests[counter].timeRange.min), float64(duration/time.Second))
	assert.GreaterOrEqual(s.T(), float64(s.tests[counter].timeRange.max), float64(duration/time.Second))
}

func (s *HubSuite) Test_ConnectionAttemptRunning() {
	ski := "test"

	s.sut.setConnectionAttemptRunning(ski, true)
	status := s.sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), true, status)
	s.sut.setConnectionAttemptRunning(ski, false)
	status = s.sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), false, status)
}

func (s *HubSuite) Test_InitiatePairingWithSKI() {
	ski := "test"

	s.sut.InitiatePairingWithSKI(ski)
	assert.Equal(s.T(), 0, len(s.sut.connections))

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	s.sut.registerConnection(con)
	s.sut.InitiatePairingWithSKI(ski)
	assert.Equal(s.T(), 1, len(s.sut.connections))
}

func (s *HubSuite) Test_CancelPairingWithSKI() {
	ski := "test"

	s.sut.CancelPairingWithSKI(ski)
	assert.Equal(s.T(), 0, len(s.sut.connections))
	assert.Equal(s.T(), 0, len(s.sut.connectionAttemptRunning))

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	s.sut.registerConnection(con)
	assert.Equal(s.T(), 1, len(s.sut.connections))

	s.sut.CancelPairingWithSKI(ski)
	assert.Equal(s.T(), 0, len(s.sut.connectionAttemptRunning))
}

func (s *HubSuite) Test_ReportMdnsEntries() {
	testski1 := "test1"
	testski2 := "test2"

	entries := make(map[string]*MdnsEntry)

	s.serviceProvider.EXPECT().VisibleMDNSRecordsUpdated(gomock.Any()).AnyTimes()
	s.sut.ReportMdnsEntries(entries)

	entries[testski1] = &MdnsEntry{
		Ski: testski1,
	}
	service1 := s.sut.ServiceForSKI(testski1)
	service1.Trusted = true
	service1.IPv4 = "127.0.0.1"

	entries[testski2] = &MdnsEntry{
		Ski: testski2,
	}
	service2 := s.sut.ServiceForSKI(testski2)
	service2.Trusted = true
	service2.IPv4 = "127.0.0.1"

	s.sut.ReportMdnsEntries(entries)
}
