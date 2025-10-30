package api

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Monitoring Unit
// UseCase: Monitoring of Power Consumption
type MuMPCInterface interface {
	api.UseCaseInterface
	// ------------------------- Getters ------------------------- //

	// Scenario 1

	// get the momentary active power consumption or production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	Power() (float64, error)

	// get the momentary active power consumption or production per phase
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	PowerPerPhase() (map[model.ElectricalConnectionPhaseNameType]float64, error)

	// Scenario 2

	// get the total feed in energy
	//
	//   - negative values are used for production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	EnergyProduced() (float64, error)

	// get the total feed in energy
	//
	//   - negative values are used for production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	EnergyConsumed() (float64, error)

	// Scenario 3

	// get the momentary phase specific current consumption or production
	//
	//   - positive values are used for consumption
	//   - negative values are used for production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	CurrentPerPhase() (map[model.ElectricalConnectionPhaseNameType]float64, error)

	// Scenario 4

	// get the phase specific voltage details
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	VoltagePerPhase() (map[model.ElectricalConnectionPhaseNameType]float64, error)

	// Scenario 5

	// get frequency
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	Frequency() (float64, error)

	// ------------------------- Setters ------------------------- //

	// use Update to update the measurement data
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
	Update(data ...UpdateMeasurementData) error

	// Scenario 1

	// use UpdateDataPowerTotal in Update to set the momentary active power consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataPowerTotal(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataPowerPhaseA in Update to set the momentary active power consumption or production per phase
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataPowerPhaseA(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataPowerPhaseB in Update to set the momentary active power consumption or production per phase
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataPowerPhaseB(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataPowerPhaseC in Update to set the momentary active power consumption or production per phase
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataPowerPhaseC(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// Scenario 2

	// use UpdateDataEnergyConsumed in Update to set the total feed in energy
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	// The evaluationStart and End are optional and can be nil (both must be set to be used)
	UpdateDataEnergyConsumed(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
		evaluationStart *time.Time,
		evaluationEnd *time.Time,
	) UpdateMeasurementData

	// use UpdateDataEnergyProduced in Update to set the total feed in energy
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	// The evaluationStart and End are optional and can be nil (both must be set to be used)
	UpdateDataEnergyProduced(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
		evaluationStart *time.Time,
		evaluationEnd *time.Time,
	) UpdateMeasurementData

	// Scenario 3

	// use UpdateDataCurrentPhaseA in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataCurrentPhaseA(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataCurrentPhaseB in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataCurrentPhaseB(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataCurrentPhaseC in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataCurrentPhaseC(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// Scenario 4

	// use UpdateDataVoltagePhaseA in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseA(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataVoltagePhaseB in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseB(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataVoltagePhaseC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseC(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataVoltagePhaseAToB in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseAToB(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataVoltagePhaseBToC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseBToC(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// use UpdateDataVoltagePhaseAToC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseAToC(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData

	// Scenario 5

	// use AcFrequency in Update to set the frequency
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataFrequency(value float64, timestamp *time.Time, valueState *model.MeasurementValueStateType) UpdateMeasurementData
}
