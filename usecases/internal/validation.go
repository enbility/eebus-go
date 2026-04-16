// Package internal provides validation utilities for EEBUS data types.
//
// The validation system uses Go generics to provide type-safe validation
// for any SPINE data type while maintaining high performance with no reflection.
//
// Basic Usage:
//
//	validator := internal.NewValidator[*model.MeasurementDataType]().
//		WithName("Power Validator").
//		WithRule(internal.RequireMeasurementId()).
//		WithRule(internal.RequireMeasurementValue())
//
//	err := validator.Validate(measurement)
//
// For custom data types, use generic rules:
//
//	validator := internal.NewValidator[*model.SetpointDataType]().
//		WithRule(internal.RequireField(func(s *model.SetpointDataType) *model.SetpointIdType {
//			return s.SetpointId
//		}, "SetpointId")).
//		WithRule(internal.ValidateRange(func(s *model.SetpointDataType) *model.ScaledNumberType {
//			return s.Value
//		}, 0, 100, "Percentage"))
package internal

import (
	"errors"
	"fmt"
	"github.com/enbility/spine-go/model"
)

// ErrSkipMeasurement indicates that a measurement should be skipped during validation
// This is used for MGCP-003 compliance where measurements with "error" or "outOfRange"
// states should be ignored by the Monitoring Appliance
var ErrSkipMeasurement = errors.New("measurement should be skipped")

// ========================================
// Generic Validation System
// ========================================

// Validator is the base interface for all validators.
//
// It provides methods to validate single items, find the first valid item
// from a slice, or validate all items with detailed error reporting.
type Validator[T any] interface {
	Validate(item T) error
	FindFirstValidItem(items []T) (T, error)
	ValidateAll(items []T) ([]T, []error)
	WithName(name string) Validator[T]
}

// ValidationRule is a function that validates a specific type.
//
// Rules are composable and can be combined using CombineRules or
// applied conditionally using ConditionalRule.
type ValidationRule[T any] func(T) error

// BaseValidator provides common validation functionality for any type.
//
// It implements the Validator interface and supports method chaining
// for building complex validators with multiple rules.
type BaseValidator[T any] struct {
	rules []ValidationRule[T]
	name  string
}

// NewValidator creates a new validator for type T.
//
// Example:
//
//	validator := NewValidator[*model.MeasurementDataType]().
//		WithName("Power Validator").
//		WithRule(RequireMeasurementId()).
//		WithRule(RequireMeasurementValue())
//
//	err := validator.Validate(measurement)
func NewValidator[T any]() *BaseValidator[T] {
	return &BaseValidator[T]{}
}

// WithName sets the validator name for better error messages.
//
// The name appears in error messages to help identify which validator failed.
//
// Example:
//
//	validator.WithName("Power Measurement")
//	// Error: "Power Measurement validation failed at rule 2: Value is required"
func (v *BaseValidator[T]) WithName(name string) *BaseValidator[T] {
	v.name = name
	return v
}

// WithRule adds a validation rule to the validator.
//
// Rules are executed in the order they are added. The first rule that fails
// stops validation and returns an error.
//
// Example:
//
//	validator.WithRule(RequireField(getter, "FieldName")).
//		WithRule(ValidateRange(getter, 0, 100, "Value"))
func (v *BaseValidator[T]) WithRule(rule ValidationRule[T]) *BaseValidator[T] {
	v.rules = append(v.rules, rule)
	return v
}

// Validate applies all rules to a single item.
//
// Returns nil if all rules pass, or the first error encountered.
// Error messages include the validator name (if set) and rule position.
func (v *BaseValidator[T]) Validate(item T) error {
	for i, rule := range v.rules {
		if err := rule(item); err != nil {
			if v.name != "" {
				return fmt.Errorf("%s validation failed at rule %d: %w", v.name, i+1, err)
			}
			return fmt.Errorf("validation failed at rule %d: %w", i+1, err)
		}
	}
	return nil
}

// FindFirstValidItem returns the first item in the slice that passes all validation rules.
//
// This is useful when you have multiple items but only need one valid one.
//
// Example:
//
//	measurements := []model.MeasurementDataType{...}
//	valid, err := validator.FindFirstValidItem(measurements)
//	if err != nil {
//		return api.ErrDataNotAvailable
//	}
func (v *BaseValidator[T]) FindFirstValidItem(items []T) (T, error) {
	var zero T
	for _, item := range items {
		if err := v.Validate(item); err == nil {
			return item, nil
		}
	}
	return zero, fmt.Errorf("no valid item found")
}

// ValidateAll validates all items and returns valid ones with errors for invalid ones.
//
// Unlike FindFirstValidItem, this validates every item and returns both
// the valid items and detailed errors for each invalid item.
//
// Example:
//
//	validItems, errors := validator.ValidateAll(allItems)
//	if len(errors) > 0 {
//		for _, err := range errors {
//			log.Printf("Validation error: %v", err)
//		}
//	}
//	// Process validItems...
func (v *BaseValidator[T]) ValidateAll(items []T) ([]T, []error) {
	valid := make([]T, 0, len(items))
	errors := make([]error, 0)

	for i, item := range items {
		if err := v.Validate(item); err != nil {
			errors = append(errors, fmt.Errorf("item %d: %w", i, err))
		} else {
			valid = append(valid, item)
		}
	}

	return valid, errors
}

