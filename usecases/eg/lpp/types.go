package lpp

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "eg-lpp-UseCaseSupportUpdate"

	// Load control obligation limit data updated
	//
	// Use `ProductionLimit` to get the current data
	//
	// Use Case LPP, Scenario 1
	DataUpdateLimit api.EventType = "eg-lpp-DataUpdateLimit"

	// Failsafe limit for the produced active (real) power of the
	// Controllable System data updated
	//
	// Use `FailsafeProductionActivePowerLimit` to get the current data
	//
	// Use Case LPP, Scenario 2
	DataUpdateFailsafeProductionActivePowerLimit api.EventType = "eg-lpp-DataUpdateFailsafeProductionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data updated
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPP, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "eg-lpp-DataUpdateFailsafeDurationMinimum"

	// Indicates a notify heartbeat event the application should care of.
	// E.g. going into or out of the Failsafe state
	//
	// Use Case LPP, Scenario 3
	DataUpdateHeartbeat api.EventType = "eg-lpp-DataUpdateHeartbeat"
)
