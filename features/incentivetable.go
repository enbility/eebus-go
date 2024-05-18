package features

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type IncentiveTable struct {
	*Feature
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
		Feature: feature,
	}

	return i, nil
}

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

// return current values for Time Series
func (i *IncentiveTable) GetValues() ([]model.IncentiveTableType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.IncentiveTableDataType](i.featureRemote, model.FunctionTypeIncentiveTableData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.IncentiveTable, nil
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

// return list of descriptions
func (i *IncentiveTable) GetDescriptions() ([]model.IncentiveTableDescriptionType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.IncentiveTableDescriptionDataType](i.featureRemote, model.FunctionTypeIncentiveTableDescriptionData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.IncentiveTableDescription, nil
}

// return list of descriptions
func (i *IncentiveTable) GetDescriptionsForScope(scope model.ScopeTypeType) ([]model.IncentiveTableDescriptionType, error) {
	data, err := i.GetDescriptions()
	if err != nil {
		return nil, err
	}

	var result []model.IncentiveTableDescriptionType
	for _, item := range data {
		if item.TariffDescription != nil && item.TariffDescription.ScopeType != nil && *item.TariffDescription.ScopeType == scope {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// return list of constraints
func (i *IncentiveTable) GetConstraints() ([]model.IncentiveTableConstraintsType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.IncentiveTableConstraintsDataType](i.featureRemote, model.FunctionTypeIncentiveTableConstraintsData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.IncentiveTableConstraints, nil
}
