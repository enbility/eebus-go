package model

type SetpointIdType uint

type SetpointTypeType string

const (
	SetpointTypeTypeValueAbsolute SetpointTypeType = "valueAbsolute"
	SetpointTypeTypeValueRelative SetpointTypeType = "valueRelative"
)

type SetpointDataType struct {
	SetpointId               *SetpointIdType   `json:"setpointId,omitempty"`
	Value                    *ScaledNumberType `json:"value,omitempty"`
	ValueMin                 *ScaledNumberType `json:"valueMin,omitempty"`
	ValueMax                 *ScaledNumberType `json:"valueMax,omitempty"`
	ValueToleranceAbsolute   *ScaledNumberType `json:"valueToleranceAbsolute,omitempty"`
	ValueTolerancePercentage *ScaledNumberType `json:"valueTolerancePercentage,omitempty"`
}

type SetpointDataElementsType struct {
	SetpointId               *ElementTagType           `json:"setpointId,omitempty"`
	Value                    *ScaledNumberElementsType `json:"value,omitempty"`
	ValueMin                 *ScaledNumberElementsType `json:"valueMin,omitempty"`
	ValueMax                 *ScaledNumberElementsType `json:"valueMax,omitempty"`
	ValueToleranceAbsolute   *ScaledNumberElementsType `json:"valueToleranceAbsolute,omitempty"`
	ValueTolerancePercentage *ScaledNumberElementsType `json:"valueTolerancePercentage,omitempty"`
}

type SetpointListDataType struct {
	SetpointData []SetpointDataType `json:"setpointData,omitempty"`
}

type SetpointListDataSelectorsType struct {
	SetpointId *SetpointIdType `json:"setpointId,omitempty"`
}

type SetpointConstraintsDataType struct {
	SetpointId       *SetpointIdType   `json:"setpointId,omitempty"`
	SetpointRangeMin *ScaledNumberType `json:"setpointRangeMin,omitempty"`
	SetpointRangeMax *ScaledNumberType `json:"setpointRangeMax,omitempty"`
	SetpointStepSize *ScaledNumberType `json:"setpointStepSize,omitempty"`
}

type SetpointConstraintsDataElementsType struct {
	SetpointId       *ElementTagType           `json:"setpointId,omitempty"`
	SetpointRangeMin *ScaledNumberElementsType `json:"setpointRangeMin,omitempty"`
	SetpointRangeMax *ScaledNumberElementsType `json:"setpointRangeMax,omitempty"`
	SetpointStepSize *ScaledNumberElementsType `json:"setpointStepSize,omitempty"`
}

type SetpointConstraintsListDataType struct {
	SetpointConstraintsData []SetpointConstraintsDataType `json:"setpointConstraintsData,omitempty"`
}

type SetpointConstraintsListDataSelectorsType struct {
	SetpointId *SetpointIdType `json:"setpointId,omitempty"`
}

type SetpointDescriptionDataType struct {
	SetpointId    *SetpointIdType   `json:"setpointId,omitempty"`
	MeasurementId *SetpointIdType   `json:"measurementId,omitempty"`
	TimeTableId   *SetpointIdType   `json:"timeTableId,omitempty"`
	SetpointType  *SetpointIdType   `json:"setpointType,omitempty"`
	Unit          *ScaledNumberType `json:"unit,omitempty"`
	ScopeType     *ScaledNumberType `json:"scopeType,omitempty"`
	Label         *ScaledNumberType `json:"label,omitempty"`
	Description   *ScaledNumberType `json:"description,omitempty"`
}

type SetpointDescriptionDataElementsType struct {
	SetpointId    *ElementTagType `json:"setpointId,omitempty"`
	MeasurementId *ElementTagType `json:"measurementId,omitempty"`
	TimeTableId   *ElementTagType `json:"timeTableId,omitempty"`
	SetpointType  *ElementTagType `json:"setpointType,omitempty"`
	Unit          *ElementTagType `json:"unit,omitempty"`
	ScopeType     *ElementTagType `json:"scopeType,omitempty"`
	Label         *ElementTagType `json:"label,omitempty"`
	Description   *ElementTagType `json:"description,omitempty"`
}

type SetpointDescriptionListDataType struct {
	SetpointDescriptionData []SetpointDescriptionDataType `json:"setpointDescriptionData,omitempty"`
}

type SetpointDescriptionListDataSelectorsType struct {
	SetpointId    *SetpointIdType   `json:"setpointId,omitempty"`
	MeasurementId *SetpointIdType   `json:"measurementId,omitempty"`
	TimeTableId   *SetpointIdType   `json:"timeTableId,omitempty"`
	SetpointType  *SetpointIdType   `json:"setpointType,omitempty"`
	ScopeType     *ScaledNumberType `json:"scopeType,omitempty"`
}
