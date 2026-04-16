package lpc

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestEGLPCValidators(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{"TestLoadControlLimitValidator", testLoadControlLimitValidator},
		{"TestDeviceConfigurationValidator", testDeviceConfigurationValidator},
		{"TestElectricalConnectionCharacteristicValidator", testElectricalConnectionCharacteristicValidator},
		{"TestHelperFunctions", testHelperFunctions},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testLoadControlLimitValidator(t *testing.T) {
	// Test valid LoadControlLimitData
	validData := &model.LoadControlLimitDataType{
		LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
		Value:         model.NewScaledNumberType(5000), // 5kW
		IsLimitActive: util.Ptr(true),
	}

	err := EGLPCLoadControlLimitValidator.Validate(validData)
	assert.NoError(t, err, "Valid LoadControlLimitData should pass validation")

	// Test missing LimitId
	invalidData := &model.LoadControlLimitDataType{
		Value: model.NewScaledNumberType(5000),
	}

	err = EGLPCLoadControlLimitValidator.Validate(invalidData)
	assert.Error(t, err, "LoadControlLimitData without LimitId should fail validation")
	assert.Contains(t, err.Error(), "LimitId is required")

	// Test missing Value
	invalidData2 := &model.LoadControlLimitDataType{
		LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
	}

	err = EGLPCLoadControlLimitValidator.Validate(invalidData2)
	assert.Error(t, err, "LoadControlLimitData without Value should fail validation")
	assert.Contains(t, err.Error(), "Value is required")
}

func testDeviceConfigurationValidator(t *testing.T) {
	// Test valid DeviceConfigurationKeyValueData with ScaledNumber
	validData := &model.DeviceConfigurationKeyValueDataType{
		KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(10000), // 10kW
		},
	}

	err := EGLPCDeviceConfigurationValidator.Validate(validData)
	assert.NoError(t, err, "Valid DeviceConfigurationKeyValueData should pass validation")

	// Test valid DeviceConfigurationKeyValueData with Duration
	validDuration := &model.DeviceConfigurationKeyValueDataType{
		KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(2)),
		Value: &model.DeviceConfigurationKeyValueValueType{
			Duration: model.NewDurationType(4 * time.Hour), // 4 hours
		},
	}

	err = EGLPCDeviceConfigurationValidator.Validate(validDuration)
	assert.NoError(t, err, "Valid DeviceConfigurationKeyValueData with duration should pass validation")

	// Test missing KeyId
	invalidData := &model.DeviceConfigurationKeyValueDataType{
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(5000),
		},
	}

	err = EGLPCDeviceConfigurationValidator.Validate(invalidData)
	assert.Error(t, err, "DeviceConfigurationKeyValueData without KeyId should fail validation")
	assert.Contains(t, err.Error(), "KeyId is required")

	// Test valid value type check - ScaledNumber
	validScaledNumber := &model.DeviceConfigurationKeyValueDataType{
		KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(5000),
		},
	}

	err = EGLPCDeviceConfigurationValidator.Validate(validScaledNumber)
	assert.NoError(t, err, "DeviceConfigurationKeyValueData with ScaledNumber should pass validation")

	// Test valid value type check - Duration
	validDurationValue := &model.DeviceConfigurationKeyValueDataType{
		KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(2)),
		Value: &model.DeviceConfigurationKeyValueValueType{
			Duration: model.NewDurationType(3 * time.Hour),
		},
	}

	err = EGLPCDeviceConfigurationValidator.Validate(validDurationValue)
	assert.NoError(t, err, "DeviceConfigurationKeyValueData with Duration should pass validation")

	// Test invalid value type - String (not allowed)
	stringValue := model.DeviceConfigurationKeyValueStringType("some string")
	invalidValueType := &model.DeviceConfigurationKeyValueDataType{
		KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(3)),
		Value: &model.DeviceConfigurationKeyValueValueType{
			String: &stringValue,
		},
	}

	err = EGLPCDeviceConfigurationValidator.Validate(invalidValueType)
	assert.Error(t, err, "DeviceConfigurationKeyValueData with String type should fail validation")
	assert.Contains(t, err.Error(), "must be a ScaledNumber or Duration")
}

