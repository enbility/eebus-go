package spine

type MeasurementIdType uint

type MeasurementTypeType MeasurementTypeEnumType

type MeasurementTypeEnumType string

const (
	MeasurementTypeEnumTypeAcceleration        MeasurementTypeEnumType = "acceleration"
	MeasurementTypeEnumTypeAngle               MeasurementTypeEnumType = "angle"
	MeasurementTypeEnumTypeAngularVelocity     MeasurementTypeEnumType = "angularVelocity"
	MeasurementTypeEnumTypeArea                MeasurementTypeEnumType = "area"
	MeasurementTypeEnumTypeAtmosphericPressure MeasurementTypeEnumType = "atmosphericPressure"
	MeasurementTypeEnumTypeCapacity            MeasurementTypeEnumType = "capacity"
	MeasurementTypeEnumTypeConcentration       MeasurementTypeEnumType = "concentration"
	MeasurementTypeEnumTypeCount               MeasurementTypeEnumType = "count"
	MeasurementTypeEnumTypeCurrent             MeasurementTypeEnumType = "current"
	MeasurementTypeEnumTypeDensity             MeasurementTypeEnumType = "density"
	MeasurementTypeEnumTypeDistance            MeasurementTypeEnumType = "distance"
	MeasurementTypeEnumTypeElectricField       MeasurementTypeEnumType = "electricField"
	MeasurementTypeEnumTypeEnergy              MeasurementTypeEnumType = "energy"
	MeasurementTypeEnumTypeForce               MeasurementTypeEnumType = "force"
	MeasurementTypeEnumTypeFrequency           MeasurementTypeEnumType = "frequency"
	MeasurementTypeEnumTypeHarmonicDistortion  MeasurementTypeEnumType = "harmonicDistortion"
	MeasurementTypeEnumTypeHeat                MeasurementTypeEnumType = "heat"
	MeasurementTypeEnumTypeHeatFlux            MeasurementTypeEnumType = "heatFlux"
	MeasurementTypeEnumTypeIlluminance         MeasurementTypeEnumType = "illuminance"
	MeasurementTypeEnumTypeImpulse             MeasurementTypeEnumType = "impulse"
	MeasurementTypeEnumTypeLevel               MeasurementTypeEnumType = "level"
	MeasurementTypeEnumTypeMagneticField       MeasurementTypeEnumType = "magneticField"
	MeasurementTypeEnumTypeMass                MeasurementTypeEnumType = "mass"
	MeasurementTypeEnumTypeMassFlow            MeasurementTypeEnumType = "massFlow"
	MeasurementTypeEnumTypeParticles           MeasurementTypeEnumType = "particles"
	MeasurementTypeEnumTypePercentage          MeasurementTypeEnumType = "percentage"
	MeasurementTypeEnumTypePower               MeasurementTypeEnumType = "power"
	MeasurementTypeEnumTypePowerFactor         MeasurementTypeEnumType = "powerFactor"
	MeasurementTypeEnumTypePressure            MeasurementTypeEnumType = "pressure"
	MeasurementTypeEnumTypeRadonActivity       MeasurementTypeEnumType = "radonActivity"
	MeasurementTypeEnumTypeRelativeHumidity    MeasurementTypeEnumType = "relativeHumidity"
	MeasurementTypeEnumTypeResistance          MeasurementTypeEnumType = "resistance"
	MeasurementTypeEnumTypeSolarRadiation      MeasurementTypeEnumType = "solarRadiation"
	MeasurementTypeEnumTypeSpeed               MeasurementTypeEnumType = "speed"
	MeasurementTypeEnumTypeTemperature         MeasurementTypeEnumType = "temperature"
	MeasurementTypeEnumTypeTime                MeasurementTypeEnumType = "time"
	MeasurementTypeEnumTypeTorque              MeasurementTypeEnumType = "torque"
	MeasurementTypeEnumTypeUnknown             MeasurementTypeEnumType = "unknown"
	MeasurementTypeEnumTypeVelocity            MeasurementTypeEnumType = "velocity"
	MeasurementTypeEnumTypeVoltage             MeasurementTypeEnumType = "voltage"
	MeasurementTypeEnumTypeVolume              MeasurementTypeEnumType = "volume"
	MeasurementTypeEnumTypeVolumetricFlow      MeasurementTypeEnumType = "volumetricFlow"
)

