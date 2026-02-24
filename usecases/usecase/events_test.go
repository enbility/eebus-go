package usecase

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *UseCaseSuite) Test_HandleEvent() {
	payload := spineapi.EventPayload{}
	s.uc.HandleEvent(payload)

	payload.Ski = "test"
	payload.Device = s.remoteDevice
	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.uc.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.Entity = s.remoteDevice.Entities()[0]
	s.uc.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = &model.NodeManagementUseCaseDataType{}
	s.uc.HandleEvent(payload)

	payload.Data = &model.NodeManagementDetailedDiscoveryDataType{}
	s.uc.HandleEvent(payload)
}

func (s *UseCaseSuite) Test_useCaseDataUpdate() {
	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.remoteDevice.Entities()[0],
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.uc.useCaseDataUpdate(payload)

	result := s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)

	feature := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)
	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeHeatPump,
		useCaseName,
		"1.0.0",
		"release",
		false,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		false,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.True(s.T(), result)

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.True(s.T(), result)

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseScenarios[0].ServerFeatures = []model.FeatureTypeType{model.FeatureTypeTypeSmartEnergyManagementPs}

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)
}

func (s *UseCaseSuite) Test_OnScenariosChanged() {
	// Track callback invocations
	var callbackEntity spineapi.EntityRemoteInterface
	var callbackScenarios []uint
	callbackCount := 0

	s.uc.OnScenariosChanged = func(entity spineapi.EntityRemoteInterface, scenarios []uint) {
		callbackEntity = entity
		callbackScenarios = scenarios
		callbackCount++
	}

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.remoteDevice.Entities()[0],
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}

	// Set up use case data with valid scenarios
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	// Callback should have been called with non-empty scenarios
	assert.Equal(s.T(), 1, callbackCount)
	assert.NotNil(s.T(), callbackEntity)
	assert.NotEmpty(s.T(), callbackScenarios)

	// Update with same scenarios should not trigger callback again
	s.uc.useCaseDataUpdate(payload)
	assert.Equal(s.T(), 1, callbackCount)

	// Update with different scenarios should trigger callback
	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)
	assert.Equal(s.T(), 2, callbackCount)
}

func (s *UseCaseSuite) Test_OnScenariosChanged_EmptyScenarios() {
	callbackCount := 0
	s.uc.OnScenariosChanged = func(entity spineapi.EntityRemoteInterface, scenarios []uint) {
		callbackCount++
	}

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.remoteDevice.Entities()[0],
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}

	// Set up use case data with wrong actor type (results in empty scenarios)
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeHeatPump, // wrong actor type
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	// Callback should NOT have been called because no scenarios matched
	assert.Equal(s.T(), 0, callbackCount)
}

func (s *UseCaseSuite) Test_useCaseDataUpdate_UseCaseRemoved() {
	// First, set up valid scenarios
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.remoteDevice.Entities()[0],
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.uc.useCaseDataUpdate(payload)

	// Scenarios should be available
	assert.True(s.T(), s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1))

	// Now replace with a completely different use case
	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		model.UseCaseNameTypeCoordinatedEVCharging, // different use case
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	// Scenarios should now be cleared
	assert.False(s.T(), s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1))
	assert.False(s.T(), s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 2))
}

func (s *UseCaseSuite) Test_useCaseDataUpdate_UseCaseUnavailable() {
	// First, set up valid scenarios
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.remoteDevice.Entities()[0],
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.uc.useCaseDataUpdate(payload)

	assert.True(s.T(), s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1))

	// Now set the use case as unavailable
	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		useCaseName,
		"1.0.0",
		"release",
		false, // unavailable
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	// Scenarios should be cleared
	assert.False(s.T(), s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1))
}

func (s *UseCaseSuite) Test_useCaseDataUpdate_PMCP() {
	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.remoteDevice.Entities()[0],
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.uc.useCaseDataUpdate(payload)

	result := s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)

	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	// the PMCP device reports the wrong entity address for the EV use cases :(
	data := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Address: &model.FeatureAddressType{
					Device: util.Ptr(model.AddressDeviceType("d:_i:19667_PorscheEVSE_0012345")),
					Entity: []model.AddressEntityType{1},
				},
				Actor: util.Ptr(model.UseCaseActorTypeEVSE),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:     util.Ptr(model.UseCaseNameType("evseCommissioningAndConfiguration")),
						UseCaseVersion:  util.Ptr(model.SpecificationVersionType("1.0.0")),
						ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2},
					},
				},
			},
			{
				Address: &model.FeatureAddressType{
					Device: util.Ptr(model.AddressDeviceType("d:_i:19667_PorscheEVSE_0012345")),
					Entity: []model.AddressEntityType{1}, // this should include the subentity 1 to point to the EV entity
				},
				Actor: util.Ptr(model.UseCaseActorTypeEV),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:     util.Ptr(model.UseCaseNameType("evCommissioningAndConfiguration")),
						UseCaseVersion:  util.Ptr(model.SpecificationVersionType("1.0.0")),
						ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3, 6, 7, 8},
					},
					{
						UseCaseName:     util.Ptr(model.UseCaseNameType("overloadProtectionByEvChargingCurrentCurtailment")),
						UseCaseVersion:  util.Ptr(model.SpecificationVersionType("1.0.0")),
						ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3},
					},
					{
						UseCaseName:     util.Ptr(model.UseCaseNameType("measurementOfElectricityDuringEvCharging")),
						UseCaseVersion:  util.Ptr(model.SpecificationVersionType("1.0.0")),
						ScenarioSupport: []model.UseCaseScenarioSupportType{1},
					},
				},
			},
		},
	}
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.True(s.T(), result)
}
