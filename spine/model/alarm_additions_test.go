package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestAlarmListDataType_Update(t *testing.T) {
	sut := AlarmListDataType{
		AlarmListData: []AlarmDataType{
			{
				AlarmId:     util.Ptr(AlarmIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				AlarmId:     util.Ptr(AlarmIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := AlarmListDataType{
		AlarmListData: []AlarmDataType{
			{
				AlarmId:     util.Ptr(AlarmIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.AlarmListData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.AlarmId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.AlarmId))
	assert.Equal(t, "new", string(*item2.Description))
}
