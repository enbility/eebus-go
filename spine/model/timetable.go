package model

type TimeTableIdType uint

type TimeSlotIdType uint

type TimeSlotCountType TimeSlotIdType

type TimeSlotTimeModeType string

const (
	TimeSlotTimeModeTypeAbsolute  TimeSlotTimeModeType = "absolute"
	TimeSlotTimeModeTypeRecurring TimeSlotTimeModeType = "recurring"
	TimeSlotTimeModeTypeBoth      TimeSlotTimeModeType = "both"
)

type TimeTableDataType struct {
	TimeTableId           *TimeTableIdType             `json:"timeTableId,omitempty" eebus:"key"`
	TimeSlotId            *TimeSlotIdType              `json:"timeSlotId,omitempty"`
	RecurrenceInformation *RecurrenceInformationType   `json:"recurrenceInformation,omitempty"`
	StartTime             *AbsoluteOrRecurringTimeType `json:"startTime,omitempty"`
	EndTime               *AbsoluteOrRecurringTimeType `json:"endTime,omitempty"`
}

type TimeTableDataElementsType struct {
	TimeTableId           *ElementTagType                      `json:"timeTableId,omitempty"`
	TimeSlotId            *ElementTagType                      `json:"timeSlotId,omitempty"`
	RecurrenceInformation *RecurrenceInformationElementsType   `json:"recurrenceInformation,omitempty"`
	StartTime             *AbsoluteOrRecurringTimeElementsType `json:"startTime,omitempty"`
	EndTime               *AbsoluteOrRecurringTimeElementsType `json:"endTime,omitempty"`
}

type TimeTableListDataType struct {
	TimeTableData []TimeTableDataType `json:"timeTableData,omitempty"`
}

type TimeTableListDataSelectorsType struct {
	TimeTableId *TimeTableIdType `json:"timeTableId,omitempty"`
	TimeSlotId  *TimeSlotIdType  `json:"timeSlotId,omitempty"`
}

type TimeTableConstraintsDataType struct {
	TimeTableId          *TimeTableIdType   `json:"timeTableId,omitempty" eebus:"key"`
	SlotCountMin         *TimeSlotCountType `json:"slotCountMin,omitempty"`
	SlotCountMax         *TimeSlotCountType `json:"slotCountMax,omitempty"`
	SlotDurationMin      *DurationType      `json:"slotDurationMin,omitempty"`
	SlotDurationMax      *DurationType      `json:"slotDurationMax,omitempty"`
	SlotDurationStepSize *DurationType      `json:"slotDurationStepSize,omitempty"`
	SlotShiftStepSize    *DurationType      `json:"slotShiftStepSize,omitempty"`
	FirstSlotBeginsAt    *TimeType          `json:"firstSlotBeginsAt,omitempty"`
}

type TimeTableConstraintsDataElementsType struct {
	TimeTableId          *ElementTagType `json:"timeTableId,omitempty"`
	SlotCountMin         *ElementTagType `json:"slotCountMin,omitempty"`
	SlotCountMax         *ElementTagType `json:"slotCountMax,omitempty"`
	SlotDurationMin      *ElementTagType `json:"slotDurationMin,omitempty"`
	SlotDurationMax      *ElementTagType `json:"slotDurationMax,omitempty"`
	SlotDurationStepSize *ElementTagType `json:"slotDurationStepSize,omitempty"`
	SlotShiftStepSize    *ElementTagType `json:"slotShiftStepSize,omitempty"`
	FirstSlotBeginsAt    *ElementTagType `json:"firstSlotBeginsAt,omitempty"`
}

type TimeTableConstraintsListDataType struct {
	TimeTableConstraintsData []TimeTableConstraintsDataType `json:"timeTableConstraintsData,omitempty"`
}

type TimeTableConstraintsListDataSelectorsType struct {
	TimeTableId *TimeTableIdType `json:"timeTableId,omitempty"`
}

type TimeTableDescriptionDataType struct {
	TimeTableId             *TimeTableIdType      `json:"timeTableId,omitempty" eebus:"key"`
	TimeSlotCountChangeable *bool                 `json:"timeSlotCountChangeable,omitempty"`
	TimeSlotTimesChangeable *bool                 `json:"timeSlotTimesChangeable,omitempty"`
	TimeSlotTimeMode        *TimeSlotTimeModeType `json:"timeSlotTimeMode,omitempty"`
	Label                   *LabelType            `json:"label,omitempty"`
	Description             *DescriptionType      `json:"description,omitempty"`
}

type TimeTableDescriptionDataElementsType struct {
	TimeTableId             *ElementTagType `json:"timeTableId,omitempty"`
	TimeSlotCountChangeable *ElementTagType `json:"timeSlotCountChangeable,omitempty"`
	TimeSlotTimesChangeable *ElementTagType `json:"timeSlotTimesChangeable,omitempty"`
	TimeSlotTimeMode        *ElementTagType `json:"timeSlotTimeMode,omitempty"`
	Label                   *ElementTagType `json:"label,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
}

type TimeTableDescriptionListDataType struct {
	TimeTableDescriptionData []TimeTableDescriptionDataType `json:"timeTableDescriptionData,omitempty"`
}

type TimeTableDescriptionListDataSelectorsType struct {
	TimeTableId *TimeTableIdType `json:"timeTableId,omitempty"`
}
