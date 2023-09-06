package model

type StateInformationIdType uint

type StateInformationType string

// union  StateInformationFunctionalityType StateInformationFunctionalityType
const (
	StateInformationTypeExternalOverrideFromGrid      StateInformationType = "externalOverrideFromGrid"
	StateInformationTypeAutonomousGridSupport         StateInformationType = "autonomousGridSupport"
	StateInformationTypeIslandingMode                 StateInformationType = "islandingMode"
	StateInformationTypeBalancing                     StateInformationType = "balancing"
	StateInformationTypeTrickleCharging               StateInformationType = "trickleCharging"
	StateInformationTypeCalibration                   StateInformationType = "calibration"
	StateInformationTypeCommissioningMissing          StateInformationType = "commissioningMissing"
	StateInformationTypeSleeping                      StateInformationType = "sleeping"
	StateInformationTypeStarting                      StateInformationType = "starting"
	StateInformationTypeMppt                          StateInformationType = "mppt"
	StateInformationTypeThrottled                     StateInformationType = "throttled"
	StateInformationTypeShuttingDown                  StateInformationType = "shuttingDown"
	StateInformationTypeManualShutdown                StateInformationType = "manualShutdown"
	StateInformationTypeInverterDefective             StateInformationType = "inverterDefective"
	StateInformationTypeBatteryOvercurrentProtection  StateInformationType = "batteryOvercurrentProtection"
	StateInformationTypePvStringOvercurrentProtection StateInformationType = "pvStringOvercurrentProtection"
	StateInformationTypeGridFault                     StateInformationType = "gridFault"
	StateInformationTypeGroundFault                   StateInformationType = "groundFault"
	StateInformationTypeAcDisconnected                StateInformationType = "acDisconnected"
	StateInformationTypeDcDisconnected                StateInformationType = "dcDisconnected"
	StateInformationTypeCabinetOpen                   StateInformationType = "cabinetOpen"
	StateInformationTypeOverTemperature               StateInformationType = "overTemperature"
	StateInformationTypeUnderTemperature              StateInformationType = "underTemperature"
	StateInformationTypeFrequencyAboveLimit           StateInformationType = "frequencyAboveLimit"
	StateInformationTypeFrequencyBelowLimit           StateInformationType = "frequencyBelowLimit"
	StateInformationTypeAcVoltageAboveLimit           StateInformationType = "acVoltageAboveLimit"
	StateInformationTypeAcVoltageBelowLimit           StateInformationType = "acVoltageBelowLimit"
	StateInformationTypeDcVoltageAboveLimit           StateInformationType = "dcVoltageAboveLimit"
	StateInformationTypeDcVoltageBelowLimit           StateInformationType = "dcVoltageBelowLimit"
	StateInformationTypeHardwareTestFailure           StateInformationType = "hardwareTestFailure"
	StateInformationTypeGenericInternalError          StateInformationType = "genericInternalError"
)

type StateInformationFunctionalityType string

const (
	StateInformationFunctionalityTypeExternalOverrideFromGrid StateInformationFunctionalityType = "externalOverrideFromGrid"
	StateInformationFunctionalityTypeAutonomousGridSupport    StateInformationFunctionalityType = "autonomousGridSupport"
	StateInformationFunctionalityTypeIslandingMode            StateInformationFunctionalityType = "islandingMode"
	StateInformationFunctionalityTypeBalancing                StateInformationFunctionalityType = "balancing"
	StateInformationFunctionalityTypeTrickleCharging          StateInformationFunctionalityType = "trickleCharging"
	StateInformationFunctionalityTypeCalibration              StateInformationFunctionalityType = "calibration"
	StateInformationFunctionalityTypeCommissioningMissing     StateInformationFunctionalityType = "commissioningMissing"
	StateInformationFunctionalityTypeSleeping                 StateInformationFunctionalityType = "sleeping"
	StateInformationFunctionalityTypeStarting                 StateInformationFunctionalityType = "starting"
	StateInformationFunctionalityTypeMppt                     StateInformationFunctionalityType = "mppt"
	StateInformationFunctionalityTypeThrottled                StateInformationFunctionalityType = "throttled"
	StateInformationFunctionalityTypeShuttingDown             StateInformationFunctionalityType = "shuttingDown"
	StateInformationFunctionalityTypeManualShutdown           StateInformationFunctionalityType = "manualShutdown"
)

