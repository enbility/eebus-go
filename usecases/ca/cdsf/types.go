package cdsf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ca-cdsf-UseCaseSupportUpdate"

	// Room DHW operation mode data updated
	//
	// Use `CurrentOperationMode` to get the current data
	//
	// Use Case CDSF, Scenario 1
	DataUpdateOperationMode api.EventType = "ca-cdsf-DataUpdateOperationMode"

	// DHW overrun data updated
	//
	// Use Case CDSF, Scenarios 2 and 3
	DataUpdateOverrun api.EventType = "ca-cdsf-DataUpdateOverrun"
)
