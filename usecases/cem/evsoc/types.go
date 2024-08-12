package evsoc

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-evsoc-UseCaseSupportUpdate"

	// EV state of charge data was updated
	//
	// Use `StateOfCharge` to get the current data
	//
	// Use Case EVSOC, Scenario 1
	DataUpdateStateOfCharge api.EventType = "cem-evsoc-DataUpdateStateOfCharge"
)
