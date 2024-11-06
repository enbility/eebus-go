package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
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

// GetHvacOperationModeDescriptions returns the operation mode descriptions
func (h *HvacCommon) GetHvacOperationModeDescriptions() ([]model.HvacOperationModeDescriptionDataType, error) {
	function := model.FunctionTypeHvacOperationModeDescriptionListData
	operationModeDescriptions := make([]model.HvacOperationModeDescriptionDataType, 0)

	data, err := featureDataCopyOfType[model.HvacOperationModeDescriptionListDataType](h.featureLocal, h.featureRemote, function)
	if err == nil || data != nil {
		operationModeDescriptions = append(operationModeDescriptions, data.HvacOperationModeDescriptionData...)
	}

	return operationModeDescriptions, nil
}

// GetHvacSystemFunctionSetpointRelations returns the operation mode relations for a given system function id
func (h *HvacCommon) GetHvacSystemFunctionSetpointRelationsForSystemFunctionId(
	id model.HvacSystemFunctionIdType,
) ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionSetPointRelationListData
	filter := model.HvacSystemFunctionSetpointRelationDataType{
		SystemFunctionId: &id,
	}

	data, err := featureDataCopyOfType[model.HvacSystemFunctionSetpointRelationListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacSystemFunctionSetpointRelationData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacSystemFunctionSetpointRelationDataType](data.HvacSystemFunctionSetpointRelationData, filter)
	if len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// GetHvacSystemFunctionDescriptions returns the system function descriptions for a given filter
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
