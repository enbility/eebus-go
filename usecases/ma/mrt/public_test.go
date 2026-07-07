package mrt

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMRTSuite) Test_Temperature() {
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
				CommodityType:   util.Ptr(model.CommodityTypeTypeAir),
				ScopeType:       util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
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
				Value:         model.NewScaledNumberType(20),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	// announced in degC, requested in degC (conversion math covered by Test_convertTemperature)
	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 20.0, data)

	// requested in an unsupported unit
	_, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypeW)
	assert.NotNil(s.T(), err)

	// announced in degF, requested in degC
	descData.MeasurementDescriptionData[0].Unit = util.Ptr(model.UnitOfMeasurementTypedegF)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	fahrenheitData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(68),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, fahrenheitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Temperature(s.monitoredEntity, model.UnitOfMeasurementTypedegC)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 20.0, data)

	// value state not normal means the value must be ignored
	invalidData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(20),
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

func (s *MaMRTSuite) Test_convertTemperature() {
	tc := []struct {
		from, to model.UnitOfMeasurementType
		in, out  float64
		err      bool
	}{
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypedegC, 20, 20, false},
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypedegF, 20, 68, false},
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypeK, 0, 273.15, false},
		{model.UnitOfMeasurementTypedegF, model.UnitOfMeasurementTypedegC, 68, 20, false},
		{model.UnitOfMeasurementTypeK, model.UnitOfMeasurementTypedegC, 273.15, 0, false},
		{model.UnitOfMeasurementTypeW, model.UnitOfMeasurementTypedegC, 1, 0, true},
		{model.UnitOfMeasurementTypedegC, model.UnitOfMeasurementTypeW, 1, 0, true},
	}

	for _, tc := range tc {
		out, err := convertTemperature(tc.in, tc.from, tc.to)
		if tc.err {
			assert.NotNil(s.T(), err)
			continue
		}
		assert.Nil(s.T(), err)
		assert.InDelta(s.T(), tc.out, out, 1e-9)
	}
}
