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
}

var _ ucapi.GcpMGCPInterface = (*MGCP)(nil)

// Add support for the Monitoring of Grid Connection Point (MGCP) use case
// as a Grid Connection Point actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewMGCP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MGCP {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeGridConnectionPoint}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCEM,
		model.EntityTypeTypeGridConnectionPointOfPremises,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(3),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(4),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(5),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(6),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(7),
			Mandatory: true,
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeGridConnectionPoint,
		model.UseCaseNameTypeMonitoringOfGridConnectionPoint,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &MGCP{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MGCP) AddFeatures() {
	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)

	if dcs, err := server.NewDeviceConfiguration(e.LocalEntity); err == nil {
		dcs.AddKeyValueDescription(
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
				Unit:      util.Ptr(model.UnitOfMeasurementTypepct),
			},
		)

		value := &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(0),
		}
		_ = dcs.UpdateKeyValueDataForFilter(
			model.DeviceConfigurationKeyValueDataType{
				Value:             value,
				IsValueChangeable: util.Ptr(true),
			},
			nil,
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
			},
		)
	}

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	descDataList := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(3)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(4)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(5)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(6)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(7)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(8)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(9)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
			},
		},
	}

	f.SetData(model.FunctionTypeMeasurementDescriptionListData, descDataList)

	if _, err := server.NewMeasurement(e.LocalEntity); err == nil {
	}

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(5)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(6)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(7)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(8)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(9)),
				Value:         model.NewScaledNumberType(0),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
	}

	f.SetData(model.FunctionTypeMeasurementListData, measData)

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	descList := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				AcConnectedPhases:       util.Ptr(uint(3)),
			},
		},
	}

	f.SetData(model.FunctionTypeElectricalConnectionDescriptionListData, descList)

	paramDescList := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(0)), // acPowerTotal
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)), // gridFeedIn
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)), // gridConsumption
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(3)), // current a
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(4)), // current b
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(5)), // current c
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(6)), // voltage a
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(7)), // voltage b
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(8)), // voltage c
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(9)), // frequency
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	f.SetData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDescList)

	if _, err := server.NewElectricalConnection(e.LocalEntity); err == nil {
	}
}
