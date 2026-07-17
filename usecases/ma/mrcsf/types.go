package mrcsf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mrcsf-UseCaseSupportUpdate"

	// Room cooling operation mode data updated
	//
	// Use `CurrentOperationMode` to get the current data
	//
	// Use Case MRCSF, Scenario 1
	DataUpdateOperationMode api.EventType = "ma-mrcsf-DataUpdateOperationMode"
)
