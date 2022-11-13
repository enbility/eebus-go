package model

type ElectricalConnectionIdType uint

type ElectricalConnectionParameterIdType uint

type ElectricalConnectionMeasurandVariantType string

const (
	ElectricalConnectionMeasurandVariantTypeAmplitude     ElectricalConnectionMeasurandVariantType = "amplitude"
	ElectricalConnectionMeasurandVariantTypeRms           ElectricalConnectionMeasurandVariantType = "rms"
	ElectricalConnectionMeasurandVariantTypeInstantaneous ElectricalConnectionMeasurandVariantType = "instantaneous"
	ElectricalConnectionMeasurandVariantTypeAngle         ElectricalConnectionMeasurandVariantType = "angle"
	ElectricalConnectionMeasurandVariantTypeCosphi        ElectricalConnectionMeasurandVariantType = "cosPhi"
)

type ElectricalConnectionVoltageTypeType string

const (
	ElectricalConnectionVoltageTypeTypeAc ElectricalConnectionVoltageTypeType = "ac"
	ElectricalConnectionVoltageTypeTypeDc ElectricalConnectionVoltageTypeType = "dc"
)

type ElectricalConnectionAcMeasurementTypeType string

const (
	ElectricalConnectionAcMeasurementTypeTypeReal     ElectricalConnectionAcMeasurementTypeType = "real"
	ElectricalConnectionAcMeasurementTypeTypeReactive ElectricalConnectionAcMeasurementTypeType = "reactive"
	ElectricalConnectionAcMeasurementTypeTypeApparent ElectricalConnectionAcMeasurementTypeType = "apparent"
	ElectricalConnectionAcMeasurementTypeTypePhase    ElectricalConnectionAcMeasurementTypeType = "phase"
)

type ElectricalConnectionPhaseNameType string

const (
	ElectricalConnectionPhaseNameTypeA       ElectricalConnectionPhaseNameType = "a"
	ElectricalConnectionPhaseNameTypeB       ElectricalConnectionPhaseNameType = "b"
	ElectricalConnectionPhaseNameTypeC       ElectricalConnectionPhaseNameType = "c"
	ElectricalConnectionPhaseNameTypeAb      ElectricalConnectionPhaseNameType = "ab"
	ElectricalConnectionPhaseNameTypeBc      ElectricalConnectionPhaseNameType = "bc"
	ElectricalConnectionPhaseNameTypeAc      ElectricalConnectionPhaseNameType = "ac"
	ElectricalConnectionPhaseNameTypeAbc     ElectricalConnectionPhaseNameType = "abc"
	ElectricalConnectionPhaseNameTypeNeutral ElectricalConnectionPhaseNameType = "neutral"
	ElectricalConnectionPhaseNameTypeGround  ElectricalConnectionPhaseNameType = "ground"
	ElectricalConnectionPhaseNameTypeNone    ElectricalConnectionPhaseNameType = "none"
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
	ElectricalConnectionId  *ElectricalConnectionIdType                `json:"electricalConnectionId,omitempty" eebus:"key"`
	ParameterId             *ElectricalConnectionParameterIdType       `json:"parameterId,omitempty" eebus:"key"`
	MeasurementId           *MeasurementIdType                         `json:"measurementId,omitempty" eebus:"key"`
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
	ElectricalConnectionId *ElectricalConnectionIdType          `json:"electricalConnectionId,omitempty" eebus:"key"`
	ParameterId            *ElectricalConnectionParameterIdType `json:"parameterId,omitempty" eebus:"key"`
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
	ElectricalConnectionId *ElectricalConnectionIdType `json:"electricalConnectionId,omitempty" eebus:"key"`
	Timestamp              *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	CurrentEnergyMode      *EnergyModeType             `json:"currentEnergyMode,omitempty"`
	ConsumptionTime        *DurationType               `json:"consumptionTime,omitempty"`
	ProductionTime         *DurationType               `json:"productionTime,omitempty"`
	TotalConsumptionTime   *DurationType               `json:"totalConsumptionTime,omitempty"`
	TotalProductionTime    *DurationType               `json:"totalProductionTime,omitempty"`
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
	ElectricalConnectionId  *ElectricalConnectionIdType          `json:"electricalConnectionId,omitempty" eebus:"key"`
	PowerSupplyType         *ElectricalConnectionVoltageTypeType `json:"powerSupplyType,omitempty"`
	AcConnectedPhases       *uint                                `json:"acConnectedPhases,omitempty"`
	AcRmsPeriodDuration     *DurationType                        `json:"acRmsPeriodDuration,omitempty"`
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
