package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type TimeSeries struct {
	*Feature

	*internal.TimeSeriesCommon
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
		Feature:          feature,
		TimeSeriesCommon: internal.NewRemoteTimeSeries(feature.featureRemote),
	}

	return t, nil
}

var _ api.TimeSeriesClientInterface = (*TimeSeries)(nil)

// request FunctionTypeTimeSeriesDescriptionListData from a remote entity
func (t *TimeSeries) RequestDescriptions(
	selector *model.TimeSeriesDescriptionListDataSelectorsType,
	elements *model.TimeSeriesDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesDescriptionListData, selector, elements)
}

// request FunctionTypeTimeSeriesConstraintsListData from a remote entity
func (t *TimeSeries) RequestConstraints(
	selector *model.TimeSeriesConstraintsListDataSelectorsType,
	elements *model.TimeSeriesConstraintsDataElementsType,
) (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesConstraintsListData, selector, elements)
}

// request FunctionTypeTimeSeriesListData from a remote device
func (t *TimeSeries) RequestData(
	selector *model.TimeSeriesListDataSelectorsType,
	elements *model.TimeSeriesDataElementsType,
) (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesListData, selector, elements)
}

// write Time Series values
// returns an error if this failed
func (t *TimeSeries) WriteData(data []model.TimeSeriesDataType) (*model.MsgCounterType, error) {
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
