package crct

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ca-crct-UseCaseSupportUpdate"

	// Room cooling temperature setpoint data updated
	//
	// Use `Setpoints` to get the current data
	//
	// Use Case CRCT, Scenario 1
	DataUpdateSetpoints api.EventType = "ca-crct-DataUpdateSetpoints"

	// Room cooling temperature setpoint constraints updated
	//
	// Use `SetpointConstraints` to get the current data
	//
	// Use Case CRCT, Scenario 1
	DataUpdateSetpointConstraints api.EventType = "ca-crct-DataUpdateSetpointConstraints"
)
