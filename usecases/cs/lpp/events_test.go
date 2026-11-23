package lpp

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CsLPPSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity:    s.mockRemoteEntity,
		EventType: spineapi.EventTypeSubscriptionChange,
	}
	s.sut.HandleEvent(payload)

	payload.Device = s.monitoredEntity.Device()
	payload.Entity = s.monitoredEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	// test scenario gate rejects data changes without scenarios
	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.CmdClassifier = util.Ptr(model.CmdClassifierTypeWrite)
	payload.Data = &model.LoadControlLimitListDataType{}
	s.sut.HandleEvent(payload)

	s.setUpUseCaseScenarios()

	payload.Data = nil
	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.CmdClassifier = util.Ptr(model.CmdClassifierTypeWrite)
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Function = model.FunctionTypeLoadControlLimitListData
	payload.Data = util.Ptr(model.LoadControlLimitListDataType{})
	s.sut.HandleEvent(payload)

	payload.LocalFeature = s.loadControlFeature
	s.sut.HandleEvent(payload)

	payload.Function = model.FunctionTypeDeviceConfigurationKeyValueListData
	payload.Data = util.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)

	payload.LocalFeature = s.deviceConfigurationFeature
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeBindingChange
	payload.ChangeType = spineapi.ElementChangeAdd
	payload.LocalFeature = s.loadControlFeature
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Function = model.FunctionTypeDeviceDiagnosisHeartbeatData
	payload.LocalFeature = s.deviceDiagnosisFeature
	payload.CmdClassifier = util.Ptr(model.CmdClassifierTypeNotify)
	payload.Data = util.Ptr(model.DeviceDiagnosisHeartbeatDataType{})
	s.sut.HandleEvent(payload)

	payload.Function = model.FunctionTypeDeviceConfigurationKeyValueListData
	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *CsLPPSuite) Test_subscribeHeartbeat() {
	// test heartbeat subscription with a compatible entity
	s.sut.subscribeHeartbeat(s.monitoredEntity)
}

func (s *CsLPPSuite) Test_loadControlLimitDataUpdate() {
	localDevice := s.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
	}
	s.sut.loadControlLimitDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	lFeature := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	lFeature.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, descData)

	s.sut.loadControlLimitDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	payload.Data = data

	s.sut.loadControlLimitDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
		},
	}

	payload.Data = data

	s.sut.loadControlLimitDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *CsLPPSuite) Test_configurationDataUpdate() {
	localDevice := s.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)
	lFeature := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

	payload := spineapi.EventPayload{
		Ski:          remoteSki,
		Device:       s.remoteDevice,
		Entity:       s.monitoredEntity,
		LocalFeature: lFeature,
	}

	s.sut.configurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			},
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(2)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			},
		},
	}

	lFeature.SetData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData)

	s.eventCalled = false
	s.sut.configurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{},
	}

	payload.Data = data

	s.eventCalled = false
	s.sut.configurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: util.Ptr(model.DeviceConfigurationKeyValueValueType{}),
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(2)),
				Value: util.Ptr(model.DeviceConfigurationKeyValueValueType{}),
			},
		},
	}

	payload.Data = data

	s.eventCalled = false
	s.sut.configurationDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *CsLPPSuite) setUpUseCaseScenarios() {
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEnergyGuard,
		model.UseCaseNameTypeLimitationOfPowerProduction,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.monitoredEntity,
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.sut.UseCaseBase.HandleEvent(payload)
}