type MeasurementValueTypeType MeasurementValueTypeEnumType

type MeasurementValueTypeEnumType string

const (
	MeasurementValueTypeEnumTypeValue             MeasurementValueTypeEnumType = "value"
	MeasurementValueTypeEnumTypeAverageValue      MeasurementValueTypeEnumType = "averageValue"
	MeasurementValueTypeEnumTypeMinvValue         MeasurementValueTypeEnumType = "minValue"
	MeasurementValueTypeEnumTypeMaxvVlue          MeasurementValueTypeEnumType = "maxValue"
	MeasurementValueTypeEnumTypeStandardDeviation MeasurementValueTypeEnumType = "standardDeviation"
)

type MeasurementValueSourceType MeasurementValueSourceEnumType

type MeasurementValueSourceEnumType string

const (
	MeasurementValueSourceEnumTypeMeasuredValue   MeasurementValueSourceEnumType = "measuredValue"
	MeasurementValueSourceEnumTypeCalculatedValue MeasurementValueSourceEnumType = "calculatedValue"
	MeasurementValueSourceEnumTypeEmpiricalValue  MeasurementValueSourceEnumType = "empiricalValue"
)

type MeasurementValueTendencyType MeasurementValueTendencyEnumType

type MeasurementValueTendencyEnumType string

const (
	MeasurementValueTendencyEnumTypeRising  MeasurementValueTendencyEnumType = "rising"
	MeasurementValueTendencyEnumTypeStable  MeasurementValueTendencyEnumType = "stable"
	MeasurementValueTendencyEnumTypeFalling MeasurementValueTendencyEnumType = "falling"
)

type MeasurementValueStateType MeasurementValueStateEnumType

type MeasurementValueStateEnumType string

const (
	MeasurementValueStateEnumTypeNormal     MeasurementValueStateEnumType = "normal"
	MeasurementValueStateEnumTypeOutofrange MeasurementValueStateEnumType = "outOfRange"
	MeasurementValueStateEnumTypeError      MeasurementValueStateEnumType = "error"
)

type MeasurementDataType struct {
	MeasurementId    *MeasurementIdType            `json:"measurementId,omitempty"`
	ValueType        *MeasurementValueTypeType     `json:"valueType,omitempty"`
	Timestamp        *string                       `json:"timestamp,omitempty"`
	Value            *ScaledNumberType             `json:"value,omitempty"`
	EvaluationPeriod *TimePeriodType               `json:"evaluationPeriod,omitempty"`
	ValueSource      *MeasurementValueSourceType   `json:"valueSource,omitempty"`
	ValueTendency    *MeasurementValueTendencyType `json:"valueTendency,omitempty"`
	ValueState       *MeasurementValueStateType    `json:"valueState,omitempty"`
}

type MeasurementDataElementsType struct {
	MeasurementId    *ElementTagType `json:"measurementId,omitempty"`
	ValueType        *ElementTagType `json:"valueType,omitempty"`
	Timestamp        *ElementTagType `json:"timestamp,omitempty"`
	Value            *ElementTagType `json:"value,omitempty"`
	EvaluationPeriod *ElementTagType `json:"evaluationPeriod,omitempty"`
	ValueSource      *ElementTagType `json:"valueSource,omitempty"`
	ValueTendency    *ElementTagType `json:"valueTendency,omitempty"`
	ValueState       *ElementTagType `json:"valueState,omitempty"`
}

type MeasurementListDataType struct {
	MeasurementData []MeasurementDataType `json:"measurementData,omitempty"`
}

