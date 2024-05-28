package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type TimeSeriesCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalTimeSeries(featureLocal spineapi.FeatureLocalInterface) *TimeSeriesCommon {
	return &TimeSeriesCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteTimeSeries(featureRemote spineapi.FeatureRemoteInterface) *TimeSeriesCommon {
	return &TimeSeriesCommon{
		featureRemote: featureRemote,
	}
}

var _ api.TimeSeriesCommonInterface = (*TimeSeriesCommon)(nil)

// return list of descriptions for a given filter
func (t *TimeSeriesCommon) GetDescriptionsForFilter(
	filter model.TimeSeriesDescriptionDataType,
) ([]model.TimeSeriesDescriptionDataType, error) {
	function := model.FunctionTypeTimeSeriesDescriptionListData

	data, err := featureDataCopyOfType[model.TimeSeriesDescriptionListDataType](t.featureLocal, t.featureRemote, function)
	if err != nil || data == nil || data.TimeSeriesDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.TimeSeriesDescriptionDataType](data.TimeSeriesDescriptionData, filter)
	return result, nil
}

// return current constraints for Time Series
func (t *TimeSeriesCommon) GetConstraints() ([]model.TimeSeriesConstraintsDataType, error) {
	function := model.FunctionTypeTimeSeriesConstraintsListData

	data, err := featureDataCopyOfType[model.TimeSeriesConstraintsListDataType](t.featureLocal, t.featureRemote, function)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.TimeSeriesConstraintsData, nil
}

// return current data for Time Series for a given filter
func (t *TimeSeriesCommon) GetDataForFilter(filter model.TimeSeriesDescriptionDataType) ([]model.TimeSeriesDataType, error) {
	function := model.FunctionTypeTimeSeriesListData

	descriptions, err := t.GetDescriptionsForFilter(filter)
	if err != nil || len(descriptions) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	data, err := featureDataCopyOfType[model.TimeSeriesListDataType](t.featureLocal, t.featureRemote, function)
	if err != nil || data == nil || data.TimeSeriesData == nil {
		return nil, api.ErrDataNotAvailable
	}

	var result []model.TimeSeriesDataType

	for _, desc := range descriptions {
		filter2 := model.TimeSeriesDataType{
			TimeSeriesId: desc.TimeSeriesId,
		}

		elements := searchFilterInList[model.TimeSeriesDataType](data.TimeSeriesData, filter2)
		result = append(result, elements...)
	}
	return result, nil
}
