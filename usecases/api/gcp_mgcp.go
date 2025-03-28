package api

import (
	"github.com/enbility/eebus-go/api"
)

// Actor: Grid Connection Point
// UseCase: Monitoring of Grid Connection Point
type GcpMGCPInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// set the current power limitation factor
	SetPowerLimitationFactor(value float64) (resultErr error)

	// Scenario 2

	// set the momentary power consumption or production at the grid connection point
	//
	// parameters:
	//   - positive values are used for consumption
	//   - negative values are used for production
	SetPower(value float64) (resultErr error)

	// Scenario 3

	// set the total feed in energy at the grid connection point
	//
	// parameters:
	//   - negative values are used for production
	SetEnergyFeedIn(value float64) (resultErr error)

	// Scenario 4

	// set the total consumption energy at the grid connection point
	//
	// parameters:
	//   - positive values are used for consumption
	SetEnergyConsumed(value float64) (resultErr error)

	// Scenario 5

	// set the current phase details at the grid connection point
	SetCurrentPerPhase(value []float64) (resultErr error)

	// Scenario 6

	// set the voltage phase details at the grid connection point
	SetVoltagePerPhase(value []float64) (resultErr error)

	// Scenario 7

	// set frequency at the grid connection point
	SetFrequency(value float64) (resultErr error)
}
