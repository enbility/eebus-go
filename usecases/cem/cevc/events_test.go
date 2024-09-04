package cevc

import (
	"time"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_Events() {
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

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.TimeSeriesDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.TimeSeriesListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.IncentiveTableDescriptionDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.IncentiveTableConstraintsDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.IncentiveDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *CemCEVCSuite) Test_Failures() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.evConnected(s.mockRemoteEntity)

	s.sut.evTimeSeriesDescriptionDataUpdate(payload)

	s.sut.evTimeSeriesDataUpdate(payload)

	s.sut.evIncentiveTableDescriptionDataUpdate(payload)

	s.sut.evCheckTimeSeriesDescriptionConstraintsUpdateRequired(s.mockRemoteEntity)

	s.sut.evCheckIncentiveTableDescriptionUpdateRequired(s.mockRemoteEntity)
}

func (s *CemCEVCSuite) Test_evTimeSeriesDescriptionDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evTimeSeriesDescriptionDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.sut.evTimeSeriesDescriptionDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	timeDesc := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
				UpdateRequired: util.Ptr(true),
			},
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
			},
		},
	}

	rTimeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr := rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, timeDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evTimeSeriesDescriptionDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						MinValue:         model.NewScaledNumberType(1000),
						Value:            model.NewScaledNumberType(10000),
						MaxValue:         model.NewScaledNumberType(100000),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err := s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1000.0, demand.MinDemand)
	assert.Equal(s.T(), 10000.0, demand.OptDemand)
	assert.Equal(s.T(), 100000.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	s.sut.evTimeSeriesDescriptionDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
	s.eventCalled = false

	constData := &model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId:         util.Ptr(model.TimeSeriesIdType(0)),
				SlotCountMin:         util.Ptr(model.TimeSeriesSlotCountType(1)),
				SlotCountMax:         util.Ptr(model.TimeSeriesSlotCountType(10)),
				SlotDurationMin:      model.NewDurationType(1 * time.Minute),
				SlotDurationMax:      model.NewDurationType(60 * time.Minute),
				SlotDurationStepSize: model.NewDurationType(1 * time.Minute),
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesConstraintsListData, constData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evTimeSeriesDescriptionDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
	s.eventCalled = false

	incConstData := &model.IncentiveTableConstraintsDataType{
		IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
			{
				IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
					SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
					SlotCountMax: util.Ptr(model.TimeSlotCountType(10)),
				},
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeIncentiveTableConstraintsData, incConstData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evTimeSeriesDescriptionDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}
