# Validation System Documentation

The `usecases/internal` package provides a comprehensive, type-safe validation framework for EEBUS data types. This system uses Go generics to provide maximum code reuse while maintaining compile-time type safety.

## Architecture Overview

The validation system consists of two main components:

1. **Generic Validation Framework** (`validation.go`) - Works with any type
2. **Measurement-Specific Validation** (`measurement.go`) - Specialized for measurement data

### File Structure

```
usecases/internal/
├── validation.go           # Generic validation framework
├── validation_test.go      # Comprehensive tests for validation.go
├── measurement.go          # Measurement utilities + validation
├── measurement_test.go     # Tests for measurement-specific code
└── <other files>           # Other internal utilities
```

## Generic Validation Framework

### Core Interfaces

```go
// Validator is the base interface for all validators
type Validator[T any] interface {
    Validate(item T) error
    ValidateFirst(items []T) (T, error)
    ValidateAll(items []T) ([]T, []error)
    WithName(name string) Validator[T]
}

// ValidationRule is a function that validates a specific type
type ValidationRule[T any] func(T) error
```

### Creating Validators

```go
// Create a validator for any type
validator := internal.NewValidator[*model.YourDataType]().
    WithName("Your Validator").
    WithRule(internal.RequireField(func(item *model.YourDataType) *FieldType { 
        return item.SomeField 
    }, "SomeField")).
    WithRule(internal.ValidateRange(func(item *model.YourDataType) *model.ScaledNumberType { 
        return item.Value 
    }, 0, 1000, "Value"))
```

### Using Validators

```go
// Validate a single item
err := validator.Validate(item)

// Find first valid item from slice
validItem, err := validator.ValidateFirst(items)

// Validate all items, get valid ones and errors for invalid ones
validItems, errors := validator.ValidateAll(items)
```

## Generic Validation Rules

### Field Validation

```go
// Require any field to be non-nil
RequireField(func(item T) *FieldType { return item.Field }, "FieldName")

// Require ScaledNumberType to be non-nil
RequireScaledNumber(func(item T) *model.ScaledNumberType { return item.Value }, "Value")
```

### Numeric Validation

```go
// Validate value is within fixed range
ValidateRange(getter, min, max, fieldName)

// Validate value respects min/max fields on the same item
ValidateMinMax(valueGetter, minGetter, maxGetter, fieldName)
```

### Enum Validation

```go
// Validate field is one of allowed values
ValidateEnum(getter, []AllowedType{value1, value2}, "FieldName")

// Empty allowed list means "allow everything"
ValidateEnum(getter, []AllowedType{}, "FieldName") // Allows any value
```

### Collection Validation

```go
// Validate slice is not empty
ValidateSliceNotEmpty(func(item T) []SliceType { return item.Items }, "Items")
```

### Advanced Rules

```go
// Custom validation logic
ValidateCustom(func(item T) error {
    if item.SomeCondition {
        return errors.New("custom error")
    }
    return nil
})

// Combine multiple rules into one
CombineRules(rule1, rule2, rule3)

// Conditional validation (only apply rule if condition met)
ConditionalRule(
    func(item T) bool { return item.IsActive }, // condition
    RequireField(getter, "RequiredWhenActive"),   // rule to apply
)
```

## Measurement-Specific Validation

### Quick Start

For common measurement validation, use predefined validators:

```go
import internal "github.com/enbility/eebus-go/usecases/internal"

// Use predefined validators
value, err := internal.GetMeasurementValue(measurements, internal.EnergyMeasurementValidator)
value, err := internal.GetMeasurementValue(measurements, internal.PowerMeasurementValidator)
value, err := internal.GetMeasurementValue(measurements, internal.CurrentMeasurementValidator)
value, err := internal.GetMeasurementValue(measurements, internal.VoltageMeasurementValidator)
value, err := internal.GetMeasurementValue(measurements, internal.FrequencyMeasurementValidator)
```

### Custom Measurement Validators

