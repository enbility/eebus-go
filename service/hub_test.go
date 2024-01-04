package service

import (
	"errors"
	"testing"
	"time"

	"github.com/enbility/eebus-go/ship"
	"github.com/enbility/eebus-go/spine/model"
	gomock "github.com/golang/mock/gomock"
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
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
	}
	ski := "test"

	paired := sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)

	// mark it as connected, so mDNS is not triggered
	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	sut.registerConnection(con)
	sut.RegisterRemoteSKI(ski, true)

	paired = sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), true, paired)

	// remove the connection, so the test doesn't try to close it
	delete(sut.connections, ski)
	sut.RegisterRemoteSKI(ski, false)
	paired = sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)
}

func (s *HubSuite) Test_HandleConnecitonClosed() {
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
	}
	ski := "test"

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}

	sut.HandleConnectionClosed(con, false)

	sut.registerConnection(con)

	sut.HandleConnectionClosed(con, true)

	assert.Equal(s.T(), 0, len(sut.connections))
}

func (s *HubSuite) Test_Mdns() {
	localService := ServiceDetails{
		DeviceType: model.DeviceTypeTypeElectricitySupplySystem,
	}
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		remoteServices:           make(map[string]*ServiceDetails),
		localService:             &localService,
		mdns:                     s.mdnsService,
		serviceProvider:          s.serviceProvider,
	}
	sut.checkRestartMdnsSearch()

	pairedServices := sut.numberPairedServices()
	assert.Equal(s.T(), 0, len(sut.connections))
	assert.Equal(s.T(), 0, pairedServices)

	ski := "testski"

	sut.RegisterRemoteSKI(ski, true)
	pairedServices = sut.numberPairedServices()
	assert.Equal(s.T(), 0, len(sut.connections))
	assert.Equal(s.T(), 1, pairedServices)

	sut.StartBrowseMdnsSearch()

	sut.StopBrowseMdnsSearch()
}

func (s *HubSuite) Test_Ship() {
	localService := ServiceDetails{
		DeviceType: model.DeviceTypeTypeElectricitySupplySystem,
	}
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		remoteServices:           make(map[string]*ServiceDetails),
		localService:             &localService,
		mdns:                     s.mdnsService,
		serviceProvider:          s.serviceProvider,
	}

	ski := "testski"

	sut.HandleShipHandshakeStateUpdate(ski, ship.ShipState{
		State: ship.SmeStateError,
		Error: errors.New("test"),
	})

	sut.HandleShipHandshakeStateUpdate(ski, ship.ShipState{
		State: ship.SmeHelloStateOk,
	})

	sut.ReportServiceShipID(ski, "test")

	trust := sut.AllowWaitingForTrust(ski)
	assert.Equal(s.T(), true, trust)

	trust = sut.AllowWaitingForTrust("test")
	assert.Equal(s.T(), false, trust)

	detail := sut.PairingDetailForSki(ski)
	assert.NotNil(s.T(), detail)

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	sut.registerConnection(con)

	detail = sut.PairingDetailForSki(ski)
	assert.NotNil(s.T(), detail)
}

func (s *HubSuite) Test_MapShipMessageExchangeState() {
	sut := connectionsHubImpl{}

	ski := "test"

	state := sut.mapShipMessageExchangeState(ship.CmiStateInitStart, ski)
	assert.Equal(s.T(), ConnectionStateQueued, state)

	state = sut.mapShipMessageExchangeState(ship.CmiStateClientSend, ski)
	assert.Equal(s.T(), ConnectionStateInitiated, state)

	state = sut.mapShipMessageExchangeState(ship.SmeHelloStateReadyInit, ski)
	assert.Equal(s.T(), ConnectionStateInProgress, state)

	state = sut.mapShipMessageExchangeState(ship.SmeHelloStatePendingInit, ski)
	assert.Equal(s.T(), ConnectionStateReceivedPairingRequest, state)

	state = sut.mapShipMessageExchangeState(ship.SmeHelloStateOk, ski)
	assert.Equal(s.T(), ConnectionStateTrusted, state)

	state = sut.mapShipMessageExchangeState(ship.SmeHelloStateAbort, ski)
	assert.Equal(s.T(), ConnectionStateNone, state)

	state = sut.mapShipMessageExchangeState(ship.SmeHelloStateRemoteAbortDone, ski)
	assert.Equal(s.T(), ConnectionStateRemoteDeniedTrust, state)

	state = sut.mapShipMessageExchangeState(ship.SmePinStateCheckInit, ski)
	assert.Equal(s.T(), ConnectionStatePin, state)

	state = sut.mapShipMessageExchangeState(ship.SmeAccessMethodsRequest, ski)
	assert.Equal(s.T(), ConnectionStateInProgress, state)

	state = sut.mapShipMessageExchangeState(ship.SmeStateComplete, ski)
	assert.Equal(s.T(), ConnectionStateCompleted, state)

	state = sut.mapShipMessageExchangeState(ship.SmeStateError, ski)
	assert.Equal(s.T(), ConnectionStateError, state)

	state = sut.mapShipMessageExchangeState(ship.SmeProtHStateTimeout, ski)
	assert.Equal(s.T(), ConnectionStateInProgress, state)
}

func (s *HubSuite) Test_DisconnectSKI() {
	sut := connectionsHubImpl{
		connections: make(map[string]*ship.ShipConnection),
	}
	ski := "test"
	sut.DisconnectSKI(ski, "none")
}

