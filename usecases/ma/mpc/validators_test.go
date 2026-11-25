package mpc

import (
	"errors"
	"testing"

	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func validMeasurement() *model.MeasurementDataType {
	return &model.MeasurementDataType{
		MeasurementId: util.Ptr(model.MeasurementIdType(1)),
		Value:         model.NewScaledNumberType(100),
		ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
	}
}

func mpcValidators() map[string]internal.MeasurementValidator {
	return map[string]internal.MeasurementValidator{
		"power":     powerValidator,
		"energy":    energyValidator,
		"current":   currentValidator,
		"voltage":   voltageValidator,
		"frequency": frequencyValidator,
	}
}

func TestMPCValidators_AcceptValidMeasurements(t *testing.T) {
	for name, v := range mpcValidators() {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, v.Validate(validMeasurement()))
		})
	}
}

func TestMPCValidators_AcceptAllAllowedValueSources(t *testing.T) {
	sources := []model.MeasurementValueSourceType{
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
		model.MeasurementValueSourceTypeEmpiricalValue,
	}
	for name, v := range mpcValidators() {
		for _, src := range sources {
			t.Run(name+"/"+string(src), func(t *testing.T) {
				m := validMeasurement()
				m.ValueSource = util.Ptr(src)
				assert.NoError(t, v.Validate(m))
			})
		}
	}
}

func TestMPCValidators_RejectMissingRequiredFields(t *testing.T) {
	for name, v := range mpcValidators() {
		t.Run(name+"/MeasurementId", func(t *testing.T) {
			m := validMeasurement()
			m.MeasurementId = nil
			err := v.Validate(m)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "MeasurementId")
		})
		t.Run(name+"/Value", func(t *testing.T) {
			m := validMeasurement()
			m.Value = nil
			err := v.Validate(m)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Value")
		})
		t.Run(name+"/ValueType", func(t *testing.T) {
			m := validMeasurement()
			m.ValueType = nil
			err := v.Validate(m)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ValueType")
		})
		t.Run(name+"/ValueSource", func(t *testing.T) {
			m := validMeasurement()
			m.ValueSource = nil
			err := v.Validate(m)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ValueSource is required")
		})
	}
}

func TestMPCValidators_RejectWrongValueType(t *testing.T) {
	for name, v := range mpcValidators() {
		t.Run(name, func(t *testing.T) {
			m := validMeasurement()
			m.ValueType = util.Ptr(model.MeasurementValueTypeTypeAverageValue)
			err := v.Validate(m)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ValueType must be")
		})
	}
}

func TestMPCValidators_SkipValueStateErrorOrOutOfRange(t *testing.T) {
	// MPC-003: measurements with state error or outOfRange SHALL be ignored.
	for _, state := range []model.MeasurementValueStateType{
		model.MeasurementValueStateTypeError,
		model.MeasurementValueStateTypeOutofrange,
	} {
		for name, v := range mpcValidators() {
			t.Run(name+"/"+string(state), func(t *testing.T) {
				m := validMeasurement()
				m.ValueState = util.Ptr(state)
				err := v.Validate(m)
				assert.Error(t, err)
				assert.True(t, errors.Is(err, internal.ErrSkipMeasurement), "expected ErrSkipMeasurement, got %v", err)
			})
		}
	}
}

func TestMPCValidators_AcceptMissingValueState(t *testing.T) {
	// ValueState is optional; its absence must not fail validation.
	for name, v := range mpcValidators() {
		t.Run(name, func(t *testing.T) {
			m := validMeasurement()
			m.ValueState = nil
			assert.NoError(t, v.Validate(m))
		})
	}
}
