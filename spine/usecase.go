package spine

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/eebus-go/spine/model"
)

var entityTypeActorMap = map[model.EntityTypeType]model.UseCaseActorType{
	model.EntityTypeTypeEV:                            model.UseCaseActorTypeEV,
	model.EntityTypeTypeEVSE:                          model.UseCaseActorTypeEVSE,
	model.EntityTypeTypeCEM:                           model.UseCaseActorTypeCEM,
	model.EntityTypeTypeGridConnectionPointOfPremises: model.UseCaseActorTypeMonitoringAppliance,
	model.EntityTypeTypeElectricityStorageSystem:      model.UseCaseActorTypeBatterySystem,
	model.EntityTypeTypeElectricityGenerationSystem:   model.UseCaseActorTypePVSystem,
}

var useCaseValidActorsMap = map[model.UseCaseNameType][]model.UseCaseActorType{
	model.UseCaseNameTypeCoordinatedEVCharging:                            {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVSECommissioningAndConfiguration:                {model.UseCaseActorTypeEVSE, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVChargingSummary:                                {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVCommissioningAndConfiguration:                  {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVStateOfCharge:                                  {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging:         {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging:    {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment: {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeMonitoringOfPowerConsumption:                     {model.UseCaseActorTypeCEM, model.UseCaseActorTypeHeatPump},
	model.UseCaseNameTypeMonitoringAndControlOfSmartGridReadyConditions:   {model.UseCaseActorTypeCEM, model.UseCaseActorTypeHeatPump},
	model.UseCaseNameTypeMonitoringOfGridConnectionPoint:                  {model.UseCaseActorTypeCEM, model.UseCaseActorTypeMonitoringAppliance},
	model.UseCaseNameTypeVisualizationOfAggregatedBatteryData:             {model.UseCaseActorTypeCEM, model.UseCaseActorTypeBatterySystem, model.UseCaseActorTypeVisualizationAppliance},
	model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData:        {model.UseCaseActorTypeCEM, model.UseCaseActorTypePVSystem, model.UseCaseActorTypeVisualizationAppliance},
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

// returns a UseCaseImpl with a default mapping of entity to actor
func NewUseCase(entity *EntityLocalImpl, ucEnumType model.UseCaseNameType, useCaseVersion model.SpecificationVersionType, useCaseAvailable bool, scenarioSupport []model.UseCaseScenarioSupportType) *UseCaseImpl {
	actor := entityTypeActorMap[entity.EntityType()]

	return NewUseCaseWithActor(entity, actor, ucEnumType, useCaseVersion, useCaseAvailable, scenarioSupport)
}

// returns a UseCaseImpl with specific entity and actor
func NewUseCaseWithActor(entity *EntityLocalImpl, actor model.UseCaseActorType, ucEnumType model.UseCaseNameType, useCaseVersion model.SpecificationVersionType, useCaseAvailable bool, scenarioSupport []model.UseCaseScenarioSupportType) *UseCaseImpl {
	checkArguments(*entity.EntityImpl, ucEnumType)

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

func checkArguments(entity EntityImpl, ucEnumType model.UseCaseNameType) {
	actor := entityTypeActorMap[entity.EntityType()]
	if actor == "" {
		panic(fmt.Errorf("cannot derive actor for entity type '%s'", entity.EntityType()))
	}

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
