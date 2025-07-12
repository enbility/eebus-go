package mpc

import (
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
)

// Validators for MPC (Monitoring Appliance Power Consumption) use case measurements

var (
	// powerSourceTypes are the allowed value sources for power measurements
	powerSourceTypes = []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
		model.MeasurementValueSourceTypeEmpiricalValue,
	}

	// energySourceTypes are the allowed value sources for energy measurements
	energySourceTypes = []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
	}

	// currentSourceTypes are the allowed value sources for current measurements
	currentSourceTypes = []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
	}

	// voltageSourceTypes are the allowed value sources for voltage measurements
	voltageSourceTypes = []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
	}

	// frequencySourceTypes are the allowed value sources for frequency measurements
	frequencySourceTypes = []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
	}
)

// powerValidator validates power measurements
var powerValidator = internal.NewMeasurementValidator().
	WithName("MPC Power").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(internal.RequireValueSource(powerSourceTypes...)).
	WithRule(internal.ValidateValueState(model.MeasurementValueStateTypeNormal, false))

// energyValidator validates energy measurements (consumption and production)
var energyValidator = internal.NewMeasurementValidator().
	WithName("MPC Energy").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(internal.RequireValueSource(energySourceTypes...)).
	WithRule(internal.ValidateValueState(model.MeasurementValueStateTypeNormal, false))

// currentValidator validates current measurements
var currentValidator = internal.NewMeasurementValidator().
	WithName("MPC Current").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(internal.RequireValueSource(currentSourceTypes...)).
	WithRule(internal.ValidateValueState("", true)) // Any state required per spec

// voltageValidator validates voltage measurements
var voltageValidator = internal.NewMeasurementValidator().
	WithName("MPC Voltage").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(internal.RequireValueSource(voltageSourceTypes...)).
	WithRule(internal.ValidateMeasurementRange(0, 1000)) // 0-1000V per spec

// frequencyValidator validates frequency measurements
var frequencyValidator = internal.NewMeasurementValidator().
	WithName("MPC Frequency").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(internal.RequireValueSource(frequencySourceTypes...)).
	WithRule(internal.ValidateValueState(model.MeasurementValueStateTypeNormal, false)) // Reject error states

// getMeasurementValue is a helper that validates and extracts the value from measurements
// DEPRECATED: Use MeasurementPhaseSpecificDataForFilter with validators instead
func getMeasurementValue(measurements []model.MeasurementDataType, validator *internal.MeasurementValidator) (float64, error) {
	return internal.GetMeasurementValue(measurements, validator)
}