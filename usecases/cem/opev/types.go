package opev

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-opev-UseCaseSupportUpdate"

	// EV current limits
	//
	// Use `CurrentLimits` to get the current data
	DataUpdateCurrentLimits api.EventType = "cem-opev-DataUpdateCurrentLimits"

	// EV load control obligation limit data updated
	//
	// Use `LoadControlLimits` to get the current data
	DataUpdateLimit api.EventType = "cem-opev-DataUpdateLimit"
)
