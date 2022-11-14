package model

type TimeSeriesIdType uint

type TimeSeriesSlotIdType uint

type TimeSeriesSlotCountType TimeSeriesSlotIdType

type TimeSeriesTypeType string

const (
	TimeSeriesTypeTypePlan         TimeSeriesTypeType = "plan"
	TimeSeriesTypeTypeSingleDemand TimeSeriesTypeType = "singleDemand"
	TimeSeriesTypeTypeConstraints  TimeSeriesTypeType = "constraints"
)

type TimeSeriesSlotType struct {
	TimeSeriesSlotId      *TimeSeriesSlotIdType        `json:"timeSeriesSlotId,omitempty"`
	TimePeriod            *TimePeriodType              `json:"timePeriod,omitempty"`
	Duration              *DurationType                `json:"duration,omitempty"`
	RecurrenceInformation *AbsoluteOrRecurringTimeType `json:"recurrenceInformation,omitempty"`
	Value                 *ScaledNumberType            `json:"value,omitempty"`
	MinValue              *ScaledNumberType            `json:"minValue,omitempty"`
	MaxValue              *ScaledNumberType            `json:"maxValue,omitempty"`
}

type TimeSeriesSlotElementsType struct {
	TimeSeriesSlotId      *ElementTagType                      `json:"timeSeriesSlotId,omitempty"`
	TimePeriod            *ElementTagType                      `json:"timePeriod,omitempty"`
	Duration              *ElementTagType                      `json:"duration,omitempty"`
	RecurrenceInformation *AbsoluteOrRecurringTimeElementsType `json:"recurrenceInformation,omitempty"`
	Value                 *ScaledNumberElementsType            `json:"value,omitempty"`
	MinValue              *ScaledNumberElementsType            `json:"minValue,omitempty"`
	MaxValue              *ScaledNumberElementsType            `json:"maxValue,omitempty"`
}

type TimeSeriesDataType struct {
	TimeSeriesId   *TimeSeriesIdType    `json:"timeSeriesId,omitempty" eebus:"key"`
	TimePeriod     *TimePeriodType      `json:"timePeriod,omitempty"`
	TimeSeriesSlot []TimeSeriesSlotType `json:"timeSeriesSlot"`
}

type TimeSeriesDataElementsType struct {
	TimeSeriesId   *ElementTagType             `json:"timeSeriesId,omitempty"`
	TimePeriod     *TimePeriodElementsType     `json:"timePeriod,omitempty"`
	TimeSeriesSlot *TimeSeriesSlotElementsType `json:"timeSeriesSlot"`
}

type TimeSeriesListDataType struct {
	TimeSeriesData []TimeSeriesDataType `json:"timeSeriesData,omitempty"`
}

type TimeSeriesListDataSelectorsType struct {
	TimeSeriesId     *TimeSeriesIdType     `json:"timeSeriesId,omitempty"`
	TimeSeriesSlotId *TimeSeriesSlotIdType `json:"timeSeriesSlotId,omitempty"`
}

type TimeSeriesDescriptionDataType struct {
	TimeSeriesId        *TimeSeriesIdType      `json:"timeSeriesId,omitempty" eebus:"key"`
	TimeSeriesType      *TimeSeriesTypeType    `json:"timeSeriesType,omitempty"`
	TimeSeriesWriteable *bool                  `json:"timeSeriesWriteable,omitempty"`
	UpdateRequired      *bool                  `json:"updateRequired,omitempty"`
	MeasurementId       *MeasurementIdType     `json:"measurementId,omitempty" eebus:"key"`
	Currency            *CurrencyType          `json:"currency,omitempty"`
	Unit                *UnitOfMeasurementType `json:"unit,omitempty"`
	Label               *LabelType             `json:"label,omitempty"`
	Description         *DescriptionType       `json:"description,omitempty"`
	ScopeType           *ScopeTypeType         `json:"scopeType,omitempty"`
}

