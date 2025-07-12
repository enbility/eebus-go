package lpc

import (
	"fmt"

	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/spine-go/model"
)

// Validators for EG LPC (Energy Guard Load Limit Control Point) use case
//
// This file provides validators for data structures that the Energy Guard
// reads from Controllable Systems (CS) in the LPC use case. These validators
// ensure data integrity and compliance with EG LPC specification requirements.

// EG LPC-specific validation rules for LoadControlLimitData

// Note: The EEBus LPC specification only requires that "A limit lower than 0W SHALL be rejected."
// This validation is handled by the CS when EG writes data. When EG reads data from CS,
// we assume the CS has already validated its own data according to the spec.

// EG LPC LoadControlLimitData validator for Scenario 1
// Validates consumption limits received from CS devices
// Note: Direction and scope validation happens at the description level,
// this validator focuses on the actual limit data structure.
var EGLPCLoadControlLimitValidator = internal.NewValidator[*model.LoadControlLimitDataType]().
	WithName("EG LPC LoadControl Limit").
	WithRule(internal.RequireField(func(data *model.LoadControlLimitDataType) *model.LoadControlLimitIdType {
		return data.LimitId
	}, "LimitId")).
	WithRule(internal.RequireScaledNumber(func(data *model.LoadControlLimitDataType) *model.ScaledNumberType {
		return data.Value
	}, "Value"))

// EG LPC-specific validation rules for DeviceConfigurationKeyValue

// Note: The EEBus LPC specification requires that "The Active Power Consumption Limit and 
// the Failsafe Consumption Active Power Limit SHALL always be greater than or equal to zero."
// This validation is handled by the CS when EG writes data. When EG reads data from CS,
// we assume the CS has already validated its own data according to the spec.

// Note: The EEBus LPC specification requires that the Failsafe Duration Minimum
// "SHALL be pre-configured by the CS's vendor in the range of 2 hours to 24 hours"
// and "The Energy Guard SHALL choose a value between 2 hours and 24 hours" when writing.
// This validation is handled by the CS when EG writes data. When EG reads data from CS,
// we assume the CS has already validated its own data according to the spec.

// ValidateConfigurationValue creates a rule that validates the value exists and is of correct type.
// For EG LPC, only ScaledNumber (for power limits) and Duration (for failsafe duration) are valid.
func ValidateConfigurationValue() internal.ValidationRule[*model.DeviceConfigurationKeyValueDataType] {
	return func(data *model.DeviceConfigurationKeyValueDataType) error {
		if data.Value == nil {
			return nil // Value is optional in the data structure
		}
		
		// For EG LPC, only ScaledNumber (power limit) or Duration (failsafe duration) are valid
		hasValidValue := (data.Value.ScaledNumber != nil) || (data.Value.Duration != nil)
		
		// If other value types are present, skip this measurement
		if data.Value.String != nil || data.Value.Boolean != nil || data.Value.DateTime != nil {
			return fmt.Errorf("DeviceConfiguration Value must be a ScaledNumber or Duration for EG LPC")
		}
		
		if !hasValidValue {
			return internal.ErrSkipMeasurement // Skip if no actual value present
		}
		
		return nil
	}
}

// EG LPC DeviceConfigurationKeyValue validator for Scenario 2
// Validates failsafe configuration data received from CS devices
// Note: KeyName validation happens at the description level when filtering,
// this validator focuses on the actual configuration data structure.
var EGLPCDeviceConfigurationValidator = internal.NewValidator[*model.DeviceConfigurationKeyValueDataType]().
	WithName("EG LPC DeviceConfiguration").
	WithRule(internal.RequireField(func(data *model.DeviceConfigurationKeyValueDataType) *model.DeviceConfigurationKeyIdType {
		return data.KeyId
	}, "KeyId")).
	WithRule(ValidateConfigurationValue())

// EG LPC-specific validation rules for ElectricalConnectionCharacteristic

// RequireConsumptionCharacteristics creates a rule that validates characteristics are consumption-related.
// EG LPC monitors nominal and contractual consumption characteristics from CS devices.
func RequireConsumptionCharacteristics() internal.ValidationRule[*model.ElectricalConnectionCharacteristicDataType] {
	return func(data *model.ElectricalConnectionCharacteristicDataType) error {
		if data.CharacteristicType == nil {
			return nil // CharacteristicType is optional, let other validators handle required fields
		}
		
		// Check if it's a consumption-related characteristic
		if *data.CharacteristicType != model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax &&
		   *data.CharacteristicType != model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax {
			return fmt.Errorf("CharacteristicType must be a consumption characteristic type (PowerConsumptionNominalMax or ContractualConsumptionNominalMax)")
		}
		
		return nil
	}
}

// Note: The specification does not require validation of characteristic values or contexts.
// The only requirement is to filter for PowerConsumptionNominalMax or ContractualConsumptionNominalMax
// characteristic types, which is handled by RequireConsumptionCharacteristics().

// EG LPC ElectricalConnectionCharacteristic validator for Scenario 4
// Validates nominal consumption characteristics received from CS devices
var EGLPCElectricalConnectionCharacteristicValidator = internal.NewValidator[*model.ElectricalConnectionCharacteristicDataType]().
	WithName("EG LPC ElectricalConnectionCharacteristic").
	WithRule(internal.RequireField(func(data *model.ElectricalConnectionCharacteristicDataType) *model.ElectricalConnectionCharacteristicIdType {
		return data.CharacteristicId
	}, "CharacteristicId")).
	WithRule(RequireConsumptionCharacteristics())

// Helper function to validate LoadControlLimitData using EG LPC validator
// This provides a convenient interface for validating limit data in EG LPC context
func ValidateLoadControlLimit(data *model.LoadControlLimitDataType) error {
	return EGLPCLoadControlLimitValidator.Validate(data)
}

// Helper function to validate DeviceConfigurationKeyValueData using EG LPC validator
// This provides a convenient interface for validating configuration data in EG LPC context
func ValidateDeviceConfiguration(data *model.DeviceConfigurationKeyValueDataType) error {
	return EGLPCDeviceConfigurationValidator.Validate(data)
}

// Helper function to validate ElectricalConnectionCharacteristicData using EG LPC validator
// This provides a convenient interface for validating characteristic data in EG LPC context
func ValidateElectricalConnectionCharacteristic(data *model.ElectricalConnectionCharacteristicDataType) error {
	return EGLPCElectricalConnectionCharacteristicValidator.Validate(data)
}