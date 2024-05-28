package opev

import "github.com/enbility/eebus-go/api"

const (
	// EV current limits
	//
	// Use `CurrentLimits` to get the current data
	DataUpdateCurrentLimits api.EventType = "cem-opev-DataUpdateCurrentLimits"

	// EV load control obligation limit data updated
	//
	// Use `LoadControlLimits` to get the current data
	DataUpdateLimit api.EventType = "cem-opev-DataUpdateLimit"
)
