package mpc

import (
	"testing"

	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMPCValidators(t *testing.T) {
	t.Run("powerValidator accepts valid power measurement", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(1000), // 1kW
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		}

		err := powerValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("powerValidator accepts empirical values", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(1000),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeEmpiricalValue), // Allowed for power
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		}

		err := powerValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("energyValidator rejects empirical values", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(1000),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeEmpiricalValue), // NOT allowed for energy
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		}

		err := energyValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueSource must be one of allowed values")
	})

	t.Run("currentValidator accepts any value state", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(10), // 10A
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeError), // Should be accepted
		}

		err := currentValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("currentValidator requires value state", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(10),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    nil, // Missing required field
		}

		err := currentValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueState is required")
	})

	t.Run("voltageValidator validates range", func(t *testing.T) {
		// Valid voltage
		validMeasurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(230), // 230V - within range
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}

		err := voltageValidator.Validate(validMeasurement)
		assert.NoError(t, err)

		// Invalid voltage - out of range
		invalidMeasurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(1500), // 1500V - out of range
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}

		err = voltageValidator.Validate(invalidMeasurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 0.00 and 1000.00")
	})

	t.Run("frequencyValidator accepts any frequency value", func(t *testing.T) {
		// Frequency values should be accepted regardless of range since it's not in spec
		testCases := []float64{
			50,   // Normal 50Hz
			60,   // Normal 60Hz
			44,   // Below typical range
			66,   // Above typical range
			100,  // Very high frequency
			25,   // Very low frequency
		}

		for _, freq := range testCases {
			measurement := &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(freq),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			}

			err := frequencyValidator.Validate(measurement)
			assert.NoError(t, err, "frequency %.0fHz should be accepted", freq)
		}
	})

	t.Run("frequencyValidator uses correct source types", func(t *testing.T) {
		// FIXED: frequencyValidator should no longer use powerSourceTypes
		// This test ensures empirical values are rejected for frequency
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(50),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeEmpiricalValue), // Should be rejected
		}

		err := frequencyValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueSource must be one of allowed values")
	})

	t.Run("all validators require measurement ID and value", func(t *testing.T) {
		validators := []*internal.MeasurementValidator{
			powerValidator,
			energyValidator,
			currentValidator,
			voltageValidator,
			frequencyValidator,
		}

		invalidMeasurement := &model.MeasurementDataType{
			MeasurementId: nil, // Missing required field
			Value:         model.NewScaledNumberType(100),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		}

		for _, validator := range validators {
			err := validator.Validate(invalidMeasurement)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "MeasurementId is required")
		}
	})

	t.Run("all validators require value type", func(t *testing.T) {
		validators := []*internal.MeasurementValidator{
			powerValidator,
			energyValidator,
			currentValidator,
			voltageValidator,
			frequencyValidator,
		}

		invalidMeasurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(100),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeAverageValue), // Wrong type
		}

		for _, validator := range validators {
			err := validator.Validate(invalidMeasurement)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ValueType must be value")
		}
	})
}

func TestGetMeasurementValue(t *testing.T) {
	t.Run("extracts value from valid measurement", func(t *testing.T) {
		measurements := []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(42),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		}

		value, err := getMeasurementValue(measurements, powerValidator)
		assert.NoError(t, err)
		assert.Equal(t, 42.0, value)
	})

	t.Run("returns error for invalid measurement", func(t *testing.T) {
		measurements := []model.MeasurementDataType{
			{
				MeasurementId: nil, // Invalid
				Value:         model.NewScaledNumberType(42),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			},
		}

		_, err := getMeasurementValue(measurements, powerValidator)
		assert.Error(t, err)
	})
}