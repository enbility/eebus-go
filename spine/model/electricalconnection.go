package model

type ElectricalConnectionIdType uint

type ElectricalConnectionParameterIdType uint

type ElectricalConnectionMeasurandVariantType ElectricalConnectionMeasurandVariantEnumType

type ElectricalConnectionMeasurandVariantEnumType string

const (
	ElectricalConnectionMeasurandVariantEnumTypeAmplitude     ElectricalConnectionMeasurandVariantEnumType = "amplitude"
	ElectricalConnectionMeasurandVariantEnumTypeRms           ElectricalConnectionMeasurandVariantEnumType = "rms"
	ElectricalConnectionMeasurandVariantEnumTypeInstantaneous ElectricalConnectionMeasurandVariantEnumType = "instantaneous"
	ElectricalConnectionMeasurandVariantEnumTypeAngle         ElectricalConnectionMeasurandVariantEnumType = "angle"
	ElectricalConnectionMeasurandVariantEnumTypeCosphi        ElectricalConnectionMeasurandVariantEnumType = "cosPhi"
)

type ElectricalConnectionVoltageTypeType ElectricalConnectionVoltageTypeEnumType

type ElectricalConnectionVoltageTypeEnumType string

const (
	ElectricalConnectionVoltageTypeEnumTypeAc ElectricalConnectionVoltageTypeEnumType = "ac"
	ElectricalConnectionVoltageTypeEnumTypeDc ElectricalConnectionVoltageTypeEnumType = "dc"
)

type ElectricalConnectionAcMeasurementTypeType ElectricalConnectionAcMeasurementTypeEnumType

type ElectricalConnectionAcMeasurementTypeEnumType string

const (
	ElectricalConnectionAcMeasurementTypeEnumTypeReal     ElectricalConnectionAcMeasurementTypeEnumType = "real"
	ElectricalConnectionAcMeasurementTypeEnumTypeReactive ElectricalConnectionAcMeasurementTypeEnumType = "reactive"
	ElectricalConnectionAcMeasurementTypeEnumTypeApparent ElectricalConnectionAcMeasurementTypeEnumType = "apparent"
	ElectricalConnectionAcMeasurementTypeEnumTypePhase    ElectricalConnectionAcMeasurementTypeEnumType = "phase"
)

type ElectricalConnectionPhaseNameType ElectricalConnectionPhaseNameEnumType

type ElectricalConnectionPhaseNameEnumType string

const (
	ElectricalConnectionPhaseNameEnumTypeA       ElectricalConnectionPhaseNameEnumType = "a"
	ElectricalConnectionPhaseNameEnumTypeB       ElectricalConnectionPhaseNameEnumType = "b"
	ElectricalConnectionPhaseNameEnumTypeC       ElectricalConnectionPhaseNameEnumType = "c"
	ElectricalConnectionPhaseNameEnumTypeAb      ElectricalConnectionPhaseNameEnumType = "ab"
	ElectricalConnectionPhaseNameEnumTypeBc      ElectricalConnectionPhaseNameEnumType = "bc"
	ElectricalConnectionPhaseNameEnumTypeAc      ElectricalConnectionPhaseNameEnumType = "ac"
	ElectricalConnectionPhaseNameEnumTypeAbc     ElectricalConnectionPhaseNameEnumType = "abc"
	ElectricalConnectionPhaseNameEnumTypeNeutral ElectricalConnectionPhaseNameEnumType = "neutral"
	ElectricalConnectionPhaseNameEnumTypeGround  ElectricalConnectionPhaseNameEnumType = "ground"
	ElectricalConnectionPhaseNameEnumTypeNone    ElectricalConnectionPhaseNameEnumType = "none"
)

type ElectricalConnectionConnectionPointType string

const (
	ElectricalConnectionConnectionPointTypeGrid  ElectricalConnectionConnectionPointType = "grid"
	ElectricalConnectionConnectionPointTypeHome  ElectricalConnectionConnectionPointType = "home"
	ElectricalConnectionConnectionPointTypePv    ElectricalConnectionConnectionPointType = "pv"
	ElectricalConnectionConnectionPointTypeSd    ElectricalConnectionConnectionPointType = "sd"
	ElectricalConnectionConnectionPointTypeOther ElectricalConnectionConnectionPointType = "other"
)

