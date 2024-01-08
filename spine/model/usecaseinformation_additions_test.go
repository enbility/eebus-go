package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUseCaseInformationDataTypeSuite(t *testing.T) {
	suite.Run(t, new(UseCaseInformationDataTypeSuite))
}

type UseCaseInformationDataTypeSuite struct {
	suite.Suite
}

func (s *UseCaseInformationDataTypeSuite) SetupSuite()   {}
func (s *UseCaseInformationDataTypeSuite) TearDownTest() {}

func (s *UseCaseInformationDataTypeSuite) BeforeTest(suiteName, testName string) {}

func (s *UseCaseInformationDataTypeSuite) Test_AdditionsAndRemovals() {
	ucs := &UseCaseInformationDataType{}
	assert.NotNil(s.T(), ucs)
	assert.Equal(s.T(), 0, len(ucs.UseCaseSupport))

	uc := UseCaseSupportType{}
	ucs.Add(uc)
	assert.Equal(s.T(), 0, len(ucs.UseCaseSupport))

	uc = UseCaseSupportType{
		UseCaseName: util.Ptr(UseCaseNameTypeControlOfBattery),
	}
	ucs.Add(uc)
	assert.Equal(s.T(), 1, len(ucs.UseCaseSupport))

	ucs.Add(uc)
	assert.Equal(s.T(), 1, len(ucs.UseCaseSupport))

	ucs.Remove(UseCaseNameTypeCoordinatedEVCharging)
	assert.Equal(s.T(), 1, len(ucs.UseCaseSupport))

	ucs.Remove(UseCaseNameTypeControlOfBattery)
	assert.Equal(s.T(), 0, len(ucs.UseCaseSupport))
}
