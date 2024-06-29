package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of Grid Connection Point
type MaMGCPInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current power limitation factor
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	PowerLimitationFactor(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2

	// return the momentary power consumption or production at the grid connection point
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	//
	// return values:
	//   - positive values are used for consumption
	//   - negative values are used for production
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return the total feed in energy at the grid connection point
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	//
	// return values:
	//   - negative values are used for production
	EnergyFeedIn(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 4

	// return the total consumption energy at the grid connection point
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	//
	// return values:
	//   - positive values are used for consumption
	EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 5

	// return the momentary current consumption or production at the grid connection point
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	//
	// return values:
	//   - positive values are used for consumption
	//   - negative values are used for production
	CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 6

	// return the voltage phase details at the grid connection point
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 7

	// return frequency at the grid connection point
	//
	// parameters:
	//   - entity: the entity of the device (e.g. SMGW)
	Frequency(entity spineapi.EntityRemoteInterface) (float64, error)
}
