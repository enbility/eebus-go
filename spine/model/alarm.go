package model

type AlarmIdType uint

type AlarmTypeType string

const (
	AlarmTypeTypeAlarmCancelled AlarmTypeType = "alarmCancelled"
	AlarmTypeTypeUnderThreshold AlarmTypeType = "underThreshold"
	AlarmTypeTypeOverThreshold  AlarmTypeType = "overThreshold"
)

type AlarmDataType struct {
	AlarmId          *AlarmIdType                `json:"alarmId,omitempty"`
	ThresholdId      *ThresholdIdType            `json:"thresholdId,omitempty"`
	Timestamp        *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	AlarmType        *AlarmTypeType              `json:"alarmType,omitempty"`
	MeasuredValue    *ScaledNumberType           `json:"measuredValue,omitempty"`
	EvaluationPeriod *TimePeriodType             `json:"evaluationPeriod,omitempty"`
	ScopeType        *ScopeTypeType              `json:"scopeType,omitempty"`
	Label            *LabelType                  `json:"label,omitempty"`
	Description      *DescriptionType            `json:"description,omitempty"`
}

type AlarmDataElementsType struct {
	AlarmId          *ElementTagType           `json:"alarmId,omitempty"`
	ThresholdId      *ElementTagType           `json:"thresholdId,omitempty"`
	Timestamp        *ElementTagType           `json:"timestamp,omitempty"`
	AlarmType        *ElementTagType           `json:"alarmType,omitempty"`
	MeasuredValue    *ScaledNumberElementsType `json:"measuredValue,omitempty"`
	EvaluationPeriod *TimePeriodElementsType   `json:"evaluationPeriod,omitempty"`
	ScopeType        *ElementTagType           `json:"scopeType,omitempty"`
	Label            *ElementTagType           `json:"label,omitempty"`
	Description      *ElementTagType           `json:"description,omitempty"`
}

type AlarmListDataType struct {
	AlarmListData *AlarmListDataType `json:"alarmListData,omitempty"`
}

type AlarmListDataSelectorsType struct {
	AlarmId   *AlarmIdType   `json:"alarmId,omitempty"`
	ScopeType *ScopeTypeType `json:"scopeType,omitempty"`
}
