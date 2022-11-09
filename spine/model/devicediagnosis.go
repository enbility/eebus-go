package model

type VendorStateCodeType string

type LastErrorCodeType string

type DeviceDiagnosisOperatingStateType string

const (
	DeviceDiagnosisOperatingStateTypeNormalOperation  DeviceDiagnosisOperatingStateType = "normalOperation"
	DeviceDiagnosisOperatingStateTypeStandby          DeviceDiagnosisOperatingStateType = "standby"
	DeviceDiagnosisOperatingStateTypeFailure          DeviceDiagnosisOperatingStateType = "failure"
	DeviceDiagnosisOperatingStateTypeServiceNeeded    DeviceDiagnosisOperatingStateType = "serviceNeeded"
	DeviceDiagnosisOperatingStateTypeOverrideDetected DeviceDiagnosisOperatingStateType = "overrideDetected"
	DeviceDiagnosisOperatingStateTypeInAlarm          DeviceDiagnosisOperatingStateType = "inAlarm"
	DeviceDiagnosisOperatingStateTypeNotReachable     DeviceDiagnosisOperatingStateType = "notReachable"
	DeviceDiagnosisOperatingStateTypeFinished         DeviceDiagnosisOperatingStateType = "finished"
)

type PowerSupplyConditionType string

const (
	PowerSupplyConditionTypeGood     PowerSupplyConditionType = "good"
	PowerSupplyConditionTypeLow      PowerSupplyConditionType = "low"
	PowerSupplyConditionTypeCritical PowerSupplyConditionType = "critical"
	PowerSupplyConditionTypeUnknown  PowerSupplyConditionType = "unknown"
	PowerSupplyConditionTypeError    PowerSupplyConditionType = "error"
)

type DeviceDiagnosisStateDataType struct {
	Timestamp            *string                            `json:"timestamp,omitempty"`
	OperatingState       *DeviceDiagnosisOperatingStateType `json:"operatingState,omitempty"`
	VendorStateCode      *VendorStateCodeType               `json:"vendorStateCode,omitempty"`
	LastErrorCode        *LastErrorCodeType                 `json:"lastErrorCode,omitempty"`
	UpTime               *DurationType                      `json:"upTime,omitempty"`
	TotalUpTime          *DurationType                      `json:"totalUpTime,omitempty"`
	PowerSupplyCondition *PowerSupplyConditionType          `json:"powerSupplyCondition,omitempty"`
}

type DeviceDiagnosisStateDataElementsType struct {
	Timestamp            *ElementTagType `json:"timestamp,omitempty"`
	OperatingState       *ElementTagType `json:"operatingState,omitempty"`
	VendorStateCode      *ElementTagType `json:"vendorStateCode,omitempty"`
	LastErrorCode        *ElementTagType `json:"lastErrorCode,omitempty"`
	UpTime               *ElementTagType `json:"upTime,omitempty"`
	TotalUpTime          *ElementTagType `json:"totalUpTime,omitempty"`
	PowerSupplyCondition *ElementTagType `json:"powerSupplyCondition,omitempty"`
}

type DeviceDiagnosisHeartbeatDataType struct {
	Timestamp        *string       `json:"timestamp,omitempty"`
	HeartbeatCounter *uint64       `json:"heartbeatCounter,omitempty"`
	HeartbeatTimeout *DurationType `json:"heartbeatTimeout,omitempty"`
}

type DeviceDiagnosisHeartbeatDataElementsType struct {
	Timestamp        *ElementTagType `json:"timestamp,omitempty"`
	HeartbeatCounter *ElementTagType `json:"heartbeatCounter,omitempty"`
	HeartbeatTimeout *ElementTagType `json:"heartbeatTimeout,omitempty"`
}

type DeviceDiagnosisServiceDataType struct {
	Timestamp        *string `json:"timestamp,omitempty"`
	InstallationTime *string `json:"installationTime,omitempty"`
	BootCounter      *uint64 `json:"bootCounter,omitempty"`
	NextService      *string `json:"nextService,omitempty"`
}

type DeviceDiagnosisServiceDataElementsType struct {
	Timestamp        *ElementTagType `json:"timestamp,omitempty"`
	InstallationTime *ElementTagType `json:"installationTime,omitempty"`
	BootCounter      *ElementTagType `json:"bootCounter,omitempty"`
	NextService      *ElementTagType `json:"nextService,omitempty"`
}
