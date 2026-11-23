package evsoc

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemEVSOCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.evEntity
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
	payload.Data = &model.MeasurementListDataType{}
	s.sut.HandleEvent(payload)

	s.setUpUseCaseScenarios()

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.MeasurementListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *CemEVSOCSuite) Test_Failures() {
	s.sut.evConnected(s.mockRemoteEntity)
}

func (s *CemEVSOCSuite) Test_evMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.eventCalled = false
	s.sut.evMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeStateOfHealth),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeTravelRange),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.eventCalled = false

	s.sut.evMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(200),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(3000),
			},
		},
	}

	payload.Data = data
	s.eventCalled = false

	s.sut.evMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *CemEVSOCSuite) Test_evConnected() {
	s.sut.evConnected(s.evEntity)
}

func (s *CemEVSOCSuite) setUpUseCaseScenarios() {
	address := &model.FeatureAddressType{
		Device:  s.evEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeEV,
		model.UseCaseNameTypeEVStateOfCharge,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.evEntity,
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.sut.UseCaseBase.HandleEvent(payload)
}
