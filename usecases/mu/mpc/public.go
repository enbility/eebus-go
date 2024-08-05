package mpc

import (
	"github.com/enbility/eebus-go/api"
)

// Scenario 1

// set the momentary active power consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) SetPower(power float64) error {
	if e.acPowerTotal == nil {
		return api.ErrMissingData
	}

	err := e.setMeasurementDataForId(e.acPowerTotal, power)
	if err != nil {
		return err
	}

	return nil
}

// set the momentary active power consumption or production per phase
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) SetPowerPerPhase(phaseA, phaseB, phaseC float64) error {
	if e.acPower[0] == nil || e.acPower[1] == nil || e.acPower[2] == nil {
		return api.ErrMissingData
	}

	err := e.setMeasurementDataForId(e.acPower[0], phaseA)
	if err != nil {
		return err
	}

	err = e.setMeasurementDataForId(e.acPower[1], phaseB)
	if err != nil {
		return err
	}

	err = e.setMeasurementDataForId(e.acPower[2], phaseC)
	if err != nil {
		return err
	}

	return nil
}

// Scenario 2

// set the total consumption energy
//
//   - positive values are used for consumption
func (e *MPC) SetEnergyConsumed(energy float64) error {
	if e.acEnergyConsumed == nil {
		return api.ErrMissingData
	}

	err := e.setMeasurementDataForId(e.acEnergyConsumed, energy)
	if err != nil {
		return err
	}

	return nil
}

// set the total feed in energy
//
//   - negative values are used for production
func (e *MPC) SetEnergyProduced(energy float64) error {
	if e.acEnergyProduced == nil {
		return api.ErrMissingData
	}

	err := e.setMeasurementDataForId(e.acEnergyProduced, energy)
	if err != nil {
		return err
	}

	return nil
}

// Scenario 3

// set the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
func (e *MPC) SetCurrentPerPhase(phaseA, phaseB, phaseC float64) error {
	if e.acCurrent[0] == nil || e.acCurrent[1] == nil || e.acCurrent[2] == nil {
		return api.ErrMissingData
	}

	err := e.setMeasurementDataForId(e.acCurrent[0], phaseA)
	if err != nil {
		return err
	}

	err = e.setMeasurementDataForId(e.acCurrent[1], phaseB)
	if err != nil {
		return err
	}

	err = e.setMeasurementDataForId(e.acCurrent[2], phaseC)
	if err != nil {
		return err
	}

	return nil
}

// Scenario 4

// set the phase specific voltage details
func (e *MPC) SetVoltagePerPhase(phaseA, phaseB, phaseC float64) error {
	for _, id := range e.acVoltage {
		if id == nil {
			return api.ErrMissingData
		}
	}

	err := e.setMeasurementDataForId(e.acVoltage[0], phaseA)
	if err != nil {
		return err
	}

	err = e.setMeasurementDataForId(e.acVoltage[1], phaseB)
	if err != nil {
		return err
	}

	err = e.setMeasurementDataForId(e.acVoltage[2], phaseC)
	if err != nil {
		return err
	}

	return nil
}

// Scenario 5

// SetFrequency set frequency
func (e *MPC) SetFrequency(frequency float64) error {
	if e.acFrequency == nil {
		return api.ErrMissingData
	}

	err := e.setMeasurementDataForId(e.acFrequency, frequency)
	if err != nil {
		return err
	}

	return nil
}
