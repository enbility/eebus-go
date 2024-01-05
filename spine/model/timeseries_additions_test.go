package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeSeriesListDataType_Update(t *testing.T) {
	sut := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(0)),
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
					},
				},
			},
		},
	}

	newData := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
					},
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("P6D")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
						TimePeriod: &TimePeriodType{
							StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
							EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("P6D")),
						},
						MaxValue: NewScaledNumberType(10000),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(2)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(DurationType("P1DT6H46M33S")),
						MaxValue:         NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(DurationType("PT7H37M53S")),
						MaxValue:         NewScaledNumberType(4410),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(2)),
						Duration:         util.Ptr(DurationType("PT38M")),
						MaxValue:         NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(3)),
						Duration:         util.Ptr(DurationType("PT32M")),
						MaxValue:         NewScaledNumberType(4410),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(4)),
						Duration:         util.Ptr(DurationType("P1D")),
						MaxValue:         NewScaledNumberType(0),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(3)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(DurationType("P1DT15H24M57S")),
						Value:            NewScaledNumberType(44229),
						MaxValue:         NewScaledNumberType(49629),
					},
				},
			},
		},
	}

	newData := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(3)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(DurationType("P1DT15H16M50S")),
						Value:            NewScaledNumberType(11539),
						MaxValue:         NewScaledNumberType(49629),
					},
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []TimeSeriesDescriptionDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(0)),
				Description:  util.Ptr(DescriptionType("old")),
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				Description:  util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []TimeSeriesDescriptionDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				Description:  util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(0)),
				SlotValueMin: NewScaledNumberType(1),
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				SlotValueMin: NewScaledNumberType(1),
			},
		},
	}

	newData := TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				SlotValueMin: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
