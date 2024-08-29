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
	"time"
)

type MGCP struct {
	*usecase.UseCaseBase

	pvFeedInLimitationFactor *model.DeviceConfigurationKeyIdType
	acPowerTotal             *model.MeasurementIdType
	gridFeedIn               *model.MeasurementIdType
	gridConsumption          *model.MeasurementIdType
	acCurrent                [3]*model.MeasurementIdType
	acVoltage                [3]*model.MeasurementIdType // Phase to phase voltages are not supported (yet)
	acFrequency              *model.MeasurementIdType
}

var _ ucapi.GcpMGCPInterface = (*MGCP)(nil)

// At the moment the MGCP use case configures itself as a 3-phase meter by default (ABC).
func NewMGCP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MGCP {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeGridConnectionPoint}
	var validEntityTypes []model.EntityTypeType = nil // all entity types are valid
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
		model.UseCaseActorTypeGridConnectionPoint,
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

func (m *MGCP) AddFeatures() {
	// server features
	deviceConfigurationFeature := m.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	deviceConfigurationFeature.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	deviceConfigurationFeature.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, false)

	measurementFeature := m.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, true, false)
	measurementFeature.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)

	electricalConnectionFeature := m.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	electricalConnectionFeature.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)

	configuration, err := server.NewDeviceConfiguration(m.LocalEntity)
	if err != nil {
		panic(err)
	}

	m.pvFeedInLimitationFactor = configuration.AddKeyValueDescription(model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
		Unit:      util.Ptr(model.UnitOfMeasurementTypepct),
	})
	if m.pvFeedInLimitationFactor == nil {
		panic("failed to add key description")
	}

	measurement, err := server.NewMeasurement(m.LocalEntity)
	if err != nil {
		panic(err)
	}

	m.acPowerTotal = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})

	m.gridFeedIn = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
	})

	m.gridConsumption = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
	})

	for i := 0; i < len(m.acCurrent); i++ {
		m.acCurrent[i] = measurement.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
		})
	}

	for i := 0; i < len(m.acVoltage); i++ {
		m.acVoltage[i] = measurement.AddDescription(model.MeasurementDescriptionDataType{
			MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			Unit:            util.Ptr(model.UnitOfMeasurementTypeV),
			ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
		})
	}

	m.acFrequency = measurement.AddDescription(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeHz),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	})

	if m.acPowerTotal == nil || m.gridFeedIn == nil || m.gridConsumption == nil || m.acCurrent[0] == nil || m.acCurrent[1] == nil || m.acCurrent[2] == nil ||
		m.acVoltage[0] == nil || m.acVoltage[1] == nil || m.acVoltage[2] == nil || m.acFrequency == nil {
		panic("failed to add measurement description")
	}

	electricalConnection, err := server.NewElectricalConnection(m.LocalEntity)
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
		MeasurementId:           m.acPowerTotal,
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP2 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          m.gridFeedIn,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
	})

	idP3 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          m.gridConsumption,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
	})

	idP41 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          m.acCurrent[0],
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP42 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          m.acCurrent[1],
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP43 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          m.acCurrent[2],
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
		AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP51 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           m.acVoltage[0],
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP52 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           m.acVoltage[1],
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP53 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId:  util.Ptr(idEc1),
		MeasurementId:           m.acVoltage[2],
		VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
		AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
		AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeApparent),
		AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
	})

	idP6 := electricalConnection.AddParameterDescription(model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(idEc1),
		MeasurementId:          m.acFrequency,
		VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
	})

	if idP1 == nil || idP2 == nil || idP3 == nil || idP41 == nil || idP42 == nil || idP43 == nil || idP51 == nil ||
		idP52 == nil || idP53 == nil || idP6 == nil {
		panic("failed to add electrical connection parameter description")
	}

	err = m.SetPvFeedInLimitationFactor(0.0)
	if err != nil {
		panic(err)
	}

	for _, meas := range []*model.MeasurementIdType{
		m.acPowerTotal,
		m.gridFeedIn,
		m.gridConsumption,
		m.acCurrent[0], m.acCurrent[1], m.acCurrent[2],
		m.acVoltage[0], m.acVoltage[1], m.acVoltage[2],
		m.acFrequency,
	} {
		err = m.setMeasurementDataForId(meas, 0.0)
		if err != nil {
			panic(err)
		}
	}
}

func (m *MGCP) setMeasurementDataForId(id *model.MeasurementIdType, value float64) error {
	measurements, err := server.NewMeasurement(m.LocalEntity)
	if err != nil {
		return err
	}

	err = measurements.UpdateDataForId(model.MeasurementDataType{
		MeasurementId: id,
		ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(time.Now()),
		Value:         model.NewScaledNumberType(value),
		ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
	}, nil, *id)
	if err != nil {
		return err
	}

	return nil
}
