package crht

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ca-crht-UseCaseSupportUpdate"

	// Room heating temperature setpoint data updated
	//
	// Use `Setpoints` to get the current data
	//
	// Use Case CRHT, Scenario 1
	DataUpdateSetpoints api.EventType = "ca-crht-DataUpdateSetpoints"

	// Room heating temperature setpoint constraints updated
	//
	// Use `SetpointConstraints` to get the current data
	//
	// Use Case CRHT, Scenario 1
	DataUpdateSetpointConstraints api.EventType = "ca-crht-DataUpdateSetpointConstraints"
)
