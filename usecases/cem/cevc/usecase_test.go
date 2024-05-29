package cevc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *CemCEVCSuite) Test_IsUseCaseSupported() {
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
						UseCaseName:      util.Ptr(model.UseCaseNameTypeCoordinatedEVCharging),
						UseCaseAvailable: util.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{2, 3, 4, 5, 6, 7, 8},
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
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	timeDescData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeTimeSeriesDescriptionListData, timeDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.IncentiveTableDescriptionDataType{
		IncentiveTableDescription: []model.IncentiveTableDescriptionType{
			{
				TariffDescription: &model.TariffDescriptionDataType{
					TariffId:  util.Ptr(model.TariffIdType(0)),
					ScopeType: util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
				},
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeIncentiveTableDescriptionData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
