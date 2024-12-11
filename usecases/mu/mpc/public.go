package mpc

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	usecaseapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
)

// ------------------------- Getters ------------------------- //

// Scenario 1

// get the momentary active power consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) Power() (float64, error) {
	if e.acPowerTotal == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acPowerTotal)
}

// get the momentary active power consumption or production per phase
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) PowerPerPhase() ([]float64, error) {
	powerPerPhase := make([]float64, 0)

	for _, id := range e.acPower {
		if id != nil {
			power, err := e.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			powerPerPhase = append(powerPerPhase, power)
		}
	}

	return powerPerPhase, nil
}

// Scenario 2

// get the total feed in energy
//
//   - negative values are used for production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) EnergyConsumed() (float64, error) {
	if e.acEnergyConsumed == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acEnergyConsumed)
}

// get the total feed in energy
//
//   - negative values are used for production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) EnergyProduced() (float64, error) {
	if e.acEnergyProduced == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acEnergyProduced)
}

// Scenario 3

// get the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) CurrentPerPhase() ([]float64, error) {
	currentPerPhase := make([]float64, 0)

	for _, id := range e.acCurrent {
		if id != nil {
			current, err := e.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			currentPerPhase = append(currentPerPhase, current)
		}
	}

	return currentPerPhase, nil
}

// Scenario 4

// get the phase specific voltage details
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) VoltagePerPhase() ([]float64, error) {
	voltagePerPhase := make([]float64, 0)

	for _, id := range e.acVoltage {
		if id != nil {
			voltage, err := e.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			voltagePerPhase = append(voltagePerPhase, voltage)
		}
	}

	return voltagePerPhase, nil
}

// Scenario 5

// get frequency
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) Frequency() (float64, error) {
	if e.acFrequency == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acFrequency)
}

// ------------------------- Setters ------------------------- //

// use MPC.Update to update the measurement data
// use it like this:
//
//	mpc.Update(
//	  mpc.UpdateDataPowerTotal(1000, nil, nil),
//	  mpc.UpdateDataPowerPhaseA(500, nil, nil),
//	  ...
//	)
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) Update(updateData ...usecaseapi.UpdateMeasurementData) error {
	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		return err
	}

	measurementDataForIds := make([]api.MeasurementDataForID, 0)

	for _, measurementDataForId := range updateData {
		if !measurementDataForId.Supported() {
			return measurementDataForId.NotSupportedError()
		}

		measurementDataForIds = append(measurementDataForIds, measurementDataForId.MeasurementData())
	}

	return measurements.UpdateDataForIds(measurementDataForIds)
}

// Scenario 1

