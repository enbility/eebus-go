package mpc

import (
	"github.com/enbility/eebus-go/features/server"
	spinemocks "github.com/enbility/spine-go/mocks"
	"testing"
	"time"

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
)

func TestBasicSuite(t *testing.T) {
	suite.Run(t, new(MuMpcUsecaseSuite))
}

type MuMpcUsecaseSuite struct {
	suite.Suite

	service api.ServiceInterface
}

func (s *MuMpcUsecaseSuite) Event(_ string, _ spineapi.DeviceRemoteInterface, _ spineapi.EntityRemoteInterface, _ api.EventType) {
}

func (s *MuMpcUsecaseSuite) BeforeTest(_, _ string) {
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
}

func (s *MuMpcUsecaseSuite) Test_MpcOptionalParameters() {
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

func (s *MuMpcUsecaseSuite) Test_getMeasurementDataForId() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	mpc, err := NewMPC(
		localEntity,
		s.Event,
		&monitorPowerConfig,
		nil,
		nil,
		nil,
		nil,
	)
	assert.Nil(s.T(), err)

	_, err = mpc.getMeasurementDataForId(mpc.acPowerTotal)
	assert.NotNil(s.T(), err)

	err = mpc.AddFeatures()
	assert.Nil(s.T(), err)
	mpc.AddUseCase()

	_, err = mpc.getMeasurementDataForId(mpc.acPowerTotal)
	assert.NotNil(s.T(), err)

	err = mpc.Update(
		mpc.UpdateDataPowerTotal(5.0, util.Ptr(time.Now()), nil),
	)
	assert.Nil(s.T(), err)

	measurementData, err := mpc.getMeasurementDataForId(mpc.acPowerTotal)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), measurementData)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_ElectricalFeatureNilError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(nil)
	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_MeasurementFeatureNilError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(anyFeature)
	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_NewMeasurementsError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(anyFeature)
	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(anyFeature)

	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_NewElectricalConnectionError() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(anyFeature)
	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(anyFeature)

	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(anyFeature)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcUsecaseSuite) Test_configureMonitorPower() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	localEntity.EXPECT().Device().Return(nil)

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	mpc, err := NewMPC(
		localEntity,
		s.Event,
		&monitorPowerConfig,
		nil,
		nil,
		nil,
		nil,
	)
	assert.Nil(s.T(), err)

	measurements, err := server.NewMeasurement(localEntity)
	assert.Nil(s.T(), err)

	var electricalConnection api.ElectricalConnectionServerInterface
	electricalConnection, err = server.NewElectricalConnection(localEntity)
	assert.Nil(s.T(), err)

	electricalConnectionId := model.ElectricalConnectionIdType(111)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	mpc.powerConfig = nil
	err = mpc.configureMonitorPower(
		measurements,
		electricalConnection,
		&electricalConnectionId,
		&constraints,
	)
	assert.NotNil(s.T(), err) // no monitorPowerConfig

	mpc.powerConfig = &monitorPowerConfig
	electricalConnection = mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	constellationsToCheck := []*ConnectedPhases{
		util.Ptr(ConnectedPhasesA),
		util.Ptr(ConnectedPhasesB),
		util.Ptr(ConnectedPhasesC),
	}

	for _, phaseConstellation := range constellationsToCheck {
		mpc.powerConfig.ConnectedPhases = *phaseConstellation

		err = mpc.configureMonitorPower(
			measurements,
			electricalConnection,
			&electricalConnectionId,
			nil,
		)

		assert.NotNil(s.T(), err) // could not add parameter description
	}
}

func (s *MuMpcUsecaseSuite) Test_configureMonitorEnergy() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	localEntity.EXPECT().Device().Return(nil)

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	monitorEnergyConfig := MonitorEnergyConfig{
		ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	mpc, err := NewMPC(
		localEntity,
		s.Event,
		&monitorPowerConfig,
		&monitorEnergyConfig,
		nil,
		nil,
		nil,
	)
	assert.Nil(s.T(), err)

	measurements, err := server.NewMeasurement(localEntity)
	assert.Nil(s.T(), err)

	var electricalConnection api.ElectricalConnectionServerInterface
	electricalConnection, err = server.NewElectricalConnection(localEntity)
	assert.Nil(s.T(), err)

	electricalConnectionId := model.ElectricalConnectionIdType(111)
	constraints := make([]model.MeasurementConstraintsDataType, 0)
	electricalConnection = mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	err = mpc.configureMonitorEnergy(
		measurements,
		electricalConnection,
		&electricalConnectionId,
		&constraints,
	)

	assert.NotNil(s.T(), err) // could not add parameter description 1
	mpc.energyConfig.ValueConstraintsConsumption = nil

	err = mpc.configureMonitorEnergy(
		measurements,
		electricalConnection,
		&electricalConnectionId,
		&constraints,
	)

	assert.NotNil(s.T(), err) // could not add parameter description 2
}

