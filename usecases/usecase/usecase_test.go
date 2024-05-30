package usecase

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *UseCaseSuite) Test() {
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEVSE}
	uc := NewUseCaseBase(
		s.localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		nil,
		validEntityTypes,
	)

	payload := spineapi.EventPayload{}
	result := uc.IsCompatibleEntityType(payload.Entity)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	result = uc.IsCompatibleEntityType(payload.Entity)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.evseEntity,
	}
	result = uc.IsCompatibleEntityType(payload.Entity)
	assert.Equal(s.T(), true, result)

	uc.AddUseCase()
	uc.UpdateUseCaseAvailability(false)
	uc.RemoveUseCase()
}

func (s *UseCaseSuite) Test_SupportedUseCaseScenarios() {
	ucName := model.UseCaseNameTypeEVSECommissioningAndConfiguration
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEVSE}
	uc := NewUseCaseBase(
		s.localEntity,
		model.UseCaseActorTypeCEM,
		ucName,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		nil,
		validEntityTypes,
	)

	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	feature := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeControllableSystem,
		model.UseCaseNameTypeLimitationOfPowerConsumption,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEVSE,
		ucName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	result := uc.SupportedUseCaseScenarios(nil)
	assert.Equal(s.T(), 0, len(result))

	result = uc.SupportedUseCaseScenarios(s.monitoredEntity)
	assert.Equal(s.T(), []model.UseCaseScenarioSupportType{1, 2, 3}, result)

	result = uc.SupportedUseCaseScenarios(s.evseEntity)
	assert.Equal(s.T(), 0, len(result))

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEVSE,
		ucName,
		"1.0.0",
		"release",
		false,
		[]model.UseCaseScenarioSupportType{1, 2})

	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	result = uc.SupportedUseCaseScenarios(s.monitoredEntity)
	assert.Equal(s.T(), []model.UseCaseScenarioSupportType{1, 2}, result)
}

func (s *UseCaseSuite) Test_HasSupportForUseCaseScenarios() {
	ucName := model.UseCaseNameTypeEVSECommissioningAndConfiguration
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEVSE}
	uc := NewUseCaseBase(
		s.localEntity,
		model.UseCaseActorTypeCEM,
		ucName,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		nil,
		validEntityTypes,
	)

	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	feature := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeControllableSystem,
		model.UseCaseNameTypeLimitationOfPowerConsumption,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEVSE,
		ucName,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	result := uc.HasSupportForUseCaseScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2, 3})
	assert.True(s.T(), result)

	result = uc.HasSupportForUseCaseScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{2, 4})
	assert.False(s.T(), result)

	result = uc.HasSupportForUseCaseScenarios(s.monitoredEntity, nil)
	assert.False(s.T(), result)

	result = uc.HasSupportForUseCaseScenarios(s.evseEntity, []model.UseCaseScenarioSupportType{2, 4})
	assert.False(s.T(), result)

	data = &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEVSE,
		ucName,
		"1.0.0",
		"release",
		false,
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	result = uc.HasSupportForUseCaseScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2, 3})
	assert.False(s.T(), result)
}

func (s *UseCaseSuite) Test_UseCaseDataUpdate() {
	ucName := model.UseCaseNameTypeEVSECommissioningAndConfiguration
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEVSE}
	uc := NewUseCaseBase(
		s.localEntity,
		model.UseCaseActorTypeCEM,
		ucName,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		nil,
		validEntityTypes,
	)

	result := uc.hasRemoteEntity(s.monitoredEntity)
	assert.False(s.T(), result)

	uc.removeRemoteEntity(s.monitoredEntity)

	scenarios := uc.scenariosForRemoteEntity(s.monitoredEntity)
	assert.Nil(s.T(), scenarios)

	s.mux.Lock()
	cbInvoked := false
	s.mux.Unlock()

	uc.setRemoteEntityScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2, 3})
	result = uc.hasRemoteEntity(s.monitoredEntity)
	assert.True(s.T(), result)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.monitoredEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeRemove,
	}

	cb := func(
		ski string,
		device spineapi.DeviceRemoteInterface,
		entity spineapi.EntityRemoteInterface,
		event api.EventType,
	) {
		s.mux.Lock()
		cbInvoked = true
		s.mux.Unlock()
	}

	s.mux.Lock()
	cbInvoked = false
	s.mux.Unlock()

	uc.UseCaseDataUpdate(payload, cb, api.EventType("test"))
	assert.True(s.T(), cbInvoked)

	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	feature := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		*feature.Address(),
		model.UseCaseActorTypeEVSE,
		ucName,
		"1.0.0",
		"release",
		false,
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload = spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.monitoredEntity,
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
	}

	s.mux.Lock()
	cbInvoked = false
	s.mux.Unlock()

	uc.UseCaseDataUpdate(payload, cb, api.EventType("test"))
	assert.True(s.T(), cbInvoked)

	result = uc.hasRemoteEntity(s.monitoredEntity)
	assert.True(s.T(), result)

	scenarios = uc.scenariosForRemoteEntity(s.monitoredEntity)
	assert.NotNil(s.T(), scenarios)

	entities := uc.RemoteEntities()
	assert.NotNil(s.T(), entities)
	assert.Equal(s.T(), 1, len(entities))

	data = &model.NodeManagementUseCaseDataType{}
	nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	s.mux.Lock()
	cbInvoked = false
	s.mux.Unlock()

	uc.UseCaseDataUpdate(payload, cb, api.EventType("test"))
	assert.True(s.T(), cbInvoked)

	result = uc.hasRemoteEntity(s.monitoredEntity)
	assert.False(s.T(), result)
}
