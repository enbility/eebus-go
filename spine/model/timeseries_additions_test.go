package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeSeriesListDataType_Update(t *testing.T) {
	sut := model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
					},
				},
			},
		},
	}

	newData := model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
					},
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeSeriesData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeSeriesId))
	assert.Equal(t, 0, int(*item1.TimeSeriesSlot[0].TimeSeriesSlotId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeSeriesId))
	assert.Equal(t, 1, int(*item2.TimeSeriesSlot[0].TimeSeriesSlotId))
}

func TestTimeSeriesDescriptionListDataType_Update(t *testing.T) {
	sut := model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				Description:  util.Ptr(model.DescriptionType("old")),
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				Description:  util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				Description:  util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeSeriesDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeSeriesId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeSeriesId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTimeSeriesConstraintsListDataType_Update(t *testing.T) {
	sut := model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				SlotValueMin: model.NewScaledNumberType(1),
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				SlotValueMin: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				SlotValueMin: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeSeriesConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeSeriesId))
	assert.Equal(t, 1.0, item1.SlotValueMin.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeSeriesId))
	assert.Equal(t, 10.0, item2.SlotValueMin.GetValue())
}
