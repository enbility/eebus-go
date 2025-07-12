package internal

import (
	"fmt"
	"slices"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/features/server"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// According to the LPC Installation Guide, this value is relatively high to make sure it doesn't conflict with other IDs on this entity
var defaultPowerTotalMeasurementId = model.MeasurementIdType(50)

// return the phase specific measurement data
func MeasurementPhaseSpecificDataForFilter(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	measurementFilter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
	validPhaseNameTypes []model.ElectricalConnectionPhaseNameType,
	validator *MeasurementValidator, // NEW: Required validator
) ([]float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	electricalConnection, err1 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err1 != nil {
		return nil, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(measurementFilter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	// Validate validator parameter
	if validator == nil {
		return nil, fmt.Errorf("validator is required")
	}

	var result []float64

	for _, item := range data {
		// Use validator instead of basic nil checks
		if err := validator.Validate(&item); err != nil {
			continue // Skip invalid measurements, don't fail entire operation
		}

		if validPhaseNameTypes != nil {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			param, err := electricalConnection.GetParameterDescriptionsForFilter(filter)
			if err != nil || len(param) == 0 ||
				param[0].AcMeasuredPhases == nil ||
				!slices.Contains(validPhaseNameTypes, *param[0].AcMeasuredPhases) {
				continue
			}
		}

		if energyDirection != "" {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			desc, err := electricalConnection.GetDescriptionForParameterDescriptionFilter(filter)
			if err != nil || desc == nil {
				continue
			}

			// if energy direction is not consume
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != energyDirection {
				return nil, err
			}
		}

		// Note: ValueState validation is now handled by the validator
		// This removes the spec-violating behavior of returning ErrDataInvalid
		// for "error" and "outOfRange" states

		value := item.Value.GetValue()
		result = append(result, value)
	}

	// Handle case where no measurements passed validation
	if len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// GetPowerTotalMeasurementId returns the MeasurementId for the AC Power Total measurement
func GetPowerTotalMeasurementId(localEntity spineapi.EntityLocalInterface) model.MeasurementIdType {
	if localEntity == nil {
		return defaultPowerTotalMeasurementId
	}
	measurementFeat := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	if measurementFeat == nil {
		return defaultPowerTotalMeasurementId
	}
	measurement, err := server.NewMeasurement(localEntity)
	if err != nil || measurement == nil {
		return defaultPowerTotalMeasurementId
	}
	MeasurementDescriptionData, err := measurement.GetDescriptionsForFilter(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})
	if err != nil || len(MeasurementDescriptionData) != 1 || MeasurementDescriptionData[0].MeasurementId == nil {
		return defaultPowerTotalMeasurementId
	}

	return *MeasurementDescriptionData[0].MeasurementId
}

// ========================================
// Measurement-Specific Validation
// ========================================

// MeasurementValidator is a specialized validator for MeasurementDataType.
//
// It wraps BaseValidator with measurement-specific convenience methods
// and provides type-safe access to measurement validation rules.
type MeasurementValidator struct {
	*BaseValidator[*model.MeasurementDataType]
}

// NewMeasurementValidator creates a new measurement validator.
//
// Example:
//
//	validator := NewMeasurementValidator().
//		WithName("Power Measurement").
//		WithRule(RequireMeasurementId()).
//		WithRule(RequireMeasurementValue()).
//		WithRule(ValidateMeasurementRange(0, 50000))
func NewMeasurementValidator() *MeasurementValidator {
	return &MeasurementValidator{
		BaseValidator: NewValidator[*model.MeasurementDataType](),
	}
}

// WithRule adds a validation rule (type-safe wrapper)
func (v *MeasurementValidator) WithRule(rule ValidationRule[*model.MeasurementDataType]) *MeasurementValidator {
	v.BaseValidator.WithRule(rule)
	return v
}

// WithName sets the validator name (type-safe wrapper)
func (v *MeasurementValidator) WithName(name string) *MeasurementValidator {
	v.BaseValidator.WithName(name)
	return v
}

// GetMeasurementValue validates measurements and extracts the first valid value.
//
// This is the primary function for measurement validation and value extraction.
// It validates all measurements against the provided validator and returns
// the numeric value from the first valid measurement.
//
// Returns api.ErrDataNotAvailable if no valid measurements are found.
//
// Example:
//
//	measurements := []model.MeasurementDataType{...}
//	value, err := GetMeasurementValue(measurements, PowerMeasurementValidator)
//	if err != nil {
//		return 0, api.ErrDataNotAvailable
//	}
//	return value, nil
func GetMeasurementValue(measurements []model.MeasurementDataType, validator *MeasurementValidator) (float64, error) {
	// Convert slice to pointer slice for validation
	ptrMeasurements := make([]*model.MeasurementDataType, len(measurements))
	for i := range measurements {
		ptrMeasurements[i] = &measurements[i]
	}
	
	valid, err := validator.ValidateFirst(ptrMeasurements)
	if err != nil {
		return 0, api.ErrDataNotAvailable
	}
	
	if valid == nil || valid.Value == nil {
		return 0, api.ErrDataNotAvailable
	}
	
	return valid.Value.GetValue(), nil
}

// Measurement-specific validation rules

// RequireMeasurementId ensures MeasurementId is present.
//
// This is the most basic measurement validation rule.
// Almost all measurement validators should include this rule.
func RequireMeasurementId() ValidationRule[*model.MeasurementDataType] {
	return RequireField(
		func(m *model.MeasurementDataType) *model.MeasurementIdType { return m.MeasurementId },
		"MeasurementId",
	)
}

// RequireMeasurementValue ensures Value is present.
//
// This rule ensures that measurement data contains an actual value.
// Use this when you need to extract numeric values from measurements.
func RequireMeasurementValue() ValidationRule[*model.MeasurementDataType] {
	return RequireScaledNumber(
		func(m *model.MeasurementDataType) *model.ScaledNumberType { return m.Value },
		"Value",
	)
}

// RequireValueType ensures ValueType matches expected type
func RequireValueType(expected model.MeasurementValueTypeType) ValidationRule[*model.MeasurementDataType] {
	return func(m *model.MeasurementDataType) error {
		if expected == "" {
			return nil // Skip if no expected type specified
		}
		if m.ValueType == nil {
			return fmt.Errorf("ValueType is required to be %s", expected)
		}
		if *m.ValueType != expected {
			return fmt.Errorf("ValueType must be %s, got %s", expected, *m.ValueType)
		}
		return nil
	}
}

// RequireValueSource ensures ValueSource is one of allowed types
func RequireValueSource(allowed ...model.MeasurementValueSourceType) ValidationRule[*model.MeasurementDataType] {
	if len(allowed) == 0 {
		return func(m *model.MeasurementDataType) error { return nil }
	}
	
	return ValidateEnum(
		func(m *model.MeasurementDataType) *model.MeasurementValueSourceType { return m.ValueSource },
		allowed,
		"ValueSource",
	)
}

// ValidateValueState validates the value state with optional requirement
func ValidateValueState(expected model.MeasurementValueStateType, required bool) ValidationRule[*model.MeasurementDataType] {
	return func(m *model.MeasurementDataType) error {
		if m.ValueState == nil {
			if required {
				return fmt.Errorf("ValueState is required")
			}
			return nil
		}
		
		if expected != "" && *m.ValueState != expected {
			return fmt.Errorf("ValueState must be %s, got %s", expected, *m.ValueState)
		}
		
		return nil
	}
}

// ValidateMeasurementRange ensures measurement value is within range
func ValidateMeasurementRange(min, max float64) ValidationRule[*model.MeasurementDataType] {
	return ValidateRange(
		func(m *model.MeasurementDataType) *model.ScaledNumberType { return m.Value },
		min, max,
		"Measurement value",
	)
}

// Common measurement validators for reuse

// PowerMeasurementValidator validates standard power measurements.
//
// Validates power measurements with standard requirements:
// - MeasurementId and Value must be present
// - ValueType must be "value"
// - ValueSource must be measured, calculated, or empirical
// - ValueState should be normal (but not required)
//
// Example:
//
//	value, err := GetMeasurementValue(measurements, PowerMeasurementValidator)
var PowerMeasurementValidator = NewMeasurementValidator().
	WithName("Power Measurement").
	WithRule(RequireMeasurementId()).
	WithRule(RequireMeasurementValue()).
	WithRule(RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSource(
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
		model.MeasurementValueSourceTypeEmpiricalValue,
	)).
	WithRule(ValidateValueState(model.MeasurementValueStateTypeNormal, false))

// EnergyMeasurementValidator validates energy measurements.
//
// Similar to PowerMeasurementValidator but with stricter ValueSource requirements
// (no empirical values allowed for energy measurements).
//
// Example:
//
//	energy, err := GetMeasurementValue(measurements, EnergyMeasurementValidator)
var EnergyMeasurementValidator = NewMeasurementValidator().
	WithName("Energy Measurement").
	WithRule(RequireMeasurementId()).
	WithRule(RequireMeasurementValue()).
	WithRule(RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(RequireValueSource(
		model.MeasurementValueSourceTypeMeasuredValue,
		model.MeasurementValueSourceTypeCalculatedValue,
	)).
	WithRule(ValidateValueState(model.MeasurementValueStateTypeNormal, false))

// CurrentMeasurementValidator validates current measurements
var CurrentMeasurementValidator = NewMeasurementValidator().
	WithName("Current Measurement").
	WithRule(RequireMeasurementId()).
	WithRule(RequireMeasurementValue()).
	WithRule(RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(ValidateValueState("", true)) // State required but any value OK

// VoltageMeasurementValidator validates voltage measurements  
var VoltageMeasurementValidator = NewMeasurementValidator().
	WithName("Voltage Measurement").
	WithRule(RequireMeasurementId()).
	WithRule(RequireMeasurementValue()).
	WithRule(RequireValueType(model.MeasurementValueTypeTypeValue)).
	WithRule(ValidateMeasurementRange(0, 1000)) // 0-1000V reasonable range

// FrequencyMeasurementValidator validates frequency measurements
var FrequencyMeasurementValidator = NewMeasurementValidator().
	WithName("Frequency Measurement").
	WithRule(RequireMeasurementId()).
	WithRule(RequireMeasurementValue()).
	WithRule(RequireValueType(model.MeasurementValueTypeTypeValue))
