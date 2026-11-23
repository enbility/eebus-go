package vabd

import (
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *CemVABDSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *CemVABDSuite) Test_AddFeatures_GenericClient() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	assert.NotNil(s.T(), localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeGeneric, model.RoleTypeClient))

	for _, f := range localEntity.Features() {
		if f.Role() == model.RoleTypeClient {
			assert.Equal(s.T(), model.FeatureTypeTypeGeneric, f.Type())
		}
	}
}
