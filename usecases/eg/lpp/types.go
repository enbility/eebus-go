package lpp

import "github.com/enbility/eebus-go/api"

const (
	// Load control obligation limit data updated
	//
	// Use `ProductionLimit` to get the current data
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "eg-lpp-DataUpdateLimit"

	// Failsafe limit for the produced active (real) power of the
	// Controllable System data updated
	//
	// Use `FailsafeProductionActivePowerLimit` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeProductionActivePowerLimit api.EventType = "eg-lpp-DataUpdateFailsafeProductionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data updated
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "eg-lpp-DataUpdateFailsafeDurationMinimum"
)
