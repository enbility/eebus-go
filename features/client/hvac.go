package client

import (
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

// request FunctionTypeHvacSystemFunctionSetPointRelationListData from a remote device
func (h *Hvac) RequestHvacSystemFunctionSetPointRelations(
	selector *model.HvacSystemFunctionSetpointRelationListDataSelectorsType,
	elements *model.HvacSystemFunctionSetpointRelationDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacSystemFunctionSetPointRelationListData, selector, elements)
}

// request FunctionTypeHvacOperationModeDescriptionListData from a remote device
func (h *Hvac) RequestHvacOperationModeDescriptions(
	selector *model.HvacOperationModeDescriptionListDataSelectorsType,
	elements *model.HvacOperationModeDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacOperationModeDescriptionListData, selector, elements)
}
