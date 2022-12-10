package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeTableListDataType_Update(t *testing.T) {
	sut := model.TimeTableListDataType{
		TimeTableData: []model.TimeTableDataType{
			{
				TimeTableId: util.Ptr(model.TimeTableIdType(0)),
				RecurrenceInformation: &model.RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(1)),
				},
			},
			{
				TimeTableId: util.Ptr(model.TimeTableIdType(1)),
				RecurrenceInformation: &model.RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(1)),
				},
			},
		},
	}

	newData := model.TimeTableListDataType{
		TimeTableData: []model.TimeTableDataType{
			{
				TimeTableId: util.Ptr(model.TimeTableIdType(1)),
				RecurrenceInformation: &model.RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(10)),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeTableData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeTableId))
	assert.Equal(t, 1, int(*item1.RecurrenceInformation.ExecutionCount))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeTableId))
	assert.Equal(t, 10, int(*item2.RecurrenceInformation.ExecutionCount))
}

func TestTimeTableConstraintsListDataType_Update(t *testing.T) {
	sut := model.TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []model.TimeTableConstraintsDataType{
			{
				TimeTableId:  util.Ptr(model.TimeTableIdType(0)),
				SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
			},
			{
				TimeTableId:  util.Ptr(model.TimeTableIdType(1)),
				SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
			},
		},
	}

	newData := model.TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []model.TimeTableConstraintsDataType{
			{
				TimeTableId:  util.Ptr(model.TimeTableIdType(1)),
				SlotCountMin: util.Ptr(model.TimeSlotCountType(10)),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeTableConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeTableId))
	assert.Equal(t, 1, int(*item1.SlotCountMin))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeTableId))
	assert.Equal(t, 10, int(*item2.SlotCountMin))
}

func TestTimeTableDescriptionListDataType_Update(t *testing.T) {
	sut := model.TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []model.TimeTableDescriptionDataType{
			{
				TimeTableId: util.Ptr(model.TimeTableIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				TimeTableId: util.Ptr(model.TimeTableIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []model.TimeTableDescriptionDataType{
			{
				TimeTableId: util.Ptr(model.TimeTableIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeTableDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeTableId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeTableId))
	assert.Equal(t, "new", string(*item2.Description))
}