type StateInformationFailureType string

const (
	StateInformationFailureTypeInverterDefective             StateInformationFailureType = "inverterDefective"
	StateInformationFailureTypeBatteryOvercurrentProtection  StateInformationFailureType = "batteryOvercurrentProtection"
	StateInformationFailureTypePvStringOvercurrentProtection StateInformationFailureType = "pvStringOvercurrentProtection"
	StateInformationFailureTypeGridFault                     StateInformationFailureType = "gridFault"
	StateInformationFailureTypeGroundFault                   StateInformationFailureType = "groundFault"
	StateInformationFailureTypeAcDisconnected                StateInformationFailureType = "acDisconnected"
	StateInformationFailureTypeDcDisconnected                StateInformationFailureType = "dcDisconnected"
	StateInformationFailureTypeCabinetOpen                   StateInformationFailureType = "cabinetOpen"
	StateInformationFailureTypeOverTemperature               StateInformationFailureType = "overTemperature"
	StateInformationFailureTypeUnderTemperature              StateInformationFailureType = "underTemperature"
	StateInformationFailureTypeFrequencyAboveLimit           StateInformationFailureType = "frequencyAboveLimit"
	StateInformationFailureTypeFrequencyBelowLimit           StateInformationFailureType = "frequencyBelowLimit"
	StateInformationFailureTypeAcVoltageAboveLimit           StateInformationFailureType = "acVoltageAboveLimit"
	StateInformationFailureTypeAcVoltageBelowLimit           StateInformationFailureType = "acVoltageBelowLimit"
	StateInformationFailureTypeDcVoltageAboveLimit           StateInformationFailureType = "dcVoltageAboveLimit"
	StateInformationFailureTypeDcVoltageBelowLimit           StateInformationFailureType = "dcVoltageBelowLimit"
	StateInformationFailureTypeHardwareTestFailure           StateInformationFailureType = "hardwareTestFailure"
	StateInformationFailureTypeGenericInternalError          StateInformationFailureType = "genericInternalError"
)

type StateInformationCategoryType string

const (
	StateInformationCategoryTypeFunctionality StateInformationCategoryType = "functionality"
	StateInformationCategoryTypeFailure       StateInformationCategoryType = "failure"
)

type StateInformationDataType struct {
	StateInformationId *StateInformationIdType       `json:"stateInformationId,omitempty" eebus:"key"`
	StateInformation   *StateInformationType         `json:"stateInformation,omitempty"`
	IsActive           *bool                         `json:"isActive,omitempty"`
	Category           *StateInformationCategoryType `json:"category,omitempty"`
	TimeOfLastChange   *AbsoluteOrRelativeTimeType   `json:"timeOfLastChange,omitempty"`
}

type StateInformationDataElementsType struct {
	StateInformationId *ElementTagType `json:"stateInformationId,omitempty"`
	StateInformation   *ElementTagType `json:"stateInformation,omitempty"`
	IsActive           *ElementTagType `json:"isActive,omitempty"`
	Category           *ElementTagType `json:"category,omitempty"`
	TimeOfLastChange   *ElementTagType `json:"timeOfLastChange,omitempty"`
}

type StateInformationListDataType struct {
	StateInformationData []StateInformationDataType `json:"stateInformationData,omitempty"`
}

type StateInformationListDataSelectorsType struct {
	StateInformationId *StateInformationIdType       `json:"stateInformationId,omitempty"`
	StateInformation   *StateInformationType         `json:"stateInformation,omitempty"`
	IsActive           *bool                         `json:"isActive,omitempty"`
	Category           *StateInformationCategoryType `json:"category,omitempty"`
}
