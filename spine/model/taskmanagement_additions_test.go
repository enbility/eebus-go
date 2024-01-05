package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTaskManagementJobListDataType_Update(t *testing.T) {
	sut := TaskManagementJobListDataType{
		TaskManagementJobData: []TaskManagementJobDataType{
			{
				JobId:    util.Ptr(TaskManagementJobIdType(0)),
				JobState: util.Ptr(TaskManagementJobStateTypeActive),
			},
			{
				JobId:    util.Ptr(TaskManagementJobIdType(1)),
				JobState: util.Ptr(TaskManagementJobStateTypeActive),
			},
		},
	}

	newData := TaskManagementJobListDataType{
		TaskManagementJobData: []TaskManagementJobDataType{
			{
				JobId:    util.Ptr(TaskManagementJobIdType(1)),
				JobState: util.Ptr(TaskManagementJobStateTypeCompleted),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.TaskManagementJobData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, TaskManagementJobStateTypeActive, *item1.JobState)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, TaskManagementJobStateTypeCompleted, *item2.JobState)
}

func TestTaskManagementJobRelationListDataType_Update(t *testing.T) {
	sut := TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []TaskManagementJobRelationDataType{
			{
				JobId: util.Ptr(TaskManagementJobIdType(0)),
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(0)),
				},
			},
			{
				JobId: util.Ptr(TaskManagementJobIdType(1)),
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(0)),
				},
			},
		},
	}

	newData := TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []TaskManagementJobRelationDataType{
			{
				JobId: util.Ptr(TaskManagementJobIdType(1)),
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(1)),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(TaskManagementJobIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				JobId:       util.Ptr(TaskManagementJobIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(TaskManagementJobIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
