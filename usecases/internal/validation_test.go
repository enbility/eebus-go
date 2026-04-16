package internal

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

// Helper function to create pointers
func ptr[T any](v T) *T {
	return &v
}

// Test structs for generic validation testing
type TestStruct struct {
	ID          *int
	Value       *model.ScaledNumberType
	Name        *string
	Status      *string
	Items       []string
	Min         *model.ScaledNumberType
	Max         *model.ScaledNumberType
	IsActive    *bool
	LastUpdated *time.Time
}

func TestBaseValidator(t *testing.T) {
	t.Run("single rule validation", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithName("TestValidator").
			WithRule(func(ts *TestStruct) error {
				if ts.ID == nil {
					return assert.AnError
				}
				return nil
			})

		// Should fail - nil ID
		err := validator.Validate(&TestStruct{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TestValidator validation failed")

		// Should pass
		err = validator.Validate(&TestStruct{ID: ptr(1)})
		assert.NoError(t, err)
	})

	t.Run("multiple rules", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID")).
			WithRule(RequireField(func(ts *TestStruct) *string { return ts.Name }, "Name"))

		// Should fail on first rule
		err := validator.Validate(&TestStruct{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")

		// Should fail on second rule
		err = validator.Validate(&TestStruct{ID: ptr(1)})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name is required")

		// Should pass
		err = validator.Validate(&TestStruct{ID: ptr(1), Name: ptr("test")})
		assert.NoError(t, err)
	})

	t.Run("WithName sets validator name", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithName("Custom Validator").
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"))

		err := validator.Validate(&TestStruct{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Custom Validator validation failed")
	})

	t.Run("ValidateFirst", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"))

		items := []*TestStruct{
			{},           // Invalid
			{ID: ptr(1)}, // Valid
			{ID: ptr(2)}, // Also valid but shouldn't be returned
		}

		result, err := validator.ValidateFirst(items)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, *result.ID)

		// No valid items
		invalidItems := []*TestStruct{{}, {}}
		_, err = validator.ValidateFirst(invalidItems)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no valid item found")
	})

	t.Run("ValidateFirst with empty slice", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"))

		items := []*TestStruct{}
		_, err := validator.ValidateFirst(items)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no valid item found")
	})

	t.Run("ValidateAll", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"))

		items := []*TestStruct{
			{},           // Invalid
			{ID: ptr(1)}, // Valid
			{},           // Invalid
			{ID: ptr(2)}, // Valid
		}

		valid, errors := validator.ValidateAll(items)
		assert.Len(t, valid, 2)
		assert.Len(t, errors, 2)
		assert.Equal(t, 1, *valid[0].ID)
		assert.Equal(t, 2, *valid[1].ID)
		assert.Contains(t, errors[0].Error(), "item 0")
		assert.Contains(t, errors[1].Error(), "item 2")
	})

	t.Run("ValidateAll with all valid items", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"))

		items := []*TestStruct{
			{ID: ptr(1)},
			{ID: ptr(2)},
		}

		valid, errors := validator.ValidateAll(items)
		assert.Len(t, valid, 2)
		assert.Len(t, errors, 0)
	})

	t.Run("ValidateAll with all invalid items", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"))

		items := []*TestStruct{{}, {}}

		valid, errors := validator.ValidateAll(items)
		assert.Len(t, valid, 0)
		assert.Len(t, errors, 2)
	})
}

