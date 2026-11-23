package internal

import (
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *InternalSuite) Test_ConfigurationWriteRequiresApproval() {
	msg := spineapi.Message{}
	configsToApprove := map[model.DeviceConfigurationKeyNameType]struct{}{
		model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit: {},
		model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum:             {},
	}
	// Header missing
	_, err := ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.NotNil(s.T(), err)

	// MsgCounter missing
	header := model.HeaderType{}
	msg = spineapi.Message{RequestHeader: &header}
	_, err = ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.NotNil(s.T(), err)

	// DeviceConfigurationKeyValueListData missing
	header = model.HeaderType{MsgCounter: util.Ptr(model.MsgCounterType(1))}
	msg = spineapi.Message{RequestHeader: &header}
	_, err = ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.NotNil(s.T(), err)

	// DeviceConfigurationKeyValueListData.DeviceConfigurationKeyValueData is nil/length of 0
	cmd := model.CmdType{DeviceConfigurationKeyValueListData: util.Ptr(model.DeviceConfigurationKeyValueListDataType{})}
	msg = spineapi.Message{RequestHeader: &header, Cmd: cmd}
	_, err = ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.NotNil(s.T(), err)

	// Not all elements in slice of DeviceConfigurationKeyValueDataType have KeyId set
	deviceConfigList := []model.DeviceConfigurationKeyValueDataType{{KeyId: nil}}
	cmd = model.CmdType{DeviceConfigurationKeyValueListData: util.Ptr(model.DeviceConfigurationKeyValueListDataType{DeviceConfigurationKeyValueData: deviceConfigList})}
	msg = spineapi.Message{RequestHeader: &header, Cmd: cmd}
	_, err = ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.NotNil(s.T(), err)

	// Valid message but not a KeyId we care about => no approval required
	deviceConfigList = []model.DeviceConfigurationKeyValueDataType{{KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0))}}
	cmd = model.CmdType{DeviceConfigurationKeyValueListData: util.Ptr(model.DeviceConfigurationKeyValueListDataType{DeviceConfigurationKeyValueData: deviceConfigList})}
	msg = spineapi.Message{RequestHeader: &header, Cmd: cmd}
	approvalRequired, err := ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.Nil(s.T(), err)
	assert.False(s.T(), approvalRequired)

	// Valid message with KeyId we care about => approval required
	if dcs, err := server.NewDeviceConfiguration(s.localEntity); err == nil {
		dcs.AddKeyValueDescription(
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
				Unit:      util.Ptr(model.UnitOfMeasurementTypeW),
			},
		)

		value := &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(0),
		}
		_ = dcs.UpdateKeyValueDataForFilter(
			model.DeviceConfigurationKeyValueDataType{
				Value:             value,
				IsValueChangeable: util.Ptr(true),
			},
			nil,
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
			},
		)
	}
	approvalRequired, err = ConfigurationWriteRequiresApproval(&msg, s.localEntity, configsToApprove)
	assert.Nil(s.T(), err)
	assert.True(s.T(), approvalRequired)
}

