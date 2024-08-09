package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
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

	electricalConnection, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		panic(err)
	}

	err = electricalConnection.AddDescription(model.ElectricalConnectionDescriptionDataType{
		PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
	})
	if err != nil {
		panic(err)
	}

}
