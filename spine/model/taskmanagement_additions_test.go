package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTaskManagementJobListDataType_Update(t *testing.T) {
	sut := model.TaskManagementJobListDataType{
		TaskManagementJobData: []model.TaskManagementJobDataType{
			{
				JobId:    util.Ptr(model.TaskManagementJobIdType(0)),
				JobState: util.Ptr(model.TaskManagementJobStateTypeActive),
			},
			{
				JobId:    util.Ptr(model.TaskManagementJobIdType(1)),
				JobState: util.Ptr(model.TaskManagementJobStateTypeActive),
			},
		},
	}

	newData := model.TaskManagementJobListDataType{
		TaskManagementJobData: []model.TaskManagementJobDataType{
			{
				JobId:    util.Ptr(model.TaskManagementJobIdType(1)),
				JobState: util.Ptr(model.TaskManagementJobStateTypeCompleted),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TaskManagementJobData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, model.TaskManagementJobStateTypeActive, *item1.JobState)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, model.TaskManagementJobStateTypeCompleted, *item2.JobState)
}

func TestTaskManagementJobRelationListDataType_Update(t *testing.T) {
	sut := model.TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []model.TaskManagementJobRelationDataType{
			{
				JobId: util.Ptr(model.TaskManagementJobIdType(0)),
				LoadControlReleated: &model.TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(model.LoadControlEventIdType(0)),
				},
			},
			{
				JobId: util.Ptr(model.TaskManagementJobIdType(1)),
				LoadControlReleated: &model.TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(model.LoadControlEventIdType(0)),
				},
			},
		},
	}

	newData := model.TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []model.TaskManagementJobRelationDataType{
			{
				JobId: util.Ptr(model.TaskManagementJobIdType(1)),
				LoadControlReleated: &model.TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(model.LoadControlEventIdType(1)),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TaskManagementJobRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, 0, int(*item1.LoadControlReleated.EventId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, 1, int(*item2.LoadControlReleated.EventId))
}

func TestTaskManagementJobDescriptionListDataType_Update(t *testing.T) {
	sut := model.TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []model.TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(model.TaskManagementJobIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				JobId:       util.Ptr(model.TaskManagementJobIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []model.TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(model.TaskManagementJobIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TaskManagementJobDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, "new", string(*item2.Description))
}
