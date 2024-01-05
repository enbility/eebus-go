package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeTableListDataType_Update(t *testing.T) {
	sut := TimeTableListDataType{
		TimeTableData: []TimeTableDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(0)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(1)),
				},
			},
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(1)),
				},
			},
		},
	}

	newData := TimeTableListDataType{
		TimeTableData: []TimeTableDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(10)),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []TimeTableConstraintsDataType{
			{
				TimeTableId:  util.Ptr(TimeTableIdType(0)),
				SlotCountMin: util.Ptr(TimeSlotCountType(1)),
			},
			{
				TimeTableId:  util.Ptr(TimeTableIdType(1)),
				SlotCountMin: util.Ptr(TimeSlotCountType(1)),
			},
		},
	}

	newData := TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []TimeTableConstraintsDataType{
			{
				TimeTableId:  util.Ptr(TimeTableIdType(1)),
				SlotCountMin: util.Ptr(TimeSlotCountType(10)),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []TimeTableDescriptionDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []TimeTableDescriptionDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
