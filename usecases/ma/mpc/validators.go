package mpc

import (
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
)

// Validators for the MPC (Monitoring of Power Consumption) use case.
// Per the MPC spec (MPC-001/002/003):
//   - measurementId, value, valueType=value are Mandatory
//   - valueSource is Mandatory and must be measuredValue/calculatedValue/empiricalValue
//   - measurements with valueState=error or outOfRange SHALL be ignored (MPC-003)

// mpcValueSources are the allowed value sources for MPC measurements.
var mpcValueSources = []model.MeasurementValueSourceType{
	model.MeasurementValueSourceTypeMeasuredValue,
	model.MeasurementValueSourceTypeCalculatedValue,
	model.MeasurementValueSourceTypeEmpiricalValue,
}

// scenarioValidator builds the common validator shape used by every MPC scenario.
func scenarioValidator(name string) internal.MeasurementValidator {
	return internal.NewMeasurementValidator().
		WithName(name).
		WithRule(internal.RequireMeasurementId()).
		WithRule(internal.RequireMeasurementValue()).
		WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
		WithRule(internal.RequireValueSource(mpcValueSources...)).
		WithRule(internal.SkipValueState())
}

var (
	powerValidator     = scenarioValidator("MPC Power")
	energyValidator    = scenarioValidator("MPC Energy")
	currentValidator   = scenarioValidator("MPC Current")
	voltageValidator   = scenarioValidator("MPC Voltage")
	frequencyValidator = scenarioValidator("MPC Frequency")
)
