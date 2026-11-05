package opev

import (
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *CemOPEVSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *CemOPEVSuite) Test_AddFeatures_GenericClient() {
	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	assert.NotNil(s.T(), localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeGeneric, model.RoleTypeClient))

	// the test helper adds a Measurement client feature on its own, so only check the types of this use case
	for _, ft := range []model.FeatureTypeType{
		model.FeatureTypeTypeLoadControl,
		model.FeatureTypeTypeElectricalConnection,
	} {
		assert.Nil(s.T(), localEntity.FeatureOfTypeAndRole(ft, model.RoleTypeClient))
	}
}
