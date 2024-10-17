package mgcp

import (
	"errors"
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	usecaseapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"time"
)

// -------- Getters -------- //

// Scenario 1

// get the current power limitation factor
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetPowerLimitationFactor() (float64, error) {
	_configurations, err := server.NewDeviceConfiguration(m.LocalEntity)
	if err != nil {
		return 0, err
	}

	value, err := _configurations.GetKeyValueDataForKeyId(*m.pvFeedInLimitationFactor)

	if err != nil {
		return 0, err
	}

	return value.Value.ScaledNumber.GetValue(), nil
}

// Scenario 2

// get the momentary active power consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetPowerTotal() (float64, error) {
	return m.getMeasurementDataForId(m.acPowerTotal)
}

// Scenario 3

// get the total produced energy
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetEnergyFeedIn() (float64, error) {
	return m.getMeasurementDataForId(m.gridFeedIn)
}

// Scenario 4

// get the total consumed energy
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetEnergyConsumed() (float64, error) {
	return m.getMeasurementDataForId(m.gridConsumption)
}

// Scenario 5

// get the momentary phase specific current consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetCurrentPerPhase() ([]float64, error) {
	acCurrent := make([]float64, 0)

	for _, id := range m.acCurrent {
		if id != nil {
			value, err := m.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			acCurrent = append(acCurrent, value)
		}
	}

	if len(acCurrent) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return acCurrent, nil
}

// Scenario 6

// get the momentary phase specific voltage consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetVoltagePerPhase() ([]float64, error) {
	acVoltage := make([]float64, 0)

	for _, id := range m.acVoltage {
		if id != nil {
			value, err := m.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			acVoltage = append(acVoltage, value)
		}
	}

	if len(acVoltage) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return acVoltage, nil
}

// Scenario 7

// get frequency
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) GetFrequency() (float64, error) {
	return m.getMeasurementDataForId(m.acFrequency)
}

// -------- Setters -------- //

// Update the data

// use MPC.Update to update the data of the MGCP Usecase
// use it like this:
//
//	mgcp.Update(
//	  mgcp.MeasuredAcPowerTotal(1000),
//	  mgcp.MeasuredAcPowerPhaseA(500),
//	  ...
//	)
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) Update(updateValueType ...UpdateValueType) error {
	measurements := make([]api.MeasurementDataForID, 0)
	configurations := make([]model.DeviceConfigurationKeyValueDataType, 0)

	for _, update := range updateValueType {
		if update.updateValueTypeType == usecaseapi.GcpMGCPUpdateValueTypeTypeMeasurement {
			measurements = append(measurements, update.updateTypeMeasurement)
		} else if update.updateValueTypeType == usecaseapi.GcpMGCPUpdateValueTypeTypeConfiguration {
			configurations = append(configurations, update.updateTypeConfiguration)
		} else {
			return errors.New("unknown UpdateValueTypeType: " + string(rune(update.updateValueTypeType)))
		}
	}

	if len(measurements) > 0 {
		_measurements, err := server.NewMeasurement(m.LocalEntity)
		if err != nil {
			return err
		}

		err = _measurements.UpdateDataForIds(measurements)
		if err != nil {
			return err
		}
	}

	if len(configurations) == 1 {
		_configurations, err := server.NewDeviceConfiguration(m.LocalEntity)
		if err != nil {
			return err
		}

		err = _configurations.UpdateKeyValueDataForKeyId(configurations[0], nil, *configurations[0].KeyId)
		if err != nil {
			return err
		}
	} else {
		if len(configurations) > 1 {
			return errors.New("only one PowerLimitationFactor update is supported at a time")
		}
	}

	return nil
}

// Scenario 1

// Use MGCP.PowerLimitationFactor in MGCP.Update to set the current power limitation factor
func (m *MGCP) PowerLimitationFactor(pvFeedInLimitationFactor float64) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return UpdateValueType{
		updateValueTypeType: usecaseapi.GcpMGCPUpdateValueTypeTypeMeasurement,
		updateTypeConfiguration: model.DeviceConfigurationKeyValueDataType{
			KeyId: m.pvFeedInLimitationFactor,
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(pvFeedInLimitationFactor),
			},
		},
	}
}

