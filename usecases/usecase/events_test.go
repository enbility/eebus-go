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

	payload.Device = s.remoteDevice
	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.uc.HandleEvent(payload)

	payload.Entity = s.remoteDevice.Entities()[0]
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
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

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
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

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
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

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
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

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
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

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
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.uc.useCaseScenarios[0].ServerFeatures = []model.FeatureTypeType{model.FeatureTypeTypeSmartEnergyManagementPs}

	s.uc.useCaseDataUpdate(payload)

	result = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), result)
}
