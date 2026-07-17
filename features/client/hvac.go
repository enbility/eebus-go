package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Hvac struct {
	*Feature

	*internal.HvacCommon
}

// Get a new HVAC features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewHvac(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
) (*Hvac, error) {
	feature, err := NewFeature(model.FeatureTypeTypeHvac, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	hvac := &Hvac{
		Feature:    feature,
		HvacCommon: internal.NewRemoteHvac(feature.featureRemote),
	}

	return hvac, nil
}

// request FunctionTypeHvacSystemFunctionDescriptionListData from a remote device
func (h *Hvac) RequestHvacSystemFunctionDescriptions(
	selector *model.HvacSystemFunctionDescriptionListDataSelectorsType,
	elements *model.HvacSystemFunctionDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacSystemFunctionDescriptionListData, selector, elements)
}

// request FunctionTypeHvacSystemFunctionListData from a remote device
func (h *Hvac) RequestHvacSystemFunctions(
	selector *model.HvacSystemFunctionListDataSelectorsType,
	elements *model.HvacSystemFunctionDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacSystemFunctionListData, selector, elements)
}

// request FunctionTypeHvacOperationModeDescriptionListData from a remote device
func (h *Hvac) RequestHvacOperationModeDescriptions(
	selector *model.HvacOperationModeDescriptionListDataSelectorsType,
	elements *model.HvacOperationModeDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacOperationModeDescriptionListData, selector, elements)
}

// request FunctionTypeHvacSystemFunctionOperationModeRelationListData from a remote device
func (h *Hvac) RequestHvacSystemFunctionOperationModeRelations(
	selector *model.HvacSystemFunctionOperationModeRelationListDataSelectorsType,
	elements *model.HvacSystemFunctionOperationModeRelationDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, selector, elements)
}

// request FunctionTypeHvacSystemFunctionSetPointRelationListData from a remote device
func (h *Hvac) RequestHvacSystemFunctionSetpointRelations(
	selector *model.HvacSystemFunctionSetpointRelationListDataSelectorsType,
	elements *model.HvacSystemFunctionSetpointRelationDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacSystemFunctionSetPointRelationListData, selector, elements)
}

// request FunctionTypeHvacOverrunDescriptionListData from a remote device
func (h *Hvac) RequestHvacOverrunDescriptions(
	selector *model.HvacOverrunDescriptionListDataSelectorsType,
	elements *model.HvacOverrunDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacOverrunDescriptionListData, selector, elements)
}

// request FunctionTypeHvacOverrunListData from a remote device
func (h *Hvac) RequestHvacOverruns(
	selector *model.HvacOverrunListDataSelectorsType,
	elements *model.HvacOverrunDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacOverrunListData, selector, elements)
}

// write the given HVAC overrun data to the remote device,
// e.g. to start or stop an overrun
func (h *Hvac) WriteHvacOverrunListData(
	data []model.HvacOverrunDataType,
) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		HvacOverrunListData: &model.HvacOverrunListDataType{
			HvacOverrunData: data,
		},
	}

	return h.remoteDevice.Sender().Write(h.featureLocal.Address(), h.featureRemote.Address(), cmd)
}

// write the given HVAC system function data to the remote device,
// e.g. to change the current operation mode of a system function
func (h *Hvac) WriteHvacSystemFunctionListData(
	data []model.HvacSystemFunctionDataType,
) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		HvacSystemFunctionListData: &model.HvacSystemFunctionListDataType{
			HvacSystemFunctionData: data,
		},
	}

	return h.remoteDevice.Sender().Write(h.featureLocal.Address(), h.featureRemote.Address(), cmd)
}
