package mpc

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMPCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.monitoredEntity
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
	payload.Data = util.Ptr(model.MeasurementDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.MeasurementListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *MaMPCSuite) Test_Failures() {
	s.sut.deviceConnected(s.mockRemoteEntity)

	s.sut.deviceMeasurementDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *MaMPCSuite) Test_deviceMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
	}
	s.sut.deviceMeasurementDataUpdate(payload)
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
				MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(3)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(4)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(5)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(6)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(6)),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.deviceMeasurementDataUpdate(payload)
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
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(6)),
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

	s.sut.deviceMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *MaMPCSuite) Test_deviceMeasurementDataUpdate_EventCallbacks() {
	// Test missing event callback coverage lines 103,105 (PowerPerPhase), 129,131 (CurrentPerPhase), 138,140 (VoltagePerPhase)
	
	// First setup measurement descriptions to enable the public methods to succeed
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
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
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

	// Add valid measurement data to make public methods succeed
	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(1000),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(230),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Setup electrical connection for phase-specific data
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test PowerPerPhase event callback (line 103,105)
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
		Data:   measData,
	}
	
	s.eventCalled = false
	s.sut.deviceMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)

	// Test CurrentPerPhase event callback (line 129,131)
	currentData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(15),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		},
	}
	
	payload.Data = currentData
	s.eventCalled = false
	s.sut.deviceMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)

	// Test VoltagePerPhase event callback (line 138,140)
	voltageData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(240),
				ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
			},
		},
	}
	
	payload.Data = voltageData
	s.eventCalled = false
	s.sut.deviceMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}
