package model

type MeasurementIdType uint

type MeasurementTypeType string

const (
	MeasurementTypeTypeAcceleration        MeasurementTypeType = "acceleration"
	MeasurementTypeTypeAngle               MeasurementTypeType = "angle"
	MeasurementTypeTypeAngularVelocity     MeasurementTypeType = "angularVelocity"
	MeasurementTypeTypeArea                MeasurementTypeType = "area"
	MeasurementTypeTypeAtmosphericPressure MeasurementTypeType = "atmosphericPressure"
	MeasurementTypeTypeCapacity            MeasurementTypeType = "capacity"
	MeasurementTypeTypeConcentration       MeasurementTypeType = "concentration"
	MeasurementTypeTypeCount               MeasurementTypeType = "count"
	MeasurementTypeTypeCurrent             MeasurementTypeType = "current"
	MeasurementTypeTypeDensity             MeasurementTypeType = "density"
	MeasurementTypeTypeDistance            MeasurementTypeType = "distance"
	MeasurementTypeTypeElectricField       MeasurementTypeType = "electricField"
	MeasurementTypeTypeEnergy              MeasurementTypeType = "energy"
	MeasurementTypeTypeForce               MeasurementTypeType = "force"
	MeasurementTypeTypeFrequency           MeasurementTypeType = "frequency"
	MeasurementTypeTypeHarmonicDistortion  MeasurementTypeType = "harmonicDistortion"
	MeasurementTypeTypeHeat                MeasurementTypeType = "heat"
	MeasurementTypeTypeHeatFlux            MeasurementTypeType = "heatFlux"
	MeasurementTypeTypeIlluminance         MeasurementTypeType = "illuminance"
	MeasurementTypeTypeImpulse             MeasurementTypeType = "impulse"
	MeasurementTypeTypeLevel               MeasurementTypeType = "level"
	MeasurementTypeTypeMagneticField       MeasurementTypeType = "magneticField"
	MeasurementTypeTypeMass                MeasurementTypeType = "mass"
	MeasurementTypeTypeMassFlow            MeasurementTypeType = "massFlow"
	MeasurementTypeTypeParticles           MeasurementTypeType = "particles"
	MeasurementTypeTypePercentage          MeasurementTypeType = "percentage"
	MeasurementTypeTypePower               MeasurementTypeType = "power"
	MeasurementTypeTypePowerFactor         MeasurementTypeType = "powerFactor"
	MeasurementTypeTypePressure            MeasurementTypeType = "pressure"
	MeasurementTypeTypeRadonActivity       MeasurementTypeType = "radonActivity"
	MeasurementTypeTypeRelativeHumidity    MeasurementTypeType = "relativeHumidity"
	MeasurementTypeTypeResistance          MeasurementTypeType = "resistance"
	MeasurementTypeTypeSolarRadiation      MeasurementTypeType = "solarRadiation"
	MeasurementTypeTypeSpeed               MeasurementTypeType = "speed"
	MeasurementTypeTypeTemperature         MeasurementTypeType = "temperature"
	MeasurementTypeTypeTime                MeasurementTypeType = "time"
	MeasurementTypeTypeTorque              MeasurementTypeType = "torque"
	MeasurementTypeTypeUnknown             MeasurementTypeType = "unknown"
	MeasurementTypeTypeVelocity            MeasurementTypeType = "velocity"
	MeasurementTypeTypeVoltage             MeasurementTypeType = "voltage"
	MeasurementTypeTypeVolume              MeasurementTypeType = "volume"
	MeasurementTypeTypeVolumetricFlow      MeasurementTypeType = "volumetricFlow"
)

type MeasurementValueTypeType string

const (
	MeasurementValueTypeTypeValue             MeasurementValueTypeType = "value"
	MeasurementValueTypeTypeAverageValue      MeasurementValueTypeType = "averageValue"
	MeasurementValueTypeTypeMinvValue         MeasurementValueTypeType = "minValue"
	MeasurementValueTypeTypeMaxvVlue          MeasurementValueTypeType = "maxValue"
	MeasurementValueTypeTypeStandardDeviation MeasurementValueTypeType = "standardDeviation"
)

type MeasurementValueSourceType string

const (
	MeasurementValueSourceTypeMeasuredValue   MeasurementValueSourceType = "measuredValue"
	MeasurementValueSourceTypeCalculatedValue MeasurementValueSourceType = "calculatedValue"
	MeasurementValueSourceTypeEmpiricalValue  MeasurementValueSourceType = "empiricalValue"
)

type MeasurementValueTendencyType string

const (
	MeasurementValueTendencyTypeRising  MeasurementValueTendencyType = "rising"
	MeasurementValueTendencyTypeStable  MeasurementValueTendencyType = "stable"
	MeasurementValueTendencyTypeFalling MeasurementValueTendencyType = "falling"
)

type MeasurementValueStateType string

const (
	MeasurementValueStateTypeNormal     MeasurementValueStateType = "normal"
	MeasurementValueStateTypeOutofrange MeasurementValueStateType = "outOfRange"
	MeasurementValueStateTypeError      MeasurementValueStateType = "error"
)

type MeasurementDataType struct {
	MeasurementId    *MeasurementIdType            `json:"measurementId,omitempty" eebus:"key"`
	ValueType        *MeasurementValueTypeType     `json:"valueType,omitempty"`
	Timestamp        *AbsoluteOrRelativeTimeType   `json:"timestamp,omitempty"`
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
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty" eebus:"key"`
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
	MeasurementId    *MeasurementIdType     `json:"measurementId,omitempty" eebus:"key"`
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
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty" eebus:"key"`
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
