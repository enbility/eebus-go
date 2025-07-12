package mgcp

import (
	"testing"

	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMGCPValidators(t *testing.T) {
	t.Run("MGCPPowerValidator accepts valid power measurement", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500), // 2500W
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		}

		err := MGCPPowerValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("MGCPPowerValidator accepts measurement without ValueSource (recommended)", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			// ValueSource is recommended, not mandatory for power
		}

		err := MGCPPowerValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("MGCPPowerValidator requires ValueType to be value", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeAverageValue), // Wrong type
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}

		err := MGCPPowerValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueType must be value")
	})

	t.Run("MGCPEnergyValidator requires ValueSource (mandatory)", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(5000),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			// ValueSource missing - should fail for energy
		}

		err := MGCPEnergyValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueSource")
	})

	t.Run("MGCPEnergyValidator accepts all allowed ValueSource types", func(t *testing.T) {
		allowedSources := []model.MeasurementValueSourceType{
			model.MeasurementValueSourceTypeMeasuredValue,
			model.MeasurementValueSourceTypeCalculatedValue,
			model.MeasurementValueSourceTypeEmpiricalValue,
		}

		for _, source := range allowedSources {
			measurement := &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(5000),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   &source,
			}

			err := MGCPEnergyValidator.Validate(measurement)
			assert.NoError(t, err, "Energy validator should accept %s", source)
		}
	})

	t.Run("MGCP-003 rule skips measurements with error state", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeError), // Should be skipped
		}

		err := MGCPPowerValidator.Validate(measurement)
		assert.ErrorIs(t, err, internal.ErrSkipMeasurement)
	})

	t.Run("MGCP-003 rule skips measurements with outOfRange state", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeOutofrange), // Should be skipped
		}

		err := MGCPEnergyValidator.Validate(measurement)
		assert.ErrorIs(t, err, internal.ErrSkipMeasurement)
	})

	t.Run("MGCP-003 rule accepts measurements with normal state", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		}

		err := MGCPFrequencyValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("all MGCP validators require MeasurementId and Value", func(t *testing.T) {
		validators := map[string]*internal.MeasurementValidator{
			"Power":     MGCPPowerValidator,
			"Energy":    MGCPEnergyValidator,
			"Current":   MGCPCurrentValidator,
			"Voltage":   MGCPVoltageValidator,
			"Frequency": MGCPFrequencyValidator,
		}

		for name, validator := range validators {
			// Test missing MeasurementId
			measurement := &model.MeasurementDataType{
				Value:     model.NewScaledNumberType(100),
				ValueType: util.Ptr(model.MeasurementValueTypeTypeValue),
			}

			err := validator.Validate(measurement)
			assert.Error(t, err, "%s validator should reject missing MeasurementId", name)
			assert.Contains(t, err.Error(), "MeasurementId")

			// Test missing Value
			measurement = &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			}

			err = validator.Validate(measurement)
			assert.Error(t, err, "%s validator should reject missing Value", name)
		}
	})

	t.Run("all MGCP validators require ValueType to be value", func(t *testing.T) {
		validators := map[string]*internal.MeasurementValidator{
			"Power":     MGCPPowerValidator,
			"Energy":    MGCPEnergyValidator,
			"Current":   MGCPCurrentValidator,
			"Voltage":   MGCPVoltageValidator,
			"Frequency": MGCPFrequencyValidator,
		}

		for name, validator := range validators {
			measurement := &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(100),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeAverageValue), // Wrong type
			}

			err := validator.Validate(measurement)
			assert.Error(t, err, "%s validator should reject non-value ValueType", name)
			assert.Contains(t, err.Error(), "ValueType must be value")
		}
	})

	t.Run("MGCPCurrentValidator accepts measurement without ValueSource (recommended)", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(15), // 15A
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			// ValueSource is recommended, not mandatory for current
		}

		err := MGCPCurrentValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("MGCPVoltageValidator accepts measurement without ValueSource (recommended)", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(230), // 230V
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			// ValueSource is recommended, not mandatory for voltage
		}

		err := MGCPVoltageValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("MGCPFrequencyValidator accepts measurement without ValueSource (recommended)", func(t *testing.T) {
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(50), // 50Hz
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			// ValueSource is recommended, not mandatory for frequency
		}

		err := MGCPFrequencyValidator.Validate(measurement)
		assert.NoError(t, err)
	})

	t.Run("MGCP validators reject invalid ValueSource types", func(t *testing.T) {
		invalidSources := []model.MeasurementValueSourceType{
			"invalidSource",
			"approximatedValue", 
			"simulatedValue",
		}

		for _, source := range invalidSources {
			measurement := &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(100),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   &source,
			}

			// Test with power validator (recommended)
			err := MGCPPowerValidator.Validate(measurement)
			assert.Error(t, err, "Should reject invalid ValueSource %s", source)
			
			// Test with energy validator (mandatory)
			err = MGCPEnergyValidator.Validate(measurement)
			assert.Error(t, err, "Should reject invalid ValueSource %s", source)
		}
	})

	t.Run("MGCP validators handle nil fields correctly", func(t *testing.T) {
		// Test nil MeasurementId
		measurement := &model.MeasurementDataType{
			Value:     model.NewScaledNumberType(100),
			ValueType: util.Ptr(model.MeasurementValueTypeTypeValue),
		}
		err := MGCPPowerValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MeasurementId")

		// Test nil Value
		measurement = &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		}
		err = MGCPPowerValidator.Validate(measurement)
		assert.Error(t, err)

		// Test nil ValueType
		measurement = &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(100),
		}
		err = MGCPPowerValidator.Validate(measurement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValueType")
	})

	t.Run("MGCP validators handle multiple ValueType edge cases", func(t *testing.T) {
		invalidValueTypes := []model.MeasurementValueTypeType{
			model.MeasurementValueTypeTypeAverageValue,
			model.MeasurementValueTypeTypeMaxValue,
			model.MeasurementValueTypeTypeMinValue,
		}

		for _, valueType := range invalidValueTypes {
			measurement := &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(100),
				ValueType:     &valueType,
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			}

			err := MGCPEnergyValidator.Validate(measurement)
			assert.Error(t, err, "Should reject ValueType %s", valueType)
			assert.Contains(t, err.Error(), "ValueType must be value")
		}
	})
}

