package vapd

import "github.com/enbility/eebus-go/api"

const (
	// PV System total power data updated
	//
	// Use `Power` to get the current data
	//
	// Use Case VAPD, Scenario 1
	DataUpdatePower api.EventType = "cem-vapd-DataUpdatePower"

	// PV System nominal peak power data updated
	//
	// Use `PowerNominalPeak` to get the current data
	//
	// Use Case VAPD, Scenario 2
	DataUpdatePowerNominalPeak api.EventType = "cem-vapd-DataUpdatePowerNominalPeak"

	// PV System total yield data updated
	//
	// Use `PVYieldTotal` to get the current data
	//
	// Use Case VAPD, Scenario 3
	DataUpdatePVYieldTotal api.EventType = "cem-vapd-DataUpdatePVYieldTotal"
)
