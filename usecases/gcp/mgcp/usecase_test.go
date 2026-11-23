package mgcp

import (
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
	"testing"
	"time"
)

func TestMGCPUsecaseSuite(t *testing.T) {
	suite.Run(t, new(MgcpUsecaseSuite))
}

type MgcpUsecaseSuite struct {
	suite.Suite

	service api.ServiceInterface
}

func (s *MgcpUsecaseSuite) BeforeTest(_, _ string) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeGridGuard},
		9999, cert, time.Second*4, nil, nil)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()
}

func (s *MgcpUsecaseSuite) Event(_ string, _ spineapi.DeviceRemoteInterface, _ spineapi.EntityRemoteInterface, _ api.EventType) {
}

func (s *MgcpUsecaseSuite) Test_RequiredParameters() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	var monitorPowerConfig = MonitorPowerConfig{
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	var monitorEnergyConfig = MonitorEnergyConfig{
		ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	numRequiredParams := 2

	// iterate over all possible combinations of the required parameters that should not work
	for i := 0; i < (1<<numRequiredParams)-1; i++ {
		var reqPowerConfig *MonitorPowerConfig
		var reqEnergyConfig *MonitorEnergyConfig

		if i&1 != 0 {
			reqPowerConfig = &monitorPowerConfig
		}

		if i&2 != 0 {
			reqEnergyConfig = &monitorEnergyConfig
		}

		mpc, err := NewMGCP(
			localEntity,
			s.Event,
			nil,
			reqPowerConfig,
			reqEnergyConfig,
			nil,
			nil,
			nil,
		)

		assert.Nil(s.T(), mpc)
		assert.NotNil(s.T(), err)
	}
}

func (s *MgcpUsecaseSuite) Test_OptionalParameters() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)
	assert.NotNil(s.T(), localEntity)

	var monitorPowerLimitationFactor = MonitorPvFeedInPowerLimitationFactorConfig{}

	var monitorPowerConfig = MonitorPowerConfig{
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

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
			ValueStepSize: model.NewScaledNumberType(1),
			ValueRangeMin: model.NewScaledNumberType(0),
			ValueRangeMax: model.NewScaledNumberType(100),
		}),
	}

	// iterate over all possible combinations of the optional parameters
	for i := 0; i < (1 << 4); i++ {
		var optPowerLimitationFactor *MonitorPvFeedInPowerLimitationFactorConfig
		var optCurrentConfig *MonitorCurrentConfig
		var optVoltageConfig *MonitorVoltageConfig
		var optFrequencyConfig *MonitorFrequencyConfig

		if i&1 != 0 {
			optPowerLimitationFactor = &monitorPowerLimitationFactor
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

		mpc, err := NewMGCP(
			localEntity,
			s.Event,
			optPowerLimitationFactor,
			&monitorPowerConfig,
			&monitorEnergyConfig,
			optCurrentConfig,
			optVoltageConfig,
			optFrequencyConfig,
		)

		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpc)

		err = mpc.AddFeatures()
		assert.Nil(s.T(), err)
	}
}

func (s *MgcpUsecaseSuite) Test_getMeasurementForId() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	var monitorPowerConfig = MonitorPowerConfig{
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	var monitorEnergyConfig = MonitorEnergyConfig{
		ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		nil,
		&monitorPowerConfig,
		&monitorEnergyConfig,
		nil,
		nil,
		nil,
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), mgcp)

	err = mgcp.AddFeatures()

	// test with invalid id
	m, err := mgcp.getMeasurementDataForId(nil)
	assert.NotNil(s.T(), err)

	id := mgcp.acPowerTotal
	value := 43.0

	err = mgcp.Update(mgcp.UpdateDataPowerTotal(value, nil, nil))
	assert.Nil(s.T(), err)

	m, err = mgcp.getMeasurementDataForId(id)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), value, m)
}

func (s *GcpMpcgSuite) Test_AddFeaturesConfigurationNilError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()
	localEntity.EXPECT().GetOrAddFeature(mock.Anything, mock.Anything).Return(anyFeature)

	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, mock.Anything).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err) // NewDeviceConfiguration failed
}

func (s *GcpMpcgSuite) Test_AddFeaturesMeasurementNilError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil).Maybe()
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()
	localEntity.EXPECT().GetOrAddFeature(mock.Anything, mock.Anything).Return(anyFeature)

	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, mock.Anything).Return(anyFeature)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, mock.Anything).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err) // NewMeasurement failed
}

func (s *GcpMpcgSuite) Test_AddFeaturesElectricalConnectionNilError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil).Maybe()
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()
	localEntity.EXPECT().GetOrAddFeature(mock.Anything, mock.Anything).Return(anyFeature)

	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, mock.Anything).Return(anyFeature)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, mock.Anything).Return(anyFeature)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, mock.Anything).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err) // NewMeasurement failed
}

