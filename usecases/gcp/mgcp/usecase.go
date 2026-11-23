package mgcp

import (
	"errors"
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type MGCP struct {
	*usecase.UseCaseBase

	limitationConfig *MonitorPvFeedInPowerLimitationFactorConfig
	powerConfig      *MonitorPowerConfig
	energyConfig     *MonitorEnergyConfig
	currentConfig    *MonitorCurrentConfig
	voltageConfig    *MonitorVoltageConfig
	frequencyConfig  *MonitorFrequencyConfig

	pvFeedInLimitationFactor *model.DeviceConfigurationKeyIdType
	acPowerTotal             *model.MeasurementIdType
	gridFeedIn               *model.MeasurementIdType
	gridConsumption          *model.MeasurementIdType
	acCurrent                [3]*model.MeasurementIdType
	acVoltage                [6]*model.MeasurementIdType
	acFrequency              *model.MeasurementIdType
}

// creates a new MGCP usecase instance for a MonitoredUnit entity
//
// parameters:
//   - localEntity: the local entity for which the use case is created
//   - eventCB: the callback function to be called when an event is triggered
//   - monitorFeedInLimitationConfig: (optional) configuration parameters for MGCP scenario 1
//   - monitorPowerConfig: (required) configuration parameters for MGCP scenario 2
//   - monitorEnergyConfig: (required) configuration parameters for MGCP scenario 3
//   - monitorCurrentConfig: (optional) configuration parameters for MGCP scenario 4
//   - monitorVoltageConfig: (optional) configuration parameters for MGCP scenario 5
//   - monitorFrequencyConfig: (optional) configuration parameters for MGCP scenario 6
//
// possible errors:
//   - configuration error
//   - and more...

func NewMGCP(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
	monitorFeedInLimitationConfig *MonitorPvFeedInPowerLimitationFactorConfig,
	monitorPowerConfig *MonitorPowerConfig,
	monitorEnergyConfig *MonitorEnergyConfig,
	monitorCurrentConfig *MonitorCurrentConfig,
	monitorVoltageConfig *MonitorVoltageConfig,
	monitorFrequencyConfig *MonitorFrequencyConfig,
) (*MGCP, error) {
	if monitorPowerConfig == nil {
		return nil, errors.New("monitorPowerConfig must be set")
	}

	if monitorEnergyConfig == nil {
		return nil, errors.New("monitorEnergyConfig must be set")
	}

	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeMonitoringAppliance}
	useCaseScenarios := make([]api.UseCaseScenario, 0)

	if monitorFeedInLimitationConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
		})
	}

	useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
		Scenario:  model.UseCaseScenarioSupportType(2),
		Mandatory: true,
		ServerFeatures: []model.FeatureTypeType{
			model.FeatureTypeTypeMeasurement,
			model.FeatureTypeTypeElectricalConnection,
		},
	})

	useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
		Scenario:  model.UseCaseScenarioSupportType(3),
		Mandatory: true,
		ServerFeatures: []model.FeatureTypeType{
			model.FeatureTypeTypeMeasurement,
			model.FeatureTypeTypeElectricalConnection,
		},
	})

	useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
		Scenario:  model.UseCaseScenarioSupportType(4),
		Mandatory: false,
		ServerFeatures: []model.FeatureTypeType{
			model.FeatureTypeTypeMeasurement,
			model.FeatureTypeTypeElectricalConnection,
		},
	})

	if monitorCurrentConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(5),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		})
	}

	if monitorVoltageConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(6),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		})
	}

	if monitorFrequencyConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(7),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		})
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfGridConnectionPoint,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		nil,
		true,
	)

	uc := &MGCP{
		UseCaseBase: usecase,

		limitationConfig: monitorFeedInLimitationConfig,
		powerConfig:      monitorPowerConfig,
		energyConfig:     monitorEnergyConfig,
		currentConfig:    monitorCurrentConfig,
		voltageConfig:    monitorVoltageConfig,
		frequencyConfig:  monitorFrequencyConfig,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc, nil
}

