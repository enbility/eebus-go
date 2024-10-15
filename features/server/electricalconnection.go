package server

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
)

type ElectricalConnection struct {
	*Feature

	*internal.ElectricalConnectionCommon
}

func NewElectricalConnection(localEntity spineapi.EntityLocalInterface) (*ElectricalConnection, error) {
	feature, err := NewFeature(model.FeatureTypeTypeElectricalConnection, localEntity)
	if err != nil {
		return nil, err
	}

	ec := &ElectricalConnection{
		Feature:                    feature,
		ElectricalConnectionCommon: internal.NewLocalElectricalConnection(feature.featureLocal),
	}

	return ec, nil
}

// Add a new description data set
//
// NOTE: the electricalConnectionId has to be provided
//
// will return nil if the data set could not be added
func (e *ElectricalConnection) AddDescription(
	description model.ElectricalConnectionDescriptionDataType,
) error {
	if description.ElectricalConnectionId == nil {
		return errors.New("missing id data")
	}

	partial := model.NewFilterTypePartial()
	datalist := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{description},
	}

	if err := e.featureLocal.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, datalist, partial, nil); err != nil {
		return errors.New(err.String())
	}

	return nil
}

// Add a new parameter description data sett and return the parameterId
//
// NOTE: the electricalConnectionId has to be provided, parameterId may not be provided
//
// will return nil if the data set could not be added
func (e *ElectricalConnection) AddParameterDescription(
	description model.ElectricalConnectionParameterDescriptionDataType,
) *model.ElectricalConnectionParameterIdType {
	if description.ElectricalConnectionId == nil || description.ParameterId != nil {
		return nil
	}

	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: description.ElectricalConnectionId,
	}
	data, _ := e.GetParameterDescriptionsForFilter(filter)

	maxId := model.ElectricalConnectionParameterIdType(0)

	for _, item := range data {
		if item.ParameterId != nil && *item.ParameterId >= maxId {
			maxId = *item.ParameterId + 1
		}
	}

	parameterId := util.Ptr(maxId)
	description.ParameterId = parameterId

	partial := model.NewFilterTypePartial()
	datalist := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{description},
	}

	if err := e.featureLocal.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, datalist, partial, nil); err != nil {
		return nil
	}

	return parameterId
}

// Add a new characteristic data set
//
// Note: ElectricalConnectionId and ParameterId must be set, CharacteristicId will be set automatically
//
// Will return an error if the data set could not be added
func (e *ElectricalConnection) AddCharacteristic(data model.ElectricalConnectionCharacteristicDataType) (*model.ElectricalConnectionCharacteristicIdType, error) {
	if data.ElectricalConnectionId == nil ||
		data.ParameterId == nil {
		return nil, errors.New("missing id data")
	}
	if data.CharacteristicId != nil {
		return nil, errors.New("characteristic id must not be set")
	}

	maxId := model.ElectricalConnectionCharacteristicIdType(0)

	listData, err := spine.LocalFeatureDataCopyOfType[*model.ElectricalConnectionCharacteristicListDataType](e.featureLocal, model.FunctionTypeElectricalConnectionCharacteristicListData)
	if err != nil {
		listData = &model.ElectricalConnectionCharacteristicListDataType{}
	}

	for _, item := range listData.ElectricalConnectionCharacteristicData {
		if item.CharacteristicId != nil && *item.CharacteristicId >= maxId {
			maxId = *item.CharacteristicId + 1
		}
	}

	charId := util.Ptr(maxId)
	data.CharacteristicId = charId

	datalist := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{data},
	}

	partial := model.NewFilterTypePartial()
	if err := e.featureLocal.UpdateData(model.FunctionTypeElectricalConnectionCharacteristicListData, datalist, partial, nil); err != nil {
		return nil, errors.New(err.String())
	}

	return charId, nil
}

