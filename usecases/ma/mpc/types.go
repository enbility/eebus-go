package mpc

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mpc-UseCaseSupportUpdate"

	// Total momentary active power consumption or production
	//
	// Use `Power` to get the current data
	//
	// Use Case MCP, Scenario 1
	DataUpdatePower api.EventType = "ma-mpc-DataUpdatePower"

	// Phase specific momentary active power consumption or production
	//
	// Use `PowerPerPhase` to get the current data
	//
	// Use Case MCP, Scenario 1
	DataUpdatePowerPerPhase api.EventType = "ma-mpc-DataUpdatePowerPerPhase"

	// Total energy consumed
	//
	// Use `EnergyConsumed` to get the current data
	//
	// Use Case MCP, Scenario 2
	DataUpdateEnergyConsumed api.EventType = "ma-mpc-DataUpdateEnergyConsumed"

	// Total energy produced
	//
	// Use `EnergyProduced` to get the current data
	//
	// Use Case MCP, Scenario 2
	DataUpdateEnergyProduced api.EventType = "ma-mpc-DataUpdateEnergyProduced"

	// Phase specific momentary current consumption or production
	//
	// Use `CurrentPerPhase` to get the current data
	//
	// Use Case MCP, Scenario 3
	DataUpdateCurrentsPerPhase api.EventType = "ma-mpc-DataUpdateCurrentsPerPhase"

	// Phase specific voltage
	//
	// Use `VoltagePerPhase` to get the current data
	//
	// Use Case MCP, Scenario 3
	DataUpdateVoltagePerPhase api.EventType = "ma-mpc-DataUpdateVoltagePerPhase"

	// Power network frequency data updated
	//
	// Use `Frequency` to get the current data
	//
	// Use Case MCP, Scenario 3
	DataUpdateFrequency api.EventType = "ma-mpc-DataUpdateFrequency"
)
