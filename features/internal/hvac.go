package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type HvacCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

// NewLocalHvac creates a new HvacCommon helper for local entities
func NewLocalHvac(featureLocal spineapi.FeatureLocalInterface) *HvacCommon {
	return &HvacCommon{
		featureLocal: featureLocal,
	}
}

// NewRemoteHvac creates a new HvacCommon helper for remote entities
func NewRemoteHvac(featureRemote spineapi.FeatureRemoteInterface) *HvacCommon {
	return &HvacCommon{
		featureRemote: featureRemote,
	}
}

// GetHvacSystemFunctionDescriptionsForFilter returns the system function descriptions for a given filter
func (h *HvacCommon) GetHvacSystemFunctionDescriptionsForFilter(
	filter model.HvacSystemFunctionDescriptionDataType,
) ([]model.HvacSystemFunctionDescriptionDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionDescriptionListData

	data, err := featureDataCopyOfType[model.HvacSystemFunctionDescriptionListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacSystemFunctionDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacSystemFunctionDescriptionDataType](data.HvacSystemFunctionDescriptionData, filter)

	return result, nil
}

// GetHvacSystemFunctionDescriptionForId returns the system function description for a given system function ID
func (h *HvacCommon) GetHvacSystemFunctionDescriptionForId(
	id model.HvacSystemFunctionIdType,
) (*model.HvacSystemFunctionDescriptionDataType, error) {
	filter := model.HvacSystemFunctionDescriptionDataType{
		SystemFunctionId: &id,
	}

	result, err := h.GetHvacSystemFunctionDescriptionsForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetHvacSystemFunctionsForFilter returns the system function data for a given filter
func (h *HvacCommon) GetHvacSystemFunctionsForFilter(
	filter model.HvacSystemFunctionDataType,
) ([]model.HvacSystemFunctionDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionListData

	data, err := featureDataCopyOfType[model.HvacSystemFunctionListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacSystemFunctionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacSystemFunctionDataType](data.HvacSystemFunctionData, filter)

	return result, nil
}

// GetHvacSystemFunctionForId returns the system function data for a given system function ID
func (h *HvacCommon) GetHvacSystemFunctionForId(
	id model.HvacSystemFunctionIdType,
) (*model.HvacSystemFunctionDataType, error) {
	filter := model.HvacSystemFunctionDataType{
		SystemFunctionId: &id,
	}

	result, err := h.GetHvacSystemFunctionsForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetHvacOperationModeDescriptionsForFilter returns the operation mode descriptions for a given filter
func (h *HvacCommon) GetHvacOperationModeDescriptionsForFilter(
	filter model.HvacOperationModeDescriptionDataType,
) ([]model.HvacOperationModeDescriptionDataType, error) {
	function := model.FunctionTypeHvacOperationModeDescriptionListData

	data, err := featureDataCopyOfType[model.HvacOperationModeDescriptionListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacOperationModeDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacOperationModeDescriptionDataType](data.HvacOperationModeDescriptionData, filter)

	return result, nil
}

// GetHvacOperationModeDescriptionForId returns the operation mode description for a given operation mode ID
func (h *HvacCommon) GetHvacOperationModeDescriptionForId(
	id model.HvacOperationModeIdType,
) (*model.HvacOperationModeDescriptionDataType, error) {
	filter := model.HvacOperationModeDescriptionDataType{
		OperationModeId: &id,
	}

	result, err := h.GetHvacOperationModeDescriptionsForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetHvacSystemFunctionOperationModeRelationsForFilter returns the system function
// operation mode relations for a given filter
func (h *HvacCommon) GetHvacSystemFunctionOperationModeRelationsForFilter(
	filter model.HvacSystemFunctionOperationModeRelationDataType,
) ([]model.HvacSystemFunctionOperationModeRelationDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionOperationModeRelationListData

	data, err := featureDataCopyOfType[model.HvacSystemFunctionOperationModeRelationListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacSystemFunctionOperationModeRelationData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacSystemFunctionOperationModeRelationDataType](data.HvacSystemFunctionOperationModeRelationData, filter)

	return result, nil
}

// GetHvacSystemFunctionSetpointRelationsForFilter returns the system function
// setpoint relations for a given filter
func (h *HvacCommon) GetHvacSystemFunctionSetpointRelationsForFilter(
	filter model.HvacSystemFunctionSetpointRelationDataType,
) ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionSetPointRelationListData

	data, err := featureDataCopyOfType[model.HvacSystemFunctionSetpointRelationListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacSystemFunctionSetpointRelationData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacSystemFunctionSetpointRelationDataType](data.HvacSystemFunctionSetpointRelationData, filter)

	return result, nil
}

// GetHvacOverrunDescriptionsForFilter returns the overrun descriptions for a given filter
func (h *HvacCommon) GetHvacOverrunDescriptionsForFilter(
	filter model.HvacOverrunDescriptionDataType,
) ([]model.HvacOverrunDescriptionDataType, error) {
	function := model.FunctionTypeHvacOverrunDescriptionListData

	data, err := featureDataCopyOfType[model.HvacOverrunDescriptionListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacOverrunDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacOverrunDescriptionDataType](data.HvacOverrunDescriptionData, filter)

	return result, nil
}

// GetHvacOverrunDataForFilter returns the overrun data for a given filter
func (h *HvacCommon) GetHvacOverrunDataForFilter(
	filter model.HvacOverrunDataType,
) ([]model.HvacOverrunDataType, error) {
	function := model.FunctionTypeHvacOverrunListData

	data, err := featureDataCopyOfType[model.HvacOverrunListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacOverrunData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacOverrunDataType](data.HvacOverrunData, filter)

	return result, nil
}

// GetHvacOverrunForId returns the overrun data for a given overrun ID
func (h *HvacCommon) GetHvacOverrunForId(
	id model.HvacOverrunIdType,
) (*model.HvacOverrunDataType, error) {
	filter := model.HvacOverrunDataType{
		OverrunId: &id,
	}

	result, err := h.GetHvacOverrunDataForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// CheckEventPayloadDataForFilter checks if the given event payload data contains
// system function data matching the given filter
func (h *HvacCommon) CheckEventPayloadDataForFilter(payloadData any, filter model.HvacSystemFunctionDataType) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.HvacSystemFunctionListDataType)
	if !ok || data.HvacSystemFunctionData == nil {
		return false
	}

	result := searchFilterInList[model.HvacSystemFunctionDataType](data.HvacSystemFunctionData, filter)

	return len(result) > 0
}
