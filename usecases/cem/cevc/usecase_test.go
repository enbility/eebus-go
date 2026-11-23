package cevc

import (
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *CemCEVCSuite) Test_AddFeatures_GenericClient() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	assert.NotNil(s.T(), localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeGeneric, model.RoleTypeClient))

	for _, f := range localEntity.Features() {
		if f.Role() == model.RoleTypeClient {
			assert.Equal(s.T(), model.FeatureTypeTypeGeneric, f.Type())
		}
	}
}
