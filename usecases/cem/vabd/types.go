package vabd

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-vabd-UseCaseSupportUpdate"

	// Battery System (dis)charge power data updated
	//
	// Use `Power` to get the current data
	//
	// Use Case VABD, Scenario 1
	DataUpdatePower api.EventType = "cem-vabd-DataUpdatePower"

	// Battery System cumulated charge energy data updated
	//
	// Use `EnergyCharged` to get the current data
	//
	// Use Case VABD, Scenario 2
	DataUpdateEnergyCharged api.EventType = "cem-vabd-DataUpdateEnergyCharged"

	// Battery System cumulated discharge energy data updated
	//
	// Use `EnergyDischarged` to get the current data
	//
	// Use Case VABD, Scenario 3
	DataUpdateEnergyDischarged api.EventType = "cem-vabd-DataUpdateEnergyDischarged"

	// Battery System state of charge data updated
	//
	// Use `StateOfCharge` to get the current data
	//
	// Use Case VABD, Scenario 4
	DataUpdateStateOfCharge api.EventType = "cem-vabd-DataUpdateStateOfCharge"
)
