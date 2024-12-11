package mpc

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
)

type MPC struct {
	*usecase.UseCaseBase

	powerConfig     *MonitorPowerConfig
	energyConfig    *MonitorEnergyConfig
	currentConfig   *MonitorCurrentConfig
	voltageConfig   *MonitorVoltageConfig
	frequencyConfig *MonitorFrequencyConfig

	acPowerTotal     *model.MeasurementIdType
	acPower          [3]*model.MeasurementIdType
	acEnergyConsumed *model.MeasurementIdType
	acEnergyProduced *model.MeasurementIdType
	acCurrent        [3]*model.MeasurementIdType
	acVoltage        [6]*model.MeasurementIdType // Phase to phase voltages are not supported (yet)
	acFrequency      *model.MeasurementIdType
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

	_ = spine.Events.Subscribe(uc)

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
		electricalConnection *server.ElectricalConnection,
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
	electricalConnection *server.ElectricalConnection,
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

	acPowerConstraints := []*model.MeasurementConstraintsDataType{
		e.powerConfig.ValueConstraintsPhaseA,
		e.powerConfig.ValueConstraintsPhaseB,
		e.powerConfig.ValueConstraintsPhaseC,
	}

	acMeasuredPhases := []*model.ElectricalConnectionPhaseNameType{
		util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
	}

	for id := 0; id < len(e.acPower); id++ {
		if e.powerConfig.SupportsPhases(phases[id]) {
			e.acPower[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			})

			// if constraints are configured for acPower[id], set the
			// constraint id and update measurementsConstraintData
			if acPowerConstraints[id] != nil {
				acPowerConstraints[id].MeasurementId = e.acPower[id]
				*measurementsConstraintData = append(*measurementsConstraintData, *acPowerConstraints[id])
			}

			parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId:  electricalConnectionId,
				MeasurementId:           e.acPower[id],
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        acMeasuredPhases[id],
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			}

			parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
			if parameterDescriptionId == nil {
				return errors.New("could not add parameter description")
			}
		}
	}

	return nil
}

func (e *MPC) configureMonitorEnergy(
	measurements *server.Measurement,
	electricalConnection *server.ElectricalConnection,
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
	electricalConnection *server.ElectricalConnection,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.currentConfig == nil {
		return nil
	}

	acCurrentConstraints := []*model.MeasurementConstraintsDataType{
		e.currentConfig.ValueConstraintsPhaseA,
		e.currentConfig.ValueConstraintsPhaseB,
		e.currentConfig.ValueConstraintsPhaseC,
	}

	acMeasuredPhases := []*model.ElectricalConnectionPhaseNameType{
		util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
	}

	for id := 0; id < len(e.acCurrent); id++ {
		if e.powerConfig.SupportsPhases(phases[id]) {
			e.acCurrent[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			})

			// if constraints are configured for acCurrent[id], set the
			// constraint id and update measurementsConstraintData
			if acCurrentConstraints[id] != nil {
				acCurrentConstraints[id].MeasurementId = e.acCurrent[id]
				*measurementsConstraintData = append(*measurementsConstraintData, *acCurrentConstraints[id])
			}

			parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId: electricalConnectionId,
				MeasurementId:          e.acCurrent[id],
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       acMeasuredPhases[id],
				AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			}

			parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
			if parameterDescriptionId == nil {
				return errors.New("could not add parameter description")
			}
		}
	}

	return nil
}

func (e *MPC) configureMonitorVoltage(
	measurements *server.Measurement,
	electricalConnection *server.ElectricalConnection,
	electricalConnectionId *model.ElectricalConnectionIdType,
	measurementsConstraintData *[]model.MeasurementConstraintsDataType,
) error {
	if e.voltageConfig == nil {
		return nil
	}

	acVoltagePhasesFrom := []*model.ElectricalConnectionPhaseNameType{
		util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
	}

	acVoltagePhasesTo := []*model.ElectricalConnectionPhaseNameType{
		util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
	}

	acVoltageConstraints := []*model.MeasurementConstraintsDataType{
		e.voltageConfig.ValueConstraintsPhaseA,
		e.voltageConfig.ValueConstraintsPhaseB,
		e.voltageConfig.ValueConstraintsPhaseC,
		e.voltageConfig.ValueConstraintsPhaseAToB,
		e.voltageConfig.ValueConstraintsPhaseBToC,
		e.voltageConfig.ValueConstraintsPhaseCToA,
	}

	for id := 0; id < len(e.acVoltage); id++ {
		if e.powerConfig.SupportsPhases(phases[id]) {
			// skip PhaseToPhase voltages if they're not supported
			if len(phases[id]) == 2 && !e.voltageConfig.SupportPhaseToPhase {
				continue
			}

			e.acVoltage[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			})

			// if constraints are configured for acVoltage[id], set the
			// constraint id and update measurementsConstraintData
			if acVoltageConstraints[id] != nil {
				acVoltageConstraints[id].MeasurementId = e.acVoltage[id]
				*measurementsConstraintData = append(*measurementsConstraintData, *acVoltageConstraints[id])
			}

			parameterDescription := model.ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId:  electricalConnectionId,
				MeasurementId:           e.acVoltage[id],
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        acVoltagePhasesFrom[id],
				AcMeasuredInReferenceTo: acVoltagePhasesTo[id],
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			}

			parameterDescriptionId := electricalConnection.AddParameterDescription(parameterDescription)
			if parameterDescriptionId == nil {
				return errors.New("could not add parameter description")
			}
		}
	}

	return nil
}

func (e *MPC) configureMonitorFrequency(
	measurements *server.Measurement,
	electricalConnection *server.ElectricalConnection,
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

var phases = [][]string{
	{"a"},
	{"b"},
	{"c"},
	{"a", "b"},
	{"b", "c"},
	{"c", "a"},
}
