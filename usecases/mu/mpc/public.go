package mpc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
	"time"
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
func (e *MPC) Update(measurementDataForIds ...api.MeasurementDataForID) error {
	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		return err
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
) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measurementData(
			acPowerTotal,
			timestamp,
			e.powerConfig.ValueSourceTotal,
			valueState,
			nil,
			nil,
		),
		Id: *e.acPowerTotal,
	}
}

// use MPC.UpdateDataPowerPhaseA in MPC.Update to set the momentary active power consumption or production per phase
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerPhaseA(
	acPowerPhaseA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acPower[0] == nil {
		panic("acPowerPhaseA is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			acPowerPhaseA,
			timestamp,
			e.powerConfig.ValueSourcePhaseA,
			valueState,
			nil,
			nil,
		),
		Id: *e.acPower[0],
	}
}

// use MPC.UpdateDataPowerPhaseB in MPC.Update to set the momentary active power consumption or production per phase
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerPhaseB(
	acPowerPhaseB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acPower[1] == nil {
		panic("acPowerPhaseB is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			acPowerPhaseB,
			timestamp,
			e.powerConfig.ValueSourcePhaseB,
			valueState,
			nil,
			nil,
		),
		Id: *e.acPower[1],
	}
}

// use MPC.UpdateDataPowerPhaseC in MPC.Update to set the momentary active power consumption or production per phase
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataPowerPhaseC(
	acPowerPhaseC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acPower[2] == nil {
		panic("acPowerPhaseC is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			acPowerPhaseC,
			timestamp,
			e.powerConfig.ValueSourcePhaseC,
			valueState,
			nil,
			nil,
		),
		Id: *e.acPower[2],
	}
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
) api.MeasurementDataForID {
	if e.acEnergyConsumed == nil {
		panic("acEnergyConsumed is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			energyConsumed,
			timestamp,
			e.energyConfig.ValueSourceConsumption,
			valueState,
			evaluationStart,
			evaluationEnd,
		),
		Id: *e.acEnergyConsumed,
	}
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
) api.MeasurementDataForID {
	if e.acEnergyProduced == nil {
		panic("acEnergyProduced is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			energyProduced,
			timestamp,
			e.energyConfig.ValueSourceProduction,
			valueState,
			evaluationStart,
			evaluationEnd,
		),
		Id: *e.acEnergyProduced,
	}
}

// Scenario 3

// use MPC.UpdateDataCurrentPhaseA in MPC.Update to set the momentary phase specific current consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataCurrentPhaseA(
	acCurrentPhaseA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acCurrent[0] == nil {
		panic("acCurrentPhaseA is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			acCurrentPhaseA,
			timestamp,
			e.currentConfig.ValueSourcePhaseA,
			valueState,
			nil,
			nil,
		),
		Id: *e.acCurrent[0],
	}
}

// use MPC.UpdateDataCurrentPhaseB in MPC.Update to set the momentary phase specific current consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataCurrentPhaseB(
	acCurrentPhaseB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acCurrent[1] == nil {
		panic("acCurrentPhaseB is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			acCurrentPhaseB,
			timestamp,
			e.currentConfig.ValueSourcePhaseB,
			valueState,
			nil,
			nil,
		),
		Id: *e.acCurrent[1],
	}
}

// use MPC.UpdateDataCurrentPhaseC in MPC.Update to set the momentary phase specific current consumption or production
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataCurrentPhaseC(
	acCurrentPhaseC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acCurrent[2] == nil {
		panic("acCurrentPhaseC is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			acCurrentPhaseC,
			timestamp,
			e.currentConfig.ValueSourcePhaseC,
			valueState,
			nil,
			nil,
		),
		Id: *e.acCurrent[2],
	}
}

// Scenario 4

// use MPC.UpdateDataVoltagePhaseA in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseA(
	voltagePhaseA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acVoltage[0] == nil {
		panic("acVoltagePhaseA is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			voltagePhaseA,
			timestamp,
			e.voltageConfig.ValueSourcePhaseA,
			valueState,
			nil,
			nil,
		),
		Id: *e.acVoltage[0],
	}
}

// use MPC.UpdateDataVoltagePhaseB in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseB(
	voltagePhaseB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acVoltage[1] == nil {
		panic("acVoltagePhaseB is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			voltagePhaseB,
			timestamp,
			e.voltageConfig.ValueSourcePhaseB,
			valueState,
			nil,
			nil,
		),
		Id: *e.acVoltage[1],
	}
}

// use MPC.UpdateDataVoltagePhaseC in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseC(
	voltagePhaseC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acVoltage[2] == nil {
		panic("acVoltagePhaseC is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			voltagePhaseC,
			timestamp,
			e.voltageConfig.ValueSourcePhaseC,
			valueState,
			nil,
			nil,
		),
		Id: *e.acVoltage[2],
	}
}

// use MPC.UpdateDataVoltagePhaseAToB in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseAToB(
	voltagePhaseAToB float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acVoltage[3] == nil {
		panic("acVoltagePhaseAToB is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			voltagePhaseAToB,
			timestamp,
			e.voltageConfig.ValueSourcePhaseAToB,
			valueState,
			nil,
			nil,
		),
		Id: *e.acVoltage[3],
	}
}

// use MPC.UpdateDataVoltagePhaseBToC in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseBToC(
	voltagePhaseBToC float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acVoltage[4] == nil {
		panic("acVoltagePhaseBToC is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			voltagePhaseBToC,
			timestamp,
			e.voltageConfig.ValueSourcePhaseBToC,
			valueState,
			nil,
			nil,
		),
		Id: *e.acVoltage[4],
	}
}

// use MPC.UpdateDataVoltagePhaseCToA in MPC.Update to set the phase specific voltage details
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataVoltagePhaseCToA(
	voltagePhaseCToA float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acVoltage[5] == nil {
		panic("acVoltagePhaseCToA is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			voltagePhaseCToA,
			timestamp,
			e.voltageConfig.ValueSourcePhaseCToA,
			valueState,
			nil,
			nil,
		),
		Id: *e.acVoltage[5],
	}
}

// Scenario 5

// use MPC.UpdateDataFrequency in MPC.Update to set the frequency
// The timestamp is optional and can be nil
// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
func (e *MPC) UpdateDataFrequency(
	frequency float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
) api.MeasurementDataForID {
	if e.acFrequency == nil {
		panic("acFrequency is not supported, please check the configuration")
	}
	return api.MeasurementDataForID{
		Data: measurementData(
			frequency,
			timestamp,
			e.frequencyConfig.ValueSource,
			valueState,
			nil,
			nil,
		),
		Id: *e.acFrequency,
	}
}