type TimeSeriesDescriptionDataElementsType struct {
	TimeSeriesId        *ElementTagType `json:"timeSeriesId,omitempty"`
	TimeSeriesType      *ElementTagType `json:"timeSeriesType,omitempty"`
	TimeSeriesWriteable *ElementTagType `json:"timeSeriesWriteable,omitempty"`
	UpdateRequired      *ElementTagType `json:"updateRequired,omitempty"`
	MeasurementId       *ElementTagType `json:"measurementId,omitempty"`
	Currency            *ElementTagType `json:"currency,omitempty"`
	Unit                *ElementTagType `json:"unit,omitempty"`
	Label               *ElementTagType `json:"label,omitempty"`
	Description         *ElementTagType `json:"description,omitempty"`
	ScopeType           *ElementTagType `json:"scopeType,omitempty"`
}

type TimeSeriesDescriptionListDataType struct {
	TimeSeriesDescriptionData []TimeSeriesDescriptionDataType `json:"timeSeriesDescriptionData,omitempty"`
}

type TimeSeriesDescriptionListDataSelectorsType struct {
	TimeSeriesId   *TimeSeriesIdType   `json:"timeSeriesId,omitempty"`
	TimeSeriesType *TimeSeriesTypeType `json:"timeSeriesType,omitempty"`
	MeasurementId  *MeasurementIdType  `json:"measurementId,omitempty"`
	ScopeType      *ScopeTypeType      `json:"scopeType,omitempty"`
}

type TimeSeriesConstraintsDataType struct {
	TimeSeriesId                *TimeSeriesIdType           `json:"timeSeriesId,omitempty" eebus:"key"`
	SlotCountMin                *TimeSeriesSlotCountType    `json:"slotCountMin,omitempty"`
	SlotCountMax                *TimeSeriesSlotCountType    `json:"slotCountMax,omitempty"`
	SlotDurationMin             *DurationType               `json:"slotDurationMin,omitempty"`
	SlotDurationMax             *DurationType               `json:"slotDurationMax,omitempty"`
	SlotDurationStepSize        *DurationType               `json:"slotDurationStepSize,omitempty"`
	EarliestTimeSeriesStartTime *AbsoluteOrRelativeTimeType `json:"earliestTimeSeriesStartTime,omitempty"`
	LatestTimeSeriesEndTime     *AbsoluteOrRelativeTimeType `json:"latestTimeSeriesEndTime,omitempty"`
	SlotValueMin                *ScaledNumberType           `json:"slotValueMin,omitempty"`
	SlotValueMax                *ScaledNumberType           `json:"slotValueMax,omitempty"`
	SlotValueStepSize           *ScaledNumberType           `json:"slotValueStepSize,omitempty"`
}

type TimeSeriesConstraintsDataElementsType struct {
	TimeSeriesId                *ElementTagType `json:"timeSeriesId,omitempty"`
	SlotCountMin                *ElementTagType `json:"slotCountMin,omitempty"`
	SlotCountMax                *ElementTagType `json:"slotCountMax,omitempty"`
	SlotDurationMin             *ElementTagType `json:"slotDurationMin,omitempty"`
	SlotDurationMax             *ElementTagType `json:"slotDurationMax,omitempty"`
	SlotDurationStepSize        *ElementTagType `json:"slotDurationStepSize,omitempty"`
	EarliestTimeSeriesStartTime *ElementTagType `json:"earliestTimeSeriesStartTime,omitempty"`
	LatestTimeSeriesEndTime     *ElementTagType `json:"latestTimeSeriesEndTime,omitempty"`
	SlotValueMin                *ElementTagType `json:"slotValueMin,omitempty"`
	SlotValueMax                *ElementTagType `json:"slotValueMax,omitempty"`
	SlotValueStepSize           *ElementTagType `json:"slotValueStepSize,omitempty"`
}

type TimeSeriesConstraintsListDataType struct {
	TimeSeriesConstraintsData []TimeSeriesConstraintsDataType `json:"timeSeriesConstraintsData,omitempty"`
}

type TimeSeriesConstraintsListDataSelectorsType struct {
	TimeSeriesId *TimeSeriesIdType `json:"timeSeriesId,omitempty"`
}
