package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Customer Energy Management
// UseCase: Visualization of Aggregated Battery Data
type CemVABDInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current (dis)charging power
	//
	// parameters:
	//   - entity: the entity of the inverter
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2

	// return the cumulated battery system charge energy
	//
	// parameters:
	//   - entity: the entity of the inverter
	EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return the cumulated battery system discharge energy
	//
	// parameters:
	//   - entity: the entity of the inverter
	EnergyDischarged(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 4

	// return the current state of charge of the battery system
	//
	// parameters:
	//   - entity: the entity of the inverter
	StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error)
}
