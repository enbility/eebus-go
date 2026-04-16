package internal

import (
	"errors"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func ptrTest[T any](v T) *T { return &v }

func validMeasurementData() model.MeasurementDataType {
	return model.MeasurementDataType{
		MeasurementId: ptrTest(model.MeasurementIdType(1)),
		Value:         model.NewScaledNumberType(100),
		ValueType:     ptrTest(model.MeasurementValueTypeTypeValue),
		ValueSource:   ptrTest(model.MeasurementValueSourceTypeMeasuredValue),
		ValueState:    ptrTest(model.MeasurementValueStateTypeNormal),
	}
}

func testValidator() MeasurementValidator {
	return NewMeasurementValidator().
		WithName("Test Validator").
		WithRule(RequireMeasurementId()).
		WithRule(RequireMeasurementValue()).
		WithRule(RequireValueType(model.MeasurementValueTypeTypeValue))
}

func TestMeasurementValidator_AcceptsValid(t *testing.T) {
	m := validMeasurementData()
	assert.NoError(t, testValidator().Validate(&m))
}

func TestMeasurementValidator_RejectsMissingMeasurementId(t *testing.T) {
	m := validMeasurementData()
	m.MeasurementId = nil
	err := testValidator().Validate(&m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MeasurementId is required")
}

func TestRequireValueSource(t *testing.T) {
	rule := RequireValueSource(
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
	)

	t.Run("accepts allowed", func(t *testing.T) {
		m := validMeasurementData()
		m.ValueSource = ptrTest(model.MeasurementValueSourceTypeCalculatedValue)
		assert.NoError(t, rule(&m))
	})

	t.Run("rejects disallowed", func(t *testing.T) {
		m := validMeasurementData()
		m.ValueSource = ptrTest(model.MeasurementValueSourceTypeEmpiricalValue)
		assert.Error(t, rule(&m))
	})

	t.Run("rejects nil (mandatory)", func(t *testing.T) {
		m := validMeasurementData()
		m.ValueSource = nil
		err := rule(&m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueSource is required")
	})

	t.Run("panics when no allowed values", func(t *testing.T) {
		assert.Panics(t, func() { RequireValueSource() })
	})
}

func TestRequireValueType(t *testing.T) {
	rule := RequireValueType(model.MeasurementValueTypeTypeValue)

	t.Run("accepts matching", func(t *testing.T) {
		m := validMeasurementData()
		assert.NoError(t, rule(&m))
	})

	t.Run("rejects wrong type", func(t *testing.T) {
		m := validMeasurementData()
		m.ValueType = ptrTest(model.MeasurementValueTypeTypeAverageValue)
		err := rule(&m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueType must be")
	})

	t.Run("rejects nil", func(t *testing.T) {
		m := validMeasurementData()
		m.ValueType = nil
		err := rule(&m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueType is required")
	})
}

func TestValidateMeasurementRange(t *testing.T) {
	rule := ValidateMeasurementRange(0, 100)

	t.Run("accepts in range", func(t *testing.T) {
		m := validMeasurementData()
		m.Value = model.NewScaledNumberType(50)
		assert.NoError(t, rule(&m))
	})

	t.Run("rejects below", func(t *testing.T) {
		m := validMeasurementData()
		m.Value = model.NewScaledNumberType(-1)
		err := rule(&m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between")
	})

	t.Run("rejects above", func(t *testing.T) {
		m := validMeasurementData()
		m.Value = model.NewScaledNumberType(101)
		err := rule(&m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between")
	})
}

func TestValidateValueState(t *testing.T) {
	t.Run("required, matching", func(t *testing.T) {
		rule := ValidateValueState(model.MeasurementValueStateTypeNormal, true)
		m := validMeasurementData()
		assert.NoError(t, rule(&m))
	})

	t.Run("required, missing", func(t *testing.T) {
		rule := ValidateValueState(model.MeasurementValueStateTypeNormal, true)
		m := validMeasurementData()
		m.ValueState = nil
		err := rule(&m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueState is required")
	})

	t.Run("optional, missing", func(t *testing.T) {
		rule := ValidateValueState(model.MeasurementValueStateTypeNormal, false)
		m := validMeasurementData()
		m.ValueState = nil
		assert.NoError(t, rule(&m))
	})
}

func TestSkipValueState(t *testing.T) {
	rule := SkipValueState()
	for _, state := range []model.MeasurementValueStateType{
		model.MeasurementValueStateTypeError,
		model.MeasurementValueStateTypeOutofrange,
	} {
		t.Run(string(state), func(t *testing.T) {
			m := validMeasurementData()
			m.ValueState = ptrTest(state)
			err := rule(&m)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, ErrSkipMeasurement))
		})
	}

	t.Run("nil passes", func(t *testing.T) {
		m := validMeasurementData()
		m.ValueState = nil
		assert.NoError(t, rule(&m))
	})

	t.Run("normal passes", func(t *testing.T) {
		m := validMeasurementData()
		assert.NoError(t, rule(&m))
	})
}

func TestGetMeasurementValue(t *testing.T) {
	t.Run("returns first valid value, skipping invalid ones", func(t *testing.T) {
		invalid := validMeasurementData()
		invalid.MeasurementId = nil
		valid := validMeasurementData()
		valid.Value = model.NewScaledNumberType(42)

		value, err := GetMeasurementValue([]model.MeasurementDataType{invalid, valid}, testValidator())
		assert.NoError(t, err)
		assert.Equal(t, 42.0, value)
	})

	t.Run("errors when no valid measurements", func(t *testing.T) {
		invalid := validMeasurementData()
		invalid.MeasurementId = nil
		_, err := GetMeasurementValue([]model.MeasurementDataType{invalid}, testValidator())
		assert.Error(t, err)
	})
}
