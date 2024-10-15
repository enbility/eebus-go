package mdt

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mdt-UseCaseSupportUpdate"

	// DHW Temperature
	//
	// Use `Temperature` to get the current data
	//
	// Use Case MDT, Scenario 1
	DataUpdateDhwTemperature api.EventType = "ma-mdt-DataUpdateDhwTemperature"
)