func (m *MGCP) AddFeatures() error {
	// server features
	deviceConfigurationFeature := m.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	measurementFeature := m.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	electricalConnectionFeature := m.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)

	deviceConfigurationFeature.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	deviceConfigurationFeature.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, false)

	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	configuration, err := server.NewDeviceConfiguration(m.LocalEntity)
	if err != nil {
		return err
	}

	err = m.configurePvFeedInLimitationFactor(configuration)
	if err != nil {
		return err
	}

	measurement, err := server.NewMeasurement(m.LocalEntity)
	if err != nil {
		return err
	}

	electricalConnection, err := server.NewElectricalConnection(m.LocalEntity)
	if err != nil {
		return err
	}

	electricalConnectionId, err := electricalConnection.GetOrAddIdForDescription(model.ElectricalConnectionDescriptionDataType{
		PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
	})
	if err != nil {
		return err
	}

	constraints := make([]model.MeasurementConstraintsDataType, 0)

	configMethods := []func(
		measurements api.MeasurementServerInterface,
		electricalConnection api.ElectricalConnectionServerInterface,
		electricalConnectionId *model.ElectricalConnectionIdType,
		measurementsConstraintData *[]model.MeasurementConstraintsDataType,
	) error{
		m.configureMonitorPower,
		m.configureGridFeedIn,
		m.configureGridConsumption,
		m.configureMonitorCurrent,
		m.configureMonitorVoltage,
		m.configureMonitorFrequency,
	}

	for _, configMethod := range configMethods {
		if err := configMethod(measurement, electricalConnection, electricalConnectionId, &constraints); err != nil {
			return err
		}
	}

	if len(constraints) > 0 {
		measurementFeature.UpdateData(
			model.FunctionTypeMeasurementConstraintsListData,
			&model.MeasurementConstraintsListDataType{
				MeasurementConstraintsData: constraints,
			}, nil, nil,
		)
	}

	return nil
}

func (m *MGCP) configurePvFeedInLimitationFactor(
	configurations api.DeviceConfigurationServerInterface,
) error {
	if m.limitationConfig == nil {
		return nil
	}

	m.pvFeedInLimitationFactor = configurations.AddKeyValueDescription(model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
		Unit:      util.Ptr(model.UnitOfMeasurementTypepct),
	})

	if m.pvFeedInLimitationFactor == nil {
		return errors.New("failed to add key description")
	}

	return nil
}

func (m *MGCP) configureMonitorPower(
	measurements api.MeasurementServerInterface,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if m.powerConfig == nil {
		return errors.New("mgcp power config must be configured")
	}

	if m.powerConfig.ValueSource == nil {
		return errors.New("mgcp power config value source must be configured")
	}

	m.acPowerTotal = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})

	parameterDescription := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  electricalConnectionId,
		MeasurementId:           m.acPowerTotal,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	if parameterDescription == nil {
		return errors.New("failed to add parameter description")
	}

	if m.powerConfig.ValueConstraints != nil {
		m.powerConfig.ValueConstraints.MeasurementId = m.acPowerTotal
		*measurementsConstraintData = append(*measurementsConstraintData, *m.powerConfig.ValueConstraints)
	}

	return nil
}

func (m *MGCP) configureGridFeedIn(
	measurements api.MeasurementServerInterface,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if m.energyConfig == nil {
		return errors.New("mgcp energy config must be configured")
	}

	if m.energyConfig.ValueSourceProduction == nil {
		return errors.New("mgcp energy config production value source must be configured")
	}

	m.gridFeedIn = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
	})

	parameterDescription := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: electricalConnectionId,
		MeasurementId:          m.gridFeedIn,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
	})
	if parameterDescription == nil {
		return errors.New("failed to add parameter description")
	}

	if m.energyConfig.ValueConstraintsProduction != nil {
		m.energyConfig.ValueConstraintsProduction.MeasurementId = m.gridFeedIn
		*measurementsConstraintData = append(*measurementsConstraintData, *m.energyConfig.ValueConstraintsProduction)
	}

	return nil
}

func (m *MGCP) configureGridConsumption(
	measurements api.MeasurementServerInterface,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if m.energyConfig == nil {
		return errors.New("mgcp energy config must be configured")
	}

	if m.energyConfig.ValueSourceConsumption == nil {
		return errors.New("value source consumption must be set")
	}

	m.gridConsumption = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
	})

	parameterDescription := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: electricalConnectionId,
		MeasurementId:          m.gridConsumption,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
	})
	if parameterDescription == nil {
		return errors.New("failed to add parameter description")
	}

	if m.energyConfig.ValueConstraintsConsumption != nil {
		m.energyConfig.ValueConstraintsConsumption.MeasurementId = m.gridConsumption
		*measurementsConstraintData = append(*measurementsConstraintData, *m.energyConfig.ValueConstraintsConsumption)
	}

	return nil
}

