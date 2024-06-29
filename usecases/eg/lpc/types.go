package lpc

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entites supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "eg-lpc-UseCaseSupportUpdate"

	// Load control obligation limit data updated
	//
	// Use `ConsumptionLimit` to get the current data
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "eg-lpc-DataUpdateLimit"

	// Failsafe limit for the consumed active (real) power of the
	// Controllable System data updated
	//
	// Use `FailsafeConsumptionActivePowerLimit` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeConsumptionActivePowerLimit api.EventType = "eg-lpc-DataUpdateFailsafeConsumptionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data updated
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "eg-lpc-DataUpdateFailsafeDurationMinimum"
)
