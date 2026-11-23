package mpc

import (
	"errors"
	"testing"
	"time"

	"github.com/enbility/eebus-go/features/server"
	spineMocks "github.com/enbility/spine-go/mocks"

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

	service       api.ServiceInterface
	mockedService *mocks.ServiceInterface

	localEntity        spineapi.EntityLocalInterface
	mockedLocalEntity  *spineMocks.EntityLocalInterface
	mockedLocalDevice  *spineMocks.DeviceLocalInterface
	mockedLocalFeature *spineMocks.FeatureLocalInterface

	mockedRemoteEntity  *spineMocks.EntityRemoteInterface
	mockedRemoteDevice  *spineMocks.DeviceRemoteInterface
	mockedRemoteFeature *spineMocks.FeatureRemoteInterface

	mockedElectricalConnectionFeature *mocks.ElectricalConnectionServerInterface
}

func (s *MuMpcUsecaseSuite) Event(_ string, _ spineapi.DeviceRemoteInterface, _ spineapi.EntityRemoteInterface, _ api.EventType) {
}

func (s *MuMpcUsecaseSuite) BeforeTest(_, _ string) {
	cert, err := cert.CreateCertificate("test", "test", "DE", "test")
	assert.Nil(s.T(), err)
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeInverter},
		9999, cert, time.Second*4, nil, nil)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	err = s.service.Setup()
	assert.Nil(s.T(), err)

	s.mockedRemoteDevice = spineMocks.NewDeviceRemoteInterface(s.T())
	s.mockedRemoteEntity = spineMocks.NewEntityRemoteInterface(s.T())
	s.mockedRemoteFeature = spineMocks.NewFeatureRemoteInterface(s.T())
	s.mockedRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(s.mockedRemoteFeature).Maybe()
	s.mockedRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockedRemoteEntity.EXPECT().Device().Return(s.mockedRemoteDevice).Maybe()
	s.mockedRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockedRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	s.mockedRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()
	s.mockedRemoteFeature.EXPECT().Address().Return(&model.FeatureAddressType{}).Maybe()
	s.mockedRemoteFeature.EXPECT().Operations().Return(nil).Maybe()

	s.localEntity = s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)
}

func (s *MuMpcUsecaseSuite) Test_MpcOptionalParameters() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)

	// required
	var monitorPowerConfig = &MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}

	// the following 4 parameters are optional and can be nil
	var monitorEnergyConfig = MonitorEnergyConfig{
		ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
	}
	var monitorCurrentConfig = MonitorCurrentConfig{
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}
	var monitorVoltageConfig = MonitorVoltageConfig{
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeAb: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeBc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeAc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
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
			monitorPowerConfig,
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

	// test creating new mpc instance without power configuration
	{
		mpcInstance, err := NewMPC(s.localEntity, s.Event, nil, nil, nil, nil, nil)
		expectedError := "the monitor power config for the MPC-Use-Case must not be nil"
		assert.ErrorContains(s.T(), err, expectedError)
		assert.Nil(s.T(), mpcInstance)
	}
}

func (s *MuMpcUsecaseSuite) Test_MpcRequredParametersError() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)

	_, err := NewMPC(
		localEntity,
		s.Event,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	assert.NotNil(s.T(), err)
}

