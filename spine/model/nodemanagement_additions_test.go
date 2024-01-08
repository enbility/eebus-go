package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNodeManagementUseCaseDataTypeSuite(t *testing.T) {
	suite.Run(t, new(NodeManagementUseCaseDataTypeSuite))
}

type NodeManagementUseCaseDataTypeSuite struct {
	suite.Suite
}

func (s *NodeManagementUseCaseDataTypeSuite) SetupSuite()   {}
func (s *NodeManagementUseCaseDataTypeSuite) TearDownTest() {}

func (s *NodeManagementUseCaseDataTypeSuite) BeforeTest(suiteName, testName string) {}

func (s *NodeManagementUseCaseDataTypeSuite) Test_AdditionsAndRemovals() {
	ucs := &NodeManagementUseCaseDataType{}
	assert.NotNil(s.T(), ucs)
	assert.Equal(s.T(), 0, len(ucs.UseCaseInformation))

	address := FeatureAddressType{
		Device: util.Ptr(AddressDeviceType("test")),
		Entity: []AddressEntityType{1},
	}
	ucs.AddUseCaseSupport(
		address,
		UseCaseActorTypeCEM,
		UseCaseNameTypeControlOfBattery,
		SpecificationVersionType(""),
		"",
		true,
		[]UseCaseScenarioSupportType{},
	)
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation))
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation[0].UseCaseSupport))

	ucs.AddUseCaseSupport(
		address,
		UseCaseActorTypeCEM,
		UseCaseNameTypeEVSECommissioningAndConfiguration,
		SpecificationVersionType(""),
		"",
		true,
		[]UseCaseScenarioSupportType{},
	)
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation))
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation[0].UseCaseSupport))

	ucs.AddUseCaseSupport(
		address,
		UseCaseActorTypeCEM,
		UseCaseNameTypeEVSECommissioningAndConfiguration,
		SpecificationVersionType(""),
		"",
		true,
		[]UseCaseScenarioSupportType{},
	)
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation))
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation[0].UseCaseSupport))

	ucs.AddUseCaseSupport(
		address,
		UseCaseActorTypeEnergyGuard,
		UseCaseNameTypeLimitationOfPowerConsumption,
		SpecificationVersionType(""),
		"",
		true,
		[]UseCaseScenarioSupportType{},
	)
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation))
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation[0].UseCaseSupport))
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation[1].UseCaseSupport))

	ucs.RemoveUseCaseSupport(
		address,
		UseCaseActorTypeCEM,
		UseCaseNameTypeEVChargingSummary,
	)
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation))
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation[0].UseCaseSupport))

	ucs.RemoveUseCaseSupport(
		address,
		UseCaseActorTypeCEM,
		UseCaseNameTypeControlOfBattery,
	)
	assert.Equal(s.T(), 2, len(ucs.UseCaseInformation))
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation[0].UseCaseSupport))

	ucs.RemoveUseCaseSupport(
		address,
		UseCaseActorTypeCEM,
		UseCaseNameTypeEVSECommissioningAndConfiguration,
	)
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation))

	ucs.RemoveUseCaseSupport(
		address,
		"",
		"",
	)
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation))

	invalidAddress := FeatureAddressType{
		Device: util.Ptr(AddressDeviceType("test")),
		Entity: []AddressEntityType{2},
	}
	ucs.RemoveUseCaseSupport(
		invalidAddress,
		UseCaseActorTypeCEM,
		UseCaseNameTypeEVSECommissioningAndConfiguration,
	)
	assert.Equal(s.T(), 1, len(ucs.UseCaseInformation))

}
