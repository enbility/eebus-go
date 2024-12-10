package usecase

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UseCaseSuite) Test() {
	payload := spineapi.EventPayload{}
	result := s.uc.IsCompatibleEntityType(payload.Entity)
	assert.False(s.T(), result)

	payload = spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	result = s.uc.IsCompatibleEntityType(payload.Entity)
	assert.False(s.T(), result)

	payload = spineapi.EventPayload{
		Entity: s.monitoredEntity,
	}
	result = s.uc.IsCompatibleEntityType(payload.Entity)
	assert.True(s.T(), result)

	usecaseFilter := model.UseCaseFilterType{
		Actor:       useCaseActor,
		UseCaseName: useCaseName,
	}
	result = s.localEntity.HasUseCaseSupport(usecaseFilter)
	assert.False(s.T(), result)

	s.uc.AddUseCase()
	result = s.localEntity.HasUseCaseSupport(usecaseFilter)
	assert.True(s.T(), result)

	s.uc.UpdateUseCaseAvailability(false)
	result = s.localEntity.HasUseCaseSupport(usecaseFilter)
	assert.True(s.T(), result)

	s.uc.RemoveUseCase()
	result = s.localEntity.HasUseCaseSupport(usecaseFilter)
	assert.False(s.T(), result)
}

func (s *UseCaseSuite) Test_RemoveDeviceScenarios() {
	result := s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 0, len(result))

	scenarios := s.uc.AvailableScenariosForEntity(s.monitoredEntity)
	assert.Equal(s.T(), 0, len(scenarios))

	ok := s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), ok)

	s.uc.updateRemoteEntityScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2, 3})

	result = s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 1, len(result))

	scenarios = s.uc.AvailableScenariosForEntity(s.monitoredEntity)
	assert.Equal(s.T(), 3, len(scenarios))

	s.uc.removeDeviceFromAvailableEntityScenarios(s.monitoredEntity.Device())

	result = s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 0, len(result))

	scenarios = s.uc.AvailableScenariosForEntity(s.monitoredEntity)
	assert.Equal(s.T(), 0, len(scenarios))
}

func (s *UseCaseSuite) Test_AvailableScenarios() {
	result := s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 0, len(result))

	scenarios := s.uc.AvailableScenariosForEntity(s.monitoredEntity)
	assert.Equal(s.T(), 0, len(scenarios))

	ok := s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.False(s.T(), ok)

	s.uc.updateRemoteEntityScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2, 3})

	result = s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 1, len(result))

	scenarios = s.uc.AvailableScenariosForEntity(s.monitoredEntity)
	assert.Equal(s.T(), 3, len(scenarios))

	ok = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.True(s.T(), ok)

	s.uc.updateRemoteEntityScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2})

	scenarios = s.uc.AvailableScenariosForEntity(s.monitoredEntity)
	assert.Equal(s.T(), []uint{1, 2}, scenarios)

	ok = s.uc.IsScenarioAvailableAtEntity(s.monitoredEntity, 1)
	assert.True(s.T(), ok)

	s.uc.removeEntityFromAvailableEntityScenarios(s.monitoredEntity)

	result = s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 0, len(result))

	s.uc.updateRemoteEntityScenarios(s.monitoredEntity, []model.UseCaseScenarioSupportType{1, 2, 3})

	result = s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 1, len(result))

	s.uc.removeDeviceFromAvailableEntityScenarios(s.monitoredEntity.Device())

	result = s.uc.RemoteEntitiesScenarios()
	assert.Equal(s.T(), 0, len(result))
}

func (s *UseCaseSuite) Test_RequiredServerFeatures() {
	required := s.uc.requiredServerFeaturesForScenario(model.UseCaseScenarioSupportType(1))
	assert.Equal(s.T(), 1, len(required))

	required = s.uc.requiredServerFeaturesForScenario(model.UseCaseScenarioSupportType(2))
	assert.Equal(s.T(), 0, len(required))

	required = s.uc.requiredServerFeaturesForScenario(model.UseCaseScenarioSupportType(4))
	assert.Equal(s.T(), 0, len(required))
}
