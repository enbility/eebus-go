package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Customer Energy Management
// UseCase: EV Charging Electricity Measurement
type CemEVCEMInterface interface {
	api.UseCaseInterface

	// return the number of ac connected phases of the EV or 0 if it is unknown
	//
	// parameters:
	//   - entity: the entity of the EV
	PhasesConnected(entity spineapi.EntityRemoteInterface) (uint, error)

	// Scenario 1

	// return the last current measurement for each phase of the connected EV
	//
	// parameters:
	//   - entity: the entity of the EV
	CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 2

	// return the last power measurement for each phase of the connected EV
	//
	// parameters:
	//   - entity: the entity of the EV
	PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 3

	// return the charged energy measurement in Wh of the connected EV
	//
	// parameters:
	//   - entity: the entity of the EV
	EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error)
}
