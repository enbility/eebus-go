package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceConfigurationCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalDeviceConfiguration(featureLocal spineapi.FeatureLocalInterface) *DeviceConfigurationCommon {
	return &DeviceConfigurationCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteDeviceConfiguration(featureRemote spineapi.FeatureRemoteInterface) *DeviceConfigurationCommon {
	return &DeviceConfigurationCommon{
		featureRemote: featureRemote,
	}
}

var _ api.DeviceConfigurationCommonInterface = (*DeviceConfigurationCommon)(nil)

// check if spine.EventPayload Data contains data for a given filter
//
// data type will be checked for model.DeviceConfigurationKeyValueListDataType,
// filter type will be checked for model.DeviceConfigurationKeyValueDescriptionDataType
func (d *DeviceConfigurationCommon) CheckEventPayloadDataForFilter(payloadData any, filter any) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.DeviceConfigurationKeyValueListDataType)
	filterData, ok2 := filter.(model.DeviceConfigurationKeyValueDescriptionDataType)
	if !ok || !ok2 {
		return false
	}

	descs, err := d.GetKeyValueDescriptionsForFilter(filterData)
	if err != nil {
		return false
	}
	for _, desc := range descs {
		if desc.KeyId == nil {
			continue
		}

		for _, item := range data.DeviceConfigurationKeyValueData {
			if item.KeyId != nil &&
				*item.KeyId == *desc.KeyId ||
				item.Value != nil {
				return true
			}
		}
	}

	return false
}

// Get the description for a given keyId
//
// Will return nil if no matching description is found
func (d *DeviceConfigurationCommon) GetKeyValueDescriptionFoKeyId(keyId model.DeviceConfigurationKeyIdType) (
	*model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	data, err := d.GetKeyValueDescriptionsForFilter(model.DeviceConfigurationKeyValueDescriptionDataType{KeyId: &keyId})

	if err != nil || len(data) != 1 {
		return nil, api.ErrDataNotAvailable
	}

	return &data[0], nil
}

// Get the description for a given value combination
//
// Returns an error if no matching description is found
func (d *DeviceConfigurationCommon) GetKeyValueDescriptionsForFilter(
	filter model.DeviceConfigurationKeyValueDescriptionDataType,
) ([]model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	function := model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData

	data, err := featureDataCopyOfType[model.DeviceConfigurationKeyValueDescriptionListDataType](d.featureLocal, d.featureRemote, function)
	if err != nil || data == nil || data.DeviceConfigurationKeyValueDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.DeviceConfigurationKeyValueDescriptionDataType](
		data.DeviceConfigurationKeyValueDescriptionData, filter)
	return result, nil
}

// Get the key value data for a given keyId
//
// Will return nil if no data is available
func (d *DeviceConfigurationCommon) GetKeyValueDataForKeyId(keyId model.DeviceConfigurationKeyIdType) (
	*model.DeviceConfigurationKeyValueDataType, error) {
	return d.GetKeyValueDataForFilter(model.DeviceConfigurationKeyValueDescriptionDataType{KeyId: &keyId})
}

// Get key value data for a given filter
//
// Will return nil if no data is available
func (d *DeviceConfigurationCommon) GetKeyValueDataForFilter(filter model.DeviceConfigurationKeyValueDescriptionDataType) (
	*model.DeviceConfigurationKeyValueDataType, error) {
	function := model.FunctionTypeDeviceConfigurationKeyValueListData

	descriptions, err := d.GetKeyValueDescriptionsForFilter(filter)
	if err != nil || descriptions == nil || len(descriptions) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	description := descriptions[0]

	data, err := featureDataCopyOfType[model.DeviceConfigurationKeyValueListDataType](d.featureLocal, d.featureRemote, function)
	if err != nil || data == nil || data.DeviceConfigurationKeyValueData == nil {
		return nil, api.ErrDataNotAvailable
	}

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId != nil && *item.KeyId == *description.KeyId {
			return &item, nil
		}
	}

	return nil, api.ErrDataNotAvailable
}