func (s *MuMpcUsecaseSuite) Test_getMeasurementDataForId() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
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
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(nil)
	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_MeasurementFeatureNilError() {
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(anyFeature)
	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_NewMeasurementsError() {
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return()

	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(anyFeature)
	localEntity.EXPECT().GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(anyFeature)

	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(nil)

	err := s.sut.AddFeatures()
	assert.NotNil(s.T(), err)
}

func (s *MuMpcAbcSuite) Test_AddFeatures_NewElectricalConnectionError() {
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	s.sut.LocalEntity = localEntity

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
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
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	mockedDevice := spineMocks.NewDeviceLocalInterface(s.T())
	mockedEvents := spineMocks.NewEventsManagerInterface(s.T())
	mockedEvents.EXPECT().Subscribe(mock.Anything).Return(nil).Maybe()
	mockedDevice.EXPECT().Events().Return(mockedEvents).Maybe()
	localEntity.EXPECT().Device().Return(mockedDevice)

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
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

	constellationsToCheck := []model.ElectricalConnectionPhaseNameType{
		model.ElectricalConnectionPhaseNameTypeA,
		model.ElectricalConnectionPhaseNameTypeB,
		model.ElectricalConnectionPhaseNameTypeC,
	}

	for _, phaseConstellation := range constellationsToCheck {
		mpc.powerConfig.ConnectedPhases = phaseConstellation

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
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	mockedDevice := spineMocks.NewDeviceLocalInterface(s.T())
	mockedEvents := spineMocks.NewEventsManagerInterface(s.T())
	mockedEvents.EXPECT().Subscribe(mock.Anything).Return(nil).Maybe()
	mockedDevice.EXPECT().Events().Return(mockedEvents).Maybe()
	localEntity.EXPECT().Device().Return(mockedDevice)

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
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
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	mockedDevice := spineMocks.NewDeviceLocalInterface(s.T())
	mockedEvents := spineMocks.NewEventsManagerInterface(s.T())
	mockedEvents.EXPECT().Subscribe(mock.Anything).Return(nil).Maybe()
	mockedDevice.EXPECT().Events().Return(mockedEvents).Maybe()
	localEntity.EXPECT().Device().Return(mockedDevice)

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil).Maybe()
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature).Maybe()

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}

	monitorCurrentConfig := MonitorCurrentConfig{
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
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
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil).Maybe()

	constellationsToCheck := []model.ElectricalConnectionPhaseNameType{
		model.ElectricalConnectionPhaseNameTypeA,
		model.ElectricalConnectionPhaseNameTypeB,
		model.ElectricalConnectionPhaseNameTypeC,
	}

	for _, phaseConstellation := range constellationsToCheck {
		mpc.powerConfig.ConnectedPhases = phaseConstellation

		err = mpc.configureMonitorCurrent(
			measurements,
			electricalConnection,
			&electricalConnectionId,
			&constraints,
		)

		assert.NotNil(s.T(), err) // could not add parameter description
	}

	// test when using phase to phase
	monitorVoltageConfig := MonitorVoltageConfig{
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeAb: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}

	mpcInstance, err := NewMPC(s.localEntity, s.Event, &monitorPowerConfig, nil, nil, &monitorVoltageConfig, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), mpcInstance)

	err = mpcInstance.configureMonitorVoltage(measurements, electricalConnection, &electricalConnectionId, &constraints)
	assert.NotNil(s.T(), err)

}

