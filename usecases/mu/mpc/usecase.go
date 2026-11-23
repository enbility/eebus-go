package mpc

import (
	"errors"
	"fmt"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type PhaseMeasurementIdMap map[model.ElectricalConnectionPhaseNameType]*model.MeasurementIdType

type MPC struct {
	*usecase.UseCaseBase

	powerConfig     *MonitorPowerConfig
	energyConfig    *MonitorEnergyConfig
	currentConfig   *MonitorCurrentConfig
	voltageConfig   *MonitorVoltageConfig
	frequencyConfig *MonitorFrequencyConfig

	acPowerTotal      *model.MeasurementIdType
	acPowerPerPhase   PhaseMeasurementIdMap
	acEnergyConsumed  *model.MeasurementIdType
	acEnergyProduced  *model.MeasurementIdType
	acCurrentPerPhase PhaseMeasurementIdMap
	acVoltagePerPhase PhaseMeasurementIdMap
	acFrequency       *model.MeasurementIdType
}

// creates a new MPC usecase instance for a MonitoredUnit entity
//
// parameters:
//   - localEntity: the local entity for which to construct an MPC instance
//   - eventCB: the callback to notify about events for this usecase
//   - monitorPowerConfig: (required) configuration parameters for MPC scenario 1
//   - monitorEnergyConfig: (optional) configuration parameters for MPC scenario 2, nil if not supported
//   - monitorCurrentConfig: (optional) configuration parameters for MPC scenario 3, nil if not supported
//   - monitorVoltageConfig: (optional) configuration parameters for MPC scenario 4, nil if not supported
//   - monitorFrequencyConfig: (optional) configuration parameters for MPC scenario, nil if not supported
//
// possible errors:
//   - if required fields in parameters are unset
func NewMPC(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
	monitorPowerConfig *MonitorPowerConfig,
	monitorEnergyConfig *MonitorEnergyConfig,
	monitorCurrentConfig *MonitorCurrentConfig,
	monitorVoltageConfig *MonitorVoltageConfig,
	monitorFrequencyConfig *MonitorFrequencyConfig,
) (*MPC, error) {
	if monitorPowerConfig == nil {
		return nil, errors.New("the monitor power config for the MPC-Use-Case must not be nil")
	}

	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeMonitoringAppliance}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
	}

	if monitorEnergyConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		})
	}

	if monitorCurrentConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(3),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		})
	}

	if monitorVoltageConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(4),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		})
	}

	if monitorFrequencyConfig != nil {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:  model.UseCaseScenarioSupportType(5),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		})
	}

	u := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoredUnit,
		model.UseCaseNameTypeMonitoringOfPowerConsumption,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		nil,
		true,
	)

	uc := &MPC{
		UseCaseBase:     u,
		powerConfig:     monitorPowerConfig,
		energyConfig:    monitorEnergyConfig,
		currentConfig:   monitorCurrentConfig,
		voltageConfig:   monitorVoltageConfig,
		frequencyConfig: monitorFrequencyConfig,
	}
	uc.acPowerPerPhase = PhaseMeasurementIdMap{}
	uc.acCurrentPerPhase = PhaseMeasurementIdMap{}
	uc.acVoltagePerPhase = PhaseMeasurementIdMap{}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc, nil
}

func (e *MPC) AddFeatures() error {
	// server features
	electricalConnectionFeature := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	if electricalConnectionFeature == nil {
		return errors.New("could not add feature: " + string(model.FeatureTypeTypeElectricalConnection))
	}
	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	measurementFeature := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	if measurementFeature == nil {
		return errors.New("could not add feature: " + string(model.FeatureTypeTypeMeasurement))
	}
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		return err
	}

	electricalConnection, err := server.NewElectricalConnection(e.LocalEntity)
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
		measurements *server.Measurement,
		electricalConnection api.ElectricalConnectionServerInterface,
		electricalConnectionId *model.ElectricalConnectionIdType,
		measurementsConstraintData *[]model.MeasurementConstraintsDataType,
	) error{
		e.configureMonitorPower,
		e.configureMonitorEnergy,
		e.configureMonitorCurrent,
		e.configureMonitorVoltage,
		e.configureMonitorFrequency,
	}

	for _, configMethod := range configMethods {
		if err := configMethod(measurements, electricalConnection, electricalConnectionId, &constraints); err != nil {
			return err
		}
	}

	// if any of the configured measurements set constraints, update the
	// measurementFeature with those accumulated constraints
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

func (e *MPC) configureMonitorPower(
	measurements *server.Measurement,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.powerConfig == nil {
		return errors.New("mpc monitoring power must be configured")
	}

	e.acPowerTotal = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})

	// if constraints are configured for acPowerTotal, set the
	// constraint id and update measurementsConstraintData
	if e.powerConfig.ValueConstraintsTotal != nil {
		e.powerConfig.ValueConstraintsTotal.MeasurementId = e.acPowerTotal
		*measurementsConstraintData = append(*measurementsConstraintData, *e.powerConfig.ValueConstraintsTotal)
	}

	parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  electricalConnectionId,
		MeasurementId:           e.acPowerTotal,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameType(e.powerConfig.ConnectedPhases)),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	}

	parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
	if parameterDescriptionId == nil {
		return errors.New("could not add parameter description")
	}

	for phase := range e.powerConfig.ValueSourcePerPhase {
		if !e.powerConfig.SupportsPhases([]string{string(phase)}) {
			errStr := fmt.Sprintf("power configuration for phase %s is not supported, please check the configuration", phase)
			return errors.New(errStr)
		}
		e.acPowerPerPhase[phase] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
		})

		parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  electricalConnectionId,
			MeasurementId:           e.acPowerPerPhase[phase],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(phase),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
		if parameterDescriptionId == nil {
			return errors.New("could not add parameter description")
		}

		if e.powerConfig.ValueConstraintsPerPhase[phase] == nil {
			continue
		}
		e.powerConfig.ValueConstraintsPerPhase[phase].MeasurementId = e.acPowerPerPhase[phase]
		*measurementsConstraintData = append(*measurementsConstraintData, *e.powerConfig.ValueConstraintsPerPhase[phase])

	}

	return nil
}

