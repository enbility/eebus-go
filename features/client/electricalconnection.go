package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type ElectricalConnection struct {
	*Feature

	*internal.ElectricalConnectionCommon
}

// Get a new ElectricalConnection features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewElectricalConnection(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*ElectricalConnection, error) {
	feature, err := NewFeature(model.FeatureTypeTypeElectricalConnection, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	e := &ElectricalConnection{
		Feature:                    feature,
		ElectricalConnectionCommon: internal.NewRemoteElectricalConnection(feature.featureRemote),
	}

	return e, nil
}

var _ api.ElectricalConnectionClientInterface = (*ElectricalConnection)(nil)

// request ElectricalConnectionDescriptionListDataType from a remote entity
func (e *ElectricalConnection) RequestDescriptions(
	selector *model.ElectricalConnectionDescriptionListDataSelectorsType,
	elements *model.ElectricalConnectionDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionDescriptionListData, selector, elements)
}

// request FunctionTypeElectricalConnectionParameterDescriptionListData from a remote entity
func (e *ElectricalConnection) RequestParameterDescriptions(
	selector *model.ElectricalConnectionParameterDescriptionListDataSelectorsType,
	elements *model.ElectricalConnectionParameterDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, selector, elements)
}

// request FunctionTypeElectricalConnectionPermittedValueSetListData from a remote entity
func (e *ElectricalConnection) RequestPermittedValueSets(
	selector *model.ElectricalConnectionPermittedValueSetListDataSelectorsType,
	elements *model.ElectricalConnectionPermittedValueSetDataElementsType,
) (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, selector, elements)
}

// request FunctionTypeElectricalConnectionCharacteristicListData from a remote entity
func (e *ElectricalConnection) RequestCharacteristics(
	selector *model.ElectricalConnectionCharacteristicListDataSelectorsType,
	elements *model.ElectricalConnectionCharacteristicDataElementsType,
) (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionCharacteristicListData, selector, elements)
}
