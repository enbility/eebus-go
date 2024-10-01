package server

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type Measurement struct {
	*Feature

	*internal.MeasurementCommon
}

func NewMeasurement(localEntity spineapi.EntityLocalInterface) (*Measurement, error) {
	feature, err := NewFeature(model.FeatureTypeTypeMeasurement, localEntity)
	if err != nil {
		return nil, err
	}

	m := &Measurement{
		Feature:           feature,
		MeasurementCommon: internal.NewLocalMeasurement(feature.featureLocal),
	}

	return m, nil
}

var _ api.MeasurementServerInterface = (*Measurement)(nil)

// Add a new parameter description data sett and return the measurementId
//
// NOTE: the measurementId may not be provided
//
// will return nil if the data set could not be added
func (m *Measurement) AddDescription(
	description model.MeasurementDescriptionDataType,
) *model.MeasurementIdType {
	if description.MeasurementId != nil {
		return nil
	}

	data, err := m.GetDescriptionsForFilter(model.MeasurementDescriptionDataType{})
	if err != nil {
		data = []model.MeasurementDescriptionDataType{}
	}

	maxId := model.MeasurementIdType(0)

	for _, item := range data {
		if item.MeasurementId != nil && *item.MeasurementId >= maxId {
			maxId = *item.MeasurementId + 1
		}
	}

	measurementId := util.Ptr(maxId)
	description.MeasurementId = measurementId

	partial := model.NewFilterTypePartial()
	datalist := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{description},
	}

	if err := m.featureLocal.UpdateData(model.FunctionTypeMeasurementDescriptionListData, datalist, partial, nil); err != nil {
		return nil
	}

	return measurementId
}

// Set or update data set for a measurementId
// Id provided in deleteId will trigger removal of matching items from the data set before the update
// Elements provided in deleteElement will limit the fields to be removed using Id
//
// Will return an error if the data set could not be updated
func (m *Measurement) UpdateDataForIds(
	data []api.MeasurementDataForID,
	deleteId *model.MeasurementIdType,
	deleteElements *model.MeasurementDataElementsType,
) (resultErr error) {
	var filterData []api.MeasurementDataForFilter
	for _, item := range data {
		filterData = append(filterData, api.MeasurementDataForFilter{
			Data:   item.Data,
			Filter: model.MeasurementDescriptionDataType{MeasurementId: &item.Id},
		})
	}

	var deleteSelector *model.MeasurementListDataSelectorsType
	if deleteId != nil {
		deleteSelector = &model.MeasurementListDataSelectorsType{
			MeasurementId: deleteId,
		}
	}

	return m.UpdateDataForFilters(filterData, deleteSelector, deleteElements)
}

// Set or update data set for a filter
// Selector provided in deleteSelector will trigger removal of matching items from the data set before the update
// Elements provided in deleteElement will limit the fields to be removed using Selector
//
// Will return an error if the data set could not be updated
func (m *Measurement) UpdateDataForFilters(
	data []api.MeasurementDataForFilter,
	deleteSelector *model.MeasurementListDataSelectorsType,
	deleteElements *model.MeasurementDataElementsType,
) (resultErr error) {
	resultErr = api.ErrDataNotAvailable

	var measurementData []model.MeasurementDataType

	for _, item := range data {
		descriptions, err := m.GetDescriptionsForFilter(item.Filter)
		if err != nil || descriptions == nil || len(descriptions) != 1 {
			return
		}

		description := descriptions[0]
		item.Data.MeasurementId = description.MeasurementId

		measurementData = append(measurementData, item.Data)
	}

	partial := model.NewFilterTypePartial()

	datalist := &model.MeasurementListDataType{
		MeasurementData: measurementData,
	}

	var deleteFilter *model.FilterType
	if deleteSelector != nil {
		deleteFilter = &model.FilterType{
			MeasurementListDataSelectors: deleteSelector,
		}

		if deleteElements != nil {
			deleteFilter.MeasurementDataElements = deleteElements
		}
	}

	if err := m.featureLocal.UpdateData(model.FunctionTypeMeasurementListData, datalist, partial, deleteFilter); err != nil {
		return errors.New(err.String())
	}

	return nil
}
