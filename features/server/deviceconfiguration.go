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

type DeviceConfiguration struct {
	*Feature

	*internal.DeviceConfigurationCommon
}

func NewDeviceConfiguration(localEntity spineapi.EntityLocalInterface) (*DeviceConfiguration, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceConfiguration, localEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceConfiguration{
		Feature:                   feature,
		DeviceConfigurationCommon: internal.NewLocalDeviceConfiguration(feature.featureLocal),
	}

	return dc, nil
}

var _ api.DeviceConfigurationServerInterface = (*DeviceConfiguration)(nil)

// Add a new description data set and return the keyId
//
// will return nil if the data set could not be added
func (d *DeviceConfiguration) AddKeyValueDescription(
	description model.DeviceConfigurationKeyValueDescriptionDataType,
) *model.DeviceConfigurationKeyIdType {
	data, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueDescriptionListDataType](d.featureLocal, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if err != nil {
		data = &model.DeviceConfigurationKeyValueDescriptionListDataType{}
	}

	var keyId *model.DeviceConfigurationKeyIdType

	maxId := model.DeviceConfigurationKeyIdType(0)

	for _, item := range data.DeviceConfigurationKeyValueDescriptionData {
		if item.KeyId != nil && *item.KeyId >= maxId {
			maxId = *item.KeyId + 1
		}
	}

	keyId = util.Ptr(maxId)
	description.KeyId = keyId

	partial := model.NewFilterTypePartial()
	datalist := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{description},
	}

	if err := d.featureLocal.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, datalist, partial, nil); err != nil {
		return nil
	}

	return keyId
}

// Set or update data set for a keyId
// Elements provided in deleteElements will be removed from the data set before the update
//
// Will return an error if the data set could not be updated
func (d *DeviceConfiguration) UpdateKeyValueDataForKeyId(
	data model.DeviceConfigurationKeyValueDataType,
	deleteElements *model.DeviceConfigurationKeyValueDataElementsType,
	keyId model.DeviceConfigurationKeyIdType,
) (resultErr error) {
	return d.UpdateKeyValueDataForFilter(data, deleteElements, model.DeviceConfigurationKeyValueDescriptionDataType{KeyId: &keyId})
}

// Set or update data set for a filter
// Elements provided in deleteElements will be removed from the data set before the update
//
// Will return an error if the data set could not be updated
func (d *DeviceConfiguration) UpdateKeyValueDataForFilter(
	data model.DeviceConfigurationKeyValueDataType,
	deleteElements *model.DeviceConfigurationKeyValueDataElementsType,
	filter model.DeviceConfigurationKeyValueDescriptionDataType,
) (resultErr error) {
	resultErr = api.ErrDataNotAvailable

	descriptions, err := d.GetKeyValueDescriptionsForFilter(filter)
	if err != nil {
		return err
	}
	if descriptions == nil || len(descriptions) != 1 {
		return
	}

	description := descriptions[0]
	data.KeyId = description.KeyId

	datalist := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{data},
	}

	partial := model.NewFilterTypePartial()
	var deleteFilter *model.FilterType
	if deleteElements != nil {
		deleteFilter = &model.FilterType{
			DeviceConfigurationKeyValueListDataSelectors: &model.DeviceConfigurationKeyValueListDataSelectorsType{
				KeyId: description.KeyId,
			},
			DeviceConfigurationKeyValueDataElements: deleteElements,
		}
	}

	if err := d.featureLocal.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, datalist, partial, deleteFilter); err != nil {
		return errors.New(err.String())
	}

	return nil
}
