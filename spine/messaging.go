package spine

type MessagingNumberType uint

type MessagingDataTextType string

type MessagingTypeType MessagingTypeEnumType

type MessagingTypeEnumType string

const (
	MessagingTypeEnumTypeLogging     MessagingTypeEnumType = "logging"
	MessagingTypeEnumTypeInformation MessagingTypeEnumType = "information"
	MessagingTypeEnumTypeWarning     MessagingTypeEnumType = "warning"
	MessagingTypeEnumTypeAlarm       MessagingTypeEnumType = "alarm"
	MessagingTypeEnumTypeEmergency   MessagingTypeEnumType = "emergency"
	MessagingTypeEnumTypeObsolete    MessagingTypeEnumType = "obsolete"
)

type MessagingDataType struct {
	Timestamp       *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	MessagingNumber *MessagingNumberType        `json:"messagingNumber,omitempty"`
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
