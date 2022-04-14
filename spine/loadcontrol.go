package spine

type LoadControlEventIdType string

type LoadControlEventActionType LoadControlEventActionEnumType

type LoadControlEventActionEnumType string

const (
	LoadControlEventActionEnumTypePause     LoadControlEventActionEnumType = "pause"
	LoadControlEventActionEnumTypeResume    LoadControlEventActionEnumType = "resume"
	LoadControlEventActionEnumTypeReduce    LoadControlEventActionEnumType = "reduce"
	LoadControlEventActionEnumTypeIncrease  LoadControlEventActionEnumType = "increase"
	LoadControlEventActionEnumTypeEmergency LoadControlEventActionEnumType = "emergency"
	LoadControlEventActionEnumTypeNormal    LoadControlEventActionEnumType = "normal"
)

type LoadControlEventStateType LoadControlEventStateEnumType

type LoadControlEventStateEnumType string

const (
	LoadControlEventStateEnumTypeEventaccepted  LoadControlEventStateEnumType = "eventAccepted"
	LoadControlEventStateEnumTypeEventstarted   LoadControlEventStateEnumType = "eventStarted"
	LoadControlEventStateEnumTypeEventstopped   LoadControlEventStateEnumType = "eventStopped"
	LoadControlEventStateEnumTypeEventrejected  LoadControlEventStateEnumType = "eventRejected"
	LoadControlEventStateEnumTypeEventcancelled LoadControlEventStateEnumType = "eventCancelled"
	LoadControlEventStateEnumTypeEventerror     LoadControlEventStateEnumType = "eventError"
)

type LoadControlLimitIdType uint

type LoadControlLimitTypeType LoadControlLimitTypeEnumType

type LoadControlLimitTypeEnumType string

const (
	LoadControlLimitTypeEnumTypeMinvaluelimit LoadControlLimitTypeEnumType = "minValueLimit"
	LoadControlLimitTypeEnumTypeMaxvaluelimit LoadControlLimitTypeEnumType = "maxValueLimit"
)

type LoadControlCategoryType LoadControlCategoryEnumType

type LoadControlCategoryEnumType string

const (
	LoadControlCategoryEnumTypeObligation     LoadControlCategoryEnumType = "obligation"
	LoadControlCategoryEnumTypeRecommendation LoadControlCategoryEnumType = "recommendation"
	LoadControlCategoryEnumTypeOptimization   LoadControlCategoryEnumType = "optimization"
)

type LoadControlNodeDataType struct {
	IsNodeRemoteControllable *bool `json:"isNodeRemoteControllable,omitempty"`
}

type LoadControlNodeDataElementsType struct {
	IsNodeRemoteControllable *ElementTagType `json:"isNodeRemoteControllable,omitempty"`
}

type LoadControlEventDataType struct {
	Timestamp          *string                     `json:"timestamp,omitempty"`
	EventId            *LoadControlEventIdType     `json:"eventId,omitempty,omitempty"`
	EventActionConsume *LoadControlEventActionType `json:"eventActionConsume,omitempty"`
	EventActionProduce *LoadControlEventActionType `json:"eventActionProduce,omitempty"`
	TimePeriod         *TimePeriodType             `json:"timePeriod,omitempty"`
}

type LoadControlEventDataElementsType struct {
	Timestamp          *ElementTagType         `json:"timestamp,omitempty"`
	EventId            *ElementTagType         `json:"eventId,omitempty,omitempty"`
	EventActionConsume *ElementTagType         `json:"eventActionConsume,omitempty"`
	EventActionProduce *ElementTagType         `json:"eventActionProduce,omitempty"`
	TimePeriod         *TimePeriodElementsType `json:"timePeriod,omitempty"`
}

type LoadControlEventListDataType struct {
	LoadControlEventData []LoadControlEventDataType `json:"loadControlEventData,omitempty"`
}

type LoadControlEventListDataSelectorsType struct {
	TimestampInterval *TimestampIntervalType  `json:"timestampInterval,omitempty"`
	EventId           *LoadControlEventIdType `json:"eventId,omitempty"`
}

type LoadControlStateDataType struct {
	Timestamp                 *string                     `json:"timestamp"`
	EventId                   *LoadControlEventIdType     `json:"eventId,omitempty"`
	EventStateConsume         *LoadControlEventStateType  `json:"eventStateConsume"`
	AppliedEventActionConsume *LoadControlEventActionType `json:"appliedEventActionConsume"`
	EventStateProduce         *LoadControlEventStateType  `json:"eventStateProduce"`
	AppliedEventActionProduce *LoadControlEventActionType `json:"appliedEventActionProduce"`
}

type LoadControlStateDataElementsType struct {
	Timestamp                 *ElementTagType `json:"timestamp"`
	EventId                   *ElementTagType `json:"eventId,omitempty"`
	EventStateConsume         *ElementTagType `json:"eventStateConsume"`
	AppliedEventActionConsume *ElementTagType `json:"appliedEventActionConsume"`
	EventStateProduce         *ElementTagType `json:"eventStateProduce"`
	AppliedEventActionProduce *ElementTagType `json:"appliedEventActionProduce"`
}

