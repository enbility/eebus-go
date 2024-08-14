package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
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
func (m *MGCP) SetPvFeedInLimitationFactor(factor float64) error {
	configuration, err := server.NewDeviceConfiguration(m.LocalEntity)
	if err != nil {
		panic(err)
	}

	value := model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(factor),
	}

	data := model.DeviceConfigurationKeyValueDataType{
		KeyId: m.pvFeedInLimitationFactor,
		Value: &value,
	}

	return configuration.UpdateKeyValueDataForKeyId(
		data,
		nil,
		*m.pvFeedInLimitationFactor,
	)
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
	if m.acPowerTotal == nil {
		return api.ErrMissingData
	}

	return m.setMeasurementDataForId(m.acPowerTotal, power)
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
	if m.gridFeedIn == nil {
		return api.ErrMissingData
	}

	return m.setMeasurementDataForId(m.gridFeedIn, energy)
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
	if m.gridConsumption == nil {
		return api.ErrMissingData
	}
	return m.setMeasurementDataForId(m.gridConsumption, energy)
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
	for _, v := range m.acCurrent {
		if v == nil {
			return api.ErrMissingData
		}
	}

	err := m.setMeasurementDataForId(m.acCurrent[0], phaseA)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.acCurrent[1], phaseB)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.acCurrent[2], phaseC)
	if err != nil {
		return err
	}

	return nil
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
	for _, v := range m.acVoltage {
		if v == nil {
			return api.ErrMissingData
		}
	}

	err := m.setMeasurementDataForId(m.acVoltage[0], phaseA)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.acVoltage[1], phaseB)
	if err != nil {
		return err
	}
	err = m.setMeasurementDataForId(m.acVoltage[2], phaseC)
	if err != nil {
		return err
	}

	return nil
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
	if m.acFrequency == nil {
		return api.ErrMissingData
	}
	return m.setMeasurementDataForId(m.acFrequency, frequency)
}
