package internal

import (
	"fmt"
	"slices"

	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Check if an incoming device configuration write requires write approval, if yes return `true` otherwise `false`
func ConfigurationWriteRequiresApproval(msg *spineapi.Message, localEntity spineapi.EntityLocalInterface, configsToApprove map[model.DeviceConfigurationKeyNameType]struct{}) (bool, error) {
	if msg.RequestHeader == nil || msg.RequestHeader.MsgCounter == nil ||
		msg.Cmd.DeviceConfigurationKeyValueListData == nil {
		return false, fmt.Errorf("invalid message")
	}

	data := msg.Cmd.DeviceConfigurationKeyValueListData

	if len(data.DeviceConfigurationKeyValueData) == 0 {
		return false, fmt.Errorf("no data")
	}

	// all DeviceConfigurationKeyValueData must have keyId set as primary identifier
	if slices.ContainsFunc(data.DeviceConfigurationKeyValueData, func(i model.DeviceConfigurationKeyValueDataType) bool {
		return i.KeyId == nil
	}) {
		return false, fmt.Errorf("invalid message")
	}

	dc, err := server.NewDeviceConfiguration(localEntity)
	if err != nil {
		return false, err
	}

	for _, deviceKeyValueData := range data.DeviceConfigurationKeyValueData {
		description, err := dc.GetKeyValueDescriptionFoKeyId(*deviceKeyValueData.KeyId)
		if description == nil || err != nil {
			logging.Log().Debug("ConfigurationWriteRequiresApproval: no device configuration for KeyID found: ", *deviceKeyValueData.KeyId)
			continue
		}

		if description.KeyName == nil {
			logging.Log().Debugf("ConfigurationWriteRequiresApproval: invalid internal data (KeyName not set on DeviceConfigurationKeyValueDescriptionDataType with key %s)", *deviceKeyValueData.KeyId)
			continue
		}

		// Only ask for write approval if at least one of the configurations we care about is trying to be set
		if _, exists := configsToApprove[*description.KeyName]; exists {
			return true, nil
		}
	}
	return false, nil
}

// Extract the device configuration writes from each pending message and return them grouped in a map by msgCounter
func GroupPendingDeviceConfigurations(pendingDeviceConfigs map[model.MsgCounterType]*spineapi.Message, localEntity spineapi.EntityLocalInterface) map[model.MsgCounterType][]ucapi.PendingDeviceConfiguration {
	result := make(map[model.MsgCounterType][]ucapi.PendingDeviceConfiguration)

	dc, err := server.NewDeviceConfiguration(localEntity)
	if err != nil {
		logging.Log().Debugf("GroupPendingDeviceConfigurations: Error occurred when getting device configuration: %s", err.Error())
		return result
	}

	for msgCounter, msg := range pendingDeviceConfigs {
		if msg.Cmd.DeviceConfigurationKeyValueListData == nil {
			continue
		}
		for _, configKeyValueData := range msg.Cmd.DeviceConfigurationKeyValueListData.DeviceConfigurationKeyValueData {
			if configKeyValueData.KeyId == nil {
				continue
			}
			description, err := dc.GetKeyValueDescriptionFoKeyId(*configKeyValueData.KeyId)
			if err != nil || description == nil || description.KeyName == nil {
				continue
			}

			result[msgCounter] = append(result[msgCounter], ucapi.PendingDeviceConfiguration{
				Description:       *description,
				KeyName:           *description.KeyName,
				Value:             configKeyValueData.Value,
				IsValueChangeable: configKeyValueData.IsValueChangeable,
			})
		}
	}
	return result
}