// Scenario 2

// Use MGCP.MeasurementAcPowerTotal in MGCP.Update to set the current total power
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcPowerTotal(
	acPowerTotal float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcPowerTotal",
		m.acPowerTotal,
		m.powerConfig.ValueSource,
		acPowerTotal,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Scenario 3

// Use MGCP.MeasurementAcEnergyFeedIn in MGCP.Update to set the total feed in energy
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
func (m *MGCP) MeasurementAcEnergyFeedIn(
	energy float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationPeriodStart *time.Time,
	evaluationPeriodEnd *time.Time,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcEnergyFeedIn",
		m.gridFeedIn,
		m.energyConfig.ValueSourceProduction,
		energy,
		timestamp,
		valueState,
		evaluationPeriodStart,
		evaluationPeriodEnd,
	)
}

// Scenario 4

// Use MGCP.MeasurementAcEnergyConsumed in MGCP.Update to set the total feed in energy
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
func (m *MGCP) MeasurementAcEnergyConsumed(
	energy float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationPeriodStart *time.Time,
	evaluationPeriodEnd *time.Time,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcEnergyConsumed",
		m.gridConsumption,
		m.energyConfig.ValueSourceConsumption,
		energy,
		timestamp,
		valueState,
		evaluationPeriodStart,
		evaluationPeriodEnd,
	)
}

// Scenario 5

// Use MGCP.MeasurementAcCurrentPhaseA in MGCP.Update to set the current of phase A
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcCurrentPhaseA(
	current float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcCurrentPhaseA",
		m.acCurrent[0],
		m.currentConfig.ValueSourcePhaseA,
		current,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcCurrentPhaseB in MGCP.Update to set the current of phase B
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcCurrentPhaseB(
	current float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcCurrentPhaseB",
		m.acCurrent[1],
		m.currentConfig.ValueSourcePhaseB,
		current,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcCurrentPhaseC in MGCP.Update to set the current of phase C
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcCurrentPhaseC(
	current float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcCurrentPhaseC",
		m.acCurrent[2],
		m.currentConfig.ValueSourcePhaseC,
		current,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Scenario 6

// Use MGCP.MeasurementAcVoltagePhaseA in MGCP.Update to set the voltage of phase A
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcVoltagePhaseA(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcVoltagePhaseA",
		m.acVoltage[0],
		m.voltageConfig.ValueSourcePhaseA,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcVoltagePhaseB in MGCP.Update to set the voltage of phase B
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcVoltagePhaseB(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcVoltagePhaseB",
		m.acVoltage[1],
		m.voltageConfig.ValueSourcePhaseB,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcVoltagePhaseC in MGCP.Update to set the voltage of phase C
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcVoltagePhaseC(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcVoltagePhaseC",
		m.acVoltage[2],
		m.voltageConfig.ValueSourcePhaseC,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcVoltagePhaseAToB in MGCP.Update to set the voltage between phase A and B
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcVoltagePhaseAToB(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcVoltagePhaseAToB",
		m.acVoltage[3],
		m.voltageConfig.ValueSourcePhaseAToB,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcVoltagePhaseBToC in MGCP.Update to set the voltage between phase B and C
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcVoltagePhaseBToC(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcVoltagePhaseBToC",
		m.acVoltage[4],
		m.voltageConfig.ValueSourcePhaseBToC,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.MeasurementAcVoltagePhaseCToA in MGCP.Update to set the voltage between phase C and A
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcVoltagePhaseCToA(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcVoltagePhaseCToA",
		m.acVoltage[5],
		m.voltageConfig.ValueSourcePhaseCToA,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Scenario 7

// Use MGCP.MeasurementAcFrequency in MGCP.Update to set the frequency
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) MeasurementAcFrequency(
	frequency float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.GcpMGCPUpdateValueTypeInterface {
	return measurementUpdateValueType(
		"MeasurementAcFrequency",
		m.acFrequency,
		m.frequencyConfig.ValueSource,
		frequency,
		timestamp,
		valueState,
		nil,
		nil,
	)
}
