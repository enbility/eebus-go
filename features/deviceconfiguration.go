package features

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type DeviceConfiguration struct {
	*Feature
}

// Get a new DeviceConfiguration features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewDeviceConfiguration(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*DeviceConfiguration, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceConfiguration, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceConfiguration{
		Feature: feature,
	}

	return dc, nil
}

// request DeviceConfiguration data from a remote entity
func (d *DeviceConfiguration) RequestDescriptions() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, nil, nil)
}

// request DeviceConfigurationKeyValueListDataType from a remote entity
func (d *DeviceConfiguration) RequestKeyValues() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceConfigurationKeyValueListData, nil, nil)
}

// return current descriptions for Device Configuration
func (d *DeviceConfiguration) GetDescriptions() ([]model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueDescriptionListDataType](d.featureRemote, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.DeviceConfigurationKeyValueDescriptionData, nil
}

// returns the description of a provided key name
func (d *DeviceConfiguration) GetDescriptionForKeyId(keyId model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	descriptions, err := d.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range descriptions {
		if item.KeyId != nil && *item.KeyId == keyId {
			return &item, nil
		}
	}

	return nil, api.ErrDataNotAvailable
}

// returns the description of a provided key name
// returns an error if the key name was not found
func (d *DeviceConfiguration) GetDescriptionForKeyName(keyName model.DeviceConfigurationKeyNameType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	descriptions, err := d.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range descriptions {
		if item.KeyId != nil &&
			item.KeyName != nil &&
			*item.KeyName == keyName {
			return &item, nil
		}
	}

	return nil, api.ErrDataNotAvailable
}

// return current values for Device Configuration
func (d *DeviceConfiguration) GetKeyValues() ([]model.DeviceConfigurationKeyValueDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueListDataType](d.featureRemote, model.FunctionTypeDeviceConfigurationKeyValueListData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.DeviceConfigurationKeyValueData, nil
}

// write key values
// returns an error if this failed
func (d *DeviceConfiguration) WriteKeyValues(data []model.DeviceConfigurationKeyValueDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionTypeDeviceConfigurationKeyValueListData),
		Filter:   []model.FilterType{*model.NewFilterTypePartial()},
		DeviceConfigurationKeyValueListData: &model.DeviceConfigurationKeyValueListDataType{
			DeviceConfigurationKeyValueData: data,
		},
	}

	return d.remoteDevice.Sender().Write(d.featureLocal.Address(), d.featureRemote.Address(), cmd)
}

// return a pointer value for a given key and value type
func (d *DeviceConfiguration) GetKeyValueForKeyName(keyname model.DeviceConfigurationKeyNameType, valueType model.DeviceConfigurationKeyValueTypeType) (any, error) {
	values, err := d.GetKeyValues()
	if err != nil {
		return nil, err
	}

	for _, item := range values {
		if item.KeyId == nil || item.Value == nil {
			continue
		}

		desc, err := d.GetDescriptionForKeyId(*item.KeyId)
		if err != nil {
			continue
		}

		if desc.KeyName != nil && *desc.KeyName == keyname {
			switch valueType {
			case model.DeviceConfigurationKeyValueTypeTypeBoolean:
				return item.Value.Boolean, nil
			case model.DeviceConfigurationKeyValueTypeTypeDate:
				return item.Value.Date, nil
			case model.DeviceConfigurationKeyValueTypeTypeDateTime:
				return item.Value.DateTime, nil
			case model.DeviceConfigurationKeyValueTypeTypeDuration:
				return item.Value.Duration, nil
			case model.DeviceConfigurationKeyValueTypeTypeString:
				return item.Value.String, nil
			case model.DeviceConfigurationKeyValueTypeTypeTime:
				return item.Value.Time, nil
			case model.DeviceConfigurationKeyValueTypeTypeScaledNumber:
				return item.Value.ScaledNumber, nil
			default:
				return nil, api.ErrDataNotAvailable
			}
		}
	}

	return nil, api.ErrDataNotAvailable
}