func (s *MuMpcUsecaseSuite) Test_configureMonitorVoltage() {
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	mockedDevice := spineMocks.NewDeviceLocalInterface(s.T())
	mockedEvents := spineMocks.NewEventsManagerInterface(s.T())
	mockedEvents.EXPECT().Subscribe(mock.Anything).Return(nil).Maybe()
	mockedDevice.EXPECT().Events().Return(mockedEvents).Maybe()
	localEntity.EXPECT().Device().Return(mockedDevice)

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil).Maybe()
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature).Maybe()

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	}

	monitorVoltageConfig := MonitorVoltageConfig{
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeAb: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeBc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeAc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
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
	electricalConnection.(*mocks.ElectricalConnectionServerInterface).EXPECT().AddParameterDescription(mock.Anything).Return(nil).Maybe()

	constellationsToCheck := []model.ElectricalConnectionPhaseNameType{
		model.ElectricalConnectionPhaseNameTypeA,
		model.ElectricalConnectionPhaseNameTypeB,
		model.ElectricalConnectionPhaseNameTypeC,
	}

	for _, phaseConstellation := range constellationsToCheck {
		mpc.powerConfig.ConnectedPhases = phaseConstellation

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
	localEntity := spineMocks.NewEntityLocalInterface(s.T())
	mockedDevice := spineMocks.NewDeviceLocalInterface(s.T())
	mockedEvents := spineMocks.NewEventsManagerInterface(s.T())
	mockedEvents.EXPECT().Subscribe(mock.Anything).Return(nil).Maybe()
	mockedDevice.EXPECT().Events().Return(mockedEvents).Maybe()
	localEntity.EXPECT().Device().Return(mockedDevice)

	anyFeature := spineMocks.NewFeatureLocalInterface(s.T())
	anyFeature.EXPECT().DataCopy(mock.Anything).Return(nil)
	anyFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	localEntity.EXPECT().FeatureOfTypeAndRole(mock.Anything, mock.Anything).Return(anyFeature)

	monitorPowerConfig := MonitorPowerConfig{
		ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
		ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueSourcePerPhase: PhaseMeasurementSourceMap{
			model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
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

func (s *MuMpcUsecaseSuite) TestAddFeatures() {
	// Testing function AddFeatures() and cover the paths that newMeasurement will return error in it
	// Covering first path when calling newMeasurement returns nil
	{
		s.mockedService = mocks.NewServiceInterface(s.T())
		s.mockedService.EXPECT().AddUseCase(mock.Anything).Return(nil).Maybe()
		s.mockedLocalDevice = spineMocks.NewDeviceLocalInterface(s.T())
		s.mockedService.EXPECT().LocalDevice().Return(s.mockedLocalDevice).Maybe()
		s.mockedLocalEntity = spineMocks.NewEntityLocalInterface(s.T())
		s.mockedLocalDevice.EXPECT().EntityForType(mock.Anything).Return(s.mockedLocalEntity).Maybe()
		s.mockedLocalFeature = spineMocks.NewFeatureLocalInterface(s.T())
		s.mockedLocalEntity.EXPECT().GetOrAddFeature(mock.Anything, mock.Anything).Return(s.mockedLocalFeature).Maybe()
		s.mockedLocalEntity.EXPECT().Device().Return(s.mockedLocalDevice).Maybe()
		mockedEvents := spineMocks.NewEventsManagerInterface(s.T())
		mockedEvents.EXPECT().Subscribe(mock.Anything).Return(nil).Maybe()
		s.mockedLocalDevice.EXPECT().Events().Return(mockedEvents).Maybe()
		s.mockedLocalFeature.EXPECT().AddFunctionType(mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
		s.mockedLocalEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(nil).Once().Maybe()

		powerConfig := &MonitorPowerConfig{
			ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
			ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}
		mpcInstance, err := NewMPC(s.mockedLocalEntity, s.Event, powerConfig, nil, nil, nil, nil)

		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.NotNil(s.T(), err)
	}

	// Test covering when calling newElectricalConnection and return error
	{
		s.mockedLocalEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(s.mockedLocalFeature).Maybe()
		s.mockedLocalEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(nil).Once().Maybe()

		powerConfig := &MonitorPowerConfig{
			ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
			ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}
		mpcInstance, err := NewMPC(s.mockedLocalEntity, s.Event, powerConfig, nil, nil, nil, nil)

		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.NotNil(s.T(), err)
	}

	// Test covering when calling GetOrAddIdForDescription Return nil
	{
		err := errors.New("test")
		mErr := model.ErrorType{}
		s.mockedElectricalConnectionFeature = mocks.NewElectricalConnectionServerInterface(s.T())
		s.mockedElectricalConnectionFeature.EXPECT().GetOrAddIdForDescription(mock.Anything).Return(nil, err).Maybe()
		s.mockedLocalEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer).Return(s.mockedLocalFeature).Maybe()
		s.mockedLocalEntity.EXPECT().FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer).Return(s.mockedLocalFeature).Maybe()
		s.mockedLocalFeature.EXPECT().UpdateData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(mErr)).Maybe()
		s.mockedLocalFeature.EXPECT().DataCopy(mock.Anything).Return(util.Ptr(model.ElectricalConnectionDescriptionListDataType{})).Maybe()
		powerConfig := &MonitorPowerConfig{
			ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
			ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}
		mpcInstance, err := NewMPC(s.mockedLocalEntity, s.Event, powerConfig, nil, nil, nil, nil)

		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.NotNil(s.T(), err)
	}
}
