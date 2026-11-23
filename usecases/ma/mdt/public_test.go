package mdt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMDTSuite) Test_Temperature() {
	data, err := s.sut.Temperature(s.mockRemoteEntity, model.UnitOfMeasurementTypedegC)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeTemperature),
				CommodityType:   util.Ptr(model.CommodityTypeTypeDomestichotwater),
				ScopeType:       util.Ptr(model.ScopeTypeTypeDhwTemperature),
				Unit:            util.Ptr(model.UnitOfMeasurementTypedegC),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	// announced in degC, requested in degC (conversion math covered by Test_convertTemperature)
	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, data)

	// requested in an unsupported unit
	_, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypeW)
	assert.ErrorIs(s.T(), err, api.ErrNotSupported)

	// announced in degF, requested in degC
	descData.MeasurementDescriptionData[0].Unit = util.Ptr(model.UnitOfMeasurementTypedegF)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	fahrenheitData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(122),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, fahrenheitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, data)

	// value state not normal means the value must be ignored
	invalidData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, invalidData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *MaMDTSuite) Test_convertTemperature() {
	tc := []struct {
		from, to model.UnitOfMeasurementType
		in, out  float64
		err      error
	}{
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypedegC, 50, 50, nil},
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypedegF, 50, 122, nil},
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypeK, 0, 273.15, nil},
		{model.UnitOfMeasurementTypedegF, model.UnitOfMeasurementTypedegC, 122, 50, nil},
		{model.UnitOfMeasurementTypeK, model.UnitOfMeasurementTypedegC, 273.15, 0, nil},
		{model.UnitOfMeasurementTypeW, model.UnitOfMeasurementTypedegC, 1, 0, api.ErrDataInvalid},
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypeW, 1, 0, api.ErrNotSupported},
	}

	for _, tc := range tc {
		out, err := convertTemperature(tc.in, tc.from, tc.to)
		if tc.err != nil {
			assert.ErrorIs(s.T(), err, tc.err)
			continue
		}
		assert.Nil(s.T(), err)
		assert.InDelta(s.T(), tc.out, out, 1e-9)
	}
}
