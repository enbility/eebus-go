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
func (e *MPC) PowerPerPhase() (map[model.ElectricalConnectionPhaseNameType]float64, error) {
	powerPerPhase := make(map[model.ElectricalConnectionPhaseNameType]float64)

	for phase, id := range e.acPowerPerPhase {
		if id != nil {
			power, err := e.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			powerPerPhase[phase] = power
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
func (e *MPC) CurrentPerPhase() (map[model.ElectricalConnectionPhaseNameType]float64, error) {
	currentPerPhase := make(map[model.ElectricalConnectionPhaseNameType]float64)

	for phase, id := range e.acCurrentPerPhase {
		if id != nil {
			current, err := e.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			currentPerPhase[phase] = current
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
func (e *MPC) VoltagePerPhase() (map[model.ElectricalConnectionPhaseNameType]float64, error) {
	voltagePerPhase := make(map[model.ElectricalConnectionPhaseNameType]float64)

	for phase, id := range e.acVoltagePerPhase {
		if id != nil {
			voltage, err := e.getMeasurementDataForId(id)
			if err != nil {
				return nil, err
			}
			voltagePerPhase[phase] = voltage
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
		e.acPowerPerPhase[model.ElectricalConnectionPhaseNameTypeA],
		measurementData(
			acPowerPhaseA,
			timestamp,
			e.powerConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeA],
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
		e.acPowerPerPhase[model.ElectricalConnectionPhaseNameTypeB],
		measurementData(
			acPowerPhaseB,
			timestamp,
			e.powerConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeB],
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
		e.acPowerPerPhase[model.ElectricalConnectionPhaseNameTypeC],
		measurementData(
			acPowerPhaseC,
			timestamp,
			e.powerConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeC],
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
	if e.acEnergyConsumed == nil {
		return newUpdateData(
			"acEnergyConsumed is not supported, please check the configuration",
			nil,
			nil,
		)
	}
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
	if e.acEnergyProduced == nil {
		return newUpdateData(
			"acEnergyProduced is not supported, please check the configuration",
			nil,
			nil,
		)
	}
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
	// validate first if current is supported
	if e.currentConfig == nil {
		return newUpdateData(
			"acCurrent is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acCurrentPhaseA is not supported, please check the configuration",
		e.acCurrentPerPhase[model.ElectricalConnectionPhaseNameTypeA],
		measurementData(
			acCurrentPhaseA,
			timestamp,
			e.currentConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeA],
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
	// validate first if current is supported
	if e.currentConfig == nil {
		return newUpdateData(
			"acCurrent is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acCurrentPhaseB is not supported, please check the configuration",
		e.acCurrentPerPhase[model.ElectricalConnectionPhaseNameTypeB],
		measurementData(
			acCurrentPhaseB,
			timestamp,
			e.currentConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeB],
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
	// validate first if current is supported
	if e.currentConfig == nil {
		return newUpdateData(
			"acCurrent is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acCurrentPhaseC is not supported, please check the configuration",
		e.acCurrentPerPhase[model.ElectricalConnectionPhaseNameTypeC],
		measurementData(
			acCurrentPhaseC,
			timestamp,
			e.currentConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeC],
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
	// validate first if voltage is supported
	if e.voltageConfig == nil {
		return newUpdateData(
			"acVoltage is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acVoltagePhaseA is not supported, please check the configuration",
		e.acVoltagePerPhase[model.ElectricalConnectionPhaseNameTypeA],
		measurementData(
			voltagePhaseA,
			timestamp,
			e.voltageConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeA],
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
	// validate first if voltage is supported
	if e.voltageConfig == nil {
		return newUpdateData(
			"acVoltage is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acVoltagePhaseB is not supported, please check the configuration",
		e.acVoltagePerPhase[model.ElectricalConnectionPhaseNameTypeB],
		measurementData(
			voltagePhaseB,
			timestamp,
			e.voltageConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeB],
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
	// validate first if voltage is supported
	if e.voltageConfig == nil {
		return newUpdateData(
			"acVoltage is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acVoltagePhaseC is not supported, please check the configuration",
		e.acVoltagePerPhase[model.ElectricalConnectionPhaseNameTypeC],
		measurementData(
			voltagePhaseC,
			timestamp,
			e.voltageConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeC],
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
	// validate first if voltage is supported
	if e.voltageConfig == nil {
		return newUpdateData(
			"acVoltage is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acVoltagePhaseAToB is not supported, please check the configuration",
		e.acVoltagePerPhase[model.ElectricalConnectionPhaseNameTypeAb],
		measurementData(
			voltagePhaseAToB,
			timestamp,
			e.voltageConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeAb],
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
	// validate first if voltage is supported
	if e.voltageConfig == nil {
		return newUpdateData(
			"acVoltage is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acVoltagePhaseBToC is not supported, please check the configuration",
		e.acVoltagePerPhase[model.ElectricalConnectionPhaseNameTypeBc],
		measurementData(
			voltagePhaseBToC,
			timestamp,
			e.voltageConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeBc],
			valueState,
			nil,
			nil,
		),
	)
}

// use MPC.UpdateDataVoltagePhaseCToA in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseAToC(
	voltagePhaseCToA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) usecaseapi.UpdateMeasurementData {
	// validate first if voltage is supported
	if e.voltageConfig == nil {
		return newUpdateData(
			"acVoltage is not supported, please check the configuration",
			nil,
			nil,
		)
	}
	return newUpdateData(
		"acVoltagePhaseCToA is not supported, please check the configuration",
		e.acVoltagePerPhase[model.ElectricalConnectionPhaseNameTypeAc],
		measurementData(
			voltagePhaseCToA,
			timestamp,
			e.voltageConfig.ValueSourcePerPhase[model.ElectricalConnectionPhaseNameTypeAc],
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
	// Validate first if frequency is supported
	if e.frequencyConfig == nil {
		return newUpdateData(
			"acFrequency is not supported, please check the configuration",
			e.acFrequency,
			nil,
		)
	}
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
