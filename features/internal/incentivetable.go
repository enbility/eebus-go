package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type IncentiveTableCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalIncentiveTable(featureLocal spineapi.FeatureLocalInterface) *IncentiveTableCommon {
	return &IncentiveTableCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteIncentiveTable(featureRemote spineapi.FeatureRemoteInterface) *IncentiveTableCommon {
	return &IncentiveTableCommon{
		featureRemote: featureRemote,
	}
}

var _ api.IncentiveTableCommonInterface = (*IncentiveTableCommon)(nil)

// return list of descriptions for a given filter
func (i *IncentiveTableCommon) GetDescriptionsForFilter(
	filter model.TariffDescriptionDataType,
) ([]model.IncentiveTableDescriptionType, error) {
	function := model.FunctionTypeIncentiveTableDescriptionData
	data, err := featureDataCopyOfType[model.IncentiveTableDescriptionDataType](i.featureLocal, i.featureRemote, function)
	if err != nil || data == nil || data.IncentiveTableDescription == nil {
		return nil, api.ErrDataNotAvailable
	}

	var result []model.IncentiveTableDescriptionType
	for _, item := range data.IncentiveTableDescription {
		match := searchFilterInItem[model.TariffDescriptionDataType](*item.TariffDescription, filter)

		if match {
			result = append(result, item)
		}
	}

	return result, nil
}

// return list of constraints
func (i *IncentiveTableCommon) GetConstraints() ([]model.IncentiveTableConstraintsType, error) {
	function := model.FunctionTypeIncentiveTableConstraintsData
	data, err := featureDataCopyOfType[model.IncentiveTableConstraintsDataType](i.featureLocal, i.featureRemote, function)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.IncentiveTableConstraints, nil
}

// return current data for Time Series
func (i *IncentiveTableCommon) GetData() ([]model.IncentiveTableType, error) {
	function := model.FunctionTypeIncentiveTableData

	data, err := featureDataCopyOfType[model.IncentiveTableDataType](i.featureLocal, i.featureRemote, function)
	if err != nil || data == nil || data.IncentiveTable == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.IncentiveTable, nil
}
