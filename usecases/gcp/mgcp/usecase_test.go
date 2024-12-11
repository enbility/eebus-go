package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
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
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()
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

func (s *MgcpUsecaseSuite) Event(_ string, _ spineapi.DeviceRemoteInterface, _ spineapi.EntityRemoteInterface, _ api.EventType) {
}
