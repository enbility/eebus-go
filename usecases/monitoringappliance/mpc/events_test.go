package mpc

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MPCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.monitoredEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
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
}

func (s *MPCSuite) Test_Failures() {
	s.sut.deviceConnected(s.mockRemoteEntity)

	s.sut.deviceMeasurementDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *MPCSuite) Test_deviceMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
	}
	s.sut.deviceMeasurementDataUpdate(payload)

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
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.deviceMeasurementDataUpdate(payload)

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
}
