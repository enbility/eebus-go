package lpp

import (
	"fmt"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
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

func (s *CsLPPSuite) Test_deviceConnected() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}

	s.sut.deviceConnected(payload)

	// no entities
	mockRemoteDevice := mocks.NewDeviceRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().Entities().Return(nil)
	payload.Device = mockRemoteDevice
	s.sut.deviceConnected(payload)

	// one entity with one DeviceDiagnosis server
	payload.Device = s.remoteDevice
	s.sut.deviceConnected(payload)

	s.sut.subscribeHeartbeatWorkaround(payload)
}

func (s *CsLPPSuite) Test_multipleDeviceDiagServer() {
	// multiple entities each with DeviceDiagnosis server

	payload := spineapi.EventPayload{
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}

	remoteDeviceName := "remote"

	var remoteFeatures = []struct {
		featureType   model.FeatureTypeType
		role          model.RoleType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeLoadControl,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
		{model.FeatureTypeTypeDeviceConfiguration,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeDeviceDiagnosisHeartbeatData,
			},
		},
		{model.FeatureTypeTypeElectricalConnection,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
	}
	var featureInformations []model.NodeManagementDetailedDiscoveryFeatureInformationType
	// 4 entities
	for i := 1; i < 5; i++ {
		for index, feature := range remoteFeatures {
			supportedFcts := []model.FunctionPropertyType{}
			for _, fct := range feature.supportedFcts {
				supportedFct := model.FunctionPropertyType{
					Function: util.Ptr(fct),
					PossibleOperations: &model.PossibleOperationsType{
						Read: &model.PossibleOperationsReadType{},
					},
				}
				supportedFcts = append(supportedFcts, supportedFct)
			}

			featureInformation := model.NodeManagementDetailedDiscoveryFeatureInformationType{
				Description: &model.NetworkManagementFeatureDescriptionDataType{
					FeatureAddress: &model.FeatureAddressType{
						Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity:  []model.AddressEntityType{model.AddressEntityType(i)},
						Feature: util.Ptr(model.AddressFeatureType(index)),
					},
					FeatureType:       util.Ptr(feature.featureType),
					Role:              util.Ptr(feature.role),
					SupportedFunction: supportedFcts,
				},
			}
			featureInformations = append(featureInformations, featureInformation)
		}
	}

	detailedData := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeCEM),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{2},
					},
					EntityType: util.Ptr(model.EntityTypeTypeCEM),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{3},
					},
					EntityType: util.Ptr(model.EntityTypeTypeCEM),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{4},
					},
					EntityType: util.Ptr(model.EntityTypeTypeCEM),
				},
			},
		},
		FeatureInformation: featureInformations,
	}

	_, err := s.remoteDevice.AddEntityAndFeatures(true, detailedData)
	if err != nil {
		fmt.Println(err)
	}
	s.remoteDevice.UpdateDevice(detailedData.DeviceInformation.Description)

	s.sut.deviceConnected(payload)

	s.sut.subscribeHeartbeatWorkaround(payload)
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