```go
// Create custom measurement validator
var customValidator = internal.NewMeasurementValidator().
    WithName("Custom Energy").
    WithRule(internal.RequireMeasurementId()).
    WithRule(internal.RequireMeasurementValue()).
    WithRule(internal.RequireValueType(model.MeasurementValueTypeTypeValue)).
    WithRule(internal.RequireValueSource(
        model.MeasurementValueSourceTypeMeasuredValue,
        model.MeasurementValueSourceTypeCalculatedValue,
    )).
    WithRule(internal.ValidateValueState(model.MeasurementValueStateTypeNormal, false)).
    WithRule(internal.ValidateMeasurementRange(0, 999999))

// Use it
value, err := internal.GetMeasurementValue(measurements, customValidator)
```

### Measurement-Specific Rules

```go
// Measurement field validation
RequireMeasurementId()        // Requires MeasurementId field
RequireMeasurementValue()     // Requires Value field

// Measurement constraints
RequireValueType(expected)                    // Validates ValueType
RequireValueSource(allowed...)                // Validates ValueSource enum
ValidateValueState(expected, required)        // Validates ValueState
ValidateMeasurementRange(min, max)           // Validates Value range
```

## Examples by Use Case

### Example 1: Setpoint Validation

```go
// Define validator for setpoint data
var setpointValidator = internal.NewValidator[*model.SetpointDataType]().
    WithName("Setpoint Validator").
    WithRule(internal.RequireField(func(s *model.SetpointDataType) *model.SetpointIdType { 
        return s.SetpointId 
    }, "SetpointId")).
    WithRule(internal.RequireScaledNumber(func(s *model.SetpointDataType) *model.ScaledNumberType { 
        return s.Value 
    }, "Value")).
    WithRule(internal.ValidateRange(func(s *model.SetpointDataType) *model.ScaledNumberType { 
        return s.Value 
    }, 0, 100, "Percentage"))

// Use validator
err := setpointValidator.Validate(setpointData)
if err != nil {
    return fmt.Errorf("setpoint validation failed: %w", err)
}
```

### Example 2: Load Control Limits

```go
// Define validator for load control limits
var limitValidator = internal.NewValidator[*model.LoadControlLimitDataType]().
    WithName("Load Control Limit").
    WithRule(internal.RequireField(func(l *model.LoadControlLimitDataType) *model.LoadControlLimitIdType { 
        return l.LimitId 
    }, "LimitId")).
    WithRule(internal.RequireScaledNumber(func(l *model.LoadControlLimitDataType) *model.ScaledNumberType { 
        return l.Value 
    }, "Value")).
    WithRule(internal.ValidateRange(func(l *model.LoadControlLimitDataType) *model.ScaledNumberType { 
        return l.Value 
    }, 0, 50000, "Power Limit")).
    WithRule(internal.ValidateEnum(func(l *model.LoadControlLimitDataType) *model.LoadControlLimitTypeType { 
        return l.LimitType 
    }, []model.LoadControlLimitTypeType{
        model.LoadControlLimitTypeTypeSignDependentAbsValueLimit,
        model.LoadControlLimitTypeTypeAbsValueLimit,
    }, "LimitType"))

// Validate multiple limits
validLimits, errors := limitValidator.ValidateAll(limits)
```

### Example 3: Complex Multi-Rule Validation

```go
// Complex validator with conditional rules
var complexValidator = internal.NewValidator[*model.YourComplexType]().
    WithName("Complex Validator").
    WithRule(internal.RequireField(func(item *model.YourComplexType) *model.IdType { 
        return item.Id 
    }, "Id")).
    WithRule(internal.ConditionalRule(
        func(item *model.YourComplexType) bool { 
            return item.IsActive != nil && *item.IsActive 
        },
        internal.CombineRules(
            internal.RequireScaledNumber(func(item *model.YourComplexType) *model.ScaledNumberType { 
                return item.Value 
            }, "Value"),
            internal.ValidateRange(func(item *model.YourComplexType) *model.ScaledNumberType { 
                return item.Value 
            }, 0, 1000, "ActiveValue"),
        ),
    )).
    WithRule(internal.ValidateCustom(func(item *model.YourComplexType) error {
        if item.StartTime != nil && item.EndTime != nil && 
           item.StartTime.After(*item.EndTime) {
            return errors.New("StartTime must be before EndTime")
        }
        return nil
    }))
```

## Error Handling

### Error Messages

Validators provide detailed error messages with context:

```go
// Named validator errors include validator name
"Power Validator validation failed at rule 2: Value is required"

// Field-specific errors
"MeasurementId is required"
"Value must be between 0.00 and 1000.00, got 1500.00"
"Status must be one of allowed values"
```

