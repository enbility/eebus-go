package oscev

import "github.com/enbility/eebus-go/api"

const (
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