func (s *InternalSuite) Test_GroupPendingDeviceConfigurations() {
	failsafeLimitDesc := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
		Unit:      util.Ptr(model.UnitOfMeasurementTypeW),
	}
	failsafeDurationMinDesc := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
	}

	dcs, err := server.NewDeviceConfiguration(s.localEntity)
	require.NoError(s.T(), err)
	failsafeLimitDesc.KeyId = dcs.AddKeyValueDescription(failsafeLimitDesc)
	require.NotNil(s.T(), failsafeLimitDesc.KeyId)
	failsafeDurationMinDesc.KeyId = dcs.AddKeyValueDescription(failsafeDurationMinDesc)
	require.NotNil(s.T(), failsafeDurationMinDesc.KeyId)

	value := &model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(0),
	}
	deviceConfigList := []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId:             failsafeLimitDesc.KeyId,
			Value:             value,
			IsValueChangeable: util.Ptr(true),
		},
		{KeyId: failsafeDurationMinDesc.KeyId},
		{KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(100))},
	}
	cmd := model.CmdType{DeviceConfigurationKeyValueListData: util.Ptr(model.DeviceConfigurationKeyValueListDataType{DeviceConfigurationKeyValueData: deviceConfigList})}
	msg := spineapi.Message{Cmd: cmd}
	pendingDeviceConfigs := map[model.MsgCounterType]*spineapi.Message{model.MsgCounterType(1): &msg}
	groupedConfigurations := GroupPendingDeviceConfigurations(pendingDeviceConfigs, s.localEntity)
	// For one of the KeyIds no corresponding device configuration exists, that element should thus be skipped
	expected := map[model.MsgCounterType][]ucapi.PendingDeviceConfiguration{
		model.MsgCounterType(1): {
			{
				Description:       failsafeLimitDesc,
				KeyName:           model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit,
				Value:             value,
				IsValueChangeable: util.Ptr(true),
			},
			{
				Description: failsafeDurationMinDesc,
				KeyName:     model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum,
			},
		},
	}
	assert.Equal(s.T(), expected, groupedConfigurations)
}

func (s *InternalSuite) Test_GroupPendingDeviceConfigurations_SkipsInvalidEntries() {
	descriptionWithoutKeyName := model.DeviceConfigurationKeyValueDescriptionDataType{
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeString),
	}
	validDescription := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
	}

	dcs, err := server.NewDeviceConfiguration(s.localEntity)
	require.NoError(s.T(), err)
	descriptionWithoutKeyName.KeyId = dcs.AddKeyValueDescription(descriptionWithoutKeyName)
	require.NotNil(s.T(), descriptionWithoutKeyName.KeyId)
	validDescription.KeyId = dcs.AddKeyValueDescription(validDescription)
	require.NotNil(s.T(), validDescription.KeyId)

	msgWithoutDeviceConfigurationData := spineapi.Message{}
	msgWithInvalidEntries := spineapi.Message{
		Cmd: model.CmdType{
			DeviceConfigurationKeyValueListData: util.Ptr(model.DeviceConfigurationKeyValueListDataType{
				DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
					{},
					{KeyId: descriptionWithoutKeyName.KeyId},
					{KeyId: validDescription.KeyId},
				},
			}),
		},
	}
	pendingDeviceConfigs := map[model.MsgCounterType]*spineapi.Message{
		model.MsgCounterType(1): &msgWithoutDeviceConfigurationData,
		model.MsgCounterType(2): &msgWithInvalidEntries,
	}

	groupedConfigurations := GroupPendingDeviceConfigurations(pendingDeviceConfigs, s.localEntity)

	expected := map[model.MsgCounterType][]ucapi.PendingDeviceConfiguration{
		model.MsgCounterType(2): {
			{
				Description: validDescription,
				KeyName:     model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit,
			},
		},
	}
	assert.Equal(s.T(), expected, groupedConfigurations)
}

func (s *InternalSuite) Test_GroupPendingDeviceConfigurations_ReturnsEmptyResultWhenDeviceConfigurationFeatureIsMissing() {
	localEntity := spinemocks.NewEntityLocalInterface(s.T())
	localEntity.EXPECT().Device().Return(nil)
	localEntity.EXPECT().
		FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer).
		Return(nil)

	msg := spineapi.Message{
		Cmd: model.CmdType{
			DeviceConfigurationKeyValueListData: util.Ptr(model.DeviceConfigurationKeyValueListDataType{
				DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
					{KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0))},
				},
			}),
		},
	}
	pendingDeviceConfigs := map[model.MsgCounterType]*spineapi.Message{
		model.MsgCounterType(1): &msg,
	}

	groupedConfigurations := GroupPendingDeviceConfigurations(pendingDeviceConfigs, localEntity)

	assert.Empty(s.T(), groupedConfigurations)
}
