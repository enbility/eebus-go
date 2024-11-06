package cdt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	usecasesapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Setpoints returns the setpoints.
//
// Possible errors:
//   - ErrDataNotAvailable: If the mapping of operation modes to setpoints or the setpoints themselves are not available.
//   - Other errors: Any other errors encountered during the process.
func (e *CDT) Setpoints(entity spineapi.EntityRemoteInterface) ([]usecasesapi.Setpoint, error) {
	setpoints := make([]usecasesapi.Setpoint, 0)

	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	for _, setpoint := range sp.GetSetpoints() {
		var value float64 = 0
		var minValue float64 = 0
		var maxValue float64 = 0
		var timePeriod model.TimePeriodType = model.TimePeriodType{}

		if setpoint.SetpointId == nil {
			logging.Log().Error("Setpoint ID is nil")
			continue
		}

		if setpoint.Value != nil {
			value = setpoint.Value.GetValue()
		}

		if setpoint.ValueMax != nil {
			maxValue = setpoint.ValueMax.GetValue()
		}

		if setpoint.ValueMin != nil {
			minValue = setpoint.ValueMin.GetValue()
		}

		if setpoint.TimePeriod != nil {
			timePeriod = *setpoint.TimePeriod
		}

		// As per [Resource Specification] 4.3.23.4 setpointListData:
		// - isSetpointActive: If false, the setpoint is inactive; if true or omitted, it is active.
		// - isSetpointChangeable: If true, the server accepts changes; if false, it declines changes. If absent, changes are accepted.
		isActive := (setpoint.IsSetpointActive == nil || *setpoint.IsSetpointActive)
		isChangeable := (setpoint.IsSetpointChangeable == nil || *setpoint.IsSetpointChangeable)

		setpoints = append(setpoints,
			usecasesapi.Setpoint{
				Id:           uint(*setpoint.SetpointId),
				Value:        value,
				MinValue:     minValue,
				MaxValue:     maxValue,
				IsActive:     isActive,
				IsChangeable: isChangeable,
				TimePeriod:   timePeriod,
			},
		)
	}

	if len(setpoints) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return setpoints, nil
}

// SetpointConstraints returns the setpoint constraints.
//
// Possible errors:
//   - ErrDataNotAvailable: If the mapping of operation modes to setpoints or the setpoint constraints are not available.
//   - Other errors: Any other errors encountered during the process.
func (e *CDT) SetpointConstraints(entity spineapi.EntityRemoteInterface) ([]usecasesapi.SetpointConstraints, error) {
	setpointConstraints := make([]usecasesapi.SetpointConstraints, 0)

	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	for _, constraints := range sp.GetSetpointConstraints() {
		var minValue float64 = 0
		var maxValue float64 = 0
		var setSize float64 = 0

		if constraints.SetpointId == nil {
			logging.Log().Error("Setpoint ID is nil")
			continue
		}

		if constraints.SetpointRangeMin != nil {
			minValue = constraints.SetpointRangeMin.GetValue()
		}

		if constraints.SetpointRangeMax != nil {
			maxValue = constraints.SetpointRangeMax.GetValue()
		}

		if constraints.SetpointStepSize != nil {
			setSize = constraints.SetpointStepSize.GetValue()
		}

		setpointConstraints = append(setpointConstraints,
			usecasesapi.SetpointConstraints{
				Id:       uint(*constraints.SetpointId),
				MinValue: minValue,
				MaxValue: maxValue,
				StepSize: setSize,
			},
		)
	}

	if len(setpointConstraints) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return setpointConstraints, nil
}

// mapSetpointsToOperationModes maps setpoints to their respective operation modes.
func (e *CDT) mapSetpointsToModes(entity spineapi.EntityRemoteInterface) error {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return err
	}

	// Get the DHW system functionId for the DHW system function.
	filter := model.HvacSystemFunctionDescriptionDataType{
		SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
	}
	functions, _ := hvac.GetHvacSystemFunctionDescriptionsForFilter(filter)
	if len(functions) == 0 {
		return api.ErrDataNotAvailable
	}

	functionId := *functions[0].SystemFunctionId

	// Get the relations between operation modes and setpoints for the DHW system function.
	relations, _ := hvac.GetHvacSystemFunctionSetpointRelationsForSystemFunctionId(functionId)
	if len(relations) == 0 {
		return api.ErrDataNotAvailable
	}

	// Get the operation mode descriptions for the operation modes in the relations.
	descriptions, _ := hvac.GetHvacOperationModeDescriptions()
	if len(descriptions) == 0 {
		return api.ErrDataNotAvailable
	}

	// Create a mapping to get the operation mode descriptions by operation mode ID.
	modeDescriptions := make(map[model.HvacOperationModeIdType]model.HvacOperationModeTypeType)
	for _, description := range descriptions {
		modeDescriptions[*description.OperationModeId] = *description.OperationModeType
	}

	// Map the setpoints to their respective operation modes.
	for _, relation := range relations {
		if mode, found := modeDescriptions[*relation.OperationModeId]; found {
			if len(relation.SetpointId) == 0 {
				// Only the 'Off' operation mode can have no setpoint associated with it.
				if mode != model.HvacOperationModeTypeTypeOff {
					logging.Log().Errorf("Operation mode '%s' has no setpoints", mode)
				}
			} else if len(relation.SetpointId) == 1 {
				// Unique 1:1 mapping of operation mode to setpoint.
				e.modes[mode] = relation.SetpointId[0]
			} else {
				if mode != model.HvacOperationModeTypeTypeAuto {
					logging.Log().Errorf("Operation mode '%s' has multiple setpoints", mode)
				}
			}
		}
	}

	return nil
}

// WriteSetpoint sets the temperature setpoint for a specific operation mode.
//
// Possible errors:
//   - ErrDataNotAvailable: If the mapping of operation modes to setpoints is not available.
//   - ErrNotSupported: If the setpoint is not changeable.
//   - Other errors: Any other errors encountered during the process.
func (e *CDT) WriteSetpoint(
	entity spineapi.EntityRemoteInterface,
	mode usecasesapi.HvacOperationModeType,
	temperature float64,
) error {
	if model.HvacOperationModeTypeType(mode) == model.HvacOperationModeTypeTypeAuto {
		// 'Auto' mode is controlled by a timetable, meaning the current setpoint
		// for the HVAC system function changes according to the timetable.
		// Only the 'Off', 'On', and 'Eco' modes can be directly controlled by a setpoint.
		return nil
	}

	if len(e.modes) == 0 && e.mapSetpointsToModes(entity) != nil {
		return api.ErrDataNotAvailable
	}

	setpointId, found := e.modes[model.HvacOperationModeTypeType(mode)]
	if !found {
		return api.ErrDataNotAvailable
	}

	setPoint, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return err
	}

	setpointToWrite := []model.SetpointDataType{
		{
			SetpointId: &setpointId,
			Value:      model.NewScaledNumberType(temperature),
		},
	}

	if _, err = setPoint.WriteSetPointListData(setpointToWrite); err != nil {
		return err
	}

	return nil
}
