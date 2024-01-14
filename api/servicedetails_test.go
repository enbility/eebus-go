package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestServiceDetails(t *testing.T) {
	suite.Run(t, new(ServiceDetailsSuite))
}

type ServiceDetailsSuite struct {
	suite.Suite
}

func (s *ServiceDetailsSuite) Test_ServiceDetails() {
	testSki := "test"

	details := NewServiceDetails(testSki)
	assert.NotNil(s.T(), details)

	conState := NewConnectionStateDetail(ConnectionStateNone, nil)
	details.SetConnectionStateDetail(conState)

	state := details.ConnectionStateDetail()
	assert.Equal(s.T(), ConnectionStateNone, state.State())
}
