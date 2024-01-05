package spine_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceRemoteSuite(t *testing.T) {
	suite.Run(t, new(DeviceRemoteSuite))
}

type DeviceRemoteSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteDevice *spine.DeviceRemoteImpl
}

func (s *DeviceRemoteSuite) WriteSpineMessage([]byte) {}

func (s *DeviceRemoteSuite) SetupSuite() {}

func (s *DeviceRemoteSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = spine.NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	sender := spine.NewSender(s)
	s.remoteDevice = spine.NewDeviceRemoteImpl(s.localDevice, ski, sender)
	s.localDevice.AddRemoteDevice(ski, s)

	entity := spine.NewEntityRemoteImpl(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	feature := spine.NewFeatureRemoteImpl(0, entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	entity.AddFeature(feature)

	s.remoteDevice.AddEntity(entity)
}

func (s *DeviceRemoteSuite) Test_RemoveByAddress() {
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{2})
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))
}

func (s *DeviceRemoteSuite) Test_FeatureByEntityTypeAndRole() {
	entity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), entity)

	assert.Equal(s.T(), 1, len(entity.Features()))

	feature := s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	assert.Nil(s.T(), feature)

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), feature)

	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))

	_ = s.remoteDevice.Entity([]model.AddressEntityType{0})
	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{0})
	assert.Equal(s.T(), 0, len(s.remoteDevice.Entities()))

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.Nil(s.T(), feature)
}

func (s *DeviceRemoteSuite) Test_VerifyUseCaseScenariosAndFeaturesSupport() {
	result := s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	s.remoteDevice.UseCaseManager().Add(
		model.UseCaseActorTypeBatterySystem,
		model.UseCaseNameTypeControlOfBattery,
		model.SpecificationVersionType("1.0.0"),
		true,
		[]model.UseCaseScenarioSupportType{1},
	)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	s.remoteDevice.UseCaseManager().Add(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.0"),
		true,
		[]model.UseCaseScenarioSupportType{1},
	)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	s.remoteDevice.UseCaseManager().Add(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.0"),
		false,
		[]model.UseCaseScenarioSupportType{1},
	)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	s.remoteDevice.UseCaseManager().Add(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.0"),
		true,
		[]model.UseCaseScenarioSupportType{1},
	)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), true, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{2},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), true, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), false, result)

	entity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), entity)

	feature := spine.NewFeatureRemoteImpl(0, entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	entity.AddFeature(feature)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), false, result)

	feature = spine.NewFeatureRemoteImpl(0, entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	entity.AddFeature(feature)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), true, result)

	s.remoteDevice.RemoveByAddress(feature.Address().Entity)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), false, result)
}
