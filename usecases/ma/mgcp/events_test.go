package mgcp

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *GcpMGCPSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.smgwEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.DeviceConfigurationKeyValueDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.MeasurementDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.MeasurementListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *GcpMGCPSuite) Test_Failures() {
	s.sut.gridConnected(s.mockRemoteEntity)

	s.sut.gridConfigurationDescriptionDataUpdate(s.mockRemoteEntity)

	s.sut.gridMeasurementDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *GcpMGCPSuite) Test_gridConfigurationDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.smgwEntity,
	}
	s.sut.gridConfigurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

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

	s.sut.gridConfigurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

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

	payload.Data = keyData

	// Update the feature with the data so it's actually stored
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.gridConfigurationDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *GcpMGCPSuite) Test_gridMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.smgwEntity,
	}
	s.sut.gridMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

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
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(3)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(4)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(5)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Add electrical connection setup for complete validation
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(3)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(4)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(5)),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.gridMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(10),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(10),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(10),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(10),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(10),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(5)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(10),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
	}

	payload.Data = data

	// Update the feature with the data so it's actually stored
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.gridMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *GcpMGCPSuite) Test_gridMeasurementDataUpdate_PhaseEventCallbacks() {
	// Test missing event callback coverage for CurrentPerPhase and VoltagePerPhase
	
	// Setup measurement descriptions for current and voltage
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
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Setup electrical connection
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
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test CurrentPerPhase event callback
	currentData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(15),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		},
	}

	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.smgwEntity,
		Data:   currentData,
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, currentData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.eventCalled = false
	s.sut.gridMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)

	// Test VoltagePerPhase event callback
	voltageData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				Value:         model.NewScaledNumberType(230),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		},
	}

	payload.Data = voltageData
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, voltageData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.eventCalled = false
	s.sut.gridMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}