func (m *MGCP) configureMonitorCurrent(
	measurements api.MeasurementServerInterface,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if m.currentConfig == nil {
		return nil
	}

	valueSourcesOfPhases := []*model.MeasurementValueSourceType{
		m.currentConfig.ValueSourcePhaseA,
		m.currentConfig.ValueSourcePhaseB,
		m.currentConfig.ValueSourcePhaseC,
	}

	valueConstraints := []*model.MeasurementConstraintsDataType{
		m.currentConfig.ValueConstraintsPhaseA,
		m.currentConfig.ValueConstraintsPhaseB,
		m.currentConfig.ValueConstraintsPhaseC,
	}

	electricalConnectedPhases := []model.ElectricalConnectionPhaseNameType{
		model.ElectricalConnectionPhaseNameTypeA,
		model.ElectricalConnectionPhaseNameTypeB,
		model.ElectricalConnectionPhaseNameTypeC,
	}

	for i := 0; i < len(m.acCurrent); i++ {
		if valueSourcesOfPhases[i] == nil {
			continue
		}

		m.acCurrent[i] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
		})

		parameterDescription := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: electricalConnectionId,
			MeasurementId:          m.acCurrent[i],
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:       util.Ptr(electricalConnectedPhases[i]),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		})

		if parameterDescription == nil {
			return errors.New("failed to add parameter description")
		}

		if valueConstraints[i] != nil {
			valueConstraints[i].MeasurementId = m.acCurrent[i]
			*measurementsConstraintData = append(*measurementsConstraintData, *valueConstraints[i])
		}
	}

	return nil
}

func (m *MGCP) configureMonitorVoltage(
	measurements api.MeasurementServerInterface,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if m.voltageConfig == nil {
		return nil
	}

	electricalConnectionPhaseToPhase := [][]model.ElectricalConnectionPhaseNameType{
		{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeNeutral},
		{model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeNeutral},
		{model.ElectricalConnectionPhaseNameTypeC, model.ElectricalConnectionPhaseNameTypeNeutral},
		{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB},
		{model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC},
		{model.ElectricalConnectionPhaseNameTypeC, model.ElectricalConnectionPhaseNameTypeA},
	}

	valueSourcesOfPhases := []*model.MeasurementValueSourceType{
		m.voltageConfig.ValueSourcePhaseA,
		m.voltageConfig.ValueSourcePhaseB,
		m.voltageConfig.ValueSourcePhaseC,
		m.voltageConfig.ValueSourcePhaseAToB,
		m.voltageConfig.ValueSourcePhaseBToC,
		m.voltageConfig.ValueSourcePhaseCToA,
	}

	valueConstraintsOfPhases := []*model.MeasurementConstraintsDataType{
		m.voltageConfig.ValueConstraintsPhaseA,
		m.voltageConfig.ValueConstraintsPhaseB,
		m.voltageConfig.ValueConstraintsPhaseC,
		m.voltageConfig.ValueConstraintsPhaseAToB,
		m.voltageConfig.ValueConstraintsPhaseBToC,
		m.voltageConfig.ValueConstraintsPhaseCToA,
	}

	for i := 0; i < len(m.acVoltage); i++ {
		if valueSourcesOfPhases[i] == nil {
			continue
		}

		m.acVoltage[i] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
		})

		parameterDescription := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  electricalConnectionId,
			MeasurementId:           m.acVoltage[i],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(electricalConnectionPhaseToPhase[i][0]),
			AcMeasuredInReferenceTo: util.Ptr(electricalConnectionPhaseToPhase[i][1]),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		})

		if parameterDescription == nil {
			return errors.New("failed to add parameter description")
		}

		if valueConstraintsOfPhases[i] != nil {
			valueConstraintsOfPhases[i].MeasurementId = m.acVoltage[i]
			*measurementsConstraintData = append(*measurementsConstraintData, *valueConstraintsOfPhases[i])
		}
	}

	return nil
}

func (m *MGCP) configureMonitorFrequency(
	measurements api.MeasurementServerInterface,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if m.frequencyConfig == nil {
		return nil
	}

	m.acFrequency = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	})

	parameterDescription := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: electricalConnectionId,
		MeasurementId:          m.acFrequency,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
	})
	if parameterDescription == nil {
		return errors.New("failed to add parameter description")
	}

	if m.frequencyConfig.ValueConstraints != nil {
		m.frequencyConfig.ValueConstraints.MeasurementId = m.acFrequency
		*measurementsConstraintData = append(*measurementsConstraintData, *m.frequencyConfig.ValueConstraints)
	}

	return nil
}

func (m *MGCP) getMeasurementDataForId(id *model.MeasurementIdType) (float64, error) {
	if id == nil {
		return 0, api.ErrMissingData
	}

	measurements, err := server.NewMeasurement(m.LocalEntity)
	if err != nil {
		return 0, err
	}

	data, err := measurements.GetDataForId(*id)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, api.ErrDataNotAvailable
	}

	return data.Value.GetValue(), nil
}