// ========================================
// Generic Validation Rules
// ========================================
// These rules work with any type using getter functions

// RequireField creates a rule that checks if a field is not nil.
//
// This is the most commonly used validation rule for ensuring required fields are present.
//
// Example:
//
//	rule := RequireField(func(item *model.MeasurementDataType) *model.MeasurementIdType {
//		return item.MeasurementId
//	}, "MeasurementId")
//
//	// Use in validator:
//	validator.WithRule(RequireField(func(s *model.SetpointDataType) *model.SetpointIdType {
//		return s.SetpointId
//	}, "SetpointId"))
func RequireField[T any, F any](getter func(T) *F, fieldName string) ValidationRule[T] {
	return func(item T) error {
		if getter(item) == nil {
			return fmt.Errorf("%s is required", fieldName)
		}
		return nil
	}
}

// RequireScaledNumber creates a rule that checks if a ScaledNumberType field is not nil.
//
// This is a specialized version of RequireField for SPINE ScaledNumberType fields.
//
// Example:
//
//	validator.WithRule(RequireScaledNumber(func(m *model.MeasurementDataType) *model.ScaledNumberType {
//		return m.Value
//	}, "Value"))
func RequireScaledNumber[T any](getter func(T) *model.ScaledNumberType, fieldName string) ValidationRule[T] {
	return func(item T) error {
		if getter(item) == nil {
			return fmt.Errorf("%s is required", fieldName)
		}
		return nil
	}
}

// ValidateRange creates a rule that validates a numeric value is within range.
//
// If the value is nil, validation passes (use RequireScaledNumber to make it required).
// Range validation is inclusive on both ends.
//
// Example:
//
//	// Validate power is between 0 and 50kW
//	validator.WithRule(ValidateRange(func(m *model.MeasurementDataType) *model.ScaledNumberType {
//		return m.Value
//	}, 0, 50000, "Power"))
func ValidateRange[T any](
	getter func(T) *model.ScaledNumberType,
	minVal, maxVal float64,
	fieldName string,
) ValidationRule[T] {
	return func(item T) error {
		value := getter(item)
		if value == nil {
			return nil // Skip if nil, use RequireScaledNumber to make it required
		}
		val := value.GetValue()
		if val < minVal || val > maxVal {
			return fmt.Errorf("%s must be between %.2f and %.2f, got %.2f", fieldName, minVal, maxVal, val)
		}
		return nil
	}
}

// ValidateMinMax creates a rule that validates min <= value <= max relationship
func ValidateMinMax[T any](
	valueGetter func(T) *model.ScaledNumberType,
	minGetter func(T) *model.ScaledNumberType,
	maxGetter func(T) *model.ScaledNumberType,
	fieldName string,
) ValidationRule[T] {
	return func(item T) error {
		value := valueGetter(item)
		if value == nil {
			return nil // Skip if no value
		}

		val := value.GetValue()

		if minVal := minGetter(item); minVal != nil && val < minVal.GetValue() {
			return fmt.Errorf("%s %.2f is below minimum %.2f", fieldName, val, minVal.GetValue())
		}

		if maxVal := maxGetter(item); maxVal != nil && val > maxVal.GetValue() {
			return fmt.Errorf("%s %.2f is above maximum %.2f", fieldName, val, maxVal.GetValue())
		}

		return nil
	}
}

// ValidateEnum creates a rule that validates a value is one of allowed values.
// If the value pointer is nil the rule passes; combine with RequireField to reject nil.
// Passing an empty allowed list panics: that is always a programming error.
func ValidateEnum[T any, E comparable](
	getter func(T) *E,
	allowed []E,
	fieldName string,
) ValidationRule[T] {
	if len(allowed) == 0 {
		panic(fmt.Sprintf("ValidateEnum for %s requires at least one allowed value", fieldName))
	}
	return func(item T) error {
		value := getter(item)
		if value == nil {
			return nil // combine with RequireField to reject nil
		}
		for _, a := range allowed {
			if *value == a {
				return nil
			}
		}
		return fmt.Errorf("%s must be one of allowed values", fieldName)
	}
}

// ValidateSliceNotEmpty creates a rule that validates a slice is not empty
func ValidateSliceNotEmpty[T any, S any](
	getter func(T) []S,
	fieldName string,
) ValidationRule[T] {
	return func(item T) error {
		slice := getter(item)
		if len(slice) == 0 {
			return fmt.Errorf("%s cannot be empty", fieldName)
		}
		return nil
	}
}

// ValidateCustom creates a rule with custom validation logic
func ValidateCustom[T any](validate func(T) error) ValidationRule[T] {
	return validate
}

// CombineRules combines multiple validation rules into one
func CombineRules[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(item T) error {
		for _, rule := range rules {
			if err := rule(item); err != nil {
				return err
			}
		}
		return nil
	}
}

// ConditionalRule applies a rule only if a condition is met
func ConditionalRule[T any](
	condition func(T) bool,
	rule ValidationRule[T],
) ValidationRule[T] {
	return func(item T) error {
		if condition(item) {
			return rule(item)
		}
		return nil
	}
}