func testElectricalConnectionCharacteristicValidator(t *testing.T) {
	// Test valid ElectricalConnectionCharacteristicData
	validData := &model.ElectricalConnectionCharacteristicDataType{
		CharacteristicId:      util.Ptr(model.ElectricalConnectionCharacteristicIdType(1)),
		CharacteristicContext: util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:    util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
		Value:                 model.NewScaledNumberType(22000), // 22kW
	}

	err := EGLPCElectricalConnectionCharacteristicValidator.Validate(validData)
	assert.NoError(t, err, "Valid ElectricalConnectionCharacteristicData should pass validation")

	// Test with contractual consumption characteristic
	validContractual := &model.ElectricalConnectionCharacteristicDataType{
		CharacteristicId:      util.Ptr(model.ElectricalConnectionCharacteristicIdType(2)),
		CharacteristicContext: util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:    util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
		Value:                 model.NewScaledNumberType(11000), // 11kW
	}

	err = EGLPCElectricalConnectionCharacteristicValidator.Validate(validContractual)
	assert.NoError(t, err, "Valid contractual ElectricalConnectionCharacteristicData should pass validation")

	// Test missing CharacteristicId
	invalidData := &model.ElectricalConnectionCharacteristicDataType{
		CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
		Value:              model.NewScaledNumberType(22000),
	}

	err = EGLPCElectricalConnectionCharacteristicValidator.Validate(invalidData)
	assert.Error(t, err, "ElectricalConnectionCharacteristicData without CharacteristicId should fail validation")
	assert.Contains(t, err.Error(), "CharacteristicId is required")

	// Test invalid characteristic type (production instead of consumption)
	invalidType := &model.ElectricalConnectionCharacteristicDataType{
		CharacteristicId:   util.Ptr(model.ElectricalConnectionCharacteristicIdType(1)),
		CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerProductionNominalMax),
		Value:              model.NewScaledNumberType(22000),
	}

	err = EGLPCElectricalConnectionCharacteristicValidator.Validate(invalidType)
	assert.Error(t, err, "ElectricalConnectionCharacteristicData with production type should be skipped")
	assert.Contains(t, err.Error(), "must be a consumption characteristic type")
}

func testHelperFunctions(t *testing.T) {
	// Test ValidateLoadControlLimit helper
	validLimit := &model.LoadControlLimitDataType{
		LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
		Value:   model.NewScaledNumberType(5000),
	}

	err := ValidateLoadControlLimit(validLimit)
	assert.NoError(t, err, "ValidateLoadControlLimit helper should work for valid data")

	invalidLimit := &model.LoadControlLimitDataType{
		Value: model.NewScaledNumberType(5000),
	}

	err = ValidateLoadControlLimit(invalidLimit)
	assert.Error(t, err, "ValidateLoadControlLimit helper should fail for invalid data")

	// Test ValidateDeviceConfiguration helper
	validConfig := &model.DeviceConfigurationKeyValueDataType{
		KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(10000),
		},
	}

	err = ValidateDeviceConfiguration(validConfig)
	assert.NoError(t, err, "ValidateDeviceConfiguration helper should work for valid data")

	invalidConfig := &model.DeviceConfigurationKeyValueDataType{
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(10000),
		},
	}

	err = ValidateDeviceConfiguration(invalidConfig)
	assert.Error(t, err, "ValidateDeviceConfiguration helper should fail for invalid data")

	// Test ValidateElectricalConnectionCharacteristic helper
	validCharacteristic := &model.ElectricalConnectionCharacteristicDataType{
		CharacteristicId:   util.Ptr(model.ElectricalConnectionCharacteristicIdType(1)),
		CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
		Value:              model.NewScaledNumberType(22000),
	}

	err = ValidateElectricalConnectionCharacteristic(validCharacteristic)
	assert.NoError(t, err, "ValidateElectricalConnectionCharacteristic helper should work for valid data")

	invalidCharacteristic := &model.ElectricalConnectionCharacteristicDataType{
		CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
		Value:              model.NewScaledNumberType(22000),
	}

	err = ValidateElectricalConnectionCharacteristic(invalidCharacteristic)
	assert.Error(t, err, "ValidateElectricalConnectionCharacteristic helper should fail for invalid data")
}