### Validation Results

```go
// ValidateFirst returns first valid item or error
validItem, err := validator.ValidateFirst(items)
if err != nil {
    // No valid items found
    return api.ErrDataNotAvailable
}

// ValidateAll returns valid items and errors for invalid ones
validItems, errors := validator.ValidateAll(items)
for i, err := range errors {
    log.Printf("Item %d validation failed: %v", i, err)
}
```

## Best Practices

### 1. Validator Reuse

Create validators as package-level variables for reuse:

```go
// validators.go
var (
    PowerValidator = internal.NewMeasurementValidator().
        WithName("Power").
        WithRule(internal.RequireMeasurementId()).
        WithRule(internal.RequireMeasurementValue())
        
    EnergyValidator = internal.NewMeasurementValidator().
        WithName("Energy").
        WithRule(internal.RequireMeasurementId()).
        WithRule(internal.RequireMeasurementValue()).
        WithRule(internal.ValidateMeasurementRange(0, 999999))
)
```

### 2. Validator Composition

Build complex validators from simpler rules:

```go
// Base rules
var baseRules = []internal.ValidationRule[*model.MeasurementDataType]{
    internal.RequireMeasurementId(),
    internal.RequireMeasurementValue(),
}

// Specialized validators
var powerValidator = internal.NewMeasurementValidator().
    WithName("Power").
    WithRule(internal.CombineRules(baseRules...)).
    WithRule(internal.ValidateMeasurementRange(0, 50000))
```

### 3. Error Context

Use descriptive names and field names:

```go
validator.WithName("EV Charging Power")  // Appears in error messages
RequireField(getter, "ChargingCurrent")  // Clear field identification
```

### 4. Performance Considerations

- Validators use no reflection - direct field access
- Rules are evaluated in order - put most likely to fail first
- Reuse validators across requests - they're stateless

## Extending the System

### Adding New Data Type Validation

1. **Create validator using generic framework:**

```go
// For any SPINE data type
var yourValidator = internal.NewValidator[*model.YourDataType]().
    WithName("Your Validator").
    WithRule(internal.RequireField(func(item *model.YourDataType) *FieldType { 
        return item.RequiredField 
    }, "RequiredField")).
    WithRule(internal.ValidateEnum(func(item *model.YourDataType) *EnumType { 
        return item.Status 
    }, allowedValues, "Status"))
```

2. **Create helper function for common usage:**

```go
func ValidateYourData(data []model.YourDataType) (*model.YourDataType, error) {
    // Convert to pointer slice if needed
    ptrData := make([]*model.YourDataType, len(data))
    for i := range data {
        ptrData[i] = &data[i]
    }
    
    return yourValidator.ValidateFirst(ptrData)
}
```

### Adding Custom Rules

```go
// Create domain-specific validation rule
func RequireValidTimestamp[T any](getter func(T) *time.Time, fieldName string) internal.ValidationRule[T] {
    return func(item T) error {
        timestamp := getter(item)
        if timestamp == nil {
            return fmt.Errorf("%s is required", fieldName)
        }
        if timestamp.Before(time.Now().Add(-24 * time.Hour)) {
            return fmt.Errorf("%s cannot be older than 24 hours", fieldName)
        }
        return nil
    }
}

// Use custom rule
validator.WithRule(RequireValidTimestamp(func(item *YourType) *time.Time { 
    return item.Timestamp 
}, "Timestamp"))
```

## Migration Guide

### From Old Validation System

**Before:**
```go
data := internal.MeasurementWithMandatoryData(results, true, 
    model.MeasurementValueTypeTypeValue, sources, false, state)
```

**After:**
```go
value, err := internal.GetMeasurementValue(results, internal.EnergyMeasurementValidator)
```

### Migration Steps

1. **Replace old validation calls** with new validator usage
2. **Create custom validators** for specific requirements using generic rules
3. **Update error handling** to use new error types
4. **Remove old validation functions** once migration is complete

## Testing

The validation system includes comprehensive tests covering all rules and edge cases. When adding new validators:

1. **Test all validation rules** individually
2. **Test validator combinations** with valid and invalid data
3. **Test error messages** for clarity and accuracy
4. **Test edge cases** like nil values, empty slices, boundary conditions

See `validation_test.go` for comprehensive examples of testing validation logic.