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

	// test scenario gate rejects data changes without scenarios
	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = &model.MeasurementDescriptionListDataType{}
	s.sut.HandleEvent(payload)

	s.setUpUseCaseScenarios()

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
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACEnergyProduced),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(5)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(6)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACFrequency),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.deviceMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.MeasurementListDataType{
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
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(4)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(5)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(6)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	payload.Data = data

	s.sut.deviceMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *MaMPCSuite) setUpUseCaseScenarios() {
	address := &model.FeatureAddressType{
		Device:  s.monitoredEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeMonitoredUnit,
		model.UseCaseNameTypeMonitoringOfPowerConsumption,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.monitoredEntity,
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.sut.UseCaseBase.HandleEvent(payload)
}

func (s *MaMPCSuite) Test_deviceConnected() {
	s.sut.deviceConnected(s.monitoredEntity)
}
