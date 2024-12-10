package mpc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMPCSuite) Test_Power() {
	data, err := s.sut.Power(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Power(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

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

	data, err = s.sut.Power(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *MaMPCSuite) Test_PowerPerPhase() {
	data, err := s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
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

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
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

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
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

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)
}

func (s *MaMPCSuite) Test_EnergyConsumed() {
	data, err := s.sut.EnergyConsumed(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *MaMPCSuite) Test_EnergyProduced() {
	data, err := s.sut.EnergyProduced(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *MaMPCSuite) Test_CurrentPerPhase() {
	data, err := s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
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

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
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

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)
}

func (s *MaMPCSuite) Test_VoltagePerPhase() {
	data, err := s.sut.VoltagePerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(230),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(230),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(230),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
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

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{230, 230, 230}, data)
}

func (s *MaMPCSuite) Test_Frequency() {
	data, err := s.sut.Frequency(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Frequency(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, data)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}
