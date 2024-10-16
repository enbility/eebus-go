package api

import (
	"github.com/enbility/eebus-go/api"
)

// Actor: Monitored Unit
// UseCase: Monitoring of Power Consumption
type MuMPCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// set the momentary active power consumption or production
	//
	// parameters:
	//   - power: the active power
	SetPower(power float64) error

	// return the momentary active power consumption or production
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - and others
	Power() (float64, error)

	// set the momentary active phase specific power consumption or production per phase
	//
	// parameters:
	//   - phaseA: the active power of phase A
	//   - phaseB: the active power of phase B
	//   - phaseC: the active power of phase C
	SetPowerPerPhase(phaseA, phaseB, phaseC float64) error

	// return the momentary active phase specific power consumption or production per phase
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such values are (yet) available
	//   - and others
	PowerPerPhase() ([]float64, error)

	// Scenario 2

	// set the total consumption energy
	//
	// parameters:
	//  - consumed: the total consumption energy
	SetEnergyConsumed(consumed float64) error

	// return the total feed in energy
	//
	// return values:
	//   - negative values are used for production
	EnergyConsumed() (float64, error)

	// set the total feed in energy
	//
	// parameters:
	//  - produced: the total feed in energy
	SetEnergyProduced(produced float64) error

	// return the total feed in energy
	//
	// return values:
	//   - negative values are used for production
	EnergyProduced() (float64, error)

	// Scenario 3

	// set the momentary phase specific current consumption or production
	//
	// parameters:
	//   - phaseA: the current of phase A
	//   - phaseB: the current of phase B
	//   - phaseC: the current of phase C
	SetCurrentPerPhase(phaseA, phaseB, phaseC float64) error

	// return the momentary phase specific current consumption or production
	//
	// return values
	//   - positive values are used for consumption
	//   - negative values are used for production
	CurrentPerPhase() ([]float64, error)

	// Scenario 4

	// set the phase specific voltage details
	//
	// parameters:
	//   - phaseA: the voltage of phase A
	//   - phaseB: the voltage of phase B
	//   - phaseC: the voltage of phase C
	SetVoltagePerPhase(phaseA, phaseB, phaseC float64) error

	// return the phase specific voltage details
	VoltagePerPhase() ([]float64, error)

	// Scenario 5

	// set frequency
	//
	// parameters:
	//   - frequency: the frequency
	SetFrequency(frequency float64) error

	// return frequency
	Frequency() (float64, error)
}