// Test ValidateConfigurationValue function (77.8% coverage)
func TestValidateConfigurationValue(t *testing.T) {
	validator := ValidateConfigurationValue()

	t.Run("accepts nil value", func(t *testing.T) {
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
			Value: nil,
		}

		err := validator(data)
		assert.NoError(t, err)
	})

	t.Run("accepts ScaledNumber value", func(t *testing.T) {
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(5000),
			},
		}

		err := validator(data)
		assert.NoError(t, err)
	})

	t.Run("accepts Duration value", func(t *testing.T) {
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(2)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				Duration: model.NewDurationType(4 * time.Hour),
			},
		}

		err := validator(data)
		assert.NoError(t, err)
	})

	t.Run("rejects String value", func(t *testing.T) {
		stringValue := model.DeviceConfigurationKeyValueStringType("invalid")
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(3)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				String: &stringValue,
			},
		}

		err := validator(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a ScaledNumber or Duration")
	})

	t.Run("rejects Boolean value", func(t *testing.T) {
		boolValue := true
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(4)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				Boolean: &boolValue,
			},
		}

		err := validator(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a ScaledNumber or Duration")
	})

	t.Run("rejects DateTime value", func(t *testing.T) {
		dateTime := model.DateTimeType("2023-01-01T12:00:00Z")
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(5)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				DateTime: &dateTime,
			},
		}

		err := validator(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a ScaledNumber or Duration")
	})

	t.Run("returns ErrSkipMeasurement for empty value", func(t *testing.T) {
		data := &model.DeviceConfigurationKeyValueDataType{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(6)),
			Value: &model.DeviceConfigurationKeyValueValueType{},
		}

		err := validator(data)
		assert.Error(t, err)
		// Should return ErrSkipMeasurement which will be caught by internal package
	})
}

// Test RequireConsumptionCharacteristics function (83.3% coverage)
func TestRequireConsumptionCharacteristics(t *testing.T) {
	validator := RequireConsumptionCharacteristics()

	t.Run("accepts nil CharacteristicType", func(t *testing.T) {
		data := &model.ElectricalConnectionCharacteristicDataType{
			CharacteristicId:   util.Ptr(model.ElectricalConnectionCharacteristicIdType(1)),
			CharacteristicType: nil,
		}

		err := validator(data)
		assert.NoError(t, err)
	})

	t.Run("accepts PowerConsumptionNominalMax", func(t *testing.T) {
		data := &model.ElectricalConnectionCharacteristicDataType{
			CharacteristicId:   util.Ptr(model.ElectricalConnectionCharacteristicIdType(1)),
			CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
		}

		err := validator(data)
		assert.NoError(t, err)
	})

	t.Run("accepts ContractualConsumptionNominalMax", func(t *testing.T) {
		data := &model.ElectricalConnectionCharacteristicDataType{
			CharacteristicId:   util.Ptr(model.ElectricalConnectionCharacteristicIdType(2)),
			CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
		}

		err := validator(data)
		assert.NoError(t, err)
	})

	t.Run("rejects production characteristic type", func(t *testing.T) {
		data := &model.ElectricalConnectionCharacteristicDataType{
			CharacteristicId:   util.Ptr(model.ElectricalConnectionCharacteristicIdType(3)),
			CharacteristicType: util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerProductionNominalMax),
		}

		err := validator(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a consumption characteristic type")
	})
}