// Update data set for a filter
// Elements provided in deleteElements will be removed from the data set before the update
//
// // ElectricalConnectionId, ParameterId and CharacteristicId must be set
//
// Will return an error if the data set could not be updated
func (e *ElectricalConnection) UpdateCharacteristic(
	data model.ElectricalConnectionCharacteristicDataType,
	deleteElements *model.ElectricalConnectionCharacteristicDataElementsType,
) error {
	if data.CharacteristicId == nil ||
		data.ElectricalConnectionId == nil ||
		data.ParameterId == nil {
		return errors.New("missing id data")
	}

	filter := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: data.ElectricalConnectionId,
		ParameterId:            data.ParameterId,
		CharacteristicId:       data.CharacteristicId,
	}
	chars, err := e.GetCharacteristicsForFilter(filter)
	if err != nil || chars == nil || len(chars) != 1 {
		return errors.New("no matching element found")
	}

	partial := model.NewFilterTypePartial()
	var deleteFilter *model.FilterType
	if deleteElements != nil {
		deleteFilter = &model.FilterType{
			ElectricalConnectionCharacteristicListDataSelectors: &model.ElectricalConnectionCharacteristicListDataSelectorsType{
				CharacteristicId: data.CharacteristicId,
			},
			ElectricalConnectionCharacteristicDataElements: deleteElements,
		}
	}

	datalist := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{data},
	}

	if err := e.featureLocal.UpdateData(model.FunctionTypeElectricalConnectionCharacteristicListData, datalist, partial, deleteFilter); err != nil {
		return errors.New(err.String())
	}

	return nil
}

// Set or update data set for a electricalConnectiontId
// Id provided in deleteId will trigger removal of matching items from the data set before the update
// Elements provided in deleteElement will limit the fields to be removed using Id
//
// Will return an error if the data set could not be updated
func (e *ElectricalConnection) UpdatePermittedValueSetForIds(
	data []api.ElectricalConnectionPermittedValueSetForID,
) (resultErr error) {
	var filterData []api.ElectricalConnectionPermittedValueSetForFilter
	for index, item := range data {
		filterData = append(filterData, api.ElectricalConnectionPermittedValueSetForFilter{
			Data: item.Data,
			Filter: model.ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId: &data[index].ElectricalConnectionId,
				ParameterId:            &data[index].ParameterId,
			},
		})
	}

	return e.UpdatePermittedValueSetForFilters(filterData, nil, nil)
}

// Set or update data set for a filter
// deleteSelector will trigger removal of matching items from the data set before the update
// deleteElement will limit the fields to be removed using Id
//
// Will return an error if the data set could not be updated
func (e *ElectricalConnection) UpdatePermittedValueSetForFilters(
	data []api.ElectricalConnectionPermittedValueSetForFilter,
	deleteSelector *model.ElectricalConnectionPermittedValueSetListDataSelectorsType,
	deleteElements *model.ElectricalConnectionPermittedValueSetDataElementsType,
) (resultErr error) {
	resultErr = api.ErrDataNotAvailable

	var permittedData []model.ElectricalConnectionPermittedValueSetDataType

	for _, item := range data {
		descriptions, err := e.GetParameterDescriptionsForFilter(item.Filter)
		if err != nil || descriptions == nil || len(descriptions) != 1 {
			return
		}

		description := descriptions[0]
		item.Data.ElectricalConnectionId = description.ElectricalConnectionId
		item.Data.ParameterId = description.ParameterId

		permittedData = append(permittedData, item.Data)
	}

	partial := model.NewFilterTypePartial()

	datalist := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: permittedData,
	}

	var deleteFilter *model.FilterType
	if deleteSelector != nil {
		deleteFilter = &model.FilterType{
			ElectricalConnectionPermittedValueSetListDataSelectors: deleteSelector,
		}

		if deleteElements != nil {
			deleteFilter.ElectricalConnectionPermittedValueSetDataElements = deleteElements
		}
	}

	if err := e.featureLocal.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, datalist, partial, deleteFilter); err != nil {
		return errors.New(err.String())
	}

	return nil
}
