package model

type UseCaseActorType string

const (
	UseCaseActorTypeBatterySystem          UseCaseActorType = "BatterySystem"
	UseCaseActorTypeCEM                    UseCaseActorType = "CEM"
	UseCaseActorTypeConfigurationAppliance UseCaseActorType = "ConfigurationAppliance"
	UseCaseActorTypeCompressor             UseCaseActorType = "Compressor"
	UseCaseActorTypeControllableSystem     UseCaseActorType = "ControllableSystem"
	UseCaseActorTypeDHWCircuit             UseCaseActorType = "DHWCircuit"
	UseCaseActorTypeEnergyGuard            UseCaseActorType = "EnergyGuard"
	UseCaseActorTypeEVSE                   UseCaseActorType = "EVSE"
	UseCaseActorTypeEV                     UseCaseActorType = "EV"
	UseCaseActorTypeGridConnectionPoint    UseCaseActorType = "GridConnectionPoint"
	UseCaseActorTypeHeatPump               UseCaseActorType = "HeatPump"
	UseCaseActorTypeHeatingCircuit         UseCaseActorType = "HeatingCircuit"
	UseCaseActorTypeHeatingZone            UseCaseActorType = "HeatingZone"
	UseCaseActorTypeHVACRoom               UseCaseActorType = "HVACRoom"
	UseCaseActorTypeMonitoredUnit          UseCaseActorType = "MonitoredUnit"
	UseCaseActorTypeMonitoringAppliance    UseCaseActorType = "MonitoringAppliance"
	UseCaseActorTypePVSystem               UseCaseActorType = "PVSystem"
	UseCaseActorTypeVisualizationAppliance UseCaseActorType = "VisualizationAppliance"
)

type UseCaseNameType string

const (
	UseCaseNameTypeConfigurationOfDhwSystemFunction                             UseCaseNameType = "configurationOfDhwSystemFunction"
	UseCaseNameTypeConfigurationOfDhwTemperature                                UseCaseNameType = "configurationOfDhwTemperature"
	UseCaseNameTypeConfigurationOfRoomCoolingSystemFunction                     UseCaseNameType = "configurationOfRoomCoolingSystemFunction"
	UseCaseNameTypeConfigurationOfRoomCoolingTemperature                        UseCaseNameType = "configurationOfRoomCoolingTemperature"
	UseCaseNameTypeConfigurationOfRoomHeatingSystemFunction                     UseCaseNameType = "configurationOfRoomHeatingSystemFunction"
	UseCaseNameTypeConfigurationOfRoomHeatingTemperature                        UseCaseNameType = "configurationOfRoomHeatingTemperature"
	UseCaseNameTypeControlOfBattery                                             UseCaseNameType = "controlOfBattery"
	UseCaseNameTypeCoordinatedEVCharging                                        UseCaseNameType = "coordinatedEvCharging"
	UseCaseNameTypeEVChargingSummary                                            UseCaseNameType = "evChargingSummary"
	UseCaseNameTypeEVCommissioningAndConfiguration                              UseCaseNameType = "evCommissioningAndConfiguration"
	UseCaseNameTypeEVSECommissioningAndConfiguration                            UseCaseNameType = "evseCommissioningAndConfiguration"
	UseCaseNameTypeEVStateOfCharge                                              UseCaseNameType = "evStateOfCharge"
	UseCaseNameTypeLimitationOfPowerConsumption                                 UseCaseNameType = "limitationOfPowerConsumption"
	UseCaseNameTypeLimitationOfPowerProduction                                  UseCaseNameType = "limitationOfPowerProduction"
	UseCaseNameTypeIncentiveTableBasedPowerConsumptionManagement                UseCaseNameType = "incentiveTableBasedPowerConsumptionManagement"
	UseCaseNameTypeMeasurementOfElectricityDuringEVCharging                     UseCaseNameType = "measurementOfElectricityDuringEvCharging"
	UseCaseNameTypeMonitoringAndControlOfSmartGridReadyConditions               UseCaseNameType = "monitoringAndControlOfSmartGridReadyConditions"
	UseCaseNameTypeMonitoringOfBattery                                          UseCaseNameType = "monitoringOfBattery"
	UseCaseNameTypeMonitoringOfDhwSystemFunction                                UseCaseNameType = "monitoringOfDhwSystemFunction"
	UseCaseNameTypeMonitoringOfDhwTemperature                                   UseCaseNameType = "monitoringOfDhwTemperature"
	UseCaseNameTypeMonitoringOfGridConnectionPoint                              UseCaseNameType = "monitoringOfGridConnectionPoint"
	UseCaseNameTypeMonitoringOfInverter                                         UseCaseNameType = "monitoringOfInverter"
	UseCaseNameTypeMonitoringOfOutdoorTemperature                               UseCaseNameType = "monitoringOfOutdoorTemperature"
	UseCaseNameTypeMonitoringOfPowerConsumption                                 UseCaseNameType = "monitoringOfPowerConsumption"
	UseCaseNameTypeMonitoringOfPvString                                         UseCaseNameType = "monitoringOfPvString"
	UseCaseNameTypeMonitoringOfRoomCoolingSystemFunction                        UseCaseNameType = "monitoringOfRoomCoolingSystemFunction"
	UseCaseNameTypeMonitoringOfRoomHeatingSystemFunction                        UseCaseNameType = "monitoringOfRoomHeatingSystemFunction"
	UseCaseNameTypeMonitoringOfRoomTemperature                                  UseCaseNameType = "monitoringOfRoomTemperature"
	UseCaseNameTypeOptimizationOfSelfConsumptionByHeatPumpCompressorFlexibility UseCaseNameType = "optimizationOfSelfConsumptionByHeatPumpCompressorFlexibility"
	UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging                UseCaseNameType = "optimizationOfSelfConsumptionDuringEvCharging"
	UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment             UseCaseNameType = "overloadProtectionByEvChargingCurrentCurtailment"
	UseCaseNameTypeVisualizationOfAggregatedBatteryData                         UseCaseNameType = "visualizationOfAggregatedBatteryData"
	UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData                    UseCaseNameType = "visualizationOfAggregatedPhotovoltaicData"
	UseCaseNameTypeVisualizationOfHeatingAreaName                               UseCaseNameType = "visualizationOfHeatingAreaName"
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
