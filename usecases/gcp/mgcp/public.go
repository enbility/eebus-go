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
func (m *MGCP) PowerLimitationFactor() (float64, error) {
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
func (m *MGCP) PowerTotal() (float64, error) {
	return m.getMeasurementDataForId(m.acPowerTotal)
}

// Scenario 3

// get the total produced energy
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) EnergyFeedIn() (float64, error) {
	return m.getMeasurementDataForId(m.gridFeedIn)
}

// Scenario 4

// get the total consumed energy
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) EnergyConsumed() (float64, error) {
	return m.getMeasurementDataForId(m.gridConsumption)
}

// Scenario 5

// get the momentary phase specific current consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) CurrentPerPhase() ([]float64, error) {
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
func (m *MGCP) VoltagePerPhase() ([]float64, error) {
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
func (m *MGCP) Frequency() (float64, error) {
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
func (m *MGCP) Update(updateValueType ...usecaseapi.UpdateData) error {
	measurements := make([]api.MeasurementDataForID, 0)
	configurations := make([]model.DeviceConfigurationKeyValueDataType, 0)

	for _, update := range updateValueType {
		switch update.(type) {
		case usecaseapi.UpdateMeasurementData:
			measurements = append(measurements, update.(usecaseapi.UpdateMeasurementData).MeasurementData())
		case usecaseapi.UpdateConfigurationData:
			configurations = append(configurations, update.(usecaseapi.UpdateConfigurationData).ConfigurationData())
		default:
			return errors.New("unsupported updateValueType")
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

// Use MGCP.UpdateDataPowerLimitationFactor in MGCP.Update to set the current power limitation factor
func (m *MGCP) UpdateDataPowerLimitationFactor(pvFeedInLimitationFactor float64) usecaseapi.UpdateData {
	if m.pvFeedInLimitationFactor == nil {
		return &UpdateData{
			supported:         false,
			notSupportedError: errors.New("id is nil: UpdateDataPowerLimitationFactor, please check the mgcp configuration"),
		}
	}

	return &UpdateConfigurationData{
		UpdateData: UpdateData{
			supported: true,
		},
		configurationData: model.DeviceConfigurationKeyValueDataType{
			KeyId: m.pvFeedInLimitationFactor,
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(pvFeedInLimitationFactor),
			},
		},
	}
}

// Scenario 2

// Use MGCP.UpdateDataPowerTotal in MGCP.Update to set the current total power
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataPowerTotal(
	acPowerTotal float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataPowerTotal",
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

// Use MGCP.UpdateDataEnergyFeedIn in MGCP.Update to set the total feed in energy
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
func (m *MGCP) UpdateDataEnergyFeedIn(
	energy float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationPeriodStart *time.Time,
	evaluationPeriodEnd *time.Time,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataEnergyFeedIn",
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

// Use MGCP.UpdateDataEnergyConsumed in MGCP.Update to set the total feed in energy
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
func (m *MGCP) UpdateDataEnergyConsumed(
	energy float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationPeriodStart *time.Time,
	evaluationPeriodEnd *time.Time,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataEnergyConsumed",
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

// Use MGCP.UpdateDataCurrentPhaseA in MGCP.Update to set the current of phase A
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataCurrentPhaseA(
	current float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataCurrentPhaseA",
		m.acCurrent[0],
		m.currentConfig.ValueSourcePhaseA,
		current,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataCurrentPhaseB in MGCP.Update to set the current of phase B
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataCurrentPhaseB(
	current float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataCurrentPhaseB",
		m.acCurrent[1],
		m.currentConfig.ValueSourcePhaseB,
		current,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataCurrentPhaseC in MGCP.Update to set the current of phase C
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataCurrentPhaseC(
	current float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataCurrentPhaseC",
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

// Use MGCP.UpdateDataVoltagePhaseA in MGCP.Update to set the voltage of phase A
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataVoltagePhaseA(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataVoltagePhaseA",
		m.acVoltage[0],
		m.voltageConfig.ValueSourcePhaseA,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataVoltagePhaseB in MGCP.Update to set the voltage of phase B
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataVoltagePhaseB(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataVoltagePhaseB",
		m.acVoltage[1],
		m.voltageConfig.ValueSourcePhaseB,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataVoltagePhaseC in MGCP.Update to set the voltage of phase C
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataVoltagePhaseC(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataVoltagePhaseC",
		m.acVoltage[2],
		m.voltageConfig.ValueSourcePhaseC,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataVoltagePhaseAToB in MGCP.Update to set the voltage between phase A and B
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataVoltagePhaseAToB(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataVoltagePhaseAToB",
		m.acVoltage[3],
		m.voltageConfig.ValueSourcePhaseAToB,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataVoltagePhaseBToC in MGCP.Update to set the voltage between phase B and C
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataVoltagePhaseBToC(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataVoltagePhaseBToC",
		m.acVoltage[4],
		m.voltageConfig.ValueSourcePhaseBToC,
		voltage,
		timestamp,
		valueState,
		nil,
		nil,
	)
}

// Use MGCP.UpdateDataVoltagePhaseCToA in MGCP.Update to set the voltage between phase C and A
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataVoltagePhaseCToA(
	voltage float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataVoltagePhaseCToA",
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

// Use MGCP.UpdateDataFrequency in MGCP.Update to set the frequency
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (m *MGCP) UpdateDataFrequency(
	frequency float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateData {
	return updateMeasurementData(
		"UpdateDataFrequency",
		m.acFrequency,
		m.frequencyConfig.ValueSource,
		frequency,
		timestamp,
		valueState,
		nil,
		nil,
	)
}
