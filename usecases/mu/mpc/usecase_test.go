package mpc

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestBasicSuite(t *testing.T) {
	suite.Run(t, new(BasicSuite))
}

type BasicSuite struct {
	suite.Suite

	service api.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
	monitoredEntity  spineapi.EntityRemoteInterface
	loadControlFeature,
	deviceDiagnosisFeature,
	deviceConfigurationFeature spineapi.FeatureLocalInterface

	eventCalled bool
}

func (s *BasicSuite) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	s.eventCalled = true
}

func (s *BasicSuite) BeforeTest(suiteName, testName string) {
	s.eventCalled = false
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeInverter},
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()

	mockRemoteDevice := spinemocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = spinemocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := spinemocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()
	mockRemoteFeature.EXPECT().Address().Return(&model.FeatureAddressType{}).Maybe()
	mockRemoteFeature.EXPECT().Operations().Return(nil).Maybe()
}

func (s *BasicSuite) Test_MpcOptionalParameters() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)

	// required
	var monitorPowerConfig = MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	// the following 4 parameters are optional and can be nil
	var monitorEnergyConfig = MonitorEnergyConfig{
		ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}
	var monitorCurrentConfig = MonitorCurrentConfig{
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}
	var monitorVoltageConfig = MonitorVoltageConfig{
		SupportPhaseToPhase:  true,
		ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}
	var monitorFrequencyConfig = MonitorFrequencyConfig{
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
			ValueRangeMin: model.NewScaledNumberType(0),
			ValueRangeMax: model.NewScaledNumberType(100),
			ValueStepSize: model.NewScaledNumberType(1),
		}),
	}

	numOptionalParams := 4

	// iterate over all permutations of nil/set
	for i := 0; i < (1 << numOptionalParams); i++ {
		// Determine which parameters to set
		var optEnergyConfig *MonitorEnergyConfig
		var optCurrentConfig *MonitorCurrentConfig
		var optVoltageConfig *MonitorVoltageConfig
		var optFrequencyConfig *MonitorFrequencyConfig
		if i&1 != 0 {
			optEnergyConfig = &monitorEnergyConfig
		}
		if i&2 != 0 {
			optCurrentConfig = &monitorCurrentConfig
		}
		if i&4 != 0 {
			optVoltageConfig = &monitorVoltageConfig
		}
		if i&8 != 0 {
			optFrequencyConfig = &monitorFrequencyConfig
		}

		mpc, err := NewMPC(
			localEntity,
			s.Event,
			&monitorPowerConfig,
			optEnergyConfig,
			optCurrentConfig,
			optVoltageConfig,
			optFrequencyConfig,
		)

		assert.Nil(s.T(), err)

		err = mpc.AddFeatures()
		assert.Nil(s.T(), err)
		mpc.AddUseCase()
	}
}