type MeasurementListDataSelectorsType struct {
	MeasurementId     *MeasurementIdType        `json:"measurementId,omitempty"`
	ValueType         *MeasurementValueTypeType `json:"valueType,omitempty"`
	TimestampInterval *TimestampIntervalType    `json:"timestampInterval,omitempty"`
}

type MeasurementConstraintsDataType struct {
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
	ValueRangeMin *ScaledNumberType  `json:"valueRangeMin,omitempty"`
	ValueRangeMax *ScaledNumberType  `json:"valueRangeMax,omitempty"`
	ValueStepSize *ScaledNumberType  `json:"valueStepSize,omitempty"`
}

type MeasurementConstraintsDataElementsType struct {
	MeasurementId *ElementTagType           `json:"measurementId,omitempty"`
	ValueRangeMin *ScaledNumberElementsType `json:"valueRangeMin,omitempty"`
	ValueRangeMax *ScaledNumberElementsType `json:"valueRangeMax,omitempty"`
	ValueStepSize *ScaledNumberElementsType `json:"valueStepSize,omitempty"`
}

type MeasurementConstraintsListDataType struct {
	MeasurementConstraintsData []MeasurementConstraintsDataType `json:"measurementConstraintsData,omitempty"`
}

type MeasurementConstraintsListDataSelectorsType struct {
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
}

type MeasurementDescriptionDataType struct {
	MeasurementId    *MeasurementIdType     `json:"measurementId,omitempty"`
	MeasurementType  *MeasurementTypeType   `json:"measurementType,omitempty"`
	CommodityType    *CommodityTypeType     `json:"commodityType,omitempty"`
	Unit             *UnitOfMeasurementType `json:"unit,omitempty"`
	CalibrationValue *ScaledNumberType      `json:"calibrationValue,omitempty"`
	ScopeType        *ScopeTypeType         `json:"scopeType,omitempty"`
	Label            *LabelType             `json:"label,omitempty"`
	Description      *DescriptionType       `json:"description,omitempty"`
}

type MeasurementDescriptionDataElementsType struct {
	MeasurementId    *ElementTagType           `json:"measurementId,omitempty"`
	MeasurementType  *ElementTagType           `json:"measurementType,omitempty"`
	CommodityType    *ElementTagType           `json:"commodityType,omitempty"`
	Unit             *ElementTagType           `json:"unit,omitempty"`
	CalibrationValue *ScaledNumberElementsType `json:"calibrationValue,omitempty"`
	ScopeType        *ElementTagType           `json:"scopeType,omitempty"`
	Label            *ElementTagType           `json:"label,omitempty"`
	Description      *ElementTagType           `json:"description,omitempty"`
}

type MeasurementDescriptionListDataType struct {
	MeasurementDescriptionData []MeasurementDescriptionDataType `json:"measurementDescriptionData,omitempty"`
}

type MeasurementDescriptionListDataSelectorsType struct {
	MeasurementId   *MeasurementIdType   `json:"measurementId,omitempty"`
	MeasurementType *MeasurementTypeType `json:"measurementType,omitempty"`
	CommodityType   *CommodityTypeType   `json:"commodityType,omitempty"`
	ScopeType       *ScopeTypeType       `json:"scopeType,omitempty"`
}

type MeasurementThresholdRelationDataType struct {
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
	ThresholdId   []ThresholdIdType  `json:"thresholdId,omitempty"`
}

type MeasurementThresholdRelationDataElementsType struct {
	MeasurementId *ElementTagType `json:"measurementId,omitempty"`
	ThresholdId   *ElementTagType `json:"thresholdId,omitempty"`
}

type MeasurementThresholdRelationListDataType struct {
	MeasurementThresholdRelationData []MeasurementThresholdRelationDataType `json:"measurementThresholdRelationData,omitempty"`
}

type MeasurementThresholdRelationListDataSelectorsType struct {
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
	ThresholdId   *ThresholdIdType   `json:"thresholdId,omitempty"`
}
