package spine

type IdentificationIdType uint

type IdentificationTypeType IdentificationTypeEnumType

type IdentificationTypeEnumType string

const (
	IdentificationTypeEnumTypeEui48       IdentificationTypeEnumType = "eui48"
	IdentificationTypeEnumTypeEui64       IdentificationTypeEnumType = "eui64"
	IdentificationTypeEnumTypeUserrfidtag IdentificationTypeEnumType = "userRfidTag"
)

type IdentificationValueType string

type IdentificationDataType struct {
	IdentificationId    *IdentificationIdType    `json:"identificationId,omitempty"`
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
