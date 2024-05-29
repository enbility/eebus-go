package evsoc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *EVSOCSuite) Test_StateOfCharge() {
	data, err := s.sut.StateOfCharge(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.StateOfCharge(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	ucData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: util.Ptr(model.UseCaseActorTypeEV),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      util.Ptr(model.UseCaseNameTypeEVStateOfCharge),
						UseCaseAvailable: util.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{1},
					},
				},
			},
		},
	}

	nodemgmtEntity := s.remoteDevice.Entity([]model.AddressEntityType{0})
	nodeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(nodemgmtEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	fErr := nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.StateOfCharge(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePercentage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.StateOfCharge(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.StateOfCharge(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.StateOfCharge(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data)
}
