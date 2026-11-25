package mgcp

import (
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
)

// Validators for the MGCP (Monitoring of Grid Connection Point) use case.
// Per the MGCP spec:
//   - measurementId, value, valueType=value are Mandatory
//   - valueSource is Mandatory and must be measuredValue/calculatedValue/empiricalValue
//   - MGCP-003: measurements with valueState=error or outOfRange SHALL be ignored

// mgcpValueSources are the allowed value sources for all MGCP measurements.
var mgcpValueSources = []model.MeasurementValueSourceType{
	model.MeasurementValueSourceTypeMeasuredValue,
	model.MeasurementValueSourceTypeCalculatedValue,
	model.MeasurementValueSourceTypeEmpiricalValue,
}

// scenarioValidator builds the common validator shape used by every MGCP scenario.
func scenarioValidator(name string) internal.MeasurementValidator {
	return internal.NewMeasurementValidator().
		WithName(name).
		WithRule(internal.RequireMeasurementId()).
		WithRule(internal.RequireMeasurementValue()).
		WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
		WithRule(internal.RequireValueSource(mgcpValueSources...)).
		WithRule(internal.SkipValueState())
}

var (
	MGCPPowerValidator     = scenarioValidator("MGCP Power")
	MGCPEnergyValidator    = scenarioValidator("MGCP Energy")
	MGCPCurrentValidator   = scenarioValidator("MGCP Current")
	MGCPVoltageValidator   = scenarioValidator("MGCP Voltage")
	MGCPFrequencyValidator = scenarioValidator("MGCP Frequency")
)
