package usecase

import (
	"slices"
	"sync"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UseCaseBase struct {
	LocalEntity spineapi.EntityLocalInterface

	UseCaseActor              model.UseCaseActorType
	UseCaseName               model.UseCaseNameType
	useCaseVersion            model.SpecificationVersionType
	useCaseDocumentSubVersion string
	useCaseScenarios          []api.UseCaseScenario

	EventCB            api.EntityEventCallback
	useCaseUpdateEvent api.EventType

	availableEntityScenarios []api.RemoteEntityScenarios // map of scenarios and their availability for each compatible remote entity

	validActorTypes  []model.UseCaseActorType // valid remote actor types for this use case
	validEntityTypes []model.EntityTypeType   // valid remote entity types for this use case

	mux sync.Mutex
}

var _ api.UseCaseBaseInterface = (*UseCaseBase)(nil)

func NewUseCaseBase(
	localEntity spineapi.EntityLocalInterface,
	usecaseActor model.UseCaseActorType,
	usecaseName model.UseCaseNameType,
	useCaseVersion string,
	useCaseDocumentSubVersion string,
	useCaseScenarios []api.UseCaseScenario,
	eventCB api.EntityEventCallback,
	useCaseUpdateEvent api.EventType,
	validActorTypes []model.UseCaseActorType,
	validEntityTypes []model.EntityTypeType,
) *UseCaseBase {
	ucb := &UseCaseBase{
		LocalEntity:               localEntity,
		UseCaseActor:              usecaseActor,
		UseCaseName:               usecaseName,
		useCaseVersion:            model.SpecificationVersionType(useCaseVersion),
		useCaseDocumentSubVersion: useCaseDocumentSubVersion,
		useCaseScenarios:          useCaseScenarios,
		EventCB:                   eventCB,
		useCaseUpdateEvent:        useCaseUpdateEvent,
		validActorTypes:           validActorTypes,
		validEntityTypes:          validEntityTypes,
	}

	_ = spine.Events.Subscribe(ucb)

	return ucb
}

func (u *UseCaseBase) AddUseCase() {
	useCaseScenarios := []model.UseCaseScenarioSupportType{}
	for _, scenario := range u.useCaseScenarios {
		useCaseScenarios = append(useCaseScenarios, scenario.Scenario)
	}

	u.LocalEntity.AddUseCaseSupport(
		u.UseCaseActor,
		u.UseCaseName,
		u.useCaseVersion,
		u.useCaseDocumentSubVersion,
		true,
		useCaseScenarios)
}

func (u *UseCaseBase) RemoveUseCase() {
	u.LocalEntity.RemoveUseCaseSupport(u.UseCaseActor, u.UseCaseName)
}

func (u *UseCaseBase) UpdateUseCaseAvailability(available bool) {
	u.LocalEntity.SetUseCaseAvailability(u.UseCaseActor, u.UseCaseName, available)
}

func (u *UseCaseBase) IsCompatibleEntityType(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil {
		return false
	}

	return slices.Contains(u.validEntityTypes, entity.EntityType())
}

// return the current list of compatible remote entities and their scenarios
func (u *UseCaseBase) RemoteEntitiesScenarios() []api.RemoteEntityScenarios {
	u.mux.Lock()
	defer u.mux.Unlock()

	return u.availableEntityScenarios
}

// return the currently available scenarios for the use case for a remote entity
func (u *UseCaseBase) AvailableScenariosForEntity(entity spineapi.EntityRemoteInterface) []uint {
	_, scenarios := u.indexAndScenariosOfEntity(entity)

	return scenarios
}

// check if the provided scenario are available at the remote entity
func (u *UseCaseBase) IsScenarioAvailableAtEntity(
	entity spineapi.EntityRemoteInterface,
	scenario uint,
) bool {
	if _, scenarios := u.indexAndScenariosOfEntity(entity); scenarios != nil {
		return slices.Contains(scenarios, scenario)
	}

	return false
}

// return the index and the scenarios of the entity in the available entity scenarios
// and return -1 and nil if not found
func (u *UseCaseBase) indexAndScenariosOfEntity(entity spineapi.EntityRemoteInterface) (int, []uint) {
	u.mux.Lock()
	defer u.mux.Unlock()

	for i, remoteEntity := range u.availableEntityScenarios {
		if entity == remoteEntity.Entity {
			return i, remoteEntity.Scenarios
		}
	}

	return -1, nil
}

// set the scenarios of a remote entity
func (u *UseCaseBase) updateRemoteEntityScenarios(
	entity spineapi.EntityRemoteInterface,
	scenarios []model.UseCaseScenarioSupportType,
) {
	updateEvent := false

	scenarioValues := []uint{}
	for _, scenario := range scenarios {
		scenarioValues = append(scenarioValues, uint(scenario))
	}

	i, _ := u.indexAndScenariosOfEntity(entity)
	if i == -1 {
		newItem := api.RemoteEntityScenarios{
			Entity:    entity,
			Scenarios: scenarioValues,
		}

		u.mux.Lock()
		u.availableEntityScenarios = append(u.availableEntityScenarios, newItem)
		u.mux.Unlock()

		updateEvent = true
	} else if i >= 0 && slices.Compare(u.availableEntityScenarios[i].Scenarios, scenarioValues) != 0 {
		u.mux.Lock()
		u.availableEntityScenarios[i].Scenarios = scenarioValues
		u.mux.Unlock()

		updateEvent = true
	}

	if updateEvent {
		u.EventCB(entity.Device().Ski(), entity.Device(), entity, u.useCaseUpdateEvent)
	}
}

// remove a remote entity from the use case
func (u *UseCaseBase) removeEntityFromAvailableEntityScenarios(entity spineapi.EntityRemoteInterface) {
	if i, _ := u.indexAndScenariosOfEntity(entity); i >= 0 {
		u.mux.Lock()
		u.availableEntityScenarios = append(u.availableEntityScenarios[:i], u.availableEntityScenarios[i+1:]...)
		u.mux.Unlock()

		u.EventCB(entity.Device().Ski(), entity.Device(), entity, u.useCaseUpdateEvent)
	}
}

// return the required server features for a use case scenario
func (u *UseCaseBase) requiredServerFeaturesForScenario(scenario model.UseCaseScenarioSupportType) []model.FeatureTypeType {
	for _, serverFeatures := range u.useCaseScenarios {
		if serverFeatures.Scenario == scenario {
			return serverFeatures.ServerFeatures
		}
	}

	return nil
}