func (s *MgcpUsecaseSuite) Test_configurePvFeedInLimitationFactorError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	deviceConfiguration := mocks.NewDeviceConfigurationServerInterface(s.T())
	deviceConfiguration.EXPECT().AddKeyValueDescription(mock.Anything).Return(nil)

	err = mgcp.configurePvFeedInLimitationFactor(deviceConfiguration)
	assert.NotNil(s.T(), err)
}

func (s *MgcpUsecaseSuite) Test_configureMonitorPowerError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	measurements := mocks.NewMeasurementServerInterface(s.T())
	electricalConnection := mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnectionId := model.ElectricalConnectionIdType(0)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	mgcp.powerConfig = nil
	err = mgcp.configureMonitorPower(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

	mgcp.powerConfig = &MonitorPowerConfig{
		ValueSource: nil,
	}
	err = mgcp.configureMonitorPower(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

	mgcp.powerConfig = &MonitorPowerConfig{
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
	}

	measurements.EXPECT().AddDescription(mock.Anything).Return(nil)
	electricalConnection.EXPECT().AddParameterDescription(mock.Anything).Return(nil)
	err = mgcp.configureMonitorPower(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)
}

func (s *MgcpUsecaseSuite) Test_configureGridFeedInError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	measurements := mocks.NewMeasurementServerInterface(s.T())
	electricalConnection := mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnectionId := model.ElectricalConnectionIdType(0)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	mgcp.energyConfig = nil
	err = mgcp.configureGridFeedIn(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

	mgcp.energyConfig = &MonitorEnergyConfig{
		ValueSourceProduction: nil,
	}
	err = mgcp.configureGridFeedIn(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

	mgcp.energyConfig = &MonitorEnergyConfig{
		ValueSourceProduction: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
	}

	measurements.EXPECT().AddDescription(mock.Anything).Return(nil)
	electricalConnection.EXPECT().AddParameterDescription(mock.Anything).Return(nil)
	err = mgcp.configureGridFeedIn(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)
}

func (s *MgcpUsecaseSuite) Test_configureGridConsumptionError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	measurements := mocks.NewMeasurementServerInterface(s.T())
	electricalConnection := mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnectionId := model.ElectricalConnectionIdType(0)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	mgcp.energyConfig = nil
	err = mgcp.configureGridConsumption(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

	mgcp.energyConfig = &MonitorEnergyConfig{
		ValueSourceConsumption: nil,
	}
	err = mgcp.configureGridConsumption(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

	mgcp.energyConfig = &MonitorEnergyConfig{
		ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
	}

	measurements.EXPECT().AddDescription(mock.Anything).Return(nil)
	electricalConnection.EXPECT().AddParameterDescription(mock.Anything).Return(nil)
	err = mgcp.configureGridConsumption(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)
}

func (s *MgcpUsecaseSuite) Test_configureMonitorCurrentError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	measurements := mocks.NewMeasurementServerInterface(s.T())
	electricalConnection := mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnectionId := model.ElectricalConnectionIdType(0)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	measurements.EXPECT().AddDescription(mock.Anything).Return(nil)
	electricalConnection.EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	currentConfigurations := []MonitorCurrentConfig{
		{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}

	for _, currentConfig := range currentConfigurations {
		mgcp.currentConfig = &currentConfig
		err = mgcp.configureMonitorCurrent(measurements, electricalConnection, &electricalConnectionId, &constraints)
		assert.NotNil(s.T(), err) // failed to add parameter description
	}
}

func (s *MgcpUsecaseSuite) Test_configureMonitorVoltageError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	measurements := mocks.NewMeasurementServerInterface(s.T())
	electricalConnection := mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnectionId := model.ElectricalConnectionIdType(0)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	measurements.EXPECT().AddDescription(mock.Anything).Return(nil)
	electricalConnection.EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	voltageConfigurations := []MonitorVoltageConfig{
		{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		{
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}

	for _, voltageConfig := range voltageConfigurations {
		mgcp.voltageConfig = &voltageConfig
		err = mgcp.configureMonitorVoltage(measurements, electricalConnection, &electricalConnectionId, &constraints)
		assert.NotNil(s.T(), err) // failed to add parameter description
	}
}

func (s *MgcpUsecaseSuite) Test_configureMonitorFrequencyError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)

	mgcp, err := NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)
	assert.Nil(s.T(), err)

	measurements := mocks.NewMeasurementServerInterface(s.T())
	electricalConnection := mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnectionId := model.ElectricalConnectionIdType(0)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	measurements.EXPECT().AddDescription(mock.Anything).Return(nil)
	electricalConnection.EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	err = mgcp.configureMonitorFrequency(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err) // failed to add parameter description
}
