package mgcp

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
		Value:         model.NewScaledNumberType(2500),
		ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
	}
}

func mgcpValidators() map[string]internal.MeasurementValidator {
	return map[string]internal.MeasurementValidator{
		"power":     MGCPPowerValidator,
		"energy":    MGCPEnergyValidator,
		"current":   MGCPCurrentValidator,
		"voltage":   MGCPVoltageValidator,
		"frequency": MGCPFrequencyValidator,
	}
}

func TestMGCPValidators_AcceptValidMeasurements(t *testing.T) {
	for name, v := range mgcpValidators() {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, v.Validate(validMeasurement()))
		})
	}
}

func TestMGCPValidators_AcceptAllAllowedValueSources(t *testing.T) {
	for _, src := range mgcpValueSources {
		for name, v := range mgcpValidators() {
			t.Run(name+"/"+string(src), func(t *testing.T) {
				m := validMeasurement()
				m.ValueSource = util.Ptr(src)
				assert.NoError(t, v.Validate(m))
			})
		}
	}
}

func TestMGCPValidators_RejectMissingRequiredFields(t *testing.T) {
	for name, v := range mgcpValidators() {
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

func TestMGCPValidators_RejectWrongValueType(t *testing.T) {
	for name, v := range mgcpValidators() {
		for _, badType := range []model.MeasurementValueTypeType{
			model.MeasurementValueTypeTypeAverageValue,
			model.MeasurementValueTypeTypeMaxValue,
			model.MeasurementValueTypeTypeMinValue,
		} {
			t.Run(name+"/"+string(badType), func(t *testing.T) {
				m := validMeasurement()
				m.ValueType = util.Ptr(badType)
				err := v.Validate(m)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "ValueType must be")
			})
		}
	}
}

func TestMGCPValidators_RejectUnknownValueSource(t *testing.T) {
	for name, v := range mgcpValidators() {
		for _, bad := range []model.MeasurementValueSourceType{"simulatedValue", "approximatedValue"} {
			t.Run(name+"/"+string(bad), func(t *testing.T) {
				m := validMeasurement()
				m.ValueSource = util.Ptr(bad)
				err := v.Validate(m)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "ValueSource")
			})
		}
	}
}

func TestMGCPValidators_SkipValueStateErrorOrOutOfRange(t *testing.T) {
	// MGCP-003: measurements with state error or outOfRange SHALL be ignored.
	for _, state := range []model.MeasurementValueStateType{
		model.MeasurementValueStateTypeError,
		model.MeasurementValueStateTypeOutofrange,
	} {
		for name, v := range mgcpValidators() {
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

func TestMGCPValidators_AcceptMissingValueState(t *testing.T) {
	for name, v := range mgcpValidators() {
		t.Run(name, func(t *testing.T) {
			m := validMeasurement()
			m.ValueState = nil
			assert.NoError(t, v.Validate(m))
		})
	}
}

func TestGetMeasurementValue_SkipsErrorState(t *testing.T) {
	measurements := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
		},
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(2)),
			Value:         model.NewScaledNumberType(3000),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		},
	}

	value, err := internal.GetMeasurementValue(measurements, MGCPPowerValidator)
	assert.NoError(t, err)
	assert.Equal(t, 3000.0, value)
}

func TestGetMeasurementValue_AllSkipped(t *testing.T) {
	measurements := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
		},
	}
	_, err := internal.GetMeasurementValue(measurements, MGCPEnergyValidator)
	assert.Error(t, err)
}
