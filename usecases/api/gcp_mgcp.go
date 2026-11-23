package api

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"time"
)

type GcpMGCPInterface interface {
	api.UseCaseInterface

	// ------------------------- Getters ------------------------- //

	// Scenario 1

	// get the current power limitation factor
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	PowerLimitationFactor() (float64, error)

	// Scenario 2

	// get the momentary power consumption or production at the grid connection point
	//
	// return values:
	//   - positive values are used for consumption
	//   - negative values are used for production
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	PowerTotal() (float64, error)

	// Scenario 3

	// get the total feed in energy at the grid connection point
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	EnergyFeedIn() (float64, error)

	// Scenario 4

	// get the total consumption energy at the grid connection point
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	EnergyConsumed() (float64, error)

	// Scenario 5

	// get the momentary phase specific current consumption or production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	CurrentPerPhase() ([]float64, error)

	// Scenario 6

	// get the momentary phase specific voltage consumption or production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	VoltagePerPhase() ([]float64, error)

	// Scenario 7

	// get frequency
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	Frequency() (float64, error)

	// ------------------------- Setters ------------------------- //

	// Update the data

	// use Update to update the data of the MGCP Usecase
	// use it like this:
	//
	//	mgcp.Update(
	//	  mgcp.MeasuredAcPowerTotal(1000, nil, nil),
	//	  mgcp.MeasuredAcPowerPhaseA(500, nil, nil),
	//	  ...
	//	)
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	Update(updateValueTypes ...UpdateData) error

	// Scenario 1

	// Use UpdateDataPowerLimitationFactor in Update to set the current power limitation factor
	UpdateDataPowerLimitationFactor(value float64) UpdateData

	// Scenario 2

	// Use UpdateDataPowerTotal in Update to set the momentary power consumption or production at the grid connection point
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataPowerTotal(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Scenario 3

	// Use UpdateDataEnergyFeedIn in Update to set the total feed in energy at the grid connection point
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
	UpdateDataEnergyFeedIn(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
		evaluationPeriodStart *time.Time,
		evaluationPeriodEnd *time.Time,
	) UpdateData

	// Scenario 4

	// Use UpdateDataEnergyConsumed in Update to set the total consumption energy at the grid connection point
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
	UpdateDataEnergyConsumed(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
		evaluationPeriodStart *time.Time,
		evaluationPeriodEnd *time.Time,
	) UpdateData

	// Scenario 5

	// Use UpdateDataCurrentPhaseA in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataCurrentPhaseA(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataCurrentPhaseB in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataCurrentPhaseB(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataCurrentPhaseC in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataCurrentPhaseC(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Scenario 6

	// Use UpdateDataVoltagePhaseA in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseA(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataVoltagePhaseB in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseB(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataVoltagePhaseC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseC(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataVoltagePhaseAToB in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseAToB(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataVoltagePhaseBToC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseBToC(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Use UpdateDataVoltagePhaseCToA in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataVoltagePhaseCToA(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData

	// Scenario 7

	// Use UpdateDataFrequency in Update to set the frequency
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	UpdateDataFrequency(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) UpdateData
}