type LoadControlStateListDataType struct {
	LoadControlStateData []LoadControlStateDataType `json:"loadControlStateData,omitempty"`
}

type LoadControlStateListDataSelectorsType struct {
	TimestampInterval *TimestampIntervalType  `json:"timestampInterval,omitempty"`
	EventId           *LoadControlEventIdType `json:"eventId,omitempty"`
}

type LoadControlLimitDataType struct {
	LimitId           *LoadControlLimitIdType `json:"limitId,omitempty"`
	IsLimitChangeable *bool                   `json:"isLimitChangeable,omitempty"`
	IsLimitActive     *bool                   `json:"isLimitActive,omitempty"`
	TimePeriod        *TimePeriodType         `json:"timePeriod,omitempty"`
	Value             *ScaledNumberType       `json:"value,omitempty"`
}

type LoadControlLimitDataElementsType struct {
	LimitId           *ElementTagType           `json:"limitId,omitempty"`
	IsLimitChangeable *ElementTagType           `json:"isLimitChangeable,omitempty"`
	IsLimitActive     *ElementTagType           `json:"isLimitActive,omitempty"`
	TimePeriod        *TimePeriodElementsType   `json:"timePeriod,omitempty"`
	Value             *ScaledNumberElementsType `json:"value,omitempty"`
}

type LoadControlLimitListDataType struct {
	LoadControlLimitData []LoadControlLimitDataType `json:"loadControlLimitData,omitempty"`
}

type LoadControlLimitListDataSelectorsType struct {
	LimitId *LoadControlLimitIdType `json:"limitId,omitempty"`
}

type LoadControlLimitConstraintsDataType struct {
	LimitId       *LoadControlLimitIdType `json:"limitId,omitempty"`
	ValueRangeMin *ScaledNumberType       `json:"valueRangeMin,omitempty"`
	ValueRangeMax *ScaledNumberType       `json:"valueRangeMax,omitempty"`
	ValueStepSize *ScaledNumberType       `json:"valueStepSize,omitempty"`
}

type LoadControlLimitConstraintsDataElementsType struct {
	LimitId       *ElementTagType           `json:"limitId,omitempty"`
	ValueRangeMin *ScaledNumberElementsType `json:"valueRangeMin,omitempty"`
	ValueRangeMax *ScaledNumberElementsType `json:"valueRangeMax,omitempty"`
	ValueStepSize *ScaledNumberElementsType `json:"valueStepSize,omitempty"`
}

type LoadControlLimitConstraintsListDataType struct {
	LoadControlLimitConstraintsData []LoadControlLimitConstraintsDataType `json:"loadControlLimitConstraintsData,omitempty"`
}

type LoadControlLimitConstraintsListDataSelectorsType struct {
	LimitId *LoadControlLimitIdType `json:"limitId,omitempty"`
}

type LoadControlLimitDescriptionDataType struct {
	LimitId        *LoadControlLimitIdType   `json:"limitId,omitempty"`
	LimitType      *LoadControlLimitTypeType `json:"limitType,omitempty"`
	LimitCategory  *LoadControlCategoryType  `json:"limitCategory,omitempty"`
	LimitDirection *EnergyDirectionType      `json:"limitDirection,omitempty"`
	MeasurementId  *MeasurementIdType        `json:"measurementId,omitempty"`
	Unit           *UnitOfMeasurementType    `json:"unit,omitempty"`
	ScopeType      *ScopeTypeType            `json:"scopeType,omitempty"`
	Label          *LabelType                `json:"label,omitempty"`
	Description    *DescriptionType          `json:"description,omitempty"`
}

type LoadControlLimitDescriptionDataElementsType struct {
	LimitId        *ElementTagType `json:"limitId,omitempty"`
	LimitType      *ElementTagType `json:"limitType,omitempty"`
	LimitCategory  *ElementTagType `json:"limitCategory,omitempty"`
	LimitDirection *ElementTagType `json:"limitDirection,omitempty"`
	MeasurementId  *ElementTagType `json:"measurementId,omitempty"`
	Unit           *ElementTagType `json:"unit,omitempty"`
	ScopeType      *ElementTagType `json:"scopeType,omitempty"`
	Label          *ElementTagType `json:"label,omitempty"`
	Description    *ElementTagType `json:"description,omitempty"`
}

type LoadControlLimitDescriptionListDataType struct {
	LoadControlLimitDescriptionData []LoadControlLimitDescriptionDataType `json:"loadControlLimitDescriptionData,omitempty"`
}

type LoadControlLimitDescriptionListDataSelectorsType struct {
	LimitId        *LoadControlLimitIdType   `json:"limitId,omitempty"`
	LimitType      *LoadControlLimitTypeType `json:"limitType,omitempty"`
	LimitDirection *EnergyDirectionType      `json:"limitDirection,omitempty"`
	MeasurementId  *MeasurementIdType        `json:"measurementId,omitempty"`
	ScopeType      *ScopeTypeType            `json:"scopeType,omitempty"`
}
