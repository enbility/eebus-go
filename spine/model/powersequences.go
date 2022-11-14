package model

type AlternativesIdType uint

type PowerSequenceIdType uint

type PowerTimeSlotNumberType uint

type PowerTimeSlotValueTypeType string

const (
	PowerTimeSlotValueTypeTypePower                   PowerTimeSlotValueTypeType = "power"
	PowerTimeSlotValueTypeTypePowerMin                PowerTimeSlotValueTypeType = "powerMin"
	PowerTimeSlotValueTypeTypePowerMax                PowerTimeSlotValueTypeType = "powerMax"
	PowerTimeSlotValueTypeTypePowerExpectedValue      PowerTimeSlotValueTypeType = "powerExpectedValue"
	PowerTimeSlotValueTypeTypePowerStandardDeviation  PowerTimeSlotValueTypeType = "powerStandardDeviation"
	PowerTimeSlotValueTypeTypePowerSkewness           PowerTimeSlotValueTypeType = "powerSkewness"
	PowerTimeSlotValueTypeTypeEnergy                  PowerTimeSlotValueTypeType = "energy"
	PowerTimeSlotValueTypeTypeEnergyMin               PowerTimeSlotValueTypeType = "energyMin"
	PowerTimeSlotValueTypeTypeEnergyMax               PowerTimeSlotValueTypeType = "energyMax"
	PowerTimeSlotValueTypeTypeEnergyExpectedValue     PowerTimeSlotValueTypeType = "energyExpectedValue"
	PowerTimeSlotValueTypeTypeEnergyStandardDeviation PowerTimeSlotValueTypeType = "energyStandardDeviation"
	PowerTimeSlotValueTypeTypeEnergySkewness          PowerTimeSlotValueTypeType = "energySkewness"
)

type PowerSequenceScopeType string

const (
	PowerSequenceScopeTypeForecast       PowerSequenceScopeType = "forecast"
	PowerSequenceScopeTypeMeasurement    PowerSequenceScopeType = "measurement"
	PowerSequenceScopeTypeRecommendation PowerSequenceScopeType = "recommendation"
)

type PowerSequenceStateType string

const (
	PowerSequenceStateTypeRunning         PowerSequenceStateType = "running"
	PowerSequenceStateTypePaused          PowerSequenceStateType = "paused"
	PowerSequenceStateTypeScheduled       PowerSequenceStateType = "scheduled"
	PowerSequenceStateTypeScheduledPaused PowerSequenceStateType = "scheduledPaused"
	PowerSequenceStateTypePending         PowerSequenceStateType = "pending"
	PowerSequenceStateTypeInactive        PowerSequenceStateType = "inactive"
	PowerSequenceStateTypeCompleted       PowerSequenceStateType = "completed"
	PowerSequenceStateTypeInvalid         PowerSequenceStateType = "invalid"
)

type PowerTimeSlotScheduleDataType struct {
	SequenceId          *PowerSequenceIdType     `json:"sequenceId,omitempty" eebus:"key"`
	SlotNumber          *PowerTimeSlotNumberType `json:"slotNumber,omitempty"`
	TimePeriod          *TimePeriodType          `json:"timePeriod,omitempty"`
	DefaultDuration     *DurationType            `json:"defaultDuration,omitempty"`
	DurationUncertainty *DurationType            `json:"durationUncertainty,omitempty"`
	SlotActivated       *bool                    `json:"slotActivated,omitempty"`
	Description         *DescriptionType         `json:"description,omitempty"`
}

type PowerTimeSlotScheduleDataElementsType struct {
	SequenceId          *ElementTagType `json:"sequenceId,omitempty"`
	SlotNumber          *ElementTagType `json:"slotNumber,omitempty"`
	TimePeriod          *ElementTagType `json:"timePeriod,omitempty"`
	DefaultDuration     *ElementTagType `json:"defaultDuration,omitempty"`
	DurationUncertainty *ElementTagType `json:"durationUncertainty,omitempty"`
	SlotActivated       *ElementTagType `json:"slotActivated,omitempty"`
	Description         *ElementTagType `json:"description,omitempty"`
}

type PowerTimeSlotScheduleListDataType struct {
	PowerTimeSlotScheduleData []PowerTimeSlotScheduleDataType `json:"powerTimeSlotScheduleData,omitempty"`
}

type PowerTimeSlotScheduleListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType     `json:"sequenceId,omitempty"`
	SlotNumber *PowerTimeSlotNumberType `json:"slotNumber,omitempty"`
}

