package cdsf

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
func (e *CDSF) OperationModes(entity spineapi.EntityRemoteInterface) ([]ucapi.HvacOperationModeType, error) {
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
func (e *CDSF) CurrentOperationMode(entity spineapi.EntityRemoteInterface) (ucapi.HvacOperationModeType, error) {
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

// set the DHW operation mode of the DHW circuit,
// returns ErrNotSupported if the operation mode is not changeable or not supported
func (e *CDSF) WriteOperationMode(entity spineapi.EntityRemoteInterface, mode ucapi.HvacOperationModeType) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return err
	}

	systemFunctionId, err := e.systemFunctionId(entity)
	if err != nil {
		return err
	}

	data, err := hvac.GetHvacSystemFunctionForId(systemFunctionId)
	if err != nil {
		return api.ErrDataNotAvailable
	}

	// only an explicit false blocks the write; an omitted flag is tolerated, as
	// some devices accept the write without advertising the changeability flag
	if data.IsOperationModeIdChangeable != nil && !*data.IsOperationModeIdChangeable {
		return api.ErrNotSupported
	}

	// resolve the requested mode through the DHW system-function relation, so a
	// mode that exists globally but is not related to this function is rejected
	relationFilter := model.HvacSystemFunctionOperationModeRelationDataType{
		SystemFunctionId: &systemFunctionId,
	}
	relations, err := hvac.GetHvacSystemFunctionOperationModeRelationsForFilter(relationFilter)
	if err != nil || len(relations) == 0 {
		return api.ErrDataNotAvailable
	}

	var modeId *model.HvacOperationModeIdType
	for _, relation := range relations {
		for _, id := range relation.OperationModeId {
			description, err := hvac.GetHvacOperationModeDescriptionForId(id)
			if err != nil || description.OperationModeType == nil {
				continue
			}
			if ucapi.HvacOperationModeType(*description.OperationModeType) == mode {
				modeId = util.Ptr(id)
				break
			}
		}
		if modeId != nil {
			break
		}
	}
	if modeId == nil {
		return api.ErrNotSupported
	}

	writeData := []model.HvacSystemFunctionDataType{
		{
			SystemFunctionId:       &systemFunctionId,
			CurrentOperationModeId: modeId,
		},
	}

	_, err = hvac.WriteHvacSystemFunctionListData(writeData)

	return err
}

// Scenario 2

// start the one-time DHW loading overrun of the DHW circuit,
// returns ErrNotSupported if the overrun status is not changeable
func (e *CDSF) StartOneTimeDhw(entity spineapi.EntityRemoteInterface) error {
	return e.writeOverrunStatus(entity, model.HvacOverrunStatusTypeActive)
}

// Scenario 3

// stop the one-time DHW loading overrun of the DHW circuit,
// returns ErrNotSupported if the overrun status is not changeable
func (e *CDSF) StopOneTimeDhw(entity spineapi.EntityRemoteInterface) error {
	return e.writeOverrunStatus(entity, model.HvacOverrunStatusTypeInactive)
}

// write the status of the one-time DHW overrun
func (e *CDSF) writeOverrunStatus(
	entity spineapi.EntityRemoteInterface,
	status model.HvacOverrunStatusType,
) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return err
	}

	overrunId, err := e.overrunId(entity)
	if err != nil {
		return err
	}

	if data, err := hvac.GetHvacOverrunForId(overrunId); err == nil &&
		data.IsOverrunStatusChangeable != nil && !*data.IsOverrunStatusChangeable {
		return api.ErrNotSupported
	}

	writeData := []model.HvacOverrunDataType{
		{
			OverrunId:     &overrunId,
			OverrunStatus: &status,
		},
	}

	_, err = hvac.WriteHvacOverrunListData(writeData)

	return err
}

// return the id of the one-time DHW overrun affecting the DHW system function
func (e *CDSF) overrunId(entity spineapi.EntityRemoteInterface) (model.HvacOverrunIdType, error) {
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
func (e *CDSF) systemFunctionId(entity spineapi.EntityRemoteInterface) (model.HvacSystemFunctionIdType, error) {
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
