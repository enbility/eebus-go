package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHubSuite(t *testing.T) {
	suite.Run(t, new(HubSuite))
}

type HubSuite struct {
	suite.Suite
}

func (s *HubSuite) Test_IncreaseConnectionAttemptCounter() {
	tests := []struct {
		counter   int
		timeRange connectionInitiationDelayTimeRange
	}{
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
	// we just need a dummy for this test
	sut := connectionsHub{}
	sut.connectionAttemptCounter = make(map[string]int)
	ski := "test"

	for _, test := range tests {
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