type PowerTimeSlotValueDataType struct {
	SequenceId *PowerSequenceIdType        `json:"sequenceId,omitempty" eebus:"key"`
	SlotNumber *PowerTimeSlotNumberType    `json:"slotNumber,omitempty"`
	ValueType  *PowerTimeSlotValueTypeType `json:"valueType,omitempty"`
	Value      *ScaledNumberType           `json:"value,omitempty"`
}

type PowerTimeSlotValueDataElementsType struct {
	SequenceId *ElementTagType           `json:"sequenceId,omitempty"`
	SlotNumber *ElementTagType           `json:"slotNumber,omitempty"`
	ValueType  *ElementTagType           `json:"valueType,omitempty"`
	Value      *ScaledNumberElementsType `json:"value,omitempty"`
}

type PowerTimeSlotValueListDataType struct {
	PowerTimeSlotValueListDataType []PowerTimeSlotValueDataType `json:"powerTimeSlotValueListData,omitempty"`
}

type PowerTimeSlotValueListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType        `json:"sequenceId,omitempty"`
	SlotNumber *PowerTimeSlotNumberType    `json:"slotNumber,omitempty"`
	ValueType  *PowerTimeSlotValueTypeType `json:"valueType,omitempty"`
}

type PowerTimeSlotScheduleConstraintsDataType struct {
	SequenceId        *PowerSequenceIdType        `json:"sequenceId,omitempty" eebus:"key"`
	SlotNumber        *PowerTimeSlotNumberType    `json:"slotNumber,omitempty"`
	EarliestStartTime *AbsoluteOrRelativeTimeType `json:"earliestStartTime,omitempty"`
	LatestEndTime     *AbsoluteOrRelativeTimeType `json:"latestEndTime,omitempty"`
	MinDuration       *DurationType               `json:"minDuration,omitempty"`
	MaxDuration       *DurationType               `json:"maxDuration,omitempty"`
	OptionalSlot      *bool                       `json:"optionalSlot,omitempty"`
}

type PowerTimeSlotScheduleConstraintsDataElementsType struct {
	SequenceId        *ElementTagType `json:"sequenceId,omitempty"`
	SlotNumber        *ElementTagType `json:"slotNumber,omitempty"`
	EarliestStartTime *ElementTagType `json:"earliestStartTime,omitempty"`
	LatestEndTime     *ElementTagType `json:"latestEndTime,omitempty"`
	MinDuration       *ElementTagType `json:"minDuration,omitempty"`
	MaxDuration       *ElementTagType `json:"maxDuration,omitempty"`
	OptionalSlot      *ElementTagType `json:"optionalSlot,omitempty"`
}

type PowerTimeSlotScheduleConstraintsListDataType struct {
	PowerTimeSlotScheduleConstraintsData []PowerTimeSlotScheduleConstraintsDataType `json:"powerTimeSlotScheduleConstraintsData,omitempty"`
}

type PowerTimeSlotScheduleConstraintsListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType     `json:"sequenceId,omitempty"`
	SlotNumber *PowerTimeSlotNumberType `json:"slotNumber,omitempty"`
}

