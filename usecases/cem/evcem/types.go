package evcem

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-evcem-UseCaseSupportUpdate"

	// EV number of connected phases data updated
	//
	// Use `PhasesConnected` to get the current data
	//
	// Use Case EVCEM, Scenario 1
	DataUpdatePhasesConnected api.EventType = "cem-evcem-DataUpdatePhasesConnected"

	// EV current measurement data updated
	//
	// Use `CurrentPerPhase` to get the current data
	//
	// Use Case EVCEM, Scenario 1
	DataUpdateCurrentPerPhase api.EventType = "cem-evcem-DataUpdateCurrentPerPhase"

	// EV power measurement data updated
	//
	// Use `PowerPerPhase` to get the current data
	//
	// Use Case EVCEM, Scenario 2
	DataUpdatePowerPerPhase api.EventType = "cem-evcem-DataUpdatePowerPerPhase"

	// EV charging energy measurement data updated
	//
	// Use `EnergyCharged` to get the current data
	//
	// Use Case EVCEM, Scenario 3
	DataUpdateEnergyCharged api.EventType = "cem-evcem-DataUpdateEnergyCharged"
)
