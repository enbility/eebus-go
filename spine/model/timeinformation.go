package model

type TimeInformationDataType struct {
	Utc          *DateTimeType     `json:"utc,omitempty"`
	UtcOffset    *DurationType     `json:"utcOffset,omitempty"`
	DayOfWeek    *DayOfWeekType    `json:"dayOfWeek,omitempty"`
	CalendarWeek *CalendarWeekType `json:"calendarWeek,omitempty"`
}

type TimeInformationDataElementsType struct {
	Utc          *ElementTagType `json:"utc,omitempty"`
	UtcOffset    *ElementTagType `json:"utcOffset,omitempty"`
	DayOfWeek    *ElementTagType `json:"dayOfWeek,omitempty"`
	CalendarWeek *ElementTagType `json:"calendarWeek,omitempty"`
}

type TimeDistributorDataType struct {
	IsTimeDistributor   *bool `json:"isTimeDistributor,omitempty"`
	DistributorPriority *uint `json:"distributorPriority,omitempty"`
}

type TimeDistributorDataElementsType struct {
	IsTimeDistributor   *ElementTagType `json:"isTimeDistributor,omitempty"`
	DistributorPriority *ElementTagType `json:"distributorPriority,omitempty"`
}

type TimePrecisionDataType struct {
	IsSynchronised *bool         `json:"isSynchronised,omitempty"`
	LastSyncAt     *DateTimeType `json:"lastSyncAt,omitempty"`
	ClockDrift     *int          `json:"clockDrift,omitempty"`
}

type TimePrecisionDataElementsType struct {
	IsSynchronised *ElementTagType `json:"isSynchronised,omitempty"`
	LastSyncAt     *ElementTagType `json:"lastSyncAt,omitempty"`
	ClockDrift     *ElementTagType `json:"clockDrift,omitempty"`
}

type TimeDistributorEnquiryCallType struct{}

type TimeDistributorEnquiryCallElementsType struct{}
