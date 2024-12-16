package internal

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *InternalSuite) Test_MeasurementPhaseSpecificDataForFilter() {
	measurementType := model.MeasurementTypeTypePower
	commodityType := model.CommodityTypeTypeElectricity
	scopeType := model.ScopeTypeTypeACPower
	energyDirection := model.EnergyDirectionTypeConsume

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: &measurementType,
		CommodityType:   &commodityType,
		ScopeType:       &scopeType,
	}

	data, err := MeasurementPhaseSpecificDataForFilter(nil, nil, filter, energyDirection, ucapi.PhaseNameMapping)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.mockRemoteEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.monitoredEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				ScopeType: util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.monitoredEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(10)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.monitoredEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{0, 0, 0}, data)

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.monitoredEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(10)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.monitoredEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *InternalSuite) Test_GetPowerTotalMeasurementId() {
	// Testing nil localEntity should return defaultPowerTotalMeasurementId
	measurementId := GetPowerTotalMeasurementId(nil)
	assert.Equal(s.T(), defaultPowerTotalMeasurementId, measurementId)

	// Testing localEntity without measurement server feature should return defaultPowerTotalMeasurementId
	// The test setup doesn't include a measurement server feature, only client features
	measurementId = GetPowerTotalMeasurementId(s.localEntity)
	assert.Equal(s.T(), defaultPowerTotalMeasurementId, measurementId)

	// Testing Create a test entity with measurement server feature
	localDevice := s.service.LocalDevice()
	testEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	// Add a measurement server feature
	measurementServerFeature := spine.NewFeatureLocal(10, testEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	measurementServerFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	testEntity.AddFeature(measurementServerFeature)

	// Initially, there should be no measurement descriptions, so it should return defaultPowerTotalMeasurementId
	measurementId = GetPowerTotalMeasurementId(testEntity)
	assert.Equal(s.T(), defaultPowerTotalMeasurementId, measurementId)

	// Testing Add measurement descriptions that don't match the filter
	descDataNoMatch := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(100)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent), // Different type
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	fErr := measurementServerFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descDataNoMatch, nil, nil)
	assert.Nil(s.T(), fErr)

	// Should still return defaultPowerTotalMeasurementId as no matching descriptions
	measurementId = GetPowerTotalMeasurementId(testEntity)
	assert.Equal(s.T(), defaultPowerTotalMeasurementId, measurementId)

	// Testing Add a matching power total measurement description
	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(200)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				// Add another measurement to ensure it only matches the correct one
				MeasurementId:   util.Ptr(model.MeasurementIdType(101)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	fErr = measurementServerFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Now it should return the correct measurement ID
	measurementId = GetPowerTotalMeasurementId(testEntity)
	assert.Equal(s.T(), model.MeasurementIdType(200), measurementId)

	// Testing: Multiple matching descriptions (though this shouldn't happen in practice)
	descDataMultiple := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(200)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(300)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	fErr = measurementServerFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descDataMultiple, nil, nil)
	assert.Nil(s.T(), fErr)

	// Should return defaultPowerTotalMeasurementId because there are multiple matching descriptions (length != 1)
	measurementId = GetPowerTotalMeasurementId(testEntity)
	assert.Equal(s.T(), defaultPowerTotalMeasurementId, measurementId)

	// Testing Measurement description without MeasurementId should return defaultPowerTotalMeasurementId
	descDataNoId := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				// MeasurementId is nil
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	fErr = measurementServerFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descDataNoId, nil, nil)
	assert.Nil(s.T(), fErr)

	measurementId = GetPowerTotalMeasurementId(testEntity)
	assert.Equal(s.T(), defaultPowerTotalMeasurementId, measurementId)
}

func (s *InternalSuite) Test_MeasurementSinglePhaseSpecificDataForFilter() {
	measurementType := model.MeasurementTypeTypePower
	commodityType := model.CommodityTypeTypeElectricity
	scopeType := model.ScopeTypeTypeACPower
	energyDirection := model.EnergyDirectionTypeConsume

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: &measurementType,
		CommodityType:   &commodityType,
		ScopeType:       &scopeType,
	}

	// set up ElectricalConnection
	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)

	_, fErr := rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := MeasurementPhaseSpecificDataForFilter(
		s.localEntity,
		s.monitoredEntity,
		filter,
		energyDirection,
		ucapi.PhaseNameMapping,
	)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{0, 10, 0}, data)

}
