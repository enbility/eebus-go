package model

type DirectControlActivityStateType string

const (
	DirectControlActivityStateTypeRunning  AlarmTypeType = "running"
	DirectControlActivityStateTypePaused   AlarmTypeType = "paused"
	DirectControlActivityStateTypeInactive AlarmTypeType = "inactive"
)

type DirectControlActivityDataType struct {
	Timestamp                 *AbsoluteOrRelativeTimeType     `json:"timestamp,omitempty"`
	ActivityState             *DirectControlActivityStateType `json:"activityState,omitempty"`
	IsActivityStateChangeable *bool                           `json:"isActivityStateChangeable,omitempty"`
	EnergyMode                *EnergyModeType                 `json:"energyMode,omitempty"`
	IsEnergyModeChangeable    *bool                           `json:"isEnergyModeChangeable,omitempty"`
	Power                     *ScaledNumberType               `json:"power,omitempty"`
	IsPowerChangeable         *bool                           `json:"isPowerChangeable,omitempty"`
	Energy                    *ScaledNumberType               `json:"energy,omitempty"`
	IsEnergyChangeable        *bool                           `json:"isEnergyChangeable,omitempty"`
	Sequence_id               *PowerSequenceIdType            `json:"sequence_id,omitempty"`
}

type DirectControlActivityDataElementsType struct {
	Timestamp                 *ElementTagType           `json:"timestamp,omitempty"`
	ActivityState             *ElementTagType           `json:"activityState,omitempty"`
	IsActivityStateChangeable *ElementTagType           `json:"isActivityStateChangeable,omitempty"`
	EnergyMode                *ElementTagType           `json:"energyMode,omitempty"`
	IsEnergyModeChangeable    *ElementTagType           `json:"isEnergyModeChangeable,omitempty"`
	Power                     *ScaledNumberElementsType `json:"power,omitempty"`
	IsPowerChangeable         *ElementTagType           `json:"isPowerChangeable,omitempty"`
	Energy                    *ScaledNumberElementsType `json:"energy,omitempty"`
	IsEnergyChangeable        *ElementTagType           `json:"isEnergyChangeable,omitempty"`
	Sequence_id               *ElementTagType           `json:"sequence_id,omitempty"`
}

type DirectControlActivityListDataType struct {
	DirectControlActivityDataElements []DirectControlActivityDataType `json:"directControlActivityDataElements,omitempty"`
}

type DirectControlActivityListDataSelectorsType struct {
	TimestampInterval *TimestampIntervalType `json:"timestampInterval,omitempty"`
}

type DirectControlDescriptionDataType struct {
	PositiveEnergyDirection *EnergyDirectionType   `json:"positiveEnergyDirection,omitempty"`
	PowerUnit               *UnitOfMeasurementType `json:"powerUnit,omitempty"`
	EnergyUnit              *UnitOfMeasurementType `json:"energyUnit,omitempty"`
}

type DirectControlDescriptionDataElementsType struct {
	PositiveEnergyDirection *ElementTagType `json:"positiveEnergyDirection,omitempty"`
	PowerUnit               *ElementTagType `json:"powerUnit,omitempty"`
	EnergyUnit              *ElementTagType `json:"energyUnit,omitempty"`
}
