package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type IncentiveTable struct {
	*FeatureImpl
}

func NewIncentiveTable(localRole, remoteRole model.RoleType, localEntity *spine.EntityLocalImpl, remoteEntity *spine.EntityRemoteImpl) (*IncentiveTable, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeIncentiveTable, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	i := &IncentiveTable{
		FeatureImpl: feature,
	}

	return i, nil
}

// request FunctionTypeIncentiveTableDescriptionData from a remote entity
func (i *IncentiveTable) RequestDescriptions() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableDescriptionData, nil, nil)
	return err
}

// request FunctionTypeIncentiveTableConstraintsData from a remote entity
func (i *IncentiveTable) RequestConstraints() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableConstraintsData, nil, nil)
	return err
}

// request FunctionTypeIncentiveTableData from a remote entity
func (i *IncentiveTable) RequestValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIncentiveTableData, nil, nil)
}

// write incentivetable descriptions
// returns an error if this failed
func (i *IncentiveTable) WriteValues(data []model.IncentiveTableType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, ErrMissingData
	}

	cmd := model.CmdType{
		IncentiveTableData: &model.IncentiveTableDataType{
			IncentiveTable: data,
		},
	}

	return i.featureRemote.Sender().Write(i.featureLocal.Address(), i.featureRemote.Address(), cmd)
}

// return current values for Time Series
func (i *IncentiveTable) GetValues() ([]model.IncentiveTableType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeIncentiveTableData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.IncentiveTableDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.IncentiveTable, nil
}

// write incentivetable descriptions
// returns an error if this failed
func (i *IncentiveTable) WriteDescriptions(data []model.IncentiveTableDescriptionType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, ErrMissingData
	}

	cmd := model.CmdType{
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{
			IncentiveTableDescription: data,
		},
	}

	return i.featureRemote.Sender().Write(i.featureLocal.Address(), i.featureRemote.Address(), cmd)
}

// return list of descriptions
func (i *IncentiveTable) GetDescriptions() ([]model.IncentiveTableDescriptionType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeIncentiveTableDescriptionData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.IncentiveTableDescriptionDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
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
		return nil, ErrDataNotAvailable
	}

	return result, nil
}

// return list of constraints
func (i *IncentiveTable) GetConstraints() ([]model.IncentiveTableConstraintsType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeIncentiveTableConstraintsData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.IncentiveTableConstraintsDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.IncentiveTableConstraints, nil
}
