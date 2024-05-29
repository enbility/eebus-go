package usecase

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UseCaseSuite) Test() {
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}
	uc := NewUseCaseBase(
		s.localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1},
		nil,
		validEntityTypes,
	)

	payload := spineapi.EventPayload{}
	result := uc.IsCompatibleEntity(payload.Entity)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	result = uc.IsCompatibleEntity(payload.Entity)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.monitoredEntity,
	}
	result = uc.IsCompatibleEntity(payload.Entity)
	assert.Equal(s.T(), true, result)

	uc.AddUseCase()
	uc.UpdateUseCaseAvailability(false)
}
