package model

type UseCaseActorType string

const (
	UseCaseActorTypeCEM  UseCaseActorType = "CEM"
	UseCaseActorTypeEVSE UseCaseActorType = "EVSE"
	UseCaseActorTypeEV   UseCaseActorType = "EV"
)

type UseCaseNameType string

const (
	UseCaseNameTypeMeasurementOfElectricityDuringEVCharging         UseCaseNameType = "measurementOfElectricityDuringEvCharging"
	UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging    UseCaseNameType = "optimizationOfSelfConsumptionDuringEvCharging"
	UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment UseCaseNameType = "overloadProtectionByEvChargingCurrentCurtailment"
	UseCaseNameTypeCoordinatedEVCharging                            UseCaseNameType = "coordinatedEvCharging"
	UseCaseNameTypeEVCommissioningAndConfiguration                  UseCaseNameType = "evCommissioningAndConfiguration"
	UseCaseNameTypeEVSECommissioningAndConfiguration                UseCaseNameType = "evseCommissioningAndConfiguration"
	UseCaseNameTypeEVChargingSummary                                UseCaseNameType = "evChargingSummary"
	UseCaseNameTypeEVStateOfCharge                                  UseCaseNameType = "evStateOfCharge"
	UseCaseNameTypeMonitoringAndControlOfSmartGridReadyConditions   UseCaseNameType = "monitoringAndControlOfSmartGridReadyConditions"
	UseCaseNameTypeMonitoringOfPowerConsumption                     UseCaseNameType = "monitoringOfPowerConsumption"
)

type UseCaseScenarioSupportType uint

type UseCaseSupportType struct {
	UseCaseName      *UseCaseNameType             `json:"useCaseName,omitempty"`
	UseCaseVersion   *SpecificationVersionType    `json:"useCaseVersion,omitempty"`
	UseCaseAvailable *bool                        `json:"useCaseAvailable,omitempty"`
	ScenarioSupport  []UseCaseScenarioSupportType `json:"scenarioSupport,omitempty"`
}

type UseCaseSupportElementsType struct {
	UseCaseName      *ElementTagType `json:"useCaseName,omitempty"`
	UseCaseVersion   *ElementTagType `json:"useCaseVersion,omitempty"`
	UseCaseAvailable *ElementTagType `json:"useCaseAvailable,omitempty"`
	ScenarioSupport  *ElementTagType `json:"scenarioSupport,omitempty"`
}

type UseCaseSupportSelectorsType struct {
	UseCaseName     *UseCaseNameType            `json:"useCaseName,omitempty"`
	UseCaseVersion  *SpecificationVersionType   `json:"useCaseVersion,omitempty"`
	ScenarioSupport *UseCaseScenarioSupportType `json:"scenarioSupport,omitempty"`
}

type UseCaseInformationDataType struct {
	Address        *FeatureAddressType  `json:"address,omitempty"`
	Actor          *UseCaseActorType    `json:"actor,omitempty"`
	UseCaseSupport []UseCaseSupportType `json:"useCaseSupport,omitempty"`
}

type UseCaseInformationDataElementsType struct {
	Address        *ElementTagType `json:"address,omitempty"`
	Actor          *ElementTagType `json:"actor,omitempty"`
	UseCaseSupport *ElementTagType `json:"useCaseSupport,omitempty"`
}

type UseCaseInformationListDataType struct {
	UseCaseInformationData []UseCaseInformationDataType `json:"useCaseInformationData,omitempty"`
}

type UseCaseInformationListDataSelectorsType struct {
	Address        *FeatureAddressType          `json:"address,omitempty"`
	Actor          *UseCaseActorType            `json:"actor,omitempty"`
	UseCaseSupport *UseCaseSupportSelectorsType `json:"useCaseSupport,omitempty"`
}