func TestRequireField(t *testing.T) {
	t.Run("nil field fails", func(t *testing.T) {
		rule := RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID")

		err := rule(&TestStruct{})
		assert.Error(t, err)
		assert.Equal(t, "ID is required", err.Error())
	})

	t.Run("non-nil field passes", func(t *testing.T) {
		rule := RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID")

		err := rule(&TestStruct{ID: ptr(42)})
		assert.NoError(t, err)
	})

	t.Run("different field types", func(t *testing.T) {
		// String field
		stringRule := RequireField(func(ts *TestStruct) *string { return ts.Name }, "Name")
		err := stringRule(&TestStruct{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name is required")

		err = stringRule(&TestStruct{Name: ptr("test")})
		assert.NoError(t, err)

		// Bool field
		boolRule := RequireField(func(ts *TestStruct) *bool { return ts.IsActive }, "IsActive")
		err = boolRule(&TestStruct{})
		assert.Error(t, err)

		err = boolRule(&TestStruct{IsActive: ptr(true)})
		assert.NoError(t, err)
	})
}

func TestRequireScaledNumber(t *testing.T) {
	t.Run("nil ScaledNumber fails", func(t *testing.T) {
		rule := RequireScaledNumber(func(ts *TestStruct) *model.ScaledNumberType { return ts.Value }, "Value")

		err := rule(&TestStruct{})
		assert.Error(t, err)
		assert.Equal(t, "Value is required", err.Error())
	})

	t.Run("non-nil ScaledNumber passes", func(t *testing.T) {
		rule := RequireScaledNumber(func(ts *TestStruct) *model.ScaledNumberType { return ts.Value }, "Value")

		err := rule(&TestStruct{Value: model.NewScaledNumberType(42)})
		assert.NoError(t, err)
	})

	t.Run("custom field name in error", func(t *testing.T) {
		rule := RequireScaledNumber(func(ts *TestStruct) *model.ScaledNumberType { return ts.Min }, "Minimum")

		err := rule(&TestStruct{})
		assert.Error(t, err)
		assert.Equal(t, "Minimum is required", err.Error())
	})
}

func TestValidateRange(t *testing.T) {
	rule := ValidateRange(
		func(ts *TestStruct) *model.ScaledNumberType { return ts.Value },
		10, 100,
		"Value",
	)

	t.Run("nil value passes (skipped)", func(t *testing.T) {
		err := rule(&TestStruct{})
		assert.NoError(t, err)
	})

	t.Run("value below min fails", func(t *testing.T) {
		err := rule(&TestStruct{Value: model.NewScaledNumberType(5)})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Value must be between 10.00 and 100.00, got 5.00")
	})

	t.Run("value above max fails", func(t *testing.T) {
		err := rule(&TestStruct{Value: model.NewScaledNumberType(150)})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Value must be between 10.00 and 100.00, got 150.00")
	})

	t.Run("value at min boundary passes", func(t *testing.T) {
		err := rule(&TestStruct{Value: model.NewScaledNumberType(10)})
		assert.NoError(t, err)
	})

	t.Run("value at max boundary passes", func(t *testing.T) {
		err := rule(&TestStruct{Value: model.NewScaledNumberType(100)})
		assert.NoError(t, err)
	})

	t.Run("value in range passes", func(t *testing.T) {
		err := rule(&TestStruct{Value: model.NewScaledNumberType(50)})
		assert.NoError(t, err)
	})
}

func TestValidateMinMax(t *testing.T) {
	rule := ValidateMinMax(
		func(ts *TestStruct) *model.ScaledNumberType { return ts.Value },
		func(ts *TestStruct) *model.ScaledNumberType { return ts.Min },
		func(ts *TestStruct) *model.ScaledNumberType { return ts.Max },
		"Value",
	)

	t.Run("nil value passes (skipped)", func(t *testing.T) {
		err := rule(&TestStruct{})
		assert.NoError(t, err)
	})

	t.Run("value below min fails", func(t *testing.T) {
		err := rule(&TestStruct{
			Value: model.NewScaledNumberType(5),
			Min:   model.NewScaledNumberType(10),
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Value 5.00 is below minimum 10.00")
	})

	t.Run("value above max fails", func(t *testing.T) {
		err := rule(&TestStruct{
			Value: model.NewScaledNumberType(100),
			Max:   model.NewScaledNumberType(50),
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Value 100.00 is above maximum 50.00")
	})

	t.Run("value within range passes", func(t *testing.T) {
		err := rule(&TestStruct{
			Value: model.NewScaledNumberType(50),
			Min:   model.NewScaledNumberType(10),
			Max:   model.NewScaledNumberType(100),
		})
		assert.NoError(t, err)
	})

	t.Run("no min/max constraints passes", func(t *testing.T) {
		err := rule(&TestStruct{
			Value: model.NewScaledNumberType(50),
		})
		assert.NoError(t, err)
	})

	t.Run("only min constraint", func(t *testing.T) {
		err := rule(&TestStruct{
			Value: model.NewScaledNumberType(50),
			Min:   model.NewScaledNumberType(10),
		})
		assert.NoError(t, err)

		err = rule(&TestStruct{
			Value: model.NewScaledNumberType(5),
			Min:   model.NewScaledNumberType(10),
		})
		assert.Error(t, err)
	})

	t.Run("only max constraint", func(t *testing.T) {
		err := rule(&TestStruct{
			Value: model.NewScaledNumberType(50),
			Max:   model.NewScaledNumberType(100),
		})
		assert.NoError(t, err)

		err = rule(&TestStruct{
			Value: model.NewScaledNumberType(150),
			Max:   model.NewScaledNumberType(100),
		})
		assert.Error(t, err)
	})
}

func TestValidateEnum(t *testing.T) {
	allowed := []string{"active", "inactive", "pending"}
	rule := ValidateEnum(
		func(ts *TestStruct) *string { return ts.Status },
		allowed,
		"Status",
	)

	t.Run("nil value passes (skipped)", func(t *testing.T) {
		err := rule(&TestStruct{})
		assert.NoError(t, err)
	})

	t.Run("allowed value passes", func(t *testing.T) {
		err := rule(&TestStruct{Status: ptr("active")})
		assert.NoError(t, err)

		err = rule(&TestStruct{Status: ptr("inactive")})
		assert.NoError(t, err)

		err = rule(&TestStruct{Status: ptr("pending")})
		assert.NoError(t, err)
	})

	t.Run("disallowed value fails", func(t *testing.T) {
		err := rule(&TestStruct{Status: ptr("unknown")})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Status must be one of allowed values")
	})

	t.Run("empty allowed list passes any value", func(t *testing.T) {
		emptyRule := ValidateEnum(
			func(ts *TestStruct) *string { return ts.Status },
			[]string{},
			"Status",
		)

		err := emptyRule(&TestStruct{Status: ptr("anything")})
		assert.NoError(t, err)
	})
}

func TestValidateSliceNotEmpty(t *testing.T) {
	rule := ValidateSliceNotEmpty(
		func(ts *TestStruct) []string { return ts.Items },
		"Items",
	)

	t.Run("empty slice fails", func(t *testing.T) {
		err := rule(&TestStruct{Items: []string{}})
		assert.Error(t, err)
		assert.Equal(t, "Items cannot be empty", err.Error())
	})

	t.Run("nil slice fails", func(t *testing.T) {
		err := rule(&TestStruct{})
		assert.Error(t, err)
		assert.Equal(t, "Items cannot be empty", err.Error())
	})

	t.Run("non-empty slice passes", func(t *testing.T) {
		err := rule(&TestStruct{Items: []string{"item1"}})
		assert.NoError(t, err)

		err = rule(&TestStruct{Items: []string{"item1", "item2"}})
		assert.NoError(t, err)
	})
}

func TestValidateCustom(t *testing.T) {
	rule := ValidateCustom(func(ts *TestStruct) error {
		if ts.ID != nil && *ts.ID < 0 {
			return assert.AnError
		}
		return nil
	})

	t.Run("custom validation passes", func(t *testing.T) {
		err := rule(&TestStruct{ID: ptr(42)})
		assert.NoError(t, err)

		err = rule(&TestStruct{})
		assert.NoError(t, err)
	})

	t.Run("custom validation fails", func(t *testing.T) {
		err := rule(&TestStruct{ID: ptr(-1)})
		assert.Error(t, err)
	})
}

func TestCombineRules(t *testing.T) {
	combinedRule := CombineRules(
		RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID"),
		RequireField(func(ts *TestStruct) *string { return ts.Name }, "Name"),
		ValidateCustom(func(ts *TestStruct) error {
			if ts.ID != nil && *ts.ID < 0 {
				return assert.AnError
			}
			return nil
		}),
	)

	t.Run("all rules pass", func(t *testing.T) {
		err := combinedRule(&TestStruct{ID: ptr(1), Name: ptr("test")})
		assert.NoError(t, err)
	})

	t.Run("first rule fails", func(t *testing.T) {
		err := combinedRule(&TestStruct{Name: ptr("test")})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("second rule fails", func(t *testing.T) {
		err := combinedRule(&TestStruct{ID: ptr(1)})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name is required")
	})

	t.Run("third rule fails", func(t *testing.T) {
		err := combinedRule(&TestStruct{ID: ptr(-1), Name: ptr("test")})
		assert.Error(t, err)
	})
}

func TestConditionalRule(t *testing.T) {
	rule := ConditionalRule(
		func(ts *TestStruct) bool { return ts.IsActive != nil && *ts.IsActive },
		RequireField(func(ts *TestStruct) *string { return ts.Name }, "Name"),
	)

	t.Run("condition not met, rule skipped", func(t *testing.T) {
		// IsActive is false, so rule should be skipped
		err := rule(&TestStruct{IsActive: ptr(false)})
		assert.NoError(t, err)

		// IsActive is nil, so rule should be skipped
		err = rule(&TestStruct{})
		assert.NoError(t, err)
	})

	t.Run("condition met, rule applied and passes", func(t *testing.T) {
		err := rule(&TestStruct{IsActive: ptr(true), Name: ptr("test")})
		assert.NoError(t, err)
	})

	t.Run("condition met, rule applied and fails", func(t *testing.T) {
		err := rule(&TestStruct{IsActive: ptr(true)})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name is required")
	})
}

func TestComplexValidationScenarios(t *testing.T) {
	t.Run("complex validator with multiple rule types", func(t *testing.T) {
		validator := NewValidator[*TestStruct]().
			WithName("Complex Validator").
			WithRule(RequireField(func(ts *TestStruct) *int { return ts.ID }, "ID")).
			WithRule(RequireScaledNumber(func(ts *TestStruct) *model.ScaledNumberType { return ts.Value }, "Value")).
			WithRule(ValidateRange(
				func(ts *TestStruct) *model.ScaledNumberType { return ts.Value },
				0, 1000,
				"Value",
			)).
			WithRule(ValidateEnum(
				func(ts *TestStruct) *string { return ts.Status },
				[]string{"active", "inactive"},
				"Status",
			)).
			WithRule(ConditionalRule(
				func(ts *TestStruct) bool { return ts.IsActive != nil && *ts.IsActive },
				RequireField(func(ts *TestStruct) *string { return ts.Name }, "Name"),
			))

		// Should pass all validations
		testData := &TestStruct{
			ID:       ptr(1),
			Value:    model.NewScaledNumberType(500),
			Status:   ptr("active"),
			IsActive: ptr(true),
			Name:     ptr("test"),
		}
		err := validator.Validate(testData)
		assert.NoError(t, err)

		// Should fail on value range
		testData.Value = model.NewScaledNumberType(2000)
		err = validator.Validate(testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 0.00 and 1000.00")

		// Should fail on status enum
		testData.Value = model.NewScaledNumberType(500)
		testData.Status = ptr("unknown")
		err = validator.Validate(testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Status must be one of allowed values")

		// Should fail on conditional rule (IsActive=true but no Name)
		testData.Status = ptr("active")
		testData.Name = nil
		err = validator.Validate(testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name is required")

		// Should pass when IsActive=false (conditional rule skipped)
		testData.IsActive = ptr(false)
		err = validator.Validate(testData)
		assert.NoError(t, err)
	})

	t.Run("validator reuse with different types", func(t *testing.T) {
		// Define a simple struct for testing
		type SimpleStruct struct {
			ID   *int
			Name *string
		}

		// Create validator for SimpleStruct
		simpleValidator := NewValidator[*SimpleStruct]().
			WithRule(RequireField(func(s *SimpleStruct) *int { return s.ID }, "ID")).
			WithRule(RequireField(func(s *SimpleStruct) *string { return s.Name }, "Name"))

		// Should work with SimpleStruct
		err := simpleValidator.Validate(&SimpleStruct{ID: ptr(1), Name: ptr("test")})
		assert.NoError(t, err)

		err = simpleValidator.Validate(&SimpleStruct{ID: ptr(1)})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Name is required")
	})
}
