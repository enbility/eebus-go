package oscev

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-oscev-UseCaseSupportUpdate"

	// EV current limits
	//
	// Use `CurrentLimits` to get the current data
	DataUpdateCurrentLimits api.EventType = "cem-oscev-DataUpdateCurrentLimits"

	// EV load control recommendation limit data updated
	//
	// Use `LoadControlLimits` to get the current data
	//
	// Use Case OSCEV, Scenario 1
	DataUpdateLimit api.EventType = "cem-oscev-DataUpdateLimit"
)
