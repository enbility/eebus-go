package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *GcpMGCPSuite) Test_PowerLimitationFactor() {
	data, err := s.sut.PowerLimitationFactor(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.PowerLimitationFactor(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerLimitationFactor(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(10),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerLimitationFactor(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *GcpMGCPSuite) Test_Power() {
	data, err := s.sut.Power(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Power(s.smgwEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
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

	data, err = s.sut.Power(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *GcpMGCPSuite) Test_Power_ErrorCases() {
	// Test case where multiple measurement data entries are returned (len(data) != 1)
	// This should trigger the missing coverage on line 77-79 in Power()
	
	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Add measurement data for both descriptions
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(20),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Add electrical connection data to enable both measurements
	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	// This should now return multiple data points and trigger the len(data) != 1 condition
	data, err := s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *GcpMGCPSuite) Test_EnergyFeedIn() {
	data, err := s.sut.EnergyFeedIn(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP mandatory for energy
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
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

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *GcpMGCPSuite) Test_EnergyConsumed() {
	data, err := s.sut.EnergyConsumed(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP mandatory for energy
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
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

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *GcpMGCPSuite) Test_CurrentPerPhase() {
	data, err := s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
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

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)
}

func (s *GcpMGCPSuite) Test_VoltagePerPhase() {
	data, err := s.sut.VoltagePerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(230),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(230),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(230),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{230, 230, 230}, data)
}

func (s *GcpMGCPSuite) Test_Frequency() {
	data, err := s.sut.Frequency(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Frequency(s.smgwEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),           // MGCP required
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue), // MGCP recommended
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),         // MGCP recommended
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.smgwEntity)
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

	data, err = s.sut.Frequency(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

// Additional comprehensive tests for MGCP specification compliance

func (s *GcpMGCPSuite) Test_PowerWithInvalidValueType() {
	// Setup measurement description first
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test with invalid ValueType (averageValue instead of value)
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(2500),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeAverageValue), // Invalid
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *GcpMGCPSuite) Test_EnergyWithMissingValueSource() {
	// Setup measurement description for energy
	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test energy measurement without ValueSource (mandatory for energy per MGCP spec)
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(5000),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				// ValueSource missing - should fail for energy
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.EnergyFeedIn(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
}

func (s *GcpMGCPSuite) Test_FrequencyWithOutOfRangeState() {
	// Setup measurement description
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test with outOfRange state (should be skipped per MGCP-003)
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeOutofrange), // Should be skipped
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.Frequency(s.smgwEntity)
	assert.NotNil(s.T(), err) // Should fail due to no valid measurements
	assert.Equal(s.T(), 0.0, data)
}

func (s *GcpMGCPSuite) Test_CurrentWithMixedStates() {
	// Setup measurement description for current
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
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Setup electrical connection description
	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Setup electrical connection parameter description with phase mapping
	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test with mixed states - one error, one normal
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(15),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError), // Should be skipped
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(12),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal), // Should be used
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.CurrentPerPhase(s.smgwEntity)
	assert.Nil(s.T(), err) // Should succeed with one valid measurement
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 12.0, data[0])
}

func (s *GcpMGCPSuite) Test_VoltageValidationCompliance() {
	// Setup measurement description for voltage
	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Setup electrical connection description
	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Setup electrical connection parameter description
	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test voltage measurement without ValueSource (recommended for voltage per MGCP spec)
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(230),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				// ValueSource is recommended, not mandatory for voltage
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.VoltagePerPhase(s.smgwEntity)
	assert.Nil(s.T(), err) // Should succeed even without ValueSource
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 230.0, data[0])
}
