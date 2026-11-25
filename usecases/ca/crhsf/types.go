package crhsf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ca-crhsf-UseCaseSupportUpdate"

	// Room heating operation mode data updated
	//
	// Use `CurrentOperationMode` to get the current data
	//
	// Use Case CRHSF, Scenario 1
	DataUpdateOperationMode api.EventType = "ca-crhsf-DataUpdateOperationMode"
)
