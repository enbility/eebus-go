package features

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type TimeSeries struct {
	*Feature
}

// Get a new TimeSeries features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewTimeSeries(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*TimeSeries, error) {
	feature, err := NewFeature(model.FeatureTypeTypeTimeSeries, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	t := &TimeSeries{
		Feature: feature,
	}

	return t, nil
}

// request FunctionTypeTimeSeriesDescriptionListData from a remote entity
func (t *TimeSeries) RequestDescriptions() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesDescriptionListData, nil, nil)
}

// request FunctionTypeTimeSeriesConstraintsListData from a remote entity
func (t *TimeSeries) RequestConstraints() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesConstraintsListData, nil, nil)
}

// request FunctionTypeTimeSeriesListData from a remote device
func (t *TimeSeries) RequestValues() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesListData, nil, nil)
}

// write Time Series values
// returns an error if this failed
func (t *TimeSeries) WriteValues(data []model.TimeSeriesDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: data,
		},
	}

	return t.remoteDevice.Sender().Write(t.featureLocal.Address(), t.featureRemote.Address(), cmd)
}

// return current values for Time Series
func (t *TimeSeries) GetValues() ([]model.TimeSeriesDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.TimeSeriesListDataType](t.featureRemote, model.FunctionTypeTimeSeriesListData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.TimeSeriesData, nil
}

// return current value for a given TimeSeriesType
// there can only be one item matching the type
func (t *TimeSeries) GetValueForType(timeSeriesType model.TimeSeriesTypeType) (*model.TimeSeriesDataType, error) {
	data, err := t.GetValues()
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.TimeSeriesId == nil {
			continue
		}

		desc, err := t.GetDescriptionForId(*item.TimeSeriesId)
		if err != nil {
			continue
		}

		if desc.TimeSeriesType == nil || *desc.TimeSeriesType != timeSeriesType {
			continue
		}

		return &item, nil
	}

	return nil, api.ErrDataNotAvailable
}

// return list of descriptions
func (t *TimeSeries) GetDescriptions() ([]model.TimeSeriesDescriptionDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.TimeSeriesDescriptionListDataType](t.featureRemote, model.FunctionTypeTimeSeriesDescriptionListData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.TimeSeriesDescriptionData, nil
}

func (t *TimeSeries) GetDescriptionForId(id model.TimeSeriesIdType) (*model.TimeSeriesDescriptionDataType, error) {
	data, err := t.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.TimeSeriesId != nil && *item.TimeSeriesId == id {
			return &item, nil
		}
	}

	return nil, api.ErrDataNotAvailable
}

func (t *TimeSeries) GetDescriptionForType(timeSeriesType model.TimeSeriesTypeType) (*model.TimeSeriesDescriptionDataType, error) {
	data, err := t.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.TimeSeriesType != nil && *item.TimeSeriesType == timeSeriesType {
			return &item, nil
		}
	}

	return nil, api.ErrDataNotAvailable
}

// return current constraints for Time Series
func (t *TimeSeries) GetConstraints() ([]model.TimeSeriesConstraintsDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.TimeSeriesConstraintsListDataType](t.featureRemote, model.FunctionTypeTimeSeriesConstraintsListData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.TimeSeriesConstraintsData, nil
}
