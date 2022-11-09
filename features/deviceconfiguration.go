package features

import (
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type DeviceConfigurationType struct {
	Key           string
	ValueBool     bool
	ValueDate     time.Time
	ValueDatetime time.Time
	ValueDuration time.Duration
	ValueString   string
	ValueTime     time.Time
	ValueFloat    float64
	Type          model.DeviceConfigurationKeyValueTypeType
	Unit          string
}

type DeviceConfiguration struct {
	*FeatureImpl
}

func NewDeviceConfiguration(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*DeviceConfiguration, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceConfiguration, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceConfiguration{
		FeatureImpl: feature,
	}

	return dc, nil
}

// request DeviceConfiguration data from a remote entity
func (d *DeviceConfiguration) Request() error {
	// request DeviceConfigurationKeyValueDescriptionListData from a remote entity
	if _, err := d.requestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData); err != nil {
		logging.Log.Error(err)
		return err
	}

	return nil
}

// request DeviceConfigurationKeyValueListDataType from a remote entity
func (d *DeviceConfiguration) RequestKeyValueList() (*model.MsgCounterType, error) {
	// request FunctionTypeDeviceConfigurationKeyValueListData from a remote entity
	msgCounter, err := d.requestData(model.FunctionTypeDeviceConfigurationKeyValueListData)
	if err != nil {
		logging.Log.Error(err)
		return nil, err
	}

	return msgCounter, nil
}

// returns if a provided scopetype in the device configuration descriptions is available or not
// returns an error if no description data is available yet
func (d *DeviceConfiguration) GetDescriptionKeyNameSupport(keyName model.DeviceConfigurationKeyNameType) (bool, error) {
	rData := d.featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if rData == nil {
		return false, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceConfigurationKeyValueDescriptionListDataType)

	for _, item := range data.DeviceConfigurationKeyValueDescriptionData {
		if item.KeyId == nil || item.KeyName == nil {
			continue
		}
		if *item.KeyName == string(keyName) {
			return true, nil
		}
	}

	return false, ErrDataNotAvailable
}

// return current SoC for measurements
func (d *DeviceConfiguration) GetEVCommunicationStandard() (*string, error) {
	descRef, err := d.deviceConfigurationKeyValueDescriptionListData()
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	rData := d.featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceConfigurationKeyValueListDataType)

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId == nil || item.Value == nil {
			continue
		}

		desc, exists := descRef[*item.KeyId]
		if !exists {
			continue
		}

		if desc.KeyName == nil {
			continue
		}

		if *desc.KeyName == string(model.DeviceConfigurationKeyNameTypeCommunicationsStandard) {
			return (*string)(item.Value.String), nil
		}
	}

	return nil, ErrDataNotAvailable
}

// return current values for Device Configuration
func (d *DeviceConfiguration) GetValues() ([]DeviceConfigurationType, error) {
	if d.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	rDescData := d.featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if rDescData == nil {
		return nil, ErrMetadataNotAvailable
	}
	descData := rDescData.(*model.DeviceConfigurationKeyValueDescriptionListDataType)

	ref := make(map[model.DeviceConfigurationKeyIdType]model.DeviceConfigurationKeyValueDescriptionDataType)
	for _, item := range descData.DeviceConfigurationKeyValueDescriptionData {
		if item.KeyName == nil || item.KeyId == nil {
			continue
		}
		ref[*item.KeyId] = item
	}

	rData := d.featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceConfigurationKeyValueListDataType)
	var resultSet []DeviceConfigurationType

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId == nil {
			continue
		}
		desc, exists := ref[*item.KeyId]
		if !exists || desc.KeyName == nil {
			continue
		}

		result := DeviceConfigurationType{
			Key: *desc.KeyName,
		}
		if desc.ValueType == nil {
			continue
		}
		result.Type = *desc.ValueType
		switch *desc.ValueType {
		case model.DeviceConfigurationKeyValueTypeTypeBoolean:
			if item.Value.Boolean != nil {
				result.ValueBool = bool(*item.Value.Boolean)
			}
		case model.DeviceConfigurationKeyValueTypeTypeDate:
			if item.Value.Date != nil {
				if value, err := model.GetDateFromString(*item.Value.Date); err == nil {
					result.ValueDate = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeDateTime:
			if item.Value.DateTime != nil {
				if value, err := model.GetDateTimeFromString(*item.Value.DateTime); err == nil {
					result.ValueDatetime = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeDuration:
			if item.Value.Duration != nil {
				if value, err := item.Value.Duration.GetTimeDuration(); err == nil {
					result.ValueDuration = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeString:
			if item.Value.String != nil {
				result.ValueString = string(*item.Value.String)
			}
		case model.DeviceConfigurationKeyValueTypeTypeTime:
			if item.Value.Time != nil {
				if value, err := model.GetTime(item.Value.Time); err != nil {
					result.ValueTime = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeScalednumber:
			if item.Value.ScaledNumber != nil {
				result.ValueFloat = item.Value.ScaledNumber.GetValue()
			}
		}
		if desc.Unit != nil {
			result.Unit = *desc.Unit
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}

// helper

type deviceConfigurationKeyValueDescriptionMap map[model.DeviceConfigurationKeyIdType]model.DeviceConfigurationKeyValueDescriptionDataType

// return a map of DeviceConfigurationKeyValueDescriptionListDataType with keyId as key
func (d *DeviceConfiguration) deviceConfigurationKeyValueDescriptionListData() (deviceConfigurationKeyValueDescriptionMap, error) {
	if d.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	data := d.featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData).(*model.DeviceConfigurationKeyValueDescriptionListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}
	ref := make(deviceConfigurationKeyValueDescriptionMap)
	for _, item := range data.DeviceConfigurationKeyValueDescriptionData {
		if item.KeyId == nil {
			continue
		}
		ref[*item.KeyId] = item
	}
	return ref, nil
}
