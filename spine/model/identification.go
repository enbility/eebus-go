package model

type IdentificationIdType uint

type IdentificationTypeType string

const (
	IdentificationTypeTypeEui48       IdentificationTypeType = "eui48"
	IdentificationTypeTypeEui64       IdentificationTypeType = "eui64"
	IdentificationTypeTypeUserrfidtag IdentificationTypeType = "userRfidTag"
)

type IdentificationValueType string

type SessionIdType uint

type IdentificationDataType struct {
	IdentificationId    *IdentificationIdType    `json:"identificationId,omitempty" eebus:"key"`
	IdentificationType  *IdentificationTypeType  `json:"identificationType,omitempty"`
	IdentificationValue *IdentificationValueType `json:"identificationValue,omitempty"`
	Authorized          *bool                    `json:"authorized,omitempty"`
}

type IdentificationDataElementsType struct {
	IdentificationId    *ElementTagType `json:"identificationId,omitempty"`
	IdentificationType  *ElementTagType `json:"identificationType,omitempty"`
	IdentificationValue *ElementTagType `json:"identificationValue,omitempty"`
	Authorized          *ElementTagType `json:"authorized,omitempty"`
}

type IdentificationListDataType struct {
	IdentificationData []IdentificationDataType `json:"identificationData,omitempty"`
}

type IdentificationListDataSelectorsType struct {
	IdentificationId   *IdentificationIdType   `json:"identificationId,omitempty"`
	IdentificationType *IdentificationTypeType `json:"identificationType,omitempty"`
}

type SessionIdentificationDataType struct {
	SessionId        *SessionIdType        `json:"sessionId,omitempty" eebus:"key"`
	IdentificationId *IdentificationIdType `json:"identificationId,omitempty"`
	IsLatestSession  *bool                 `json:"isLatestSession,omitempty"`
	TimePeriod       *TimePeriodType       `json:"timePeriod,omitempty"`
}

type SessionIdentificationDataElementsType struct {
	SessionId        *ElementTagType         `json:"sessionId,omitempty"`
	IdentificationId *ElementTagType         `json:"identificationId,omitempty"`
	IsLatestSession  *ElementTagType         `json:"isLatestSession,omitempty"`
	TimePeriod       *TimePeriodElementsType `json:"timePeriod,omitempty"`
}

type SessionIdentificationListDataType struct {
	SessionIdentificationData []SessionIdentificationDataType `json:"sessionIdentificationData,omitempty"`
}

type SessionIdentificationListDataSelectorsType struct {
	SessionId        *SessionIdType        `json:"sessionId,omitempty"`
	IdentificationId *IdentificationIdType `json:"identificationId,omitempty"`
	IsLatestSession  *bool                 `json:"isLatestSession,omitempty"`
	TimePeriod       *TimePeriodType       `json:"timePeriod,omitempty"`
}

type SessionMeasurementRelationDataType struct {
	SessionId     *SessionIdType      `json:"sessionId,omitempty" eebus:"key"`
	MeasurementId []MeasurementIdType `json:"measurementId,omitempty"`
}

type SessionMeasurementRelationDataElementsType struct {
	SessionId     *ElementTagType `json:"sessionId,omitempty"`
	MeasurementId *ElementTagType `json:"measurementId,omitempty"`
}

type SessionMeasurementRelationListDataType struct {
	SessionMeasurementRelationData []SessionMeasurementRelationDataType `json:"sessionMeasurementRelationData,omitempty"`
}

type SessionMeasurementRelationListDataSelectorsType struct {
	SessionId     *SessionIdType     `json:"sessionId,omitempty"`
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
}
