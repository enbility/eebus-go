package model

type DeviceConfigurationKeyIdType uint

type DeviceConfigurationKeyValueStringType string

type DeviceConfigurationKeyNameType DeviceConfigurationKeyNameEnumType

type DeviceConfigurationKeyNameEnumType string

const (
	DeviceConfigurationKeyNameEnumTypePeakPowerOfPVSystem         DeviceConfigurationKeyNameEnumType = "peakPowerOfPvSystem"
	DeviceConfigurationKeyNameEnumTypePvCurtailmentLimitFactor    DeviceConfigurationKeyNameEnumType = "pvCurtailmentLimitFactor"
	DeviceConfigurationKeyNameEnumTypeAsymmetricChargingSupported DeviceConfigurationKeyNameEnumType = "asymmetricChargingSupported"
	DeviceConfigurationKeyNameEnumTypeCommunicationsStandard      DeviceConfigurationKeyNameEnumType = "communicationsStandard"
)

type DeviceConfigurationKeyValueTypeType string

const (
	DeviceConfigurationKeyValueTypeTypeBoolean      DeviceConfigurationKeyValueTypeType = "boolean"
	DeviceConfigurationKeyValueTypeTypeDate         DeviceConfigurationKeyValueTypeType = "date"
	DeviceConfigurationKeyValueTypeTypeDatetime     DeviceConfigurationKeyValueTypeType = "dateTime"
	DeviceConfigurationKeyValueTypeTypeDuration     DeviceConfigurationKeyValueTypeType = "duration"
	DeviceConfigurationKeyValueTypeTypeString       DeviceConfigurationKeyValueTypeType = "string"
	DeviceConfigurationKeyValueTypeTypeTime         DeviceConfigurationKeyValueTypeType = "time"
	DeviceConfigurationKeyValueTypeTypeScalednumber DeviceConfigurationKeyValueTypeType = "scaledNumber"
)

type DeviceConfigurationKeyValueValueType struct {
	Boolean      *bool                                  `json:"boolean,omitempty"`
	Date         *string                                `json:"date,omitempty"`
	DateTime     *string                                `json:"dateTime,omitempty"`
	Duration     *string                                `json:"duration,omitempty"`
	String       *DeviceConfigurationKeyValueStringType `json:"string,omitempty"`
	Time         *string                                `json:"time,omitempty"`
	ScaledNumber *ScaledNumberType                      `json:"scaledNumber,omitempty"`
}

type DeviceConfigurationKeyValueValueElementsType struct {
	Boolean      *ElementTagType           `json:"boolean,omitempty"`
	Date         *ElementTagType           `json:"date,omitempty"`
	DateTime     *ElementTagType           `json:"dateTime,omitempty"`
	Duration     *ElementTagType           `json:"duration,omitempty"`
	String       *ElementTagType           `json:"string,omitempty"`
	Time         *ElementTagType           `json:"time,omitempty"`
	ScaledNumber *ScaledNumberElementsType `json:"scaledNumber,omitempty"`
}

type DeviceConfigurationKeyValueDataType struct {
	KeyId             *DeviceConfigurationKeyIdType         `json:"keyId,omitempty"`
	Value             *DeviceConfigurationKeyValueValueType `json:"value,omitempty"`
	IsValueChangeable *bool                                 `json:"isValueChangeable,omitempty"`
}

type DeviceConfigurationKeyValueDataElementsType struct {
	KeyId             *ElementTagType                               `json:"keyId,omitempty"`
	Value             *DeviceConfigurationKeyValueValueElementsType `json:"value,omitempty"`
	IsValueChangeable *ElementTagType                               `json:"isValueChangeable,omitempty"`
}

type DeviceConfigurationKeyValueListDataType struct {
	DeviceConfigurationKeyValueData []DeviceConfigurationKeyValueDataType `json:"deviceConfigurationKeyValueData,omitempty"`
}

type DeviceConfigurationKeyValueListDataSelectorsType struct {
	KeyId *DeviceConfigurationKeyIdType `json:"keyId,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionDataType struct {
	KeyId       *DeviceConfigurationKeyIdType        `json:"keyId,omitempty"`
	KeyName     *string                              `json:"keyName,omitempty"`
	ValueType   *DeviceConfigurationKeyValueTypeType `json:"valueType,omitempty"`
	Unit        *string                              `json:"unit,omitempty"`
	Label       *LabelType                           `json:"label,omitempty"`
	Description *DescriptionType                     `json:"description,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionDataElementsType struct {
	KeyId       *ElementTagType `json:"keyId,omitempty"`
	KeyName     *ElementTagType `json:"keyName,omitempty"`
	ValueType   *ElementTagType `json:"valueType,omitempty"`
	Unit        *ElementTagType `json:"unit,omitempty"`
	Label       *ElementTagType `json:"label,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionListDataType struct {
	DeviceConfigurationKeyValueDescriptionData []DeviceConfigurationKeyValueDescriptionDataType `json:"deviceConfigurationKeyValueDescriptionData,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionListDataSelectorsType struct {
	KeyId   *DeviceConfigurationKeyIdType `json:"keyId,omitempty"`
	KeyName *string                       `json:"keyName,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsDataType struct {
	KeyId         *DeviceConfigurationKeyIdType         `json:"keyId,omitempty"`
	ValueRangeMin *DeviceConfigurationKeyValueValueType `json:"valueRangeMin,omitempty"`
	ValueRangeMax *DeviceConfigurationKeyValueValueType `json:"valueRangeMax,omitempty"`
	ValueStepSize *DeviceConfigurationKeyValueValueType `json:"valueStepSize,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsDataElementsType struct {
	KeyId         *ElementTagType                               `json:"keyId,omitempty"`
	ValueRangeMin *DeviceConfigurationKeyValueValueElementsType `json:"valueRangeMin,omitempty"`
	ValueRangeMax *DeviceConfigurationKeyValueValueElementsType `json:"valueRangeMax,omitempty"`
	ValueStepSize *DeviceConfigurationKeyValueValueElementsType `json:"valueStepSize,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsListDataType struct {
	DeviceConfigurationKeyValueConstraintsData []DeviceConfigurationKeyValueConstraintsDataType `json:"deviceConfigurationKeyValueConstraintsData,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsListDataSelectorsType struct {
	KeyId *DeviceConfigurationKeyIdType `json:"keyId,omitempty"`
}
