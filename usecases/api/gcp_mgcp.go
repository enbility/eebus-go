package api

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/usecases/gcp/mgcp"
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
	GetPowerLimitationFactor() (float64, error)

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
	GetPowerTotal() (float64, error)

	// Scenario 3

	// get the total feed in energy at the grid connection point
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	GetEnergyFeedIn() (float64, error)

	// Scenario 4

	// get the total consumption energy at the grid connection point
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	GetEnergyConsumed() (float64, error)

	// Scenario 5

	// get the momentary phase specific current consumption or production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	GetCurrentPerPhase() ([]float64, error)

	// Scenario 6

	// get the momentary phase specific voltage consumption or production
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	GetVoltagePerPhase() ([]float64, error)

	// Scenario 7

	// get frequency
	//
	// possible errors:
	//   - ErrMissingData if the id is not available
	//   - and others
	GetFrequency() (float64, error)

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
	Update(updateValueTypes ...mgcp.UpdateValueType) error

	// Scenario 1

	// Use PowerLimitationFactor in Update to set the current power limitation factor
	PowerLimitationFactor(value float64) GcpMGCPUpdateValueTypeInterface

	// Scenario 2

	// Use MeasurementAcPowerTotal in Update to set the momentary power consumption or production at the grid connection point
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcPowerTotal(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Scenario 3

	// Use MeasurementEnergyFeedIn in Update to set the total feed in energy at the grid connection point
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
	MeasurementAcEnergyFeedIn(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
		evaluationPeriodStart *time.Time,
		evaluationPeriodEnd *time.Time,
	) GcpMGCPUpdateValueTypeInterface

	// Scenario 4

	// Use MeasurementEnergyConsumed in Update to set the total consumption energy at the grid connection point
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	// The evaluationPeriodStart and evaluationPeriodEnd are optional and can be nil (both must be set to be used)
	MeasurementAcEnergyConsumed(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
		evaluationPeriodStart *time.Time,
		evaluationPeriodEnd *time.Time,
	) GcpMGCPUpdateValueTypeInterface

	// Scenario 5

	// Use MeasurementAcCurrentPhaseA in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcCurrentPhaseA(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcCurrentPhaseB in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcCurrentPhaseB(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcCurrentPhaseC in Update to set the momentary phase specific current consumption or production
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcCurrentPhaseC(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Scenario 6

	// Use MeasurementAcVoltagePhaseA in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcVoltagePhaseA(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcVoltagePhaseB in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcVoltagePhaseB(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcVoltagePhaseC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcVoltagePhaseC(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcVoltagePhaseAToB in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcVoltagePhaseAToB(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcVoltagePhaseBToC in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcVoltagePhaseBToC(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Use MeasurementAcVoltagePhaseCToA in Update to set the phase specific voltage details
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcVoltagePhaseCToA(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface

	// Scenario 7

	// Use MeasurementAcFrequency in Update to set the frequency
	// The timestamp is optional and can be nil
	// The valueState shall be set if it differs from the normal valueState otherwise it can be nil
	MeasurementAcFrequency(
		value float64,
		timestamp *time.Time,
		valueState *model.MeasurementValueStateType,
	) GcpMGCPUpdateValueTypeInterface
}

type GcpMGCPUpdateValueTypeType int

const (
	GcpMGCPUpdateValueTypeTypeMeasurement   GcpMGCPUpdateValueTypeType = 0
	GcpMGCPUpdateValueTypeTypeConfiguration GcpMGCPUpdateValueTypeType = 1
)

type GcpMGCPUpdateValueTypeInterface interface {
	GetUpdateValueTypeType() GcpMGCPUpdateValueTypeType
	GetUpdateValueTypeMeasurement() api.MeasurementDataForID
	GetUpdateValueTypeConfiguration() model.DeviceConfigurationKeyValueDataType
}
