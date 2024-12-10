package usecase

import (
	"reflect"
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

// Adds a new use case to an entity
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - usecaseActor: The actor type of the use case
//   - usecaseName: The name of the use case
//   - useCaseVersion: The version of the use case
//   - useCaseDocumentSubVersion: The sub version of the use case document
//   - useCaseScenarios: The supported scenarios of the use case
//   - eventCB: The callback to be called when an usecase update event of a remote entity is triggered (optional, can be nil)
//   - useCaseUpdateEvent: The event type of the use case update event for the eventCB
//   - validActorTypes: The valid actor types for the use case in a remote entity
//   - validEntityTypes: The valid entity types for the use case in a remote entity
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
	u.LocalEntity.RemoveUseCaseSupports(
		[]model.UseCaseFilterType{
			{
				Actor:       u.UseCaseActor,
				UseCaseName: u.UseCaseName,
			},
		})
}

func (u *UseCaseBase) UpdateUseCaseAvailability(available bool) {
	u.LocalEntity.SetUseCaseAvailability(
		model.UseCaseFilterType{
			Actor:       u.UseCaseActor,
			UseCaseName: u.UseCaseName,
		}, available)
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
	_, scenarios := u.entitiyScenarioIndexOfEntity(entity)

	return scenarios
}

// check if the provided scenario are available at the remote entity
func (u *UseCaseBase) IsScenarioAvailableAtEntity(
	entity spineapi.EntityRemoteInterface,
	scenario uint,
) bool {
	if _, scenarios := u.entitiyScenarioIndexOfEntity(entity); scenarios != nil {
		return slices.Contains(scenarios, scenario)
	}

	return false
}

// return the indices of all entities of the device in the available entity scenarios
func (u *UseCaseBase) entityScenarioIndicesOfDevice(device spineapi.DeviceRemoteInterface) []int {
	u.mux.Lock()
	defer u.mux.Unlock()

	indices := []int{}

	for i, remoteEntityScenarios := range u.availableEntityScenarios {
		if device != nil && device.Address() != nil &&
			remoteEntityScenarios.Entity != nil &&
			remoteEntityScenarios.Entity.Device() != nil &&
			remoteEntityScenarios.Entity.Device().Address() != nil &&
			reflect.DeepEqual(device.Address(), remoteEntityScenarios.Entity.Device().Address()) {
			indices = append(indices, i)
		}
	}

	return indices
}

// return the index and the scenarios of the entity in the available entity scenarios
// and return -1 and nil if not found
func (u *UseCaseBase) entitiyScenarioIndexOfEntity(entity spineapi.EntityRemoteInterface) (int, []uint) {
	u.mux.Lock()
	defer u.mux.Unlock()

	for i, remoteEntityScenarios := range u.availableEntityScenarios {
		if entity != nil && entity.Address() != nil &&
			remoteEntityScenarios.Entity != nil && remoteEntityScenarios.Entity.Address() != nil &&
			reflect.DeepEqual(entity.Address().Device, remoteEntityScenarios.Entity.Address().Device) &&
			reflect.DeepEqual(entity.Address().Entity, remoteEntityScenarios.Entity.Address().Entity) {
			return i, remoteEntityScenarios.Scenarios
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

	i, _ := u.entitiyScenarioIndexOfEntity(entity)
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

	if updateEvent && u.EventCB != nil {
		u.EventCB(entity.Device().Ski(), entity.Device(), entity, u.useCaseUpdateEvent)
	}
}

// remove all remote entities of a device from the use case
func (u *UseCaseBase) removeDeviceFromAvailableEntityScenarios(device spineapi.DeviceRemoteInterface) {
	indicies := u.entityScenarioIndicesOfDevice(device)

	for _, i := range indicies {
		u.removeEntityIndexFromAvailableEntityScenarios(i)
	}

	if u.EventCB != nil && len(indicies) > 0 {
		u.EventCB(device.Ski(), device, nil, u.useCaseUpdateEvent)
	}
}

// remove a remote entity from the use case
func (u *UseCaseBase) removeEntityFromAvailableEntityScenarios(entity spineapi.EntityRemoteInterface) {
	if i, _ := u.entitiyScenarioIndexOfEntity(entity); i >= 0 {
		u.removeEntityIndexFromAvailableEntityScenarios(i)

		if u.EventCB != nil {
			remoteDevice := entity.Device()
			u.EventCB(remoteDevice.Ski(), remoteDevice, entity, u.useCaseUpdateEvent)
		}
	}
}

// do the actual removal of the entity from the available entity scenarios
func (u *UseCaseBase) removeEntityIndexFromAvailableEntityScenarios(index int) {
	u.mux.Lock()
	u.availableEntityScenarios = append(u.availableEntityScenarios[:index], u.availableEntityScenarios[index+1:]...)
	u.mux.Unlock()
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
