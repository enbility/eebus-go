package model

type DeviceClassificationStringType string

type PowerSourceType string

const (
	PowerSourceTypeUnknown          PowerSourceType = "unknown"
	PowerSourceTypeMainssinglephase PowerSourceType = "mainsSinglePhase"
	PowerSourceTypeMains3Phase      PowerSourceType = "mains3Phase"
	PowerSourceTypeBattery          PowerSourceType = "battery"
	PowerSourceTypeDc               PowerSourceType = "dc"
)

type DeviceClassificationManufacturerDataType struct {
	DeviceName                     *DeviceClassificationStringType `json:"deviceName,omitempty"`
	DeviceCode                     *DeviceClassificationStringType `json:"deviceCode,omitempty"`
	SerialNumber                   *DeviceClassificationStringType `json:"serialNumber,omitempty"`
	SoftwareRevision               *DeviceClassificationStringType `json:"softwareRevision,omitempty"`
	HardwareRevision               *DeviceClassificationStringType `json:"hardwareRevision,omitempty"`
	VendorName                     *DeviceClassificationStringType `json:"vendorName,omitempty"`
	VendorCode                     *DeviceClassificationStringType `json:"vendorCode,omitempty"`
	BrandName                      *DeviceClassificationStringType `json:"brandName,omitempty"`
	PowerSource                    *string                         `json:"powerSource,omitempty"`
	ManufacturerNodeIdentification *DeviceClassificationStringType `json:"manufacturerNodeIdentification,omitempty"`
	ManufacturerLabel              *LabelType                      `json:"manufacturerLabel,omitempty"`
	ManufacturerDescription        *DescriptionType                `json:"manufacturerDescription,omitempty"`
}

type DeviceClassificationManufacturerDataElementsType struct {
	DeviceName                     *ElementTagType `json:"deviceName,omitempty"`
	DeviceCode                     *ElementTagType `json:"deviceCode,omitempty"`
	SerialNumber                   *ElementTagType `json:"serialNumber,omitempty"`
	SoftwareRevision               *ElementTagType `json:"softwareRevision,omitempty"`
	HardwareRevision               *ElementTagType `json:"hardwareRevision,omitempty"`
	VendorName                     *ElementTagType `json:"vendorName,omitempty"`
	VendorCode                     *ElementTagType `json:"vendorCode,omitempty"`
	BrandName                      *ElementTagType `json:"brandName,omitempty"`
	PowerSource                    *ElementTagType `json:"powerSource,omitempty"`
	ManufacturerNodeIdentification *ElementTagType `json:"manufacturerNodeIdentification,omitempty"`
	ManufacturerLabel              *ElementTagType `json:"manufacturerLabel,omitempty"`
	ManufacturerDescription        *ElementTagType `json:"manufacturerDescription,omitempty"`
}

type DeviceClassificationUserDataType struct {
	UserNodeIdentification *DeviceClassificationStringType `json:"userNodeIdentification,omitempty"`
	UserLabel              *LabelType                      `json:"userLabel,omitempty"`
	UserDescription        *DescriptionType                `json:"userDescription,omitempty"`
}

type DeviceClassificationUserDataElementsType struct {
	UserNodeIdentification *ElementTagType `json:"userNodeIdentification,omitempty"`
	UserLabel              *ElementTagType `json:"userLabel,omitempty"`
	UserDescription        *ElementTagType `json:"userDescription,omitempty"`
}