func TestMGCPSpecCompliance(t *testing.T) {
	t.Run("Scenario 2 (Power) - ValueSource recommended", func(t *testing.T) {
		// Valid with ValueSource
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(2500),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}
		assert.NoError(t, MGCPPowerValidator.Validate(measurement))

		// Valid without ValueSource (recommended, not mandatory)
		measurement.ValueSource = nil
		assert.NoError(t, MGCPPowerValidator.Validate(measurement))
	})

	t.Run("Scenarios 3&4 (Energy) - ValueSource mandatory", func(t *testing.T) {
		// Valid with ValueSource
		measurement := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         model.NewScaledNumberType(5000),
			ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
		}
		assert.NoError(t, MGCPEnergyValidator.Validate(measurement))

		// Invalid without ValueSource (mandatory for energy)
		measurement.ValueSource = nil
		assert.Error(t, MGCPEnergyValidator.Validate(measurement))
	})

	t.Run("Scenarios 5,6,7 (Current/Voltage/Frequency) - ValueSource recommended", func(t *testing.T) {
		validators := map[string]*internal.MeasurementValidator{
			"Current":   MGCPCurrentValidator,
			"Voltage":   MGCPVoltageValidator,
			"Frequency": MGCPFrequencyValidator,
		}

		for name, validator := range validators {
			measurement := &model.MeasurementDataType{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(100),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeEmpiricalValue),
			}

			// Valid with ValueSource
			assert.NoError(t, validator.Validate(measurement), "%s should accept with ValueSource", name)

			// Valid without ValueSource (recommended, not mandatory)
			measurement.ValueSource = nil
			assert.NoError(t, validator.Validate(measurement), "%s should accept without ValueSource", name)
		}
	})
}

func TestGetMeasurementValueMGCP(t *testing.T) {
	t.Run("extracts value from valid MGCP measurement", func(t *testing.T) {
		measurements := []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(2500),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		}

		value, err := internal.GetMeasurementValue(measurements, MGCPPowerValidator)
		assert.NoError(t, err)
		assert.Equal(t, 2500.0, value)
	})

	t.Run("skips measurements with error state per MGCP-003", func(t *testing.T) {
		measurements := []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(2500),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError), // Should be skipped
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(3000),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal), // This should be used
			},
		}

		value, err := internal.GetMeasurementValue(measurements, MGCPPowerValidator)
		assert.NoError(t, err)
		assert.Equal(t, 3000.0, value) // Should get the normal state measurement
	})

	t.Run("returns error when no valid measurements found", func(t *testing.T) {
		measurements := []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(2500),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError), // Should be skipped
			},
		}

		_, err := internal.GetMeasurementValue(measurements, MGCPEnergyValidator)
		assert.Error(t, err)
	})
}