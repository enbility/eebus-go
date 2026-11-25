package cdt

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CaCDTSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.dhwCircuitEntity
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
	payload.Data = util.Ptr(model.HvacSystemFunctionDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.SetpointDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.SetpointListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.SetpointConstraintsListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *CaCDTSuite) Test_Failures() {
	s.sut.dhwCircuitConnected(s.mockRemoteEntity)

	s.sut.hvacSystemFunctionDescriptionDataUpdate(s.mockRemoteEntity)

	s.sut.setpointDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *CaCDTSuite) Test_setpointDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Entity: s.dhwCircuitEntity,
	}
	s.sut.setpointDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.SetpointListDataType{})
	s.sut.setpointDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{
			{
				SetpointId: util.Ptr(model.SetpointIdType(1)),
				Value:      model.NewScaledNumberType(21),
			},
		},
	})
	s.sut.setpointDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *CaCDTSuite) Test_setpointConstraintsDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Entity: s.dhwCircuitEntity,
	}
	s.sut.setpointConstraintsDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.SetpointConstraintsListDataType{})
	s.sut.setpointConstraintsDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.SetpointConstraintsListDataType{
		SetpointConstraintsData: []model.SetpointConstraintsDataType{
			{
				SetpointId:       util.Ptr(model.SetpointIdType(1)),
				SetpointRangeMin: model.NewScaledNumberType(16),
			},
		},
	})
	s.sut.setpointConstraintsDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}
