package mrt

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mrt-UseCaseSupportUpdate"

	// HVAC room temperature data updated
	//
	// Use `Temperature` to get the current data
	//
	// Use Case MRT, Scenario 1
	DataUpdateTemperature api.EventType = "ma-mrt-DataUpdateTemperature"
)
