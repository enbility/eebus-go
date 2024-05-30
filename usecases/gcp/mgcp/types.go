package mgcp

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entites supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "gcp-mgcp-UseCaseSupportUpdate"

	// Grid maximum allowed feed-in power as percentage value of the cumulated
	// nominal peak power of all electricity producting PV systems was updated
	//
	// Use `PowerLimitationFactor` to get the current data
	//
	// Use Case MGCP, Scenario 2
	DataUpdatePowerLimitationFactor api.EventType = "gcp-mgcp-DataUpdatePowerLimitationFactor"

	// Grid momentary power consumption/production data updated
	//
	// Use `Power` to get the current data
	//
	// Use Case MGCP, Scenario 2
	DataUpdatePower api.EventType = "gcp-mgcp-DataUpdatePower"

	// Total grid feed in energy data updated
	//
	// Use `EnergyFeedIn` to get the current data
	//
	// Use Case MGCP, Scenario 3
	DataUpdateEnergyFeedIn api.EventType = "gcp-mgcp-DataUpdateEnergyFeedIn"

	// Total grid consumed energy data updated
	//
	// Use `EnergyConsumed` to get the current data
	//
	// Use Case MGCP, Scenario 4
	DataUpdateEnergyConsumed api.EventType = "gcp-mgcp-DataUpdateEnergyConsumed"

	// Phase specific momentary current consumption/production phase detail data updated
	//
	// Use `CurrentPerPhase` to get the current data
	//
	// Use Case MGCP, Scenario 5
	DataUpdateCurrentPerPhase api.EventType = "gcp-mgcp-DataUpdateCurrentPerPhase"

	// Phase specific voltage at the grid connection point
	//
	// Use `VoltagePerPhase` to get the current data
	//
	// Use Case MGCP, Scenario 6
	DataUpdateVoltagePerPhase api.EventType = "gcp-mgcp-DataUpdateVoltagePerPhase"

	// Grid frequency data updated
	//
	// Use `Frequency` to get the current data
	//
	// Use Case MGCP, Scenario 7
	DataUpdateFrequency api.EventType = "gcp-mgcp-DataUpdateFrequency"
)
