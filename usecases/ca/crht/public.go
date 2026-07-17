package crht

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the current room heating temperature setpoints of the HVAC room entity,
// returns ErrDataNotAvailable if no such data is (yet) available
func (e *CRHT) Setpoints(entity spineapi.EntityRemoteInterface) ([]ucapi.Setpoint, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	setpointIds, err := e.setpointIds(entity)
	if err != nil {
		return nil, err
	}

	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	setpoints := make([]ucapi.Setpoint, 0)
	for _, id := range setpointIds {
		data, err := sp.GetSetpointForId(id)
		if err != nil {
			continue
		}

		setpoint := ucapi.Setpoint{
			Id: uint(id),
			// if absent, the setpoint is active and changeable
			IsActive:     data.IsSetpointActive == nil || *data.IsSetpointActive,
			IsChangeable: data.IsSetpointChangeable == nil || *data.IsSetpointChangeable,
		}
		if data.Value != nil {
			setpoint.Value = data.Value.GetValue()
		}
		if data.ValueMin != nil {
			setpoint.MinValue = data.ValueMin.GetValue()
		}
		if data.ValueMax != nil {
			setpoint.MaxValue = data.ValueMax.GetValue()
		}

		setpoints = append(setpoints, setpoint)
	}

	if len(setpoints) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return setpoints, nil
}

// return the constraints for the room heating temperature setpoints,
// returns ErrDataNotAvailable if no such data is (yet) available
func (e *CRHT) SetpointConstraints(entity spineapi.EntityRemoteInterface) ([]ucapi.SetpointConstraints, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	setpointIds, err := e.setpointIds(entity)
	if err != nil {
		return nil, err
	}

	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	constraints := make([]ucapi.SetpointConstraints, 0)
	for _, id := range setpointIds {
		data, err := sp.GetSetpointConstraintsForId(id)
		if err != nil {
			continue
		}

		constraint := ucapi.SetpointConstraints{
			Id: uint(id),
		}
		if data.SetpointRangeMin != nil {
			constraint.MinValue = data.SetpointRangeMin.GetValue()
		}
		if data.SetpointRangeMax != nil {
			constraint.MaxValue = data.SetpointRangeMax.GetValue()
		}
		if data.SetpointStepSize != nil {
			constraint.StepSize = data.SetpointStepSize.GetValue()
		}

		constraints = append(constraints, constraint)
	}

	if len(constraints) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return constraints, nil
}

// write the room heating temperature setpoint in degree Celsius for a heating
// operation mode (on, off or eco), returns ErrNotSupported if it is not changeable
func (e *CRHT) WriteSetpoint(
	entity spineapi.EntityRemoteInterface,
	mode ucapi.HvacOperationModeType,
	degC float64,
) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	// the setpoints of the "auto" mode are controlled by a timetable of the device
	if mode == ucapi.HvacOperationModeTypeAuto {
		return api.ErrNotSupported
	}

	setpointId, err := e.setpointIdForMode(entity, mode)
	if err != nil {
		return err
	}

	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return err
	}

	if data, err := sp.GetSetpointForId(setpointId); err == nil &&
		data.IsSetpointChangeable != nil && !*data.IsSetpointChangeable {
		return api.ErrNotSupported
	}

	data := []model.SetpointDataType{
		{
			SetpointId: util.Ptr(setpointId),
			Value:      model.NewScaledNumberType(degC),
		},
	}

	_, err = sp.WriteSetpointListData(data)

	return err
}

// return the ids of the setpoints related to the heating system function
func (e *CRHT) setpointIds(entity spineapi.EntityRemoteInterface) ([]model.SetpointIdType, error) {
	relations, err := e.setpointRelations(entity)
	if err != nil {
		return nil, err
	}

	var ids []model.SetpointIdType
	for _, relation := range relations {
		ids = append(ids, relation.SetpointId...)
	}

	if len(ids) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return ids, nil
}

// return the id of the setpoint related to the given heating operation mode
func (e *CRHT) setpointIdForMode(
	entity spineapi.EntityRemoteInterface,
	mode ucapi.HvacOperationModeType,
) (model.SetpointIdType, error) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	modeFilter := model.HvacOperationModeDescriptionDataType{
		OperationModeType: util.Ptr(model.HvacOperationModeTypeType(mode)),
	}
	modeDescriptions, err := hvac.GetHvacOperationModeDescriptionsForFilter(modeFilter)
	if err != nil || len(modeDescriptions) == 0 || modeDescriptions[0].OperationModeId == nil {
		return 0, api.ErrDataNotAvailable
	}

	relations, err := e.setpointRelations(entity)
	if err != nil {
		return 0, err
	}

	for _, relation := range relations {
		if relation.OperationModeId != nil &&
			*relation.OperationModeId == *modeDescriptions[0].OperationModeId &&
			len(relation.SetpointId) == 1 {
			return relation.SetpointId[0], nil
		}
	}

	return 0, api.ErrDataNotAvailable
}

// return the setpoint relations of the heating system function
func (e *CRHT) setpointRelations(
	entity spineapi.EntityRemoteInterface,
) ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	descFilter := model.HvacSystemFunctionDescriptionDataType{
		SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeHeating),
	}
	descriptions, err := hvac.GetHvacSystemFunctionDescriptionsForFilter(descFilter)
	if err != nil || len(descriptions) == 0 || descriptions[0].SystemFunctionId == nil {
		return nil, api.ErrDataNotAvailable
	}

	relationFilter := model.HvacSystemFunctionSetpointRelationDataType{
		SystemFunctionId: descriptions[0].SystemFunctionId,
	}

	return hvac.GetHvacSystemFunctionSetpointRelationsForFilter(relationFilter)
}
