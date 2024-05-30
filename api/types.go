package api

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// type for cem and usecase specfic event names
type EventType string

type RemoteEntityScenarios struct {
	Entity    spineapi.EntityRemoteInterface
	Scenarios []model.UseCaseScenarioSupportType
}
