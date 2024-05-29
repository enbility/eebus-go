package evcem

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *EVCEMSuite) Test_EVConnectedPhases() {
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
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
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

	fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(1), data)
}

func (s *EVCEMSuite) Test_EVCurrentPerPhase() {
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
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
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
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
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

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])
}

func (s *EVCEMSuite) Test_EVPowerPerPhase_Power() {
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
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
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
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
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

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data[0])
}

func (s *EVCEMSuite) Test_EVPowerPerPhase_Current() {
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
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
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
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
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

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2300.0, data[0])
}

func (s *EVCEMSuite) Test_EVChargedEnergy() {
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
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
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

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data)
}
