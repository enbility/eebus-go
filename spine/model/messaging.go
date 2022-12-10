package model

type MessagingNumberType uint

type MessagingDataTextType string

type MessagingTypeType string

const (
	MessagingTypeTypeLogging     MessagingTypeType = "logging"
	MessagingTypeTypeInformation MessagingTypeType = "information"
	MessagingTypeTypeWarning     MessagingTypeType = "warning"
	MessagingTypeTypeAlarm       MessagingTypeType = "alarm"
	MessagingTypeTypeEmergency   MessagingTypeType = "emergency"
	MessagingTypeTypeObsolete    MessagingTypeType = "obsolete"
)

type MessagingDataType struct {
	Timestamp       *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	MessagingNumber *MessagingNumberType        `json:"messagingNumber,omitempty" eebus:"key"`
	MessagingType   *MessagingTypeType          `json:"type,omitempty"` // xsd defines "type", but that is a reserved keyword
	Text            *MessagingDataTextType      `json:"text,omitempty"`
}

type MessagingDataElementsType struct {
	Timestamp       *ElementTagType `json:"timestamp,omitempty"`
	MessagingNumber *ElementTagType `json:"messagingNumber,omitempty"`
	MessagingType   *ElementTagType `json:"type,omitempty"`
	Text            *ElementTagType `json:"text,omitempty"`
}

type MessagingListDataType struct {
	MessagingData []MessagingDataType `json:"messagingData,omitempty"`
}

type MessagingListDataSelectorsType struct {
	TimestampInterval *TimestampIntervalType `json:"timestampInterval,omitempty"`
	MessagingNumber   *MessagingNumberType   `json:"messagingNumber,omitempty"`
}