type ElectricalConnectionParameterDescriptionDataType struct {
	ElectricalConnectionId  *ElectricalConnectionIdType                `json:"electricalConnectionId,omitempty"`
	ParameterId             *ElectricalConnectionParameterIdType       `json:"parameterId,omitempty"`
	MeasurementId           *MeasurementIdType                         `json:"measurementId,omitempty"`
	VoltageType             *ElectricalConnectionVoltageTypeType       `json:"voltageType,omitempty"`
	AcMeasuredPhases        *ElectricalConnectionPhaseNameType         `json:"acMeasuredPhases,omitempty"`
	AcMeasuredInReferenceTo *ElectricalConnectionPhaseNameType         `json:"acMeasuredInReferenceTo,omitempty"`
	AcMeasurementType       *ElectricalConnectionAcMeasurementTypeType `json:"acMeasurementType,omitempty"`
	AcMeasurementVariant    *ElectricalConnectionMeasurandVariantType  `json:"acMeasurementVariant,omitempty"`
	AcMeasuredHarmonic      *uint8                                     `json:"acMeasuredHarmonic,omitempty"`
	ScopeType               *ScopeTypeType                             `json:"scopeType,omitempty"`
	Label                   *LabelType                                 `json:"label,omitempty"`
	Description             *DescriptionType                           `json:"description,omitempty"`
}

