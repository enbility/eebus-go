package service

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/ship"
	"github.com/enbility/eebus-go/spine/model"
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

	tests []testStruct
}

func (s *HubSuite) SetupSuite() {
	s.tests = []testStruct{
		{0, connectionInitiationDelayTimeRanges[0]},
		{1, connectionInitiationDelayTimeRanges[1]},
		{2, connectionInitiationDelayTimeRanges[2]},
		{3, connectionInitiationDelayTimeRanges[3]},
		{4, connectionInitiationDelayTimeRanges[4]},
		{5, connectionInitiationDelayTimeRanges[5]},
		{6, connectionInitiationDelayTimeRanges[5]},
		{7, connectionInitiationDelayTimeRanges[5]},
		{8, connectionInitiationDelayTimeRanges[5]},
		{9, connectionInitiationDelayTimeRanges[5]},
		{10, connectionInitiationDelayTimeRanges[5]},
	}
}

// Service Provider Interface
var _ serviceProvider = (*HubSuite)(nil)

func (s *HubSuite) RemoteSKIConnected(string)          {}
func (s *HubSuite) RemoteSKIDisconnected(string)       {}
func (s *HubSuite) ReportServiceShipID(string, string) {}

func (s *HubSuite) Test_NewConnectionsHub() {
	ski := "12af9e"
	localService := &ServiceDetails{
		SKI: ski,
	}
	configuration := &Configuration{
		interfaces: []string{"en0"},
	}
	hub := newConnectionsHub(s, nil, configuration, localService)
	assert.NotNil(s.T(), hub)

	hub.start()
}

func (s *HubSuite) Test_IsRemoteSKIPaired() {
	sut := connectionsHub{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	paired := sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)

	service := ServiceDetails{
		SKI: ski,
	}
	// mark it as connected, so mDNS is not triggered
	sut.connections[ski] = &ship.ShipConnection{}
	sut.PairRemoteService(service)

	paired = sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), true, paired)

	// remove the connection, so the test doesn't try to close it
	delete(sut.connections, ski)
	err := sut.UnpairRemoteService(ski)
	assert.Nil(s.T(), err)
	paired = sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)
}

func (s *HubSuite) Test_CheckRestartMdnsSearch() {
	sut := connectionsHub{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
	}
	sut.checkRestartMdnsSearch()
	// Nothing to verify yet
}

func (s *HubSuite) Test_ReportServiceShipID() {
	sut := connectionsHub{
		serviceProvider: s,
	}
	sut.ReportServiceShipID("", "")
	// Nothing to verify yet
}

func (s *HubSuite) Test_DisconnectSKI() {
	sut := connectionsHub{
		connections: make(map[string]*ship.ShipConnection),
	}
	ski := "test"
	sut.DisconnectSKI(ski, "none")
}

func (s *HubSuite) Test_RegisterConnection() {
	ski := "12af9e"
	localService := &ServiceDetails{
		SKI:        ski,
		deviceType: model.DeviceTypeTypeEnergyManagementSystem, // this won't trigger mDNS announcement
	}

	sut := connectionsHub{
		connections:  make(map[string]*ship.ShipConnection),
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

func (s *HubSuite) Test_IncreaseConnectionAttemptCounter() {

	// we just need a dummy for this test
	sut := connectionsHub{
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
	sut := connectionsHub{
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
	sut := connectionsHub{
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
	localService := &ServiceDetails{
		SKI: ski,
	}
	sut := connectionsHub{
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
	sut := connectionsHub{
		connectionAttemptRunning: make(map[string]bool),
	}

	sut.setConnectionAttemptRunning(ski, true)
	status := sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), true, status)
	sut.setConnectionAttemptRunning(ski, false)
	status = sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), false, status)
}
