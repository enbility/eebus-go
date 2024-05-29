package opev

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/stretchr/testify/assert"
)

func (s *CemOPEVSuite) Test_Public() {
	// The actual tests of the functionality is located in the util package

	_, _, _, err := s.sut.CurrentLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, _, _, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.LoadControlLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.LoadControlLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimits(s.mockRemoteEntity, []ucapi.LoadLimitsPhase{})
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimits(s.evEntity, []ucapi.LoadLimitsPhase{})
	assert.NotNil(s.T(), err)
}
