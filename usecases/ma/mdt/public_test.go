package mdt

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMDTSuite) Test_Temperature() {
	data, err := s.sut.Temperature(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Temperature(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeTemperature),
				CommodityType:   util.Ptr(model.CommodityTypeTypeDomestichotwater),
				ScopeType:       util.Ptr(model.ScopeTypeTypeDhwTemperature),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Temperature(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(55),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Temperature(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 55.0, data)
}
