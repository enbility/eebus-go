package mpc

import (
	"errors"
	"time"

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

func (e *MPC) AddFeatures() {
	// server features
	electricalConnectionFeatrue := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	electricalConnectionFeatrue.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	electricalConnectionFeatrue.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	measurementFeature := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	var phases = [][]string{
		{"a"},
		{"b"},
		{"c"},
		{"a", "b"},
		{"b", "c"},
		{"c", "a"},
	}

	var constraints = make([]model.MeasurementConstraintsDataType, 0)

	e.acPowerTotal = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})

	if e.powerConfig.ValueConstraintsTotal != nil {
		e.powerConfig.ValueConstraintsTotal.MeasurementId = e.acPowerTotal
		constraints = append(constraints, *e.powerConfig.ValueConstraintsTotal)
	}

	acPowerConstraints := []*model.MeasurementConstraintsDataType{
		e.powerConfig.ValueConstraintsPhaseA,
		e.powerConfig.ValueConstraintsPhaseB,
		e.powerConfig.ValueConstraintsPhaseC,
	}
	for id := 0; id < len(e.acPower); id++ {
		if e.powerConfig.SupportsPhases(phases[id]) {
			e.acPower[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			})
			if acPowerConstraints[id] != nil {
				acPowerConstraints[id].MeasurementId = e.acPower[id]
				constraints = append(constraints, *acPowerConstraints[id])
			}
		}
	}

	if e.energyConfig != nil {
		if e.energyConfig.ValueSourceConsumption != nil {
			e.acEnergyConsumed = measurements.AddDescription(model.MeasurementDescriptionDataType{
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
			})
			if e.energyConfig.ValueConstraintsConsumption != nil {
				e.energyConfig.ValueConstraintsConsumption.MeasurementId = e.acEnergyConsumed
				constraints = append(constraints, *e.energyConfig.ValueConstraintsConsumption)
			}
		}

		if e.energyConfig.ValueSourceProduction != nil {
			e.acEnergyProduced = measurements.AddDescription(model.MeasurementDescriptionDataType{
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
			})
			if e.energyConfig.ValueConstraintsProduction != nil {
				e.energyConfig.ValueConstraintsProduction.MeasurementId = e.acEnergyProduced
				constraints = append(constraints, *e.energyConfig.ValueConstraintsProduction)
			}
		}
	}

	if e.currentConfig != nil {
		acCurrentConstraints := []*model.MeasurementConstraintsDataType{
			e.currentConfig.ValueConstraintsPhaseA,
			e.currentConfig.ValueConstraintsPhaseB,
			e.currentConfig.ValueConstraintsPhaseC,
		}
		for id := 0; id < len(e.acCurrent); id++ {
			if e.powerConfig.SupportsPhases(phases[id]) {
				e.acCurrent[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
					MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
					CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
					Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
					ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				})
				if acCurrentConstraints[id] != nil {
					acCurrentConstraints[id].MeasurementId = e.acCurrent[id]
					constraints = append(constraints, *acCurrentConstraints[id])
				}
			}
		}
	}

	if e.voltageConfig != nil {
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
				if len(phases[id]) == 2 && !e.voltageConfig.SupportPhaseToPhase {
					continue
				}
				e.acVoltage[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
					MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
					CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
					Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
					ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
				})
				if acVoltageConstraints[id] != nil {
					acVoltageConstraints[id].MeasurementId = e.acVoltage[id]
					constraints = append(constraints, *acVoltageConstraints[id])
				}
			}
		}
	}

	if e.frequencyConfig != nil {
		e.acFrequency = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
		})
		if e.frequencyConfig.ValueConstraints != nil {
			e.frequencyConfig.ValueConstraints.MeasurementId = e.acFrequency
			constraints = append(constraints, *e.frequencyConfig.ValueConstraints)
		}
	}

	electricalConnection, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	idEc1 := model.ElectricalConnectionIdType(0)
	ec1 := model.ElectricalConnectionDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
	}
	if err := electricalConnection.AddDescription(ec1); err != nil {
		panic(err)
	}

	p1 := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.acPowerTotal,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameType(e.powerConfig.ConnectedPhases)),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	}
	idP1 := electricalConnection.AddParameterDescription(p1)
	if idP1 == nil {
		panic("error adding parameter description")
	}

	if e.acPower[0] != nil {
		p21 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acPower[0],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP21 := electricalConnection.AddParameterDescription(p21)
		if idP21 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acPower[1] != nil {
		p22 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acPower[1],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP22 := electricalConnection.AddParameterDescription(p22)
		if idP22 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acPower[2] != nil {
		p23 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acPower[2],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP23 := electricalConnection.AddParameterDescription(p23)
		if idP23 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acEnergyConsumed != nil {
		p3 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(idEc1),
			MeasurementId:          e.acEnergyConsumed,
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		}
		idP3 := electricalConnection.AddParameterDescription(p3)
		if idP3 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acEnergyProduced != nil {
		p4 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(idEc1),
			MeasurementId:          e.acEnergyProduced,
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		}
		idP4 := electricalConnection.AddParameterDescription(p4)
		if idP4 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acCurrent[0] != nil {
		p51 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(idEc1),
			MeasurementId:          e.acCurrent[0],
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP51 := electricalConnection.AddParameterDescription(p51)
		if idP51 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acCurrent[1] != nil {
		p52 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(idEc1),
			MeasurementId:          e.acCurrent[1],
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP52 := electricalConnection.AddParameterDescription(p52)
		if idP52 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acCurrent[2] != nil {
		p53 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(idEc1),
			MeasurementId:          e.acCurrent[2],
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP53 := electricalConnection.AddParameterDescription(p53)
		if idP53 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acVoltage[0] != nil {
		p61 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acVoltage[0],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP61 := electricalConnection.AddParameterDescription(p61)
		if idP61 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acVoltage[1] != nil {
		p62 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acVoltage[1],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP62 := electricalConnection.AddParameterDescription(p62)
		if idP62 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acVoltage[2] != nil {
		p63 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acVoltage[2],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP63 := electricalConnection.AddParameterDescription(p63)
		if idP63 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acVoltage[3] != nil {
		p64 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acVoltage[3],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP64 := electricalConnection.AddParameterDescription(p64)
		if idP64 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acVoltage[4] != nil {
		p65 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acVoltage[4],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP65 := electricalConnection.AddParameterDescription(p65)
		if idP65 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acVoltage[5] != nil {
		p66 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(idEc1),
			MeasurementId:           e.acVoltage[5],
			VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
			AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
		}
		idP66 := electricalConnection.AddParameterDescription(p66)
		if idP66 == nil {
			panic("error adding parameter description")
		}
	}

	if e.acFrequency != nil {
		p7 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(idEc1),
			MeasurementId:          e.acFrequency,
			VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		}
		idP7 := electricalConnection.AddParameterDescription(p7)
		if idP7 == nil {
			panic("error adding parameter description")
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

func measurementData(
	value float64,
	timestamp *time.Time,
	valueSource *model.MeasurementValueSourceType,
	valueState *model.MeasurementValueStateType,
	evaluationStart *time.Time,
	evaluationEnd *time.Time,
) model.MeasurementDataType {
	measurement := model.MeasurementDataType{
		ValueType:   util.Ptr(model.MeasurementValueTypeTypeValue),
		Value:       model.NewScaledNumberType(value),
		ValueSource: valueSource,
		ValueState:  valueState,
	}

	if timestamp != nil {
		measurement.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(*timestamp)
	}

	if evaluationStart != nil && evaluationEnd != nil {
		measurement.EvaluationPeriod = &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationStart),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationEnd),
		}
	}

	return measurement
}
