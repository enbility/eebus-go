package evcem

import (
	"time"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemEVCEMSuite) Test_EVConnectedPhases() {
	data, err := s.sut.PhasesConnected(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	descData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	descData = &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				AcConnectedPhases:      util.Ptr(uint(1)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(1), data)
}

func (s *CemEVCEMSuite) Test_EVCurrentPerPhase() {
	data, err := s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

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

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])

	now := time.Now().Add(-50 * time.Second)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(now),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])

	now = now.Add(-1 * time.Hour)
	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(now),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])

	now = now.Add(-1 * time.Hour)
	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(now),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])
}

func (s *CemEVCEMSuite) Test_EVCurrentPerPhase_AudiConnect() {
	data, err := s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(1)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(4)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(7)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(8)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(4)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(7)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
				ScopeType:       util.Ptr(model.ScopeTypeTypeCharge),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])

	now := time.Now().Add(-50 * time.Second)

	measData = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(now),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
				Value:         model.NewScaledNumberType(10),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(now),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(7)),
				Value:         model.NewScaledNumberType(10),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(now),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 10.0, data[0])
}

func (s *CemEVCEMSuite) Test_EVPowerPerPhase_Power() {
	data, err := s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPower),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data[0])
}

func (s *CemEVCEMSuite) Test_EVPowerPerPhase_Current() {
	data, err := s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(4)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(3)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(5)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(4)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(6)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(5)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(7)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(6)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(3)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(4)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(5)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(6)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeCharge),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(5.09),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(4.04),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(5.09),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(5)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(6)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
}

func (s *CemEVCEMSuite) Test_EVChargedEnergy() {
	data, err := s.sut.EnergyCharged(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeCharge),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data)
}

func (s *CemEVCEMSuite) Test_EVChargedEnergy_ElliGen1() {
	data, err := s.sut.EnergyCharged(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(3)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(4)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(5)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(6)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeCharge),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeWh),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(5.09),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(4.04),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(5.09),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(5)),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(6)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}
