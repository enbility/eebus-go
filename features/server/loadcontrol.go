package server

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type LoadControl struct {
	*Feature

	*internal.LoadControlCommon
}

func NewLoadControl(localEntity spineapi.EntityLocalInterface) (*LoadControl, error) {
	feature, err := NewFeature(model.FeatureTypeTypeLoadControl, localEntity)
	if err != nil {
		return nil, err
	}

	lc := &LoadControl{
		Feature:           feature,
		LoadControlCommon: internal.NewLocalLoadControl(feature.featureLocal),
	}

	return lc, nil
}

// Add a new description data set and return the limitId
//
// NOTE: the limitId may not be provided
//
// will return nil if the data set could not be added
func (l *LoadControl) AddLimitDescription(
	description model.LoadControlLimitDescriptionDataType,
) *model.LoadControlLimitIdType {
	if description.LimitId != nil {
		return nil
	}

	data, err := l.GetLimitDescriptionsForFilter(model.LoadControlLimitDescriptionDataType{})
	if err != nil {
		data = []model.LoadControlLimitDescriptionDataType{}
	}

	maxId := model.LoadControlLimitIdType(0)

	for _, item := range data {
		if item.LimitId != nil && *item.LimitId >= maxId {
			maxId = *item.LimitId + 1
		}
	}

	limitId := util.Ptr(maxId)
	description.LimitId = limitId

	partial := model.NewFilterTypePartial()
	datalist := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{description},
	}

	if err := l.featureLocal.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, datalist, partial, nil); err != nil {
		return nil
	}

	return limitId
}

// Set or update data set for a limitId
//
// Will return an error if the data set could not be updated
func (l *LoadControl) UpdateLimitDataForIds(
	data []api.LoadControlLimitDataForID,
) (resultErr error) {
	var filterData []api.LoadControlLimitDataForFilter
	for index, item := range data {
		filterData = append(filterData, api.LoadControlLimitDataForFilter{
			Data:   item.Data,
			Filter: model.LoadControlLimitDescriptionDataType{LimitId: &data[index].Id},
		})
	}

	return l.UpdateLimitDataForFilters(filterData, nil, nil)
}

// Set or update data set for a filter
// deleteSelector will trigger removal of matching items from the data set before the update
// deleteElement will limit the fields to be removed using Id
//
// Will return an error if the data set could not be updated
func (l *LoadControl) UpdateLimitDataForFilters(
	data []api.LoadControlLimitDataForFilter,
	deleteSelector *model.LoadControlLimitListDataSelectorsType,
	deleteElements *model.LoadControlLimitDataElementsType,
) (resultErr error) {
	resultErr = api.ErrDataNotAvailable

	var limitData []model.LoadControlLimitDataType

	for _, item := range data {
		descriptions, err := l.GetLimitDescriptionsForFilter(item.Filter)
		if err != nil || descriptions == nil || len(descriptions) != 1 {
			return
		}

		description := descriptions[0]
		item.Data.LimitId = description.LimitId

		limitData = append(limitData, item.Data)
	}

	partial := model.NewFilterTypePartial()

	datalist := &model.LoadControlLimitListDataType{
		LoadControlLimitData: limitData,
	}

	var deleteFilter *model.FilterType
	if deleteSelector != nil {
		deleteFilter = &model.FilterType{
			LoadControlLimitListDataSelectors: deleteSelector,
		}

		if deleteElements != nil {
			deleteFilter.LoadControlLimitDataElements = deleteElements
		}
	}

	if err := l.featureLocal.UpdateData(model.FunctionTypeLoadControlLimitListData, datalist, partial, deleteFilter); err != nil {
		return errors.New(err.String())
	}

	return nil
}
