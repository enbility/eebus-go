package cevc

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-cevc-UseCaseSupportUpdate"

	// Scenario 1

	// EV provided an energy demand
	//
	// Use `EnergyDemand` to get the current data
	DataUpdateEnergyDemand api.EventType = "cem-cevc-DataUpdateEnergyDemand"

	// Scenario 2

	// EV provided a charge plan constraints
	//
	// Use `TimeSlotConstraints` to get the current data
	DataUpdateTimeSlotConstraints api.EventType = "cem-cevc-DataUpdateTimeSlotConstraints"

	// Scenario 3

	// EV incentive table data updated
	//
	// Use `IncentiveConstraints` to get the current data
	DataUpdateIncentiveTable api.EventType = "cem-cevc-DataUpdateIncentiveTable"

	// EV requested an incentive table, call to WriteIncentiveTableDescriptions required
	DataRequestedIncentiveTableDescription api.EventType = "cem-cevc-DataRequestedIncentiveTableDescription"

	// Scenario 2 & 3

	// EV requested power limits, call to WritePowerLimits and WriteIncentives required
	DataRequestedPowerLimitsAndIncentives api.EventType = "cem-cevc-DataRequestedPowerLimitsAndIncentives"

	// Scenario 4

	// EV provided a charge plan
	//
	// Use `ChargePlanConstraints` to get the current data
	DataUpdateChargePlanConstraints api.EventType = "cem-cevc-DataUpdateChargePlanConstraints"

	// EV provided a charge plan
	//
	// Use `ChargePlan` to get the current data
	DataUpdateChargePlan api.EventType = "cem-cevc-DataUpdateChargePlan"
)
