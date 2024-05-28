package server

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	"github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
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
// Elements provided in deleteElements will be removed from the data set before the update
//
// Will return an error if the data set could not be updated
func (m *Measurement) UpdateDataForId(
	data model.MeasurementDataType,
	deleteElements *model.MeasurementDataElementsType,
	measurementId model.MeasurementIdType,
) (resultErr error) {
	return m.UpdateDataForFilter(data, deleteElements, model.MeasurementDescriptionDataType{MeasurementId: &measurementId})
}

// Set or update data set for a filter
// Elements provided in deleteElements will be removed from the data set before the update
//
// Will return an error if the data set could not be updated
func (m *Measurement) UpdateDataForFilter(
	data model.MeasurementDataType,
	deleteElements *model.MeasurementDataElementsType,
	filter model.MeasurementDescriptionDataType,
) (resultErr error) {
	resultErr = api.ErrDataNotAvailable

	descriptions, err := m.GetDescriptionsForFilter(filter)
	if err != nil || descriptions == nil || len(descriptions) != 1 {
		return
	}

	description := descriptions[0]
	data.MeasurementId = description.MeasurementId

	datalist := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{data},
	}

	partial := model.NewFilterTypePartial()
	var delete *model.FilterType
	if deleteElements != nil {
		delete = &model.FilterType{
			MeasurementListDataSelectors: &model.MeasurementListDataSelectorsType{
				MeasurementId: description.MeasurementId,
			},
			MeasurementDataElements: deleteElements,
		}
	}

	if err := m.featureLocal.UpdateData(model.FunctionTypeMeasurementListData, datalist, partial, delete); err != nil {
		return errors.New(err.String())
	}

	return nil
}
