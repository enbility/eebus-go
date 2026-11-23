package ohpcf

import "github.com/enbility/eebus-go/api"

const (
	UseCaseSupportUpdate api.EventType = "cem-ohpcf-UseCaseSupportUpdate"

	// Scenario 1

	DataUpdateRequestedPowerEstimate api.EventType = "cem-ohpcf-DataUpdateRequestedPowerEstimate"

	DataUpdateRequestedPowerMax api.EventType = "cem-ohpcf-DataUpdateRequestedPowerMax"

	DataUpdateConsumptionIsStoppable api.EventType = "cem-ohpcf-DataUpdateConsumptionIsStoppable"

	DataUpdateConsumptionIsPausable api.EventType = "cem-ohpcf-DataUpdateConsumptionIsPausable"

	DataUpdateConsumptionStartTime api.EventType = "cem-ohpcf-DataUpdateConsumptionStartTime"

	DataUpdateConsumptionState api.EventType = "cem-ohpcf-DataUpdateConsumptionState"

	DataUpdateMinimalRunDuration api.EventType = "cem-ohpcf-DataUpdateMinimalRunDuration"

	DataUpdateMinimalPauseDuration api.EventType = "cem-ohpcf-DataUpdateMinimalPauseDuration"
)
