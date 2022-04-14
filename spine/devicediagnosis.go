package spine

type VendorStateCodeType string

type LastErrorCodeType string

type DeviceDiagnosisOperatingStateType DeviceDiagnosisOperatingStateEnumType

type DeviceDiagnosisOperatingStateEnumType string

const (
	DeviceDiagnosisOperatingStateEnumTypeNormalOperation  DeviceDiagnosisOperatingStateEnumType = "normalOperation"
	DeviceDiagnosisOperatingStateEnumTypeStandby          DeviceDiagnosisOperatingStateEnumType = "standby"
	DeviceDiagnosisOperatingStateEnumTypeFailure          DeviceDiagnosisOperatingStateEnumType = "failure"
	DeviceDiagnosisOperatingStateEnumTypeServiceNeeded    DeviceDiagnosisOperatingStateEnumType = "serviceNeeded"
	DeviceDiagnosisOperatingStateEnumTypeOverrideDetected DeviceDiagnosisOperatingStateEnumType = "overrideDetected"
	DeviceDiagnosisOperatingStateEnumTypeInAlarm          DeviceDiagnosisOperatingStateEnumType = "inAlarm"
	DeviceDiagnosisOperatingStateEnumTypeNotReachable     DeviceDiagnosisOperatingStateEnumType = "notReachable"
	DeviceDiagnosisOperatingStateEnumTypeFinished         DeviceDiagnosisOperatingStateEnumType = "finished"
)

type PowerSupplyConditionType PowerSupplyConditionEnumType

type PowerSupplyConditionEnumType string

const (
	PowerSupplyConditionEnumTypeGood     PowerSupplyConditionEnumType = "good"
	PowerSupplyConditionEnumTypeLow      PowerSupplyConditionEnumType = "low"
	PowerSupplyConditionEnumTypeCritical PowerSupplyConditionEnumType = "critical"
	PowerSupplyConditionEnumTypeUnknown  PowerSupplyConditionEnumType = "unknown"
	PowerSupplyConditionEnumTypeError    PowerSupplyConditionEnumType = "error"
)

type DeviceDiagnosisStateDataType struct {
	Timestamp            *string                            `json:"timestamp,omitempty"`
	OperatingState       *DeviceDiagnosisOperatingStateType `json:"operatingState,omitempty"`
	VendorStateCode      *VendorStateCodeType               `json:"vendorStateCode,omitempty"`
	LastErrorCode        *LastErrorCodeType                 `json:"lastErrorCode,omitempty"`
	UpTime               *string                            `json:"upTime,omitempty"`
	TotalUpTime          *string                            `json:"totalUpTime,omitempty"`
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
	Timestamp        *string `json:"timestamp,omitempty"`
	HeartbeatCounter *uint64 `json:"heartbeatCounter,omitempty"`
	HeartbeatTimeout *string `json:"heartbeatTimeout,omitempty"`
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
