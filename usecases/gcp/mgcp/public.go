package mgcp

import (
	"errors"
	"github.com/enbility/eebus-go/api"
)

// Scenario 1

// set the current power limitation factor
//
// parameters:
//   - factor: the factor to set
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (m *MGCP) SetPowerLimitationFactor(factor float64) error {
	return errors.New("not implemented")
}

// return the current power limitation factor
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (m *MGCP) PowerLimitationFactor() (float64, error) {
	return 0, errors.New("not implemented")
}

// Scenario 2

// set the momentary power consumption or production at the grid connection point
//
// parameters:
//   - power: the power to set
//   - positive values are used for consumption
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (m *MGCP) SetPower(power float64) error {
	if m.idM1 == nil {
		return api.ErrMissingData
	}

	return m.setMeasurementDataForId(m.idM1, power)
}

// return the momentary power consumption or production at the grid connection point
//
// return values:
//   - positive values are used for consumption
//   - negative values are used for production
func (m *MGCP) Power() (float64, error) {
	if m.idM1 == nil {
		return 0, api.ErrMissingData
	}

	return m.getMeasurementDataForId(m.idM1)
}

// Scenario 3

// set the total feed in energy at the grid connection point
//
// parameters:
//   - energy: the energy to set
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (m *MGCP) SetEnergyFeedIn(energy float64) error {
	if m.idM2 == nil {
		return api.ErrMissingData
	}

	return m.setMeasurementDataForId(m.idM2, energy)
}

// return the total feed in energy at the grid connection point
//
// return values:
//   - negative values are used for production
func (m *MGCP) EnergyFeedIn() (float64, error) {
	if m.idM2 == nil {
		return 0, api.ErrMissingData
	}

	return m.getMeasurementDataForId(m.idM2)
}

// Scenario 4

// set the total consumption energy at the grid connection point
//
// parameters:
//   - energy: the energy to set
//   - positive values are used for consumption
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (m *MGCP) SetEnergyConsumed(energy float64) error {
	if m.idM3 == nil {
		return api.ErrMissingData
	}

	return m.setMeasurementDataForId(m.idM3, energy)
}

// return the total consumption energy at the grid connection point
//
// return values:
//   - positive values are used for consumption
func (m *MGCP) EnergyConsumed() (float64, error) {
	if m.idM3 == nil {
		return 0, api.ErrMissingData
	}

	return m.getMeasurementDataForId(m.idM3)
}

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
func (m *MGCP) SetCurrentPerPhase(phaseA, phaseB, phaseC float64) error {
	if (m.idM41 == nil) || (m.idM42 == nil) || (m.idM43 == nil) {
		return api.ErrMissingData
	}

	err := m.setMeasurementDataForId(m.idM41, phaseA)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM42, phaseB)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM43, phaseC)
	if err != nil {
		return err
	}

	return nil
}

// return the momentary current consumption or production at the grid connection point
//
// return values:
//   - positive values are used for consumption
//   - negative values are used for production
func (m *MGCP) CurrentPerPhase() ([]float64, error) {
	if (m.idM41 == nil) || (m.idM42 == nil) || (m.idM43 == nil) {
		return []float64{}, api.ErrMissingData
	}

	valueA, err := m.getMeasurementDataForId(m.idM41)
	if err != nil {
		return []float64{}, err
	}
	valueB, err := m.getMeasurementDataForId(m.idM42)
	if err != nil {
		return []float64{}, err
	}
	valueC, err := m.getMeasurementDataForId(m.idM43)
	if err != nil {
		return []float64{}, err
	}

	return []float64{valueA, valueB, valueC}, nil
}

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
func (m *MGCP) SetVoltagePerPhase(phaseA, phaseB, phaseC float64) error {
	if (m.idM51 == nil) || (m.idM52 == nil) || (m.idM53 == nil) || (m.idM54 == nil) || (m.idM55 == nil) || (m.idM56 == nil) {
		return api.ErrMissingData
	}

	err := m.setMeasurementDataForId(m.idM51, phaseA)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM52, phaseB)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM53, phaseC)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM54, phaseA-phaseB)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM55, phaseB-phaseC)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.idM56, phaseC-phaseA)
	if err != nil {
		return err
	}

	return nil
}

// return the voltage phase details at the grid connection point
func (m *MGCP) VoltagePerPhase() ([]float64, error) {
	if (m.idM51 == nil) || (m.idM52 == nil) || (m.idM53 == nil) || (m.idM54 == nil) || (m.idM55 == nil) || (m.idM56 == nil) {
		return []float64{}, api.ErrMissingData
	}

	valueA, err := m.getMeasurementDataForId(m.idM51)
	if err != nil {
		return []float64{}, err
	}
	valueB, err := m.getMeasurementDataForId(m.idM52)
	if err != nil {
		return []float64{}, err
	}
	valueC, err := m.getMeasurementDataForId(m.idM53)
	if err != nil {
		return []float64{}, err
	}
	valueAB, err := m.getMeasurementDataForId(m.idM54)
	if err != nil {
		return []float64{}, err
	}
	valueBC, err := m.getMeasurementDataForId(m.idM55)
	if err != nil {
		return []float64{}, err
	}
	valueCA, err := m.getMeasurementDataForId(m.idM56)
	if err != nil {
		return []float64{}, err
	}

	return []float64{valueA, valueB, valueC, valueAB, valueBC, valueCA}, nil
}

// Scenario 7

// set the frequency at the grid connection point
//
// parameters:
//   - frequency: the frequency to set
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (m *MGCP) SetFrequency(frequency float64) error {
	if m.idM6 == nil {
		return api.ErrMissingData
	}

	return m.setMeasurementDataForId(m.idM6, frequency)
}

// return frequency at the grid connection point
func (m *MGCP) Frequency() (float64, error) {
	if m.idM6 == nil {
		return 0, api.ErrMissingData
	}

	return m.getMeasurementDataForId(m.idM6)
}
