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
func (t *TimeSeries) RequestDescriptions() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesDescriptionListData, nil, nil)
}

// request FunctionTypeTimeSeriesConstraintsListData from a remote entity
func (t *TimeSeries) RequestConstraints() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesConstraintsListData, nil, nil)
}

// request FunctionTypeTimeSeriesListData from a remote device
func (t *TimeSeries) RequestData() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesListData, nil, nil)
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