func (s *HubSuite) Test_RegisterConnection() {
	ski := "12af9e"
	localService := NewServiceDetails(ski)

	sut := connectionsHubImpl{
		connections:  make(map[string]*ship.ShipConnection),
		mdns:         s.mdnsService,
		localService: localService,
	}

	ski = "test"
	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	sut.registerConnection(con)
	assert.Equal(s.T(), 1, len(sut.connections))
	con = sut.connectionForSKI(ski)
	assert.NotNil(s.T(), con)
}

func (s *HubSuite) Test_Shutdown() {
	sut := connectionsHubImpl{
		connections: make(map[string]*ship.ShipConnection),
		mdns:        s.mdnsService,
	}
	s.mdnsService.EXPECT().ShutdownMdnsService()
	sut.Shutdown()
}

func (s *HubSuite) Test_IncreaseConnectionAttemptCounter() {

	// we just need a dummy for this test
	sut := connectionsHubImpl{
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	for _, test := range s.tests {
		sut.increaseConnectionAttemptCounter(ski)

		sut.muxConAttempt.Lock()
		counter, exists := sut.connectionAttemptCounter[ski]
		timeRange := connectionInitiationDelayTimeRanges[counter]
		sut.muxConAttempt.Unlock()

		assert.Equal(s.T(), true, exists)
		assert.Equal(s.T(), test.timeRange.min, timeRange.min)
		assert.Equal(s.T(), test.timeRange.max, timeRange.max)
	}
}

func (s *HubSuite) Test_RemoveConnectionAttemptCounter() {
	// we just need a dummy for this test
	sut := connectionsHubImpl{
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	sut.increaseConnectionAttemptCounter(ski)
	_, exists := sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), true, exists)

	sut.removeConnectionAttemptCounter(ski)
	_, exists = sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), false, exists)
}

func (s *HubSuite) Test_GetCurrentConnectionAttemptCounter() {
	// we just need a dummy for this test
	sut := connectionsHubImpl{
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	sut.increaseConnectionAttemptCounter(ski)
	_, exists := sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), exists, true)
	sut.increaseConnectionAttemptCounter(ski)

	value, exists := sut.getCurrentConnectionAttemptCounter(ski)
	assert.Equal(s.T(), 1, value)
	assert.Equal(s.T(), true, exists)
}

func (s *HubSuite) Test_GetConnectionInitiationDelayTime() {
	// we just need a dummy for this test
	ski := "12af9e"
	localService := NewServiceDetails(ski)
	sut := connectionsHubImpl{
		localService:             localService,
		connectionAttemptCounter: make(map[string]int),
	}

	counter, duration := sut.getConnectionInitiationDelayTime(ski)
	assert.Equal(s.T(), 0, counter)
	assert.LessOrEqual(s.T(), float64(s.tests[counter].timeRange.min), float64(duration/time.Second))
	assert.GreaterOrEqual(s.T(), float64(s.tests[counter].timeRange.max), float64(duration/time.Second))
}

func (s *HubSuite) Test_ConnectionAttemptRunning() {
	// we just need a dummy for this test
	ski := "test"
	sut := connectionsHubImpl{
		connectionAttemptRunning: make(map[string]bool),
	}

	sut.setConnectionAttemptRunning(ski, true)
	status := sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), true, status)
	sut.setConnectionAttemptRunning(ski, false)
	status = sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), false, status)
}

func (s *HubSuite) Test_InitiatePairingWithSKI() {
	ski := "test"
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptRunning: make(map[string]bool),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
		mdns:                     s.mdnsService,
	}

	sut.InitiatePairingWithSKI(ski)
	assert.Equal(s.T(), 0, len(sut.connections))

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	sut.registerConnection(con)
	sut.InitiatePairingWithSKI(ski)
	assert.Equal(s.T(), 1, len(sut.connections))
}

func (s *HubSuite) Test_CancelPairingWithSKI() {
	ski := "test"
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptRunning: make(map[string]bool),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
		mdns:                     s.mdnsService,
	}

	sut.CancelPairingWithSKI(ski)
	assert.Equal(s.T(), 0, len(sut.connections))
	assert.Equal(s.T(), 0, len(sut.connectionAttemptRunning))

	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	sut.registerConnection(con)
	assert.Equal(s.T(), 1, len(sut.connections))

	sut.CancelPairingWithSKI(ski)
	assert.Equal(s.T(), 0, len(sut.connectionAttemptRunning))
}

func (s *HubSuite) Test_ReportMdnsEntries() {
	localService := &ServiceDetails{
		SKI: "localSKI",
	}
	sut := connectionsHubImpl{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		connectionAttemptRunning: make(map[string]bool),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
		localService:             localService,
	}

	testski1 := "test1"
	testski2 := "test2"

	entries := make(map[string]*MdnsEntry)

	s.serviceProvider.EXPECT().VisibleMDNSRecordsUpdated(gomock.Any()).AnyTimes()
	sut.ReportMdnsEntries(entries)

	entries[testski1] = &MdnsEntry{
		Ski: testski1,
	}
	service1 := sut.ServiceForSKI(testski1)
	service1.Trusted = true
	service1.IPv4 = "127.0.0.1"

	entries[testski2] = &MdnsEntry{
		Ski: testski2,
	}
	service2 := sut.ServiceForSKI(testski2)
	service2.Trusted = true
	service2.IPv4 = "127.0.0.1"

	sut.ReportMdnsEntries(entries)
}
