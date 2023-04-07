package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type TimeSeries struct {
	*FeatureImpl
}

func NewTimeSeries(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*TimeSeries, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeTimeSeries, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	t := &TimeSeries{
		FeatureImpl: feature,
	}

	return t, nil
}

// request FunctionTypeTimeSeriesDescriptionListData from a remote entity
func (t *TimeSeries) RequestDescriptions() error {
	_, err := t.requestData(model.FunctionTypeTimeSeriesDescriptionListData, nil, nil)
	return err
}

// request FunctionTypeTimeSeriesConstraintsListData from a remote entity
func (t *TimeSeries) RequestConstraints() error {
	_, err := t.requestData(model.FunctionTypeTimeSeriesConstraintsListData, nil, nil)
	return err
}

// request FunctionTypeTimeSeriesListData from a remote device
func (t *TimeSeries) RequestValues() (*model.MsgCounterType, error) {
	return t.requestData(model.FunctionTypeTimeSeriesListData, nil, nil)
}

// write Time Series values
// returns an error if this failed
func (t *TimeSeries) WriteValues(data []model.TimeSeriesDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, ErrMissingData
	}

	cmd := model.CmdType{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: data,
		},
	}

	return t.featureRemote.Sender().Write(t.featureLocal.Address(), t.featureRemote.Address(), cmd)
}

// return current values for Time Series
func (t *TimeSeries) GetValues() ([]model.TimeSeriesDataType, error) {
	rData := t.featureRemote.Data(model.FunctionTypeTimeSeriesListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.TimeSeriesListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
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

	return nil, ErrDataNotAvailable
}

// return list of descriptions
func (t *TimeSeries) GetDescriptions() ([]model.TimeSeriesDescriptionDataType, error) {
	rData := t.featureRemote.Data(model.FunctionTypeTimeSeriesDescriptionListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.TimeSeriesDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
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

	return nil, ErrDataNotAvailable
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

	return nil, ErrDataNotAvailable
}

// return current constraints for Time Series
func (t *TimeSeries) GetConstraints() ([]model.TimeSeriesConstraintsDataType, error) {
	rData := t.featureRemote.Data(model.FunctionTypeTimeSeriesConstraintsListData)
	switch constraintsData := rData.(type) {
	case *model.TimeSeriesConstraintsListDataType:
		if constraintsData == nil {
			return nil, ErrDataNotAvailable
		}
	}

	data := rData.(*model.TimeSeriesConstraintsListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.TimeSeriesConstraintsData, nil
}
