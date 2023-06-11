package spine

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/eebus-go/spine/model"
)

// a default mapping of a given EntityTypeType to a UseCaseActorType
var entityTypeActorMap = map[model.EntityTypeType]model.UseCaseActorType{
	model.EntityTypeTypeBattery:                       model.UseCaseActorTypeBattery,
	model.EntityTypeTypeCEM:                           model.UseCaseActorTypeCEM,
	model.EntityTypeTypeCompressor:                    model.UseCaseActorTypeCompressor,
	model.EntityTypeTypeElectricityStorageSystem:      model.UseCaseActorTypeBatterySystem,
	model.EntityTypeTypeElectricityGenerationSystem:   model.UseCaseActorTypePVSystem,
	model.EntityTypeTypeEV:                            model.UseCaseActorTypeEV,
	model.EntityTypeTypeEVSE:                          model.UseCaseActorTypeEVSE,
	model.EntityTypeTypeDHWCircuit:                    model.UseCaseActorTypeDHWCircuit,
	model.EntityTypeTypeHeatingCircuit:                model.UseCaseActorTypeHeatingCircuit,
	model.EntityTypeTypeHeatPumpAppliance:             model.UseCaseActorTypeHeatPump,
	model.EntityTypeTypeHvacRoom:                      model.UseCaseActorTypeHVACRoom,
	model.EntityTypeTypeInverter:                      model.UseCaseActorTypeInverter,
	model.EntityTypeTypeSmartEnergyAppliance:          model.UseCaseActorTypeControllableSystem,
	model.EntityTypeTypeSubMeterElectricity:           model.UseCaseActorTypeControllableSystem,
	model.EntityTypeTypeGridConnectionPointOfPremises: model.UseCaseActorTypeGridConnectionPoint,
}

