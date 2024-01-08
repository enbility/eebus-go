package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type DeviceConfiguration struct {
	*FeatureImpl
}

func NewDeviceConfiguration(localRole, remoteRole model.RoleType, localEntity *spine.EntityLocalImpl, remoteEntity *spine.EntityRemoteImpl) (*DeviceConfiguration, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceConfiguration, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceConfiguration{
		FeatureImpl: feature,
	}

	return dc, nil
}

// request DeviceConfiguration data from a remote entity
func (d *DeviceConfiguration) RequestDescriptions() error {
	_, err := d.requestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, nil, nil)
	return err
}

// request DeviceConfigurationKeyValueListDataType from a remote entity
func (d *DeviceConfiguration) RequestKeyValues() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceConfigurationKeyValueListData, nil, nil)
}

// return current descriptions for Device Configuration
func (d *DeviceConfiguration) GetDescriptions() ([]model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	rData := d.featureRemote.DataCopy(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceConfigurationKeyValueDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
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

	return nil, ErrDataNotAvailable
}

// returns the description of a provided key name
// returns an error if the key name was not found
func (d *DeviceConfiguration) GetDescriptionForKeyName(keyName model.DeviceConfigurationKeyNameType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	descriptions, err := d.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range descriptions {
		if item.KeyId == nil || item.KeyName == nil {
			continue
		}
		if *item.KeyName == keyName {
			return &item, nil
		}
	}

	return nil, ErrDataNotAvailable
}

// return current values for Device Configuration
func (d *DeviceConfiguration) GetKeyValues() ([]model.DeviceConfigurationKeyValueDataType, error) {
	rData := d.featureRemote.DataCopy(model.FunctionTypeDeviceConfigurationKeyValueListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceConfigurationKeyValueListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.DeviceConfigurationKeyValueData, nil
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

		if desc.KeyName == nil {
			continue
		}

		if *desc.KeyName == keyname {
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
				return nil, ErrDataNotAvailable
			}
		}
	}

	return nil, ErrDataNotAvailable
}
