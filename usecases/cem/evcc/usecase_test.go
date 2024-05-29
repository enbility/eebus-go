package evcc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemEVCCSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *CemEVCCSuite) Test_IsUseCaseSupported() {
	data, err := s.sut.IsUseCaseSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.IsUseCaseSupported(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	ucData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: util.Ptr(model.UseCaseActorTypeEV),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      util.Ptr(model.UseCaseNameTypeEVCommissioningAndConfiguration),
						UseCaseAvailable: util.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{1, 3, 4, 5, 6, 7, 8},
					},
				},
			},
		},
	}

	nodemgmtEntity := s.remoteDevice.Entity([]model.AddressEntityType{0})
	nodeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(nodemgmtEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	fErr := nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	ucData = &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: util.Ptr(model.UseCaseActorTypeEV),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      util.Ptr(model.UseCaseNameTypeEVCommissioningAndConfiguration),
						UseCaseAvailable: util.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
					},
				},
			},
		},
	}

	fErr = nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
