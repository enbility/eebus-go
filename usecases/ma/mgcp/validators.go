package mgcp

import (
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
)

// Validators for MGCP (Monitoring Appliance Grid Connection Point) use case measurements

var (
	// mgcpValueSources are the allowed value sources for all MGCP measurements per specification
	// All MGCP scenarios allow: measuredValue, calculatedValue, empiricalValue
	mgcpValueSources = []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
		model.MeasurementValueSourceTypeEmpiricalValue,
	}
)

// MGCP-specific validation rules

// SkipValueState creates a validation rule that implements MGCP-003:
// Values with state "outOfRange" or "error" SHALL be ignored by the Monitoring Appliance
func SkipValueState() internal.ValidationRule[*model.MeasurementDataType] {
	return func(m *model.MeasurementDataType) error {
		if m.ValueState != nil {
			if *m.ValueState == model.MeasurementValueStateTypeError ||
				*m.ValueState == model.MeasurementValueStateTypeOutofrange {
				return internal.ErrSkipMeasurement // Custom error to indicate skipping
			}
		}
		return nil
	}
}

// RequireValueSourceMGCP validates ValueSource when mandatory per MGCP spec
func RequireValueSourceMGCP(mandatory bool) internal.ValidationRule[*model.MeasurementDataType] {
	if mandatory {
		return internal.RequireValueSourceMandatory(mgcpValueSources...)
	}
	// For recommended, we validate if present but don't require
	return internal.RequireValueSource(mgcpValueSources...)
}

// Scenario-specific validators

// MGCPPowerValidator validates power measurements for Scenario 2
// - valueType: M (Mandatory) = "value"
// - valueSource: R (Recommended) = measuredValue|calculatedValue|empiricalValue
// - valueState: R (Recommended) with MGCP-003 rule
var MGCPPowerValidator = internal.NewMeasurementValidator().
	WithName("MGCP Power").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSourceMGCP(false)). // Recommended, not mandatory
	WithRule(SkipValueState())               // MGCP-003 rule

// MGCPEnergyValidator validates energy measurements for Scenarios 3&4
// - valueType: M (Mandatory) = "value"
// - valueSource: M (Mandatory) = measuredValue|calculatedValue|empiricalValue
// - valueState: R (Recommended) with MGCP-003 rule
var MGCPEnergyValidator = internal.NewMeasurementValidator().
	WithName("MGCP Energy").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSourceMGCP(true)). // Mandatory for energy
	WithRule(SkipValueState())              // MGCP-003 rule

// MGCPCurrentValidator validates current measurements for Scenario 5
// - valueType: M (Mandatory) = "value"
// - valueSource: R (Recommended) = measuredValue|calculatedValue|empiricalValue
// - valueState: R (Recommended) with MGCP-003 rule
var MGCPCurrentValidator = internal.NewMeasurementValidator().
	WithName("MGCP Current").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSourceMGCP(false)). // Recommended, not mandatory
	WithRule(SkipValueState())               // MGCP-003 rule

// MGCPVoltageValidator validates voltage measurements for Scenario 6
// - valueType: M (Mandatory) = "value"
// - valueSource: R (Recommended) = measuredValue|calculatedValue|empiricalValue
// - valueState: R (Recommended) with MGCP-003 rule
var MGCPVoltageValidator = internal.NewMeasurementValidator().
	WithName("MGCP Voltage").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSourceMGCP(false)). // Recommended, not mandatory
	WithRule(SkipValueState())               // MGCP-003 rule

// MGCPFrequencyValidator validates frequency measurements for Scenario 7
// - valueType: M (Mandatory) = "value"
// - valueSource: R (Recommended) = measuredValue|calculatedValue|empiricalValue
// - valueState: R (Recommended) with MGCP-003 rule
var MGCPFrequencyValidator = internal.NewMeasurementValidator().
	WithName("MGCP Frequency").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSourceMGCP(false)). // Recommended, not mandatory
	WithRule(SkipValueState())               // MGCP-003 rule

// Legacy validators for backward compatibility (deprecated for MGCP use)

// basicMGCPValidator provides minimal validation for backward compatibility
// Deprecated: Use scenario-specific validators (MGCPPowerValidator, etc.) for proper MGCP compliance
var basicMGCPValidator = internal.NewMeasurementValidator().
	WithName("MGCP Basic").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue())
