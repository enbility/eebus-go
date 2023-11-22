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

func TestTimeSeriesListDataType_Update_02(t *testing.T) {
	sut := model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimePeriod: &model.TimePeriodType{
					StartTime: util.Ptr(model.AbsoluteOrRelativeTimeType("PT0S")),
					EndTime:   util.Ptr(model.AbsoluteOrRelativeTimeType("P6D")),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						TimePeriod: &model.TimePeriodType{
							StartTime: util.Ptr(model.AbsoluteOrRelativeTimeType("PT0S")),
							EndTime:   util.Ptr(model.AbsoluteOrRelativeTimeType("P6D")),
						},
						MaxValue: model.NewScaledNumberType(10000),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: util.Ptr(model.AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("P1DT6H46M33S")),
						MaxValue:         model.NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("PT7H37M53S")),
						MaxValue:         model.NewScaledNumberType(4410),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(2)),
						Duration:         util.Ptr(model.DurationType("PT38M")),
						MaxValue:         model.NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(3)),
						Duration:         util.Ptr(model.DurationType("PT32M")),
						MaxValue:         model.NewScaledNumberType(4410),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(4)),
						Duration:         util.Ptr(model.DurationType("P1D")),
						MaxValue:         model.NewScaledNumberType(0),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
				TimePeriod: &model.TimePeriodType{
					StartTime: util.Ptr(model.AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P1DT15H24M57S")),
						Value:            model.NewScaledNumberType(44229),
						MaxValue:         model.NewScaledNumberType(49629),
					},
				},
			},
		},
	}

	newData := model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
				TimePeriod: &model.TimePeriodType{
					StartTime: util.Ptr(model.AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P1DT15H16M50S")),
						Value:            model.NewScaledNumberType(11539),
						MaxValue:         model.NewScaledNumberType(49629),
					},
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TimeSeriesData
	// check the non changing items
	assert.Equal(t, 3, len(data))
	item1 := data[0]
	assert.Equal(t, 1, int(*item1.TimeSeriesId))
	assert.Equal(t, 0, int(*item1.TimeSeriesSlot[0].TimeSeriesSlotId))
	// check properties of updated item
	item2 := data[2]
	assert.Equal(t, 3, int(*item2.TimeSeriesId))
	assert.Equal(t, 1, int(*item2.TimeSeriesSlot[0].TimeSeriesSlotId))
	assert.Equal(t, 11539, int(item2.TimeSeriesSlot[0].Value.GetValue()))
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
