package vabd

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemVABDSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.batteryEntity
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

func (s *CemVABDSuite) Test_Failures() {
	s.sut.inverterConnected(s.mockRemoteEntity)

	s.sut.inverterMeasurementDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *CemVABDSuite) Test_inverterMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.batteryEntity,
	}
	s.sut.inverterMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeCharge),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeDischarge),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(3)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.batteryEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.inverterMeasurementDataUpdate(payload)
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
		},
	}

	payload.Data = data

	s.sut.inverterMeasurementDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *CemVABDSuite) setUpUseCaseScenarios() {
	address := &model.FeatureAddressType{
		Device:  s.batteryEntity.Device().Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	nodeFeature := s.remoteDevice.FeatureByAddress(address)

	data := &model.NodeManagementUseCaseDataType{}
	data.AddUseCaseSupport(
		model.FeatureAddressType{},
		model.UseCaseActorTypeBatterySystem,
		model.UseCaseNameTypeVisualizationOfAggregatedBatteryData,
		"1.0.0",
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
	_, _ = nodeFeature.UpdateData(true, model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	payload := spineapi.EventPayload{
		Device:     s.remoteDevice,
		Entity:     s.batteryEntity,
		EventType:  spineapi.EventTypeDataChange,
		ChangeType: spineapi.ElementChangeUpdate,
		Data:       &model.NodeManagementUseCaseDataType{},
	}
	s.sut.UseCaseBase.HandleEvent(payload)
}

func (s *CemVABDSuite) Test_inverterConnected() {
	s.sut.inverterConnected(s.batteryEntity)
}
