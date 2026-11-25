package cdt

import (
	"math"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Setpoints returns the current DHW temperature setpoints of the DHW circuit entity.
func (e *CDT) Setpoints(entity spineapi.EntityRemoteInterface) ([]ucapi.Setpoint, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	setpointIDs, err := e.setpointIDs(entity)
	if err != nil {
		return nil, err
	}

	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	setpoints := make([]ucapi.Setpoint, 0, len(setpointIDs))
	for _, id := range setpointIDs {
		data, err := sp.GetSetpointForId(id)
		if err != nil || data.Value == nil {
			continue
		}

		setpoint := ucapi.Setpoint{
			Id:           uint(id),
			Value:        data.Value.GetValue(),
			IsActive:     data.IsSetpointActive == nil || *data.IsSetpointActive,
			IsChangeable: sp.IsSetpointListDataWritable() && (data.IsSetpointChangeable == nil || *data.IsSetpointChangeable),
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

// SetpointConstraints returns complete constraints for the DHW temperature setpoints.
func (e *CDT) SetpointConstraints(entity spineapi.EntityRemoteInterface) ([]ucapi.SetpointConstraints, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	setpointIDs, err := e.setpointIDs(entity)
	if err != nil {
		return nil, err
	}
	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	constraints := make([]ucapi.SetpointConstraints, 0, len(setpointIDs))
	for _, id := range setpointIDs {
		data, err := sp.GetSetpointConstraintsForId(id)
		if err != nil || !validConstraints(data) {
			continue
		}
		constraints = append(constraints, ucapi.SetpointConstraints{
			Id:       uint(id),
			MinValue: data.SetpointRangeMin.GetValue(),
			MaxValue: data.SetpointRangeMax.GetValue(),
			StepSize: data.SetpointStepSize.GetValue(),
		})
	}

	if len(constraints) == 0 {
		return nil, api.ErrDataNotAvailable
	}
	return constraints, nil
}

// WriteSetpoint writes a DHW temperature setpoint and preserves all other
// entries from the remote SetpointListData cache for devices requiring full-list writes.
func (e *CDT) WriteSetpoint(
	entity spineapi.EntityRemoteInterface,
	mode ucapi.HvacOperationModeType,
	degC float64,
	resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType),
) (*model.MsgCounterType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	setpointID, err := e.setpointIDForMode(entity, mode)
	if err != nil {
		return nil, err
	}
	sp, err := client.NewSetpoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}
	if !sp.IsSetpointListDataWritable() {
		return nil, api.ErrNotSupported
	}

	current, err := sp.GetSetpointForId(setpointID)
	if err != nil {
		return nil, err
	}
	if current.IsSetpointChangeable != nil && !*current.IsSetpointChangeable {
		return nil, api.ErrNotSupported
	}

	constraints, err := sp.GetSetpointConstraintsForId(setpointID)
	if err != nil || !validConstraints(constraints) || !valueFitsConstraints(degC, constraints) {
		return nil, api.ErrDataInvalid
	}

	all, err := sp.GetSetpointDataForFilter(model.SetpointDataType{})
	if err != nil {
		return nil, err
	}
	writeData := append([]model.SetpointDataType(nil), all...)
	matches := 0
	for i := range writeData {
		if writeData[i].SetpointId != nil && *writeData[i].SetpointId == setpointID {
			writeData[i].Value = model.NewScaledNumberType(degC)
			matches++
		}
	}
	if matches != 1 {
		return nil, api.ErrDataInvalid
	}

	msgCounter, err := sp.WriteSetpointListData(writeData)
	if err != nil || msgCounter == nil {
		if err != nil {
			return nil, err
		}
		return nil, api.ErrDataInvalid
	}
	if resultCB != nil {
		counter := *msgCounter
		if err := sp.AddResponseCallback(counter, func(msg spineapi.ResponseMessage) {
			if result, ok := msg.Data.(*model.ResultDataType); ok {
				if result.ErrorNumber != nil && *result.ErrorNumber == model.ErrorNumberTypeNoError {
					_, _ = sp.RequestSetpoints(nil, nil)
				}
				resultCB(*result, counter)
			}
		}); err != nil {
			return msgCounter, err
		}
	}
	return msgCounter, nil
}

func validConstraints(data *model.SetpointConstraintsDataType) bool {
	if data == nil || data.SetpointRangeMin == nil || data.SetpointRangeMax == nil || data.SetpointStepSize == nil {
		return false
	}
	minValue := data.SetpointRangeMin.GetValue()
	maxValue := data.SetpointRangeMax.GetValue()
	step := data.SetpointStepSize.GetValue()
	return !math.IsNaN(minValue) && !math.IsNaN(maxValue) && !math.IsNaN(step) &&
		!math.IsInf(minValue, 0) && !math.IsInf(maxValue, 0) && !math.IsInf(step, 0) &&
		minValue <= maxValue && step > 0
}

func valueFitsConstraints(value float64, data *model.SetpointConstraintsDataType) bool {
	if !validConstraints(data) || math.IsNaN(value) || math.IsInf(value, 0) {
		return false
	}
	minValue := data.SetpointRangeMin.GetValue()
	maxValue := data.SetpointRangeMax.GetValue()
	step := data.SetpointStepSize.GetValue()
	if value < minValue || value > maxValue {
		return false
	}
	steps := (value - minValue) / step
	return math.Abs(steps-math.Round(steps)) <= 1e-9
}

func (e *CDT) setpointIDs(entity spineapi.EntityRemoteInterface) ([]model.SetpointIdType, error) {
	relations, err := e.setpointRelations(entity)
	if err != nil {
		return nil, err
	}
	seen := make(map[model.SetpointIdType]struct{})
	ids := make([]model.SetpointIdType, 0)
	for _, relation := range relations {
		for _, id := range relation.SetpointId {
			if _, exists := seen[id]; !exists {
				seen[id] = struct{}{}
				ids = append(ids, id)
			}
		}
	}
	if len(ids) == 0 {
		return nil, api.ErrDataNotAvailable
	}
	return ids, nil
}

func (e *CDT) setpointIDForMode(entity spineapi.EntityRemoteInterface, mode ucapi.HvacOperationModeType) (model.SetpointIdType, error) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}
	relations, err := e.setpointRelations(entity)
	if err != nil {
		return 0, err
	}

	ids := make(map[model.SetpointIdType]struct{})
	for _, relation := range relations {
		if relation.OperationModeId == nil {
			continue
		}
		descriptions, err := hvac.GetHvacOperationModeDescriptionsForFilter(model.HvacOperationModeDescriptionDataType{
			OperationModeId: relation.OperationModeId,
		})
		if err != nil || len(descriptions) != 1 || descriptions[0].OperationModeType == nil ||
			*descriptions[0].OperationModeType != model.HvacOperationModeTypeType(mode) {
			continue
		}
		for _, id := range relation.SetpointId {
			ids[id] = struct{}{}
		}
	}
	if len(ids) != 1 {
		return 0, api.ErrDataNotAvailable
	}
	for id := range ids {
		return id, nil
	}
	return 0, api.ErrDataNotAvailable
}

func (e *CDT) setpointRelations(entity spineapi.EntityRemoteInterface) ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}
	descriptions, err := hvac.GetHvacSystemFunctionDescriptionsForFilter(model.HvacSystemFunctionDescriptionDataType{
		SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
	})
	if err != nil || len(descriptions) != 1 || descriptions[0].SystemFunctionId == nil {
		return nil, api.ErrDataNotAvailable
	}
	relations, err := hvac.GetHvacSystemFunctionSetpointRelationsForFilter(model.HvacSystemFunctionSetpointRelationDataType{
		SystemFunctionId: descriptions[0].SystemFunctionId,
	})
	if err != nil || len(relations) == 0 {
		return nil, api.ErrDataNotAvailable
	}
	return relations, nil
}
