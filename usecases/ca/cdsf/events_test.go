package cdsf

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CaCDSFSuite) Test_Events() {
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

	payload.Data = util.Ptr(model.HvacSystemFunctionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.HvacOverrunListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *CaCDSFSuite) Test_hvacOverrunDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Entity: s.dhwCircuitEntity,
	}
	s.sut.hvacOverrunDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.HvacOverrunListDataType{})
	s.sut.hvacOverrunDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.HvacOverrunListDataType{
		HvacOverrunData: []model.HvacOverrunDataType{
			{
				OverrunId:     util.Ptr(model.HvacOverrunIdType(1)),
				OverrunStatus: util.Ptr(model.HvacOverrunStatusTypeActive),
			},
		},
	})
	s.sut.hvacOverrunDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *CaCDSFSuite) Test_Failures() {
	s.sut.dhwCircuitConnected(s.mockRemoteEntity)

	s.sut.hvacSystemFunctionDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *CaCDSFSuite) Test_hvacSystemFunctionDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Entity: s.dhwCircuitEntity,
	}
	s.sut.hvacSystemFunctionDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.HvacSystemFunctionListDataType{})
	s.sut.hvacSystemFunctionDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Data = util.Ptr(model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(1)),
			},
		},
	})
	s.sut.hvacSystemFunctionDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}
