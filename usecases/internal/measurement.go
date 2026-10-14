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
	validator MeasurementValidator,
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

	var result []float64

	for _, item := range data {
		// Skip measurements that fail validation. This also covers
		// MGCP-003 / MPC-003 (values with state error or outOfRange
		// SHALL be ignored) via the SkipValueState rule in scenario validators.
		if err := validator.Validate(&item); err != nil {
			continue
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

		result = append(result, item.Value.GetValue())
	}

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
type MeasurementValidator struct {
	BaseValidator[*model.MeasurementDataType]
}

// NewMeasurementValidator creates a new measurement validator.
func NewMeasurementValidator() MeasurementValidator {
	return MeasurementValidator{}
}

// WithRule adds a validation rule.
func (v MeasurementValidator) WithRule(rule ValidationRule[*model.MeasurementDataType]) MeasurementValidator {
	v.BaseValidator.rules = append(v.BaseValidator.rules, rule)
	return v
}

// WithName sets the validator name for better error messages.
func (v MeasurementValidator) WithName(name string) MeasurementValidator {
	v.BaseValidator.name = name
	return v
}

// Validate applies all rules to a single measurement.
func (v MeasurementValidator) Validate(m *model.MeasurementDataType) error {
	return v.BaseValidator.Validate(m)
}

// FindFirstValidItem returns the first measurement in items that passes all rules.
func (v MeasurementValidator) FindFirstValidItem(items []*model.MeasurementDataType) (*model.MeasurementDataType, error) {
	return v.BaseValidator.FindFirstValidItem(items)
}

// GetMeasurementValue returns the first valid measurement's numeric value.
// Returns api.ErrDataNotAvailable if no valid measurements are found.
func GetMeasurementValue(measurements []model.MeasurementDataType, validator MeasurementValidator) (float64, error) {
	ptrMeasurements := make([]*model.MeasurementDataType, len(measurements))
	for i := range measurements {
		ptrMeasurements[i] = &measurements[i]
	}

	valid, err := validator.FindFirstValidItem(ptrMeasurements)
	if err != nil || valid == nil || valid.Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	return valid.Value.GetValue(), nil
}

// Measurement-specific validation rules

// RequireMeasurementId ensures MeasurementId is present.
func RequireMeasurementId() ValidationRule[*model.MeasurementDataType] {
	return RequireField(
		func(m *model.MeasurementDataType) *model.MeasurementIdType { return m.MeasurementId },
		"MeasurementId",
	)
}

// RequireMeasurementValue ensures Value is present.
func RequireMeasurementValue() ValidationRule[*model.MeasurementDataType] {
	return RequireScaledNumber(
		func(m *model.MeasurementDataType) *model.ScaledNumberType { return m.Value },
		"Value",
	)
}

// RequireValueType ensures ValueType matches expected type.
func RequireValueType(expected model.MeasurementValueTypeType) ValidationRule[*model.MeasurementDataType] {
	return func(m *model.MeasurementDataType) error {
		if m.ValueType == nil {
			return fmt.Errorf("ValueType is required to be %s", expected)
		}
		if *m.ValueType != expected {
			return fmt.Errorf("ValueType must be %s, got %s", expected, *m.ValueType)
		}
		return nil
	}
}

// RequireValueSource ensures ValueSource is present and one of the allowed types.
// Callers must pass at least one allowed value.
func RequireValueSource(allowed ...model.MeasurementValueSourceType) ValidationRule[*model.MeasurementDataType] {
	if len(allowed) == 0 {
		panic("RequireValueSource requires at least one allowed value")
	}
	return func(m *model.MeasurementDataType) error {
		if m.ValueSource == nil {
			return fmt.Errorf("ValueSource is required")
		}
		if slices.Contains(allowed, *m.ValueSource) {
			return nil
		}
		return fmt.Errorf("ValueSource must be one of %v, got %s", allowed, *m.ValueSource)
	}
}

// ValidateValueState validates the value state, optionally requiring presence.
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

// ValidateMeasurementRange ensures the measurement value is within range.
func ValidateMeasurementRange(minVal, maxVal float64) ValidationRule[*model.MeasurementDataType] {
	return ValidateRange(
		func(m *model.MeasurementDataType) *model.ScaledNumberType { return m.Value },
		minVal, maxVal,
		"Measurement value",
	)
}

// SkipValueState implements MGCP-003 / MPC-003: measurements with ValueState
// "outOfRange" or "error" SHALL be ignored by the Monitoring Appliance.
func SkipValueState() ValidationRule[*model.MeasurementDataType] {
	return func(m *model.MeasurementDataType) error {
		if m.ValueState == nil {
			return nil
		}
		if *m.ValueState == model.MeasurementValueStateTypeError ||
			*m.ValueState == model.MeasurementValueStateTypeOutofrange {
			return ErrSkipMeasurement
		}
		return nil
	}
}
