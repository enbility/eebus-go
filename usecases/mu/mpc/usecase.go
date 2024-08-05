package mpc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"time"
)

type MPC struct {
	*usecase.UseCaseBase

	acPowerTotal     *model.MeasurementIdType
	acPower          [3]*model.MeasurementIdType
	acEnergyConsumed *model.MeasurementIdType
	acEnergyProduced *model.MeasurementIdType
	acCurrent        [3]*model.MeasurementIdType
	acVoltage        [3]*model.MeasurementIdType // Phase to phase voltages are not supported (yet)
	acFrequency      *model.MeasurementIdType
}

var _ ucapi.MuMPCInterface = (*MPC)(nil)

// At the moment the MPC use case configures itself as a 3-phase meter by default (ABC).
func NewMPC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MPC {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeMonitoringAppliance}
	var validEntityTypes []model.EntityTypeType = nil // all entity types are valid
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(3),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(4),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(5),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
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
		validEntityTypes)

	uc := &MPC{
		UseCaseBase: u,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MPC) AddFeatures() {
	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, true, false)
	f.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	e.acPowerTotal = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})

	for id := 0; id < len(e.acPower); id++ {
		e.acPower[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
		})
	}

	e.acEnergyConsumed = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
	})

	e.acEnergyProduced = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
	})

	for id := 0; id < len(e.acCurrent); id++ {
		e.acCurrent[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
		})
	}

	for id := 0; id < len(e.acVoltage); id++ {
		e.acVoltage[id] = measurements.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
		})
	}

	e.acFrequency = measurements.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	})

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
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	}
	idP1 := electricalConnection.AddParameterDescription(p1)
	if idP1 == nil {
		panic("error adding parameter description")
	}

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

	p7 := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.acFrequency,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
	}
	idP7 := electricalConnection.AddParameterDescription(p7)
	if idP7 == nil {
		panic("error adding parameter description")
	}

	for _, id := range []*model.MeasurementIdType{
		e.acPowerTotal,
		e.acPower[0], e.acPower[1], e.acPower[2],
		e.acEnergyConsumed,
		e.acEnergyProduced,
		e.acCurrent[0], e.acCurrent[1], e.acCurrent[2],
		e.acVoltage[0], e.acVoltage[1], e.acVoltage[2],
		e.acFrequency} {
		if err := e.setMeasurementDataForId(id, 0); err != nil {
			panic(err)
		}
	}
}

func (e *MPC) setMeasurementDataForId(id *model.MeasurementIdType, measurementData float64) error {
	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		return err
	}

	err = measurements.UpdateDataForId(model.MeasurementDataType{
		MeasurementId: id,
		ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(time.Now()),
		Value:         model.NewScaledNumberType(measurementData),
		ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
	}, nil, *id)

	return err
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