type PowerSequenceAlternativesRelationDataType struct {
	AlternativeId *AlternativesIdType   `json:"alternativeId,omitempty" eebus:"key"`
	SequenceId    []PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceAlternativesRelationDataElementsType struct {
	AlternativeId *ElementTagType `json:"alternativeId,omitempty"`
	SequenceId    *ElementTagType `json:"sequenceId,omitempty"`
}

type PowerSequenceAlternativesRelationListDataType struct {
	PowerSequenceAlternativesRelationData []PowerSequenceAlternativesRelationDataType `json:"powerSequenceAlternativesRelationData,omitempty"`
}

type PowerSequenceAlternativesRelationListDataSelectorsType struct {
	AlternativeId *AlternativesIdType   `json:"alternativeId,omitempty"`
	SequenceId    []PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceDescriptionDataType struct {
	SequenceId              *PowerSequenceIdType        `json:"sequenceId,omitempty" eebus:"key"`
	Description             *DescriptionType            `json:"description,omitempty"`
	PositiveEnergyDirection *EnergyDirectionType        `json:"positiveEnergyDirection,omitempty"`
	PowerUnit               *UnitOfMeasurementType      `json:"powerUnit,omitempty"`
	EnergyUnit              *UnitOfMeasurementType      `json:"energyUnit,omitempty"`
	ValueSource             *MeasurementValueSourceType `json:"valueSource,omitempty"`
	Scope                   *PowerSequenceScopeType     `json:"scope,omitempty"`
	TaskIdentifier          *uint                       `json:"taskIdentifier,omitempty"`
	RepetitionsTotal        *uint                       `json:"repetitionsTotal,omitempty"`
}

type PowerSequenceDescriptionDataElementsType struct {
	SequenceId              *ElementTagType `json:"sequenceId,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
	PositiveEnergyDirection *ElementTagType `json:"positiveEnergyDirection,omitempty"`
	PowerUnit               *ElementTagType `json:"powerUnit,omitempty"`
	EnergyUnit              *ElementTagType `json:"energyUnit,omitempty"`
	ValueSource             *ElementTagType `json:"valueSource,omitempty"`
	Scope                   *ElementTagType `json:"scope,omitempty"`
	TaskIdentifier          *ElementTagType `json:"taskIdentifier,omitempty"`
	RepetitionsTotal        *ElementTagType `json:"repetitionsTotal,omitempty"`
}

type PowerSequenceDescriptionListDataType struct {
	PowerSequenceDescriptionData []PowerSequenceDescriptionDataType `json:"powerSequenceDescriptionData,omitempty"`
}

type PowerSequenceDescriptionListDataSelectorsType struct {
	SequenceId []PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceStateDataType struct {
	SequenceId                 *PowerSequenceIdType     `json:"sequenceId,omitempty" eebus:"key"`
	State                      *PowerSequenceStateType  `json:"state,omitempty"`
	ActiveSlotNumber           *PowerTimeSlotNumberType `json:"activeSlotNumber,omitempty"`
	ElapsedSlotTime            *DurationType            `json:"elapsedSlotTime,omitempty"`
	RemainingSlotTime          *DurationType            `json:"remainingSlotTime,omitempty"`
	SequenceRemoteControllable *bool                    `json:"sequenceRemoteControllable,omitempty"`
	ActiveRepetitionNumber     *uint                    `json:"activeRepetitionNumber,omitempty"`
	RemainingPauseTime         *DurationType            `json:"remainingPauseTime,omitempty"`
}

type PowerSequenceStateDataElementsType struct {
	SequenceId                 *ElementTagType `json:"sequenceId,omitempty"`
	State                      *ElementTagType `json:"state,omitempty"`
	ActiveSlotNumber           *ElementTagType `json:"activeSlotNumber,omitempty"`
	ElapsedSlotTime            *ElementTagType `json:"elapsedSlotTime,omitempty"`
	RemainingSlotTime          *ElementTagType `json:"remainingSlotTime,omitempty"`
	SequenceRemoteControllable *ElementTagType `json:"sequenceRemoteControllable,omitempty"`
	ActiveRepetitionNumber     *ElementTagType `json:"activeRepetitionNumber,omitempty"`
	RemainingPauseTime         *ElementTagType `json:"remainingPauseTime,omitempty"`
}

type PowerSequenceStateListDataType struct {
	PowerSequenceStateData []PowerSequenceStateDataType `json:"powerSequenceStateData,omitempty"`
}

type PowerSequenceStateListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceScheduleDataType struct {
	SequenceId *PowerSequenceIdType        `json:"sequenceId,omitempty" eebus:"key"`
	StartTime  *AbsoluteOrRelativeTimeType `json:"startTime,omitempty"`
	EndTime    *AbsoluteOrRelativeTimeType `json:"endTime,omitempty"`
}

type PowerSequenceScheduleDataElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
	StartTime  *ElementTagType `json:"startTime,omitempty"`
	EndTime    *ElementTagType `json:"endTime,omitempty"`
}

type PowerSequenceScheduleListDataType struct {
	PowerSequenceScheduleData []PowerSequenceScheduleDataType `json:"powerSequenceScheduleData,omitempty"`
}

type PowerSequenceScheduleListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceScheduleConstraintsDataType struct {
	SequenceId        *PowerSequenceIdType        `json:"sequenceId,omitempty" eebus:"key"`
	EarliestStartTime *AbsoluteOrRelativeTimeType `json:"earliestStartTime,omitempty"`
	LatestStartTime   *AbsoluteOrRelativeTimeType `json:"latestStartTime,omitempty"`
	EarliestEndTime   *AbsoluteOrRelativeTimeType `json:"earliestEndTime,omitempty"`
	LatestEndTime     *AbsoluteOrRelativeTimeType `json:"latestEndTime,omitempty"`
	OptionalSequence  *bool                       `json:"optionalSequence,omitempty"`
}

type PowerSequenceScheduleConstraintsDataElementsType struct {
	SequenceId        *ElementTagType `json:"sequenceId,omitempty"`
	EarliestStartTime *ElementTagType `json:"earliestStartTime,omitempty"`
	LatestStartTime   *ElementTagType `json:"latestStartTime,omitempty"`
	EarliestEndTime   *ElementTagType `json:"earliestEndTime,omitempty"`
	LatestEndTime     *ElementTagType `json:"latestEndTime,omitempty"`
	OptionalSequence  *ElementTagType `json:"optionalSequence,omitempty"`
}

type PowerSequenceScheduleConstraintsListDataType struct {
	PowerSequenceScheduleConstraintsData []PowerSequenceScheduleConstraintsDataType `json:"powerSequenceScheduleConstraintsData,omitempty"`
}

type PowerSequenceScheduleConstraintsListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequencePriceDataType struct {
	SequenceId         *PowerSequenceIdType        `json:"sequenceId,omitempty" eebus:"key"`
	PotentialStartTime *AbsoluteOrRelativeTimeType `json:"potentialStartTime,omitempty"`
	Price              *ScaledNumberType           `json:"price,omitempty"`
	Currency           *CurrencyType               `json:"currency,omitempty"`
}

type PowerSequencePriceDataElementsType struct {
	SequenceId         *ElementTagType `json:"sequenceId,omitempty"`
	PotentialStartTime *ElementTagType `json:"potentialStartTime,omitempty"`
	Price              *ElementTagType `json:"price,omitempty"`
	Currency           *ElementTagType `json:"currency,omitempty"`
}

type PowerSequencePriceListDataType struct {
	PowerSequencePriceData []PowerSequencePriceDataType `json:"powerSequencePriceData,omitempty"`
}

type PowerSequencePriceListDataSelectorsType struct {
	SequenceId                 *PowerSequenceIdType        `json:"sequenceId,omitempty"`
	PotentialStartTimeInterval *AbsoluteOrRelativeTimeType `json:"potentialStartTimeInterval,omitempty"`
}

type PowerSequenceSchedulePreferenceDataType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty" eebus:"key"`
	Greenest   *bool                `json:"greenest,omitempty"`
	Cheapest   *bool                `json:"cheapest,omitempty"`
}

type PowerSequenceSchedulePreferenceDataElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
	Greenest   *ElementTagType `json:"greenest,omitempty"`
	Cheapest   *ElementTagType `json:"cheapest,omitempty"`
}

type PowerSequenceSchedulePreferenceListDataType struct {
	PowerSequenceSchedulePreferenceData []PowerSequenceSchedulePreferenceDataType `json:"powerSequenceSchedulePreferenceData,omitempty"`
}

type PowerSequenceSchedulePreferenceListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceNodeScheduleInformationDataType struct {
	NodeRemoteControllable           *bool `json:"nodeRemoteControllable,omitempty"`
	SupportsSingleSlotSchedulingOnly *bool `json:"supportsSingleSlotSchedulingOnly,omitempty"`
	AlternativesCount                *uint `json:"alternativesCount,omitempty"`
	TotalSequencesCountMax           *uint `json:"totalSequencesCountMax,omitempty"`
	SupportsReselection              *bool `json:"supportsReselection,omitempty"`
}

type PowerSequenceNodeScheduleInformationDataElementsType struct {
	NodeRemoteControllable           *ElementTagType `json:"nodeRemoteControllable,omitempty"`
	SupportsSingleSlotSchedulingOnly *ElementTagType `json:"supportsSingleSlotSchedulingOnly,omitempty"`
	AlternativesCount                *ElementTagType `json:"alternativesCount,omitempty"`
	TotalSequencesCountMax           *ElementTagType `json:"totalSequencesCountMax,omitempty"`
	SupportsReselection              *ElementTagType `json:"supportsReselection,omitempty"`
}

type PowerSequenceScheduleConfigurationRequestCallType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type PowerSequenceScheduleConfigurationRequestCallElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
}

type PowerSequencePriceCalculationRequestCallType struct {
	SequenceId         *PowerSequenceIdType        `json:"sequenceId,omitempty"`
	PotentialStartTime *AbsoluteOrRelativeTimeType `json:"potentialStartTime,omitempty"`
}

type PowerSequencePriceCalculationRequestCallElementsType struct {
	SequenceId         *ElementTagType `json:"sequenceId,omitempty"`
	PotentialStartTime *ElementTagType `json:"potentialStartTime,omitempty"`
}
