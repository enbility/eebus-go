package api

import (
	"github.com/enbility/eebus-go/api"
)

// Actor: Monitored Unit
// UseCase: Monitoring of Power Consumption
type MuMPCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the momentary active power consumption or production
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetPower(value float64) (resultErr error)

	// return the momentary active phase specific power consumption or production per phase
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetPowerPerPhase(value []float64) (resultErr error)

	// Scenario 2

	// return the total consumption energy
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	//   - positive values are used for consumption
	SetEnergyConsumed(value float64) (resultErr error)

	// return the total feed in energy
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// return values:
	//   - negative values are used for production
	SetEnergyProduced(value float64) (resultErr error)

	// Scenario 3

	// return the momentary phase specific current consumption or production
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// return values
	//   - positive values are used for consumption
	//   - negative values are used for production
	SetCurrentPerPhase(value []float64) (resultErr error)

	// Scenario 4

	// return the phase specific voltage details
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	SetVoltagePerPhase(value []float64) (resultErr error)

	// Scenario 5

	// return frequency
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	SetFrequency(value float64) (resultErr error)
}