func (e *MPC) configureMonitorEnergy(
	measurements *server.Measurement,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.energyConfig == nil {
		return nil
	}

	if e.energyConfig.ValueSourceConsumption != nil {
		e.acEnergyConsumed = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
		})

		// if constraints are configured for acEnergyConsumed, set the
		// constraint id and update measurementsConstraintData
		if e.energyConfig.ValueConstraintsConsumption != nil {
			e.energyConfig.ValueConstraintsConsumption.MeasurementId = e.acEnergyConsumed
			*measurementsConstraintData = append(*measurementsConstraintData, *e.energyConfig.ValueConstraintsConsumption)
		}

		parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: electricalConnectionId,
			MeasurementId:          e.acEnergyConsumed,
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		}

		parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
		if parameterDescriptionId == nil {
			return errors.New("could not add parameter description")
		}
	}

	if e.energyConfig.ValueSourceProduction != nil {
		e.acEnergyProduced = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
		})

		// if constraints are configured for acEnergyProduced, set the
		// constraint id and update measurementsConstraintData
		if e.energyConfig.ValueConstraintsProduction != nil {
			e.energyConfig.ValueConstraintsProduction.MeasurementId = e.acEnergyProduced
			*measurementsConstraintData = append(*measurementsConstraintData, *e.energyConfig.ValueConstraintsProduction)
		}

		p4 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: electricalConnectionId,
			MeasurementId:          e.acEnergyProduced,
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		}
		idP4 := electricalConnection.AddParameterDescription(p4)
		if idP4 == nil {
			return errors.New("could not add parameter description")
		}
	}

	return nil
}

func (e *MPC) configureMonitorCurrent(
	measurements *server.Measurement,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.currentConfig == nil {
		return nil
	}

	for phase := range e.currentConfig.ValueSourcePerPhase {
		if !e.powerConfig.SupportsPhases([]string{string(phase)}) {
			errStr := fmt.Sprintf("power configuration for phase %s is not supported, please check the configuration", phase)
			return errors.New(errStr)
		}
		e.acCurrentPerPhase[phase] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
		})

		parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: electricalConnectionId,
			MeasurementId:          e.acCurrentPerPhase[phase],
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:       util.Ptr(phase),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}

		parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
		if parameterDescriptionId == nil {
			return errors.New("could not add parameter description")
		}

		if e.currentConfig.ValueConstraintsPerPhase[phase] == nil {
			continue
		}
		e.currentConfig.ValueConstraintsPerPhase[phase].MeasurementId = e.acCurrentPerPhase[phase]
		*measurementsConstraintData = append(*measurementsConstraintData, *e.currentConfig.ValueConstraintsPerPhase[phase])

	}

	return nil
}

func (e *MPC) configureMonitorVoltage(
	measurements *server.Measurement,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.voltageConfig == nil {
		return nil
	}

	for phase := range e.voltageConfig.ValueSourcePerPhase {
		switch len(string(phase)) {
		case 1:
			if !e.powerConfig.SupportsPhases([]string{string(phase)}) {
				errStr := fmt.Sprintf("power configuration for phase %s is not supported, please check the configuration", phase)
				return errors.New(errStr)
			}
		case 2:
			fromPhase := string(string(phase)[0])
			toPhase := string(string(phase)[1])
			if !e.powerConfig.SupportsPhases([]string{fromPhase, toPhase}) {
				errStr := fmt.Sprintf("power configuration for phase %s is not supported, please check the configuration", phase)
				return errors.New(errStr)
			}
		}
		e.acVoltagePerPhase[phase] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
		})
		parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{}

		switch len(string(phase)) {
		case 1:
			parameterDescription = model.ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId:  electricalConnectionId,
				MeasurementId:           e.acVoltagePerPhase[phase],
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(phase),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			}
		case 2:
			fromPhase := string(string(phase)[0])
			toPhase := string(string(phase)[1])
			parameterDescription = model.ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId:  electricalConnectionId,
				MeasurementId:           e.acVoltagePerPhase[phase],
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameType(fromPhase)),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameType(toPhase)),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			}
		}
		parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
		if parameterDescriptionId == nil {
			return errors.New("could not add parameter description")
		}
	}

	return nil
}

func (e *MPC) configureMonitorFrequency(
	measurements *server.Measurement,
	electricalConnection api.ElectricalConnectionServerInterface,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.frequencyConfig == nil {
		return nil
	}

	e.acFrequency = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	})

	// if constraints are configured for acFrequency, set the
	// constraint id and update measurementsConstraintData
	if e.frequencyConfig.ValueConstraints != nil {
		e.frequencyConfig.ValueConstraints.MeasurementId = e.acFrequency
		*measurementsConstraintData = append(*measurementsConstraintData, *e.frequencyConfig.ValueConstraints)
	}

	parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: electricalConnectionId,
		MeasurementId:          e.acFrequency,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
	}

	parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
	if parameterDescriptionId == nil {
		return errors.New("could not add parameter description")
	}

	return nil
}

func (e *MPC) getMeasurementDataForId(id *model.MeasurementIdType) (float64, error) {
	measurements, err := server.NewMeasurement(e.LocalEntity)
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
