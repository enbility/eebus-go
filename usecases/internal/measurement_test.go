package internal

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
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
	assert.Equal(s.T(), 0, len(data))

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
