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
func (e *ElectricalConnection) RequestDescriptions() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionDescriptionListData, nil, nil)
}

// request FunctionTypeElectricalConnectionParameterDescriptionListData from a remote entity
func (e *ElectricalConnection) RequestParameterDescriptions() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, nil, nil)
}

// request FunctionTypeElectricalConnectionPermittedValueSetListData from a remote entity
func (e *ElectricalConnection) RequestPermittedValueSets() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, nil, nil)
}

// request FunctionTypeElectricalConnectionCharacteristicListData from a remote entity
func (e *ElectricalConnection) RequestCharacteristics() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionCharacteristicListData, nil, nil)
}