// use MPC.UpdateDataPowerTotal in MPC.Update to set the momentary active power consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerTotal(
	acPowerTotal float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acPowerTotal is not supported, please check the configuration",
		e.acPowerTotal,
		measurementData(
			acPowerTotal,
			timestamp,
			e.powerConfig.ValueSourceTotal,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataPowerPhaseA in MPC.Update to set the momentary active power consumption or production per phase
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerPhaseA(
	acPowerPhaseA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acPowerPhaseA is not supported, please check the configuration",
		e.acPower[0],
		measurementData(
			acPowerPhaseA,
			timestamp,
			e.powerConfig.ValueSourcePhaseA,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataPowerPhaseB in MPC.Update to set the momentary active power consumption or production per phase
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerPhaseB(
	acPowerPhaseB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acPowerPhaseB is not supported, please check the configuration",
		e.acPower[1],
		measurementData(
			acPowerPhaseB,
			timestamp,
			e.powerConfig.ValueSourcePhaseB,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataPowerPhaseC in MPC.Update to set the momentary active power consumption or production per phase
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerPhaseC(
	acPowerPhaseC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acPowerPhaseC is not supported, please check the configuration",
		e.acPower[2],
		measurementData(
			acPowerPhaseC,
			timestamp,
			e.powerConfig.ValueSourcePhaseC,
			valueState,
			nil,
			nil,
		),
	)
}

// Scenario 2

// use MPC.UpdateDataEnergyConsumed in MPC.Update to set the total feed in energy
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
// The evaluationStart and End are optional and can be nil (both must be set to be used)
func (e *MPC) UpdateDataEnergyConsumed(
	energyConsumed float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationStart *time.Time,
	evaluationEnd *time.Time,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acEnergyConsumed is not supported, please check the configuration",
		e.acEnergyConsumed,
		measurementData(
			energyConsumed,
			timestamp,
			e.energyConfig.ValueSourceConsumption,
			valueState,
			evaluationStart,
			evaluationEnd,
		),
	)
}

// use MPC.MeasuredUpdateDataEnergyProduced in MPC.Update to set the total feed in energy
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
// The evaluationStart and End are optional and can be nil (both must be set to be used)
func (e *MPC) UpdateDataEnergyProduced(
	energyProduced float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationStart *time.Time,
	evaluationEnd *time.Time,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acEnergyProduced is not supported, please check the configuration",
		e.acEnergyProduced,
		measurementData(
			energyProduced,
			timestamp,
			e.energyConfig.ValueSourceProduction,
			valueState,
			evaluationStart,
			evaluationEnd,
		),
	)
}

// Scenario 3

// use MPC.UpdateDataCurrentPhaseA in MPC.Update to set the momentary phase specific current consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataCurrentPhaseA(
	acCurrentPhaseA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acCurrentPhaseA is not supported, please check the configuration",
		e.acCurrent[0],
		measurementData(
			acCurrentPhaseA,
			timestamp,
			e.currentConfig.ValueSourcePhaseA,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataCurrentPhaseB in MPC.Update to set the momentary phase specific current consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataCurrentPhaseB(
	acCurrentPhaseB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acCurrentPhaseB is not supported, please check the configuration",
		e.acCurrent[1],
		measurementData(
			acCurrentPhaseB,
			timestamp,
			e.currentConfig.ValueSourcePhaseB,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataCurrentPhaseC in MPC.Update to set the momentary phase specific current consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataCurrentPhaseC(
	acCurrentPhaseC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acCurrentPhaseC is not supported, please check the configuration",
		e.acCurrent[2],
		measurementData(
			acCurrentPhaseC,
			timestamp,
			e.currentConfig.ValueSourcePhaseC,
			valueState,
			nil,
			nil,
		),
	)
}

// Scenario 4

// use MPC.UpdateDataVoltagePhaseA in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseA(
	voltagePhaseA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acVoltagePhaseA is not supported, please check the configuration",
		e.acVoltage[0],
		measurementData(
			voltagePhaseA,
			timestamp,
			e.voltageConfig.ValueSourcePhaseA,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataVoltagePhaseB in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseB(
	voltagePhaseB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acVoltagePhaseB is not supported, please check the configuration",
		e.acVoltage[1],
		measurementData(
			voltagePhaseB,
			timestamp,
			e.voltageConfig.ValueSourcePhaseB,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataVoltagePhaseC in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseC(
	voltagePhaseC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acVoltagePhaseC is not supported, please check the configuration",
		e.acVoltage[2],
		measurementData(
			voltagePhaseC,
			timestamp,
			e.voltageConfig.ValueSourcePhaseC,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataVoltagePhaseAToB in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseAToB(
	voltagePhaseAToB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acVoltagePhaseAToB is not supported, please check the configuration",
		e.acVoltage[3],
		measurementData(
			voltagePhaseAToB,
			timestamp,
			e.voltageConfig.ValueSourcePhaseAToB,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataVoltagePhaseBToC in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseBToC(
	voltagePhaseBToC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acVoltagePhaseBToC is not supported, please check the configuration",
		e.acVoltage[4],
		measurementData(
			voltagePhaseBToC,
			timestamp,
			e.voltageConfig.ValueSourcePhaseBToC,
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataVoltagePhaseCToA in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseCToA(
	voltagePhaseCToA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acVoltagePhaseCToA is not supported, please check the configuration",
		e.acVoltage[5],
		measurementData(
			voltagePhaseCToA,
			timestamp,
			e.voltageConfig.ValueSourcePhaseCToA,
			valueState,
			nil,
			nil,
		),
	)
}

// Scenario 5

// use MPC.UpdateDataFrequency in MPC.Update to set the frequency
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataFrequency(
	frequency float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	return newUpdateData(
		"acFrequency is not supported, please check the configuration",
		e.acFrequency,
		measurementData(
			frequency,
			timestamp,
			e.frequencyConfig.ValueSource,
			valueState,
			nil,
			nil,
		),
	)
}