type ElectricalConnectionParameterDescriptionDataElementsType struct {
	ElectricalConnectionId  *ElementTagType `json:"electricalConnectionId,omitempty"`
	ParameterId             *ElementTagType `json:"parameterId,omitempty"`
	MeasurementId           *ElementTagType `json:"measurementId,omitempty"`
	VoltageType             *ElementTagType `json:"voltageType,omitempty"`
	AcMeasuredPhases        *ElementTagType `json:"acMeasuredPhases,omitempty"`
	AcMeasuredInReferenceTo *ElementTagType `json:"acMeasuredInReferenceTo,omitempty"`
	AcMeasurementType       *ElementTagType `json:"acMeasurementType,omitempty"`
	AcMeasurementVariant    *ElementTagType `json:"acMeasurementVariant,omitempty"`
	AcMeasuredHarmonic      *ElementTagType `json:"acMeasuredHarmonic,omitempty"`
	ScopeType               *ElementTagType `json:"scopeType,omitempty"`
	Label                   *ElementTagType `json:"label,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
}

type ElectricalConnectionParameterDescriptionListDataType struct {
	ElectricalConnectionParameterDescriptionData []ElectricalConnectionParameterDescriptionDataType `json:"electricalConnectionParameterDescriptionData,omitempty"`
}

type ElectricalConnectionParameterDescriptionListDataSelectorsType struct {
	ElectricalConnectionId *ElectricalConnectionIdType          `json:"electricalConnectionId,omitempty"`
	ParameterId            *ElectricalConnectionParameterIdType `json:"parameterId,omitempty"`
	MeasurementId          *MeasurementIdType                   `json:"measurementId,omitempty"`
	ScopeType              *ScopeTypeType                       `json:"scopeType,omitempty"`
}

type ElectricalConnectionPermittedValueSetDataType struct {
	ElectricalConnectionId *ElectricalConnectionIdType          `json:"electricalConnectionId,omitempty"`
	ParameterId            *ElectricalConnectionParameterIdType `json:"parameterId,omitempty"`
	PermittedValueSet      []ScaledNumberSetType                `json:"permittedValueSet,omitempty"`
}

type ElectricalConnectionPermittedValueSetDataElementsType struct {
	ElectricalConnectionId *ElementTagType `json:"electricalConnectionId,omitempty"`
	ParameterId            *ElementTagType `json:"parameterId,omitempty"`
	PermittedValueSet      *ElementTagType `json:"permittedValueSet,omitempty"`
}

type ElectricalConnectionPermittedValueSetListDataType struct {
	ElectricalConnectionPermittedValueSetData []ElectricalConnectionPermittedValueSetDataType `json:"electricalConnectionPermittedValueSetData,omitempty"`
}

type ElectricalConnectionPermittedValueSetListDataSelectorsType struct {
	ElectricalConnectionId *ElectricalConnectionIdType          `json:"electricalConnectionId,omitempty"`
	ParameterId            *ElectricalConnectionParameterIdType `json:"parameterId,omitempty"`
}

type ElectricalConnectionStateDataType struct {
	ElectricalConnectionId *ElectricalConnectionIdType `json:"electricalConnectionId,omitempty"`
	Timestamp              *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	CurrentEnergyMode      *EnergyModeType             `json:"currentEnergyMode,omitempty"`
	ConsumptionTime        *string                     `json:"consumptionTime,omitempty"`
	ProductionTime         *string                     `json:"productionTime,omitempty"`
	TotalConsumptionTime   *string                     `json:"totalConsumptionTime,omitempty"`
	TotalProductionTime    *string                     `json:"totalProductionTime,omitempty"`
}

type ElectricalConnectionStateDataElementsType struct {
	ElectricalConnectionId *ElementTagType `json:"electricalConnectionId,omitempty"`
	Timestamp              *ElementTagType `json:"timestamp,omitempty"`
	CurrentEnergyMode      *ElementTagType `json:"currentEnergyMode,omitempty"`
	ConsumptionTime        *ElementTagType `json:"consumptionTime,omitempty"`
	ProductionTime         *ElementTagType `json:"productionTime,omitempty"`
	TotalConsumptionTime   *ElementTagType `json:"totalConsumptionTime,omitempty"`
	TotalProductionTime    *ElementTagType `json:"totalProductionTime,omitempty"`
}

type ElectricalConnectionStateListDataType struct {
	ElectricalConnectionStateData []ElectricalConnectionStateDataType `json:"electricalConnectionStateData,omitempty"`
}

type ElectricalConnectionStateListDataSelectorsType struct {
	ElectricalConnectionId *ElectricalConnectionIdType `json:"electricalConnectionId,omitempty"`
}

type ElectricalConnectionDescriptionDataType struct {
	ElectricalConnectionId  *ElectricalConnectionIdType          `json:"electricalConnectionId,omitempty"`
	PowerSupplyType         *ElectricalConnectionVoltageTypeType `json:"powerSupplyType,omitempty"`
	AcConnectedPhases       *uint                                `json:"acConnectedPhases,omitempty"`
	AcRmsPeriodDuration     *string                              `json:"acRmsPeriodDuration,omitempty"`
	PositiveEnergyDirection *EnergyDirectionType                 `json:"positiveEnergyDirection,omitempty"`
	ScopeType               *ScopeTypeType                       `json:"scopeType,omitempty"`
	Label                   *LabelType                           `json:"label,omitempty"`
	Description             *DescriptionType                     `json:"description,omitempty"`
}

type ElectricalConnectionDescriptionDataElementsType struct {
	ElectricalConnectionId  *ElementTagType `json:"electricalConnectionId,omitempty"`
	PowerSupplyType         *ElementTagType `json:"powerSupplyType,omitempty"`
	AcConnectedPhases       *ElementTagType `json:"acConnectedPhases,omitempty"`
	AcRmsPeriodDuration     *ElementTagType `json:"acRmsPeriodDuration,omitempty"`
	PositiveEnergyDirection *ElementTagType `json:"positiveEnergyDirection,omitempty"`
	ScopeType               *ElementTagType `json:"scopeType,omitempty"`
	Label                   *ElementTagType `json:"label,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
}

type ElectricalConnectionDescriptionListDataType struct {
	ElectricalConnectionDescriptionData []ElectricalConnectionDescriptionDataType `json:"electricalConnectionDescriptionData,omitempty"`
}

type ElectricalConnectionDescriptionListDataSelectorsType struct {
	ElectricalConnectionId *ElectricalConnectionIdType `json:"electricalConnectionId,omitempty"`
	ScopeType              *ScopeTypeType              `json:"scopeType,omitempty"`
}
