package vhan

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case.
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "va-vhan-UseCaseSupportUpdate"

	// Heating area name data updated.
	// Use `Name` to get the current data
	DataUpdateName api.EventType = "va-vhan-DataUpdateName"
)
