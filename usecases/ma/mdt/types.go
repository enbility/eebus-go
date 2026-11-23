package mdt

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mdt-UseCaseSupportUpdate"

	// Domestic hot water temperature data updated
	//
	// Use `Temperature` to get the current data
	//
	// Use Case MDT, Scenario 1
	DataUpdateTemperature api.EventType = "ma-mdt-DataUpdateTemperature"
)
