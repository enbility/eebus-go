package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type IncentiveTable struct {
	*Feature

	*internal.IncentiveTableCommon
}

// Get a new IncentiveTable features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewIncentiveTable(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*IncentiveTable, error) {
	feature, err := NewFeature(model.FeatureTypeTypeIncentiveTable, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	i := &IncentiveTable{
		Feature:              feature,
		IncentiveTableCommon: internal.NewRemoteIncentiveTable(feature.featureRemote),
	}

	return i, nil
}

var _ api.IncentiveTableClientInterface = (*IncentiveTable)(nil)

// request FunctionTypeIncentiveTableDescriptionData from a remote entity
func (i *IncentiveTable) RequestDescriptions() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIncentiveTableDescriptionData, nil, nil)
}

// request FunctionTypeIncentiveTableConstraintsData from a remote entity
func (i *IncentiveTable) RequestConstraints() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIncentiveTableConstraintsData, nil, nil)
}

// request FunctionTypeIncentiveTableData from a remote entity
func (i *IncentiveTable) RequestValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIncentiveTableData, nil, nil)
}

// write incentivetable descriptions
// returns an error if this failed
func (i *IncentiveTable) WriteDescriptions(data []model.IncentiveTableDescriptionType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{
			IncentiveTableDescription: data,
		},
	}

	return i.remoteDevice.Sender().Write(i.featureLocal.Address(), i.featureRemote.Address(), cmd)
}

// write incentivetable descriptions
// returns an error if this failed
func (i *IncentiveTable) WriteValues(data []model.IncentiveTableType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		IncentiveTableData: &model.IncentiveTableDataType{
			IncentiveTable: data,
		},
	}

	return i.remoteDevice.Sender().Write(i.featureLocal.Address(), i.featureRemote.Address(), cmd)
}
