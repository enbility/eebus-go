package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
)

type MGCP struct {
	*usecase.UseCaseBase

	idK1  *model.DeviceConfigurationKeyIdType
	idM1  *model.MeasurementIdType
	idM2  *model.MeasurementIdType
	idM3  *model.MeasurementIdType
	idM43 *model.MeasurementIdType
	idM41 *model.MeasurementIdType
	idM42 *model.MeasurementIdType
	idM51 *model.MeasurementIdType
	idM52 *model.MeasurementIdType
	idM53 *model.MeasurementIdType
	idM54 *model.MeasurementIdType
	idM55 *model.MeasurementIdType
	idM56 *model.MeasurementIdType
	idM6  *model.MeasurementIdType
}

var _ ucapi.GcpMGCPInterface = (*MGCP)(nil)

func NewMGCP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MGCP {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeGridConnectionPoint}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCEM,
		model.EntityTypeTypeGridConnectionPointOfPremises,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(3),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(4),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(5),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(6),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(7),
			Mandatory: false,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
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
		validEntityTypes)

	uc := &MGCP{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MGCP) AddFeatures() {
	// server features
	deviceConfigurationFeature := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	measurementFeature := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	electricalConnectionFeature := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)

	deviceConfigurationFeature.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	deviceConfigurationFeature.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, false)

	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	configuration, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	e.idK1 = configuration.AddKeyValueDescription(model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
		Unit:      util.Ptr(model.UnitOfMeasurementTypepct),
	})
	if e.idK1 == nil {
		panic("failed to add key description")
	}

	measurement, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	e.idM1 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})

	e.idM2 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
	})

	e.idM3 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
	})

	e.idM41 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	})

	e.idM42 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	})

	e.idM43 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	})

	e.idM51 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	})

	e.idM52 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	})

	e.idM53 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	})

	e.idM54 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	})

	e.idM55 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	})

	e.idM56 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	})

	e.idM6 = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	})

	if e.idM1 == nil || e.idM2 == nil || e.idM3 == nil || e.idM41 == nil || e.idM42 == nil || e.idM43 == nil ||
		e.idM51 == nil || e.idM52 == nil || e.idM53 == nil || e.idM54 == nil || e.idM55 == nil || e.idM56 == nil || e.idM6 == nil {
		panic("failed to add measurement description")

	}

	electricalConnection, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	idEc1 := model.ElectricalConnectionIdType(0)
	err = electricalConnection.AddDescription(model.ElectricalConnectionDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
	})
	if err != nil {
		panic(err)
	}

	idP1 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM1,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP2 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.idM2,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
	})

	idP3 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.idM3,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
	})

	idP41 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.idM41,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP42 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.idM42,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP43 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.idM43,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP51 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM51,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP52 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM52,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP53 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM53,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP54 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM54,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP55 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM55,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP56 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           e.idM56,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP6 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          e.idM6,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
	})

	if idP1 == nil || idP2 == nil || idP3 == nil || idP41 == nil || idP42 == nil || idP43 == nil || idP51 == nil ||
		idP52 == nil || idP53 == nil || idP54 == nil || idP55 == nil || idP56 == nil || idP6 == nil {
		panic("failed to add electrical connection parameter description")
	}

	for _, m := range []*model.MeasurementIdType{
		e.idM1,
		e.idM2,
		e.idM3,
		e.idM41,
		e.idM42,
		e.idM43,
		e.idM51,
		e.idM52,
		e.idM53,
		e.idM54,
		e.idM55,
		e.idM56,
		e.idM6,
	} {
		e.setMeasurementForId(m, 0.0)
	}
}

func (e *MGCP) setMeasurementForId(id *model.MeasurementIdType, value float64) {
	// TODO
}

func (e *MGCP) getMeasurementForId(id *model.MeasurementIdType) float64 {
	// TODO
	return 0.0
}
