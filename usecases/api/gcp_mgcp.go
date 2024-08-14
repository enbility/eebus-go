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
	//
	// parameters:
	//   - factor: the factor to set
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetPvFeedInLimitationFactor(factor float64) error

	// Scenario 2

	// set the momentary power consumption or production at the grid connection point
	//
	// parameters:
	//   - power: the power to set
	//    - positive values are used for consumption
	//    - negative values are used for production
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetPower(power float64) error

	// Scenario 3

	// set the total feed in energy at the grid connection point
	//
	// parameters:
	//   - energy: the energy to set
	//    - negative values are used for production
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetEnergyFeedIn(energy float64) error

	// Scenario 4

	// set the total consumption energy at the grid connection point
	//
	// parameters:
	//   - energy: the energy to set
	//    - positive values are used for consumption
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetEnergyConsumed(energy float64) error

	// Scenario 5

	// set the momentary current consumption or production at the grid connection point
	//
	// parameters:
	//   - phaseA: the current of phase A
	//   - phaseB: the current of phase B
	//   - phaseC: the current of phase C
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetCurrentPerPhase(phaseA, phaseB, phaseC float64) error

	// Scenario 6

	// set the voltage phase details at the grid connection point
	//
	// parameters:
	//   - phaseA: the voltage of phase A
	//   - phaseB: the voltage of phase B
	//   - phaseC: the voltage of phase C
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetVoltagePerPhase(phaseA, phaseB, phaseC float64) error

	// Scenario 7

	// set the frequency at the grid connection point
	//
	// parameters:
	//   - frequency: the frequency to set
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetFrequency(frequency float64) error
}
