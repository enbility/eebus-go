package mdsf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mdsf-UseCaseSupportUpdate"

	// Room DHW operation mode data updated
	//
	// Use `CurrentOperationMode` to get the current data
	//
	// Use Case MDSF, Scenario 1
	DataUpdateOperationMode api.EventType = "ma-mdsf-DataUpdateOperationMode"

	// DHW overrun data updated
	//
	// Use `IsOverrunActive` or `OverrunStatus` to get the current data
	//
	// Use Case MDSF, Scenario 2
	DataUpdateOverrun api.EventType = "ma-mdsf-DataUpdateOverrun"
)
