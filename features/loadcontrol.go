package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type LoadControl struct {
	*FeatureImpl
}

func NewLoadControl(localRole, remoteRole model.RoleType, localEntity *spine.EntityLocalImpl, remoteEntity *spine.EntityRemoteImpl) (*LoadControl, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeLoadControl, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	lc := &LoadControl{
		FeatureImpl: feature,
	}

	return lc, nil
}

// request FunctionTypeLoadControlLimitDescriptionListData from a remote device
func (l *LoadControl) RequestLimitDescriptions() error {
	_, err := l.requestData(model.FunctionTypeLoadControlLimitDescriptionListData, nil, nil)
	return err
}

// request FunctionTypeLoadControlLimitConstraintsListData from a remote device
func (l *LoadControl) RequestLimitConstraints() error {
	_, err := l.requestData(model.FunctionTypeLoadControlLimitConstraintsListData, nil, nil)
	return err
}

// request FunctionTypeLoadControlLimitListData from a remote device
func (l *LoadControl) RequestLimitValues() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitListData, nil, nil)
}

// returns the load control limit descriptions
// returns an error if no description data is available yet
func (l *LoadControl) GetLimitDescriptions() ([]model.LoadControlLimitDescriptionDataType, error) {
	rData := l.featureRemote.DataCopy(model.FunctionTypeLoadControlLimitDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}

	data := rData.(*model.LoadControlLimitDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.LoadControlLimitDescriptionData, nil
}

// returns the load control limit descriptions of a provided category
// returns an error if no description data for the category is available
func (l *LoadControl) GetLimitDescriptionsForCategory(category model.LoadControlCategoryType) ([]model.LoadControlLimitDescriptionDataType, error) {
	data, err := l.GetLimitDescriptions()
	if err != nil {
		return nil, err
	}

	var result []model.LoadControlLimitDescriptionDataType

	for _, item := range data {
		if item.LimitId != nil && item.LimitCategory != nil && *item.LimitCategory == category {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, ErrDataNotAvailable
	}

	return result, nil
}

// returns the load control limit descriptions for a provided measurementId
// returns an error if no description data for the measurementId is available
func (l *LoadControl) GetLimitDescriptionsForMeasurementId(measurementId model.MeasurementIdType) ([]model.LoadControlLimitDescriptionDataType, error) {
	data, err := l.GetLimitDescriptions()
	if err != nil {
		return nil, err
	}

	var result []model.LoadControlLimitDescriptionDataType

	for _, item := range data {
		if item.LimitId != nil && item.MeasurementId != nil && *item.MeasurementId == measurementId {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, ErrDataNotAvailable
	}

	return result, nil
}

// write load control limits
// returns an error if this failed
func (l *LoadControl) WriteLimitValues(data []model.LoadControlLimitDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, ErrMissingData
	}

	cmd := model.CmdType{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: data,
		},
	}

	return l.featureRemote.Sender().Write(l.featureLocal.Address(), l.featureRemote.Address(), cmd)
}

// return limit data
func (l *LoadControl) GetLimitValues() ([]model.LoadControlLimitDataType, error) {
	rData := l.featureRemote.DataCopy(model.FunctionTypeLoadControlLimitListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.LoadControlLimitListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.LoadControlLimitData, nil
}

// return limit values
func (l *LoadControl) GetLimitValueForLimitId(limitId model.LoadControlLimitIdType) (*model.LoadControlLimitDataType, error) {
	data, err := l.GetLimitValues()
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.LimitId != nil && *item.LimitId == limitId {
			return &item, nil
		}
	}

	return nil, ErrDataNotAvailable
}