func (s *MuMpcUsecaseSuite) Test_configureMonitorCurrent() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	localEntity.EXPECT().Device().Return(nil)

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	monitorCurrentConfig := MonitorCurrentConfig{
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	mpc, err := NewMPC(
		localEntity,
		s.Event,
		&monitorPowerConfig,
		nil,
		&monitorCurrentConfig,
		nil,
		nil,
	)
	assert.Nil(s.T(), err)

	measurements, err := server.NewMeasurement(localEntity)
	assert.Nil(s.T(), err)

	var electricalConnection api.ElectricalConnectionServerInterface
	electricalConnection, err = server.NewElectricalConnection(localEntity)
	assert.Nil(s.T(), err)

	electricalConnectionId := model.ElectricalConnectionIdType(111)
	constraints := make([]model.MeasurementConstraintsDataType, 0)
	electricalConnection = mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	constellationsToCheck := []*ConnectedPhases{
		util.Ptr(ConnectedPhasesA),
		util.Ptr(ConnectedPhasesB),
		util.Ptr(ConnectedPhasesC),
	}

	for _, phaseConstellation := range constellationsToCheck {
		mpc.powerConfig.ConnectedPhases = *phaseConstellation

		err = mpc.configureMonitorCurrent(
			measurements,
			electricalConnection,
			&electricalConnectionId,
			&constraints,
		)

		assert.NotNil(s.T(), err) // could not add parameter description
	}
}

func (s *MuMpcUsecaseSuite) Test_configureMonitorVoltage() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	localEntity.EXPECT().Device().Return(nil)

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	monitorVoltageConfig := MonitorVoltageConfig{
		SupportPhaseToPhase:  true,
		ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	mpc, err := NewMPC(
		localEntity,
		s.Event,
		&monitorPowerConfig,
		nil,
		nil,
		&monitorVoltageConfig,
		nil,
	)
	assert.Nil(s.T(), err)

	measurements, err := server.NewMeasurement(localEntity)
	assert.Nil(s.T(), err)

	var electricalConnection api.ElectricalConnectionServerInterface
	electricalConnection, err = server.NewElectricalConnection(localEntity)
	assert.Nil(s.T(), err)

	electricalConnectionId := model.ElectricalConnectionIdType(111)
	constraints := make([]model.MeasurementConstraintsDataType, 0)

	electricalConnection = mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	constellationsToCheck := []*ConnectedPhases{
		util.Ptr(ConnectedPhasesA),
		util.Ptr(ConnectedPhasesB),
		util.Ptr(ConnectedPhasesC),
	}

	for _, phaseConstellation := range constellationsToCheck {
		mpc.powerConfig.ConnectedPhases = *phaseConstellation

		err = mpc.configureMonitorVoltage(
			measurements,
			electricalConnection,
			&electricalConnectionId,
			&constraints,
		)

		assert.NotNil(s.T(), err) // could not add parameter description
	}
}

func (s *MuMpcUsecaseSuite) Test_configureMonitorFrequency() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	localEntity.EXPECT().Device().Return(nil)

	anyFeature := spinemocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:   ConnectedPhasesABC,
		ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}

	monitorFrequencyConfig := MonitorFrequencyConfig{
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
			ValueRangeMin: model.NewScaledNumberType(0),
			ValueRangeMax: model.NewScaledNumberType(100),
			ValueStepSize: model.NewScaledNumberType(1),
		}),
	}

	mpc, err := NewMPC(
		localEntity,
		s.Event,
		&monitorPowerConfig,
		nil,
		nil,
		nil,
		&monitorFrequencyConfig,
	)
	assert.Nil(s.T(), err)

	measurements, err := server.NewMeasurement(localEntity)
	assert.Nil(s.T(), err)

	var electricalConnection api.ElectricalConnectionServerInterface
	electricalConnection, err = server.NewElectricalConnection(localEntity)
	assert.Nil(s.T(), err)

	electricalConnectionId := model.ElectricalConnectionIdType(111)
	constraints := make([]model.MeasurementConstraintsDataType, 0)
	electricalConnection = mocks.NewElectricalConnectionServerInterface(s.T())
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil)

	err = mpc.configureMonitorFrequency(
		measurements,
		electricalConnection,
		&electricalConnectionId,
		&constraints,
	)
	assert.NotNil(s.T(), err) // could not add parameter description
}
