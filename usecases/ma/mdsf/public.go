package mdsf

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the DHW operation modes supported by the DHW circuit,
// returns ErrDataNotAvailable if no such data is (yet) available
func (e *MDSF) OperationModes(entity spineapi.EntityRemoteInterface) ([]ucapi.HvacOperationModeType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	systemFunctionId, err := e.systemFunctionId(entity)
	if err != nil {
		return nil, err
	}

	relationFilter := model.HvacSystemFunctionOperationModeRelationDataType{
		SystemFunctionId: &systemFunctionId,
	}
	relations, err := hvac.GetHvacSystemFunctionOperationModeRelationsForFilter(relationFilter)
	if err != nil || len(relations) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	modes := make([]ucapi.HvacOperationModeType, 0)
	for _, relation := range relations {
		for _, modeId := range relation.OperationModeId {
			description, err := hvac.GetHvacOperationModeDescriptionForId(modeId)
			if err != nil || description.OperationModeType == nil {
				continue
			}

			modes = append(modes, ucapi.HvacOperationModeType(*description.OperationModeType))
		}
	}

	if len(modes) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return modes, nil
}

// return the current DHW operation mode of the DHW circuit,
// returns ErrDataNotAvailable if no such data is (yet) available
func (e *MDSF) CurrentOperationMode(entity spineapi.EntityRemoteInterface) (ucapi.HvacOperationModeType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return "", api.ErrNoCompatibleEntity
	}

	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return "", err
	}

	systemFunctionId, err := e.systemFunctionId(entity)
	if err != nil {
		return "", err
	}

	data, err := hvac.GetHvacSystemFunctionForId(systemFunctionId)
	if err != nil || data.CurrentOperationModeId == nil {
		return "", api.ErrDataNotAvailable
	}

	description, err := hvac.GetHvacOperationModeDescriptionForId(*data.CurrentOperationModeId)
	if err != nil || description.OperationModeType == nil {
		return "", api.ErrDataNotAvailable
	}

	return ucapi.HvacOperationModeType(*description.OperationModeType), nil
}

// Scenario 2

// return whether an overrun affecting the DHW system function is active,
// returns ErrDataNotAvailable if no such data is (yet) available
func (e *MDSF) IsOverrunActive(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return false, err
	}

	systemFunctionId, err := e.systemFunctionId(entity)
	if err != nil {
		return false, err
	}

	data, err := hvac.GetHvacSystemFunctionForId(systemFunctionId)
	if err != nil || data.IsOverrunActive == nil {
		return false, api.ErrDataNotAvailable
	}

	return *data.IsOverrunActive, nil
}

// return the status of the one-time DHW overrun,
// returns ErrDataNotAvailable if no such data is (yet) available
func (e *MDSF) OverrunStatus(entity spineapi.EntityRemoteInterface) (model.HvacOverrunStatusType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return "", api.ErrNoCompatibleEntity
	}

	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return "", err
	}

	overrunId, err := e.overrunId(entity)
	if err != nil {
		return "", err
	}

	data, err := hvac.GetHvacOverrunForId(overrunId)
	if err != nil || data.OverrunStatus == nil {
		return "", api.ErrDataNotAvailable
	}

	return *data.OverrunStatus, nil
}

// return the id of the one-time DHW overrun affecting the DHW system function
func (e *MDSF) overrunId(entity spineapi.EntityRemoteInterface) (model.HvacOverrunIdType, error) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	systemFunctionId, err := e.systemFunctionId(entity)
	if err != nil {
		return 0, err
	}

	descFilter := model.HvacOverrunDescriptionDataType{
		OverrunType: util.Ptr(model.HvacOverrunTypeTypeOneTimeDhw),
	}
	descriptions, err := hvac.GetHvacOverrunDescriptionsForFilter(descFilter)
	if err != nil {
		return 0, api.ErrDataNotAvailable
	}

	for _, description := range descriptions {
		if description.OverrunId == nil {
			continue
		}

		for _, affectedId := range description.AffectedSystemFunctionId {
			if affectedId == systemFunctionId {
				return *description.OverrunId, nil
			}
		}
	}

	return 0, api.ErrDataNotAvailable
}

// return the id of the DHW system function of the DHW circuit
func (e *MDSF) systemFunctionId(entity spineapi.EntityRemoteInterface) (model.HvacSystemFunctionIdType, error) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	descFilter := model.HvacSystemFunctionDescriptionDataType{
		SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
	}
	descriptions, err := hvac.GetHvacSystemFunctionDescriptionsForFilter(descFilter)
	if err != nil || len(descriptions) == 0 || descriptions[0].SystemFunctionId == nil {
		return 0, api.ErrDataNotAvailable
	}

	return *descriptions[0].SystemFunctionId, nil
}
