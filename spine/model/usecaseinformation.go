package model

type UseCaseActorType UseCaseActorEnumType

type UseCaseActorEnumType string

const (
	UseCaseActorEnumTypeEV UseCaseActorEnumType = "EV"
)

type UseCaseNameType UseCaseNameEnumType

type UseCaseNameEnumType string

const (
	UseCaseNameEnumTypeMeasurementOfElectricityDuringEVCharging         UseCaseNameEnumType = "measurementOfElectricityDuringEvCharging"
	UseCaseNameEnumTypeOptimizationOfSelfConsumptionDuringEVCharging    UseCaseNameEnumType = "optimizationOfSelfConsumptionDuringEvCharging"
	UseCaseNameEnumTypeOverloadProtectionByEVChargingCurrentCurtailment UseCaseNameEnumType = "overloadProtectionByEvChargingCurrentCurtailment"
	UseCaseNameEnumTypeCoordinatedEVCharging                            UseCaseNameEnumType = "coordinatedEvCharging"
	UseCaseNameEnumTypeEVCommissioningAndConfiguration                  UseCaseNameEnumType = "evCommissioningAndConfiguration"
	UseCaseNameEnumTypeEVSECommissioningAndConfiguration                UseCaseNameEnumType = "evseCommissioningAndConfiguration"
	UseCaseNameEnumTypeEVChargingSummary                                UseCaseNameEnumType = "evChargingSummary"
	UseCaseNameEnumTypeEVStateOfCharge                                  UseCaseNameEnumType = "evStateOfCharge"
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
