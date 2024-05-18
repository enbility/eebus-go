package features

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type LoadControl struct {
	*Feature
}

// Get a new LoadControl features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewLoadControl(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*LoadControl, error) {
	feature, err := NewFeature(model.FeatureTypeTypeLoadControl, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	lc := &LoadControl{
		Feature: feature,
	}

	return lc, nil
}

// request FunctionTypeLoadControlLimitDescriptionListData from a remote device
func (l *LoadControl) RequestLimitDescriptions() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitDescriptionListData, nil, nil)
}

// request FunctionTypeLoadControlLimitConstraintsListData from a remote device
func (l *LoadControl) RequestLimitConstraints() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitConstraintsListData, nil, nil)
}

// request FunctionTypeLoadControlLimitListData from a remote device
func (l *LoadControl) RequestLimitValues() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitListData, nil, nil)
}

// returns the load control limit descriptions
// returns an error if no description data is available yet
func (l *LoadControl) GetLimitDescriptions() ([]model.LoadControlLimitDescriptionDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.LoadControlLimitDescriptionListDataType](l.featureRemote, model.FunctionTypeLoadControlLimitDescriptionListData)
	if err != nil {
		return nil, api.ErrMetadataNotAvailable
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
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// returns the load control limit descriptions of a provided type, direction and scope
// returns an error if no description data for the category is available
//
// providing an empty string for any of the params, will ignore the value in the request
func (l *LoadControl) GetLimitDescriptionsForTypeCategoryDirectionScope(
	limitType model.LoadControlLimitTypeType,
	limitCategory model.LoadControlCategoryType,
	limitDirection model.EnergyDirectionType,
	scope model.ScopeTypeType,
) ([]model.LoadControlLimitDescriptionDataType, error) {
	data, err := l.GetLimitDescriptions()
	if err != nil || len(data) == 0 {
		return nil, err
	}

	var result []model.LoadControlLimitDescriptionDataType

	for _, item := range data {
		if item.LimitId != nil &&
			(limitType == "" || (item.LimitType != nil && *item.LimitType == limitType)) &&
			(limitCategory == "" || (item.LimitCategory != nil && *item.LimitCategory == limitCategory)) &&
			(limitDirection == "" || (item.LimitDirection != nil && *item.LimitDirection == limitDirection)) &&
			(scope == "" || (item.ScopeType != nil && *item.ScopeType == scope)) {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, api.ErrDataNotAvailable
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
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// write load control limits
// returns an error if this failed
func (l *LoadControl) WriteLimitValues(data []model.LoadControlLimitDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionTypeLoadControlLimitListData),
		Filter:   []model.FilterType{*model.NewFilterTypePartial()},
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: data,
		},
	}

	return l.remoteDevice.Sender().Write(l.featureLocal.Address(), l.featureRemote.Address(), cmd)
}

// return limit data
func (l *LoadControl) GetLimitValues() ([]model.LoadControlLimitDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.LoadControlLimitListDataType](l.featureRemote, model.FunctionTypeLoadControlLimitListData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
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

	return nil, api.ErrDataNotAvailable
}