// list of known use cases and the allowed actors for each
var useCaseValidActorsMap = map[model.UseCaseNameType][]model.UseCaseActorType{
	model.UseCaseNameTypeConfigurationOfDhwSystemFunction:                             {model.UseCaseActorTypeConfigurationAppliance, model.UseCaseActorTypeDHWCircuit},
	model.UseCaseNameTypeConfigurationOfDhwTemperature:                                {model.UseCaseActorTypeConfigurationAppliance, model.UseCaseActorTypeDHWCircuit},
	model.UseCaseNameTypeConfigurationOfRoomCoolingSystemFunction:                     {model.UseCaseActorTypeConfigurationAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeConfigurationOfRoomCoolingTemperature:                        {model.UseCaseActorTypeConfigurationAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeConfigurationOfRoomHeatingSystemFunction:                     {model.UseCaseActorTypeConfigurationAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeConfigurationOfRoomHeatingTemperature:                        {model.UseCaseActorTypeConfigurationAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeControlOfBattery:                                             {model.UseCaseActorTypeInverter, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeCoordinatedEVCharging:                                        {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM, model.UseCaseActorTypeEnergyBroker},
	model.UseCaseNameTypeEVChargingSummary:                                            {model.UseCaseActorTypeEVSE, model.UseCaseActorTypeCEM, model.UseCaseActorTypeEnergyBroker},
	model.UseCaseNameTypeEVCommissioningAndConfiguration:                              {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVSECommissioningAndConfiguration:                            {model.UseCaseActorTypeEVSE, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVStateOfCharge:                                              {model.UseCaseActorTypeEV, model.UseCaseActorTypeMonitoringAppliance},
	model.UseCaseNameTypeFlexibleLoad:                                                 {model.UseCaseActorTypeEnergyConsumer, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeFlexibleStartForWhiteGoods:                                   {model.UseCaseActorTypeSmartAppliance, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeIncentiveTableBasedPowerConsumptionManagement:                {model.UseCaseActorTypeEnergyConsumer, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeLimitationOfPowerConsumption:                                 {model.UseCaseActorTypeEnergyGuard, model.UseCaseActorTypeControllableSystem},
	model.UseCaseNameTypeLimitationOfPowerProduction:                                  {model.UseCaseActorTypeEnergyGuard, model.UseCaseActorTypeControllableSystem},
	model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging:                     {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeMonitoringAndControlOfSmartGridReadyConditions:               {model.UseCaseActorTypeHeatPump, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeMonitoringOfBattery:                                          {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeBattery},
	model.UseCaseNameTypeMonitoringOfDhwSystemFunction:                                {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeDHWCircuit},
	model.UseCaseNameTypeMonitoringOfDhwTemperature:                                   {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeDHWCircuit},
	model.UseCaseNameTypeMonitoringOfGridConnectionPoint:                              {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeGridConnectionPoint},
	model.UseCaseNameTypeMonitoringOfInverter:                                         {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeInverter},
	model.UseCaseNameTypeMonitoringOfOutdoorTemperature:                               {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeOutdoorTemperatureSensor},
	model.UseCaseNameTypeMonitoringOfPowerConsumption:                                 {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeMonitoredUnit},
	model.UseCaseNameTypeMonitoringOfPvString:                                         {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypePVString},
	model.UseCaseNameTypeMonitoringOfRoomCoolingSystemFunction:                        {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeMonitoringOfRoomHeatingSystemFunction:                        {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeMonitoringOfRoomTemperature:                                  {model.UseCaseActorTypeMonitoringAppliance, model.UseCaseActorTypeHVACRoom},
	model.UseCaseNameTypeOptimizationOfSelfConsumptionByHeatPumpCompressorFlexibility: {model.UseCaseActorTypeCompressor, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging:                {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment:             {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM, model.UseCaseActorTypeEnergyGuard},
	model.UseCaseNameTypeVisualizationOfAggregatedBatteryData:                         {model.UseCaseActorTypeVisualizationAppliance, model.UseCaseActorTypeBatterySystem},
	model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData:                    {model.UseCaseActorTypeVisualizationAppliance, model.UseCaseActorTypePVSystem},
	model.UseCaseNameTypeVisualizationOfHeatingAreaName:                               {model.UseCaseActorTypeVisualizationAppliance, model.UseCaseActorTypeHeatingCircuit, model.UseCaseActorTypeHeatingZone, model.UseCaseActorTypeHVACRoom},
}

// defines a specific usecase implementation
// right now this is just used as a wrapper for supported usecases
type UseCaseImpl struct {
	Entity *EntityLocalImpl
	Actor  model.UseCaseActorType

	name             model.UseCaseNameType
	useCaseVersion   model.SpecificationVersionType
	useCaseAvailable bool
	scenarioSupport  []model.UseCaseScenarioSupportType
}

// returns a UseCaseImpl with a default mapping of entity to actor using data
func NewUseCase(entity *EntityLocalImpl, ucEnumType model.UseCaseNameType, useCaseVersion model.SpecificationVersionType, useCaseAvailable bool, scenarioSupport []model.UseCaseScenarioSupportType) *UseCaseImpl {
	checkEntityArguments(*entity.EntityImpl)

	actor := entityTypeActorMap[entity.EntityType()]

	return NewUseCaseWithActor(entity, actor, ucEnumType, useCaseVersion, useCaseAvailable, scenarioSupport)
}

// returns a UseCaseImpl with specific entity and actor
func NewUseCaseWithActor(entity *EntityLocalImpl, actor model.UseCaseActorType, ucEnumType model.UseCaseNameType, useCaseVersion model.SpecificationVersionType, useCaseAvailable bool, scenarioSupport []model.UseCaseScenarioSupportType) *UseCaseImpl {
	checkUCArguments(actor, ucEnumType)

	ucManager := entity.Device().UseCaseManager()
	ucManager.Add(actor, ucEnumType, useCaseVersion, useCaseAvailable, scenarioSupport)

	return &UseCaseImpl{
		Entity:           entity,
		Actor:            actor,
		name:             model.UseCaseNameType(ucEnumType),
		useCaseVersion:   useCaseVersion,
		useCaseAvailable: useCaseAvailable,
		scenarioSupport:  scenarioSupport,
	}
}

// check if there is an predefined mapping available
func checkEntityArguments(entity EntityImpl) {
	actor := entityTypeActorMap[entity.EntityType()]
	if actor == "" {
		panic(fmt.Errorf("cannot derive actor for entity type '%s'", entity.EntityType()))
	}
}

// check if the actor is valid for the given usecase type
func checkUCArguments(actor model.UseCaseActorType, ucEnumType model.UseCaseNameType) {
	if !linq.From(useCaseValidActorsMap[ucEnumType]).Contains(actor) {
		panic(fmt.Errorf("the actor '%s' is not valid for the use case '%s'", actor, ucEnumType))
	}
}

// Update the availability of this usecase and
// trigger a notification being sent to the remote device
func (u *UseCaseImpl) SetUseCaseAvailable(available bool) {
	u.useCaseAvailable = available

	u.Entity.Device().NotifyUseCaseData()
}

/*
// This is not yet used, might be removed?
func waitForRequest[T any](c chan T, maxDelay time.Duration) *T {
	timeout := time.After(maxDelay)

	select {
	case data := <-c:
		return &data
	case <-timeout:
		return nil
	}
}
*/
