package spine

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UsecaseSuite))
}

type UsecaseSuite struct {
	suite.Suite

	device *DeviceLocalImpl
	entity *EntityLocalImpl
}

func (s *UsecaseSuite) BeforeTest(suiteName, testName string) {
	s.device = NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	s.entity = NewEntityLocalImpl(s.device, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	s.device.AddEntity(s.entity)

}

func (s *UsecaseSuite) Test_UseCase() {
	uc := NewUseCase(
		s.entity,
		model.UseCaseNameTypeControlOfBattery,
		model.SpecificationVersionType("1.0.0"),
		true,
		[]model.UseCaseScenarioSupportType{1},
	)
	assert.NotNil(s.T(), uc)

	uc.SetUseCaseAvailable(true)
}
