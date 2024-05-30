package usecase

import (
	"slices"
	"sync"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type UseCaseBase struct {
	LocalEntity    spineapi.EntityLocalInterface
	remoteEntities map[spineapi.EntityRemoteInterface][]model.UseCaseScenarioSupportType

	UseCaseActor              model.UseCaseActorType
	UseCaseName               model.UseCaseNameType
	useCaseVersion            model.SpecificationVersionType
	useCaseDocumentSubVersion string
	useCaseScenarios          []model.UseCaseScenarioSupportType

	EventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType

	mux sync.Mutex
}

var _ api.UseCaseBaseInterface = (*UseCaseBase)(nil)

func NewUseCaseBase(
	localEntity spineapi.EntityLocalInterface,
	usecaseActor model.UseCaseActorType,
	usecaseName model.UseCaseNameType,
	useCaseVersion string,
	useCaseDocumentSubVersion string,
	useCaseScenarios []model.UseCaseScenarioSupportType,
	eventCB api.EntityEventCallback,
	validEntityTypes []model.EntityTypeType,
) *UseCaseBase {
	return &UseCaseBase{
		LocalEntity:               localEntity,
		UseCaseActor:              usecaseActor,
		UseCaseName:               usecaseName,
		useCaseVersion:            model.SpecificationVersionType(useCaseVersion),
		useCaseDocumentSubVersion: useCaseDocumentSubVersion,
		useCaseScenarios:          useCaseScenarios,
		EventCB:                   eventCB,
		validEntityTypes:          validEntityTypes,
		remoteEntities:            make(map[spineapi.EntityRemoteInterface][]model.UseCaseScenarioSupportType),
	}
}

func (u *UseCaseBase) AddUseCase() {
	u.LocalEntity.AddUseCaseSupport(
		u.UseCaseActor,
		u.UseCaseName,
		u.useCaseVersion,
		u.useCaseDocumentSubVersion,
		true,
		u.useCaseScenarios)
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

func (u *UseCaseBase) SupportedUseCaseScenarios(
	entity spineapi.EntityRemoteInterface,
) []model.UseCaseScenarioSupportType {
	if entity == nil ||
		entity.Device() == nil {
		return nil
	}

	ucs := entity.Device().UseCases()
	for _, uc := range ucs {
		// check if the use case entity address is identical to the entity address
		// the address may not exist, as it only available since SPINE 1.3
		if uc.Address != nil &&
			entity.Address() != nil &&
			slices.Compare(uc.Address.Entity, entity.Address().Entity) != 0 {
			continue
		}

		for _, support := range uc.UseCaseSupport {
			if support.UseCaseName == nil ||
				*support.UseCaseName != u.UseCaseName {
				continue
			}

			return support.ScenarioSupport
		}
	}

	return nil
}

func (u *UseCaseBase) HasSupportForUseCaseScenarios(
	entity spineapi.EntityRemoteInterface,
	scenarios []model.UseCaseScenarioSupportType,
) bool {
	if entity == nil ||
		entity.Device() == nil ||
		len(scenarios) == 0 {
		return false
	}

	ucs := entity.Device().UseCases()
	for _, uc := range ucs {
		// check if the use case entity address is identical to the entity address
		// the address may not exist, as it only available since SPINE 1.3
		if uc.Address != nil &&
			entity.Address() != nil &&
			slices.Compare(uc.Address.Entity, entity.Address().Entity) != 0 {
			continue
		}

		for _, support := range uc.UseCaseSupport {
			if support.UseCaseName == nil ||
				*support.UseCaseName != u.UseCaseName ||
				(support.UseCaseAvailable != nil && !*support.UseCaseAvailable) {
				continue
			}

			allFound := true
			for _, scenario := range scenarios {
				if !slices.Contains(support.ScenarioSupport, scenario) {
					allFound = false
					break
				}
			}
			if allFound {
				return true
			}
		}
	}

	return false
}

func (u *UseCaseBase) UseCaseDataUpdate(
	payload spineapi.EventPayload,
	eventCB api.EntityEventCallback,
	event api.EventType,
) {
	// entity was removed, so remove it from the list
	if internal.IsEntityDisconnected(payload) {
		if u.hasRemoteEntity(payload.Entity) {
			u.removeRemoteEntity(payload.Entity)
		}

		eventCB(payload.Ski, payload.Device, payload.Entity, event)

		return
	}

	// entity updated usecase data
	scenarios := u.SupportedUseCaseScenarios(payload.Entity)
	if scenarios != nil {
		curScenarios := u.scenariosForRemoteEntity(payload.Entity)
		if slices.Compare(scenarios, curScenarios) != 0 {
			u.setRemoteEntityScenarios(payload.Entity, scenarios)
			eventCB(payload.Ski, payload.Device, payload.Entity, event)
		}
	} else {
		// entity does not support the use case, maybe support was removed
		u.removeRemoteEntity(payload.Entity)
		eventCB(payload.Ski, payload.Device, payload.Entity, event)
	}
}

// return the current list of compatible remote entities and their scenarios
func (u *UseCaseBase) RemoteEntities() []api.RemoteEntityScenarios {
	u.mux.Lock()
	defer u.mux.Unlock()

	entities := make([]api.RemoteEntityScenarios, 0, len(u.remoteEntities))
	for entity, scenarios := range u.remoteEntities {
		newItem := api.RemoteEntityScenarios{
			Entity:    entity,
			Scenarios: scenarios,
		}
		entities = append(entities, newItem)
	}

	return entities
}

// set the scenarios of a remote entity
func (u *UseCaseBase) setRemoteEntityScenarios(
	entity spineapi.EntityRemoteInterface,
	scenarios []model.UseCaseScenarioSupportType,
) {
	u.mux.Lock()
	defer u.mux.Unlock()

	u.remoteEntities[entity] = scenarios
}

// check if the entity is already added
func (u *UseCaseBase) hasRemoteEntity(entity spineapi.EntityRemoteInterface) bool {
	u.mux.Lock()
	defer u.mux.Unlock()

	_, ok := u.remoteEntities[entity]
	return ok
}

// remove a remote entity from the use case
func (u *UseCaseBase) removeRemoteEntity(entity spineapi.EntityRemoteInterface) {
	if !u.hasRemoteEntity(entity) {
		return
	}

	u.mux.Lock()
	defer u.mux.Unlock()

	delete(u.remoteEntities, entity)
}

// check if the entity is already added as being compatible with the use case
func (u *UseCaseBase) scenariosForRemoteEntity(entity spineapi.EntityRemoteInterface) []model.UseCaseScenarioSupportType {
	u.mux.Lock()
	defer u.mux.Unlock()

	if value, ok := u.remoteEntities[entity]; ok {
		return value
	}

	return nil
}
