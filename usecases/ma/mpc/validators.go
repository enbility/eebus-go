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

// frequencyValidator validates frequency measurements
var frequencyValidator = internal.NewMeasurementValidator().
	WithName("MPC Frequency").
	WithRule(internal.RequireMeasurementId()).
	WithRule(internal.RequireMeasurementValue()).
	WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(internal.RequireValueSource(powerSourceTypes...)).
	WithRule(internal.ValidateValueState(model.MeasurementValueStateTypeNormal, false))

// getMeasurementValue is a helper that validates and extracts the value from measurements
func getMeasurementValue(measurements []model.MeasurementDataType, validator *internal.MeasurementValidator) (float64, error) {
	return internal.GetMeasurementValue(measurements, validator)
}