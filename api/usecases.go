package api

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// details about each use case scenario
type UseCaseScenario struct {
	Scenario       model.UseCaseScenarioSupportType // the scenario number
	Mandatory      bool                             // if this scenario is mandatory to be supported by the remote entity
	ServerFeatures []model.FeatureTypeType          // the server features required for this scenario on the remote entity
}

// contains the available scenarios of a remote entity
type RemoteEntityScenarios struct {
	Entity    spineapi.EntityRemoteInterface
	Scenarios []uint
}

// Entity event callback
//
// Used by Use Case implementations
type EntityEventCallback func(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event EventType)

type UseCaseBaseInterface interface {
	// add the use case
	AddUseCase()

	// remove the use case
	RemoveUseCase()

	// update availability of the use case
	//
	// NOTE: only allowed to be used for client side implementations
	// of a use case! Otherwise use `RemoveUseCase` and `AddUseCase`.
	UpdateUseCaseAvailability(available bool)

	// check if the entity is compatible with the use case
	IsCompatibleEntityType(entity spineapi.EntityRemoteInterface) bool

	// return the current list of compatible remote entities and their available scenarios of this use case
	RemoteEntitiesScenarios() []RemoteEntityScenarios

	// return the current list of available scenarios of this use case for the remote entity
	AvailableScenariosForEntity(entity spineapi.EntityRemoteInterface) []uint

	// check if the provided scenario are available at the remote entity
	IsScenarioAvailableAtEntity(
		entity spineapi.EntityRemoteInterface,
		scenario uint,
	) bool
}

// Implemented by each Use Case
type UseCaseInterface interface {
	UseCaseBaseInterface

	// add the features
	AddFeatures()
}
