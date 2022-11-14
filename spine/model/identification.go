package model

type IdentificationIdType uint

type IdentificationTypeType string

const (
	IdentificationTypeTypeEui48       IdentificationTypeType = "eui48"
	IdentificationTypeTypeEui64       IdentificationTypeType = "eui64"
	IdentificationTypeTypeUserrfidtag IdentificationTypeType = "userRfidTag"
)

type IdentificationValueType string

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
