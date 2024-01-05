package model

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestPowerTimeSlotScheduleListDataType_Update(t *testing.T) {
	sut := PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []PowerTimeSlotScheduleDataType{
			{
				SequenceId:    util.Ptr(PowerSequenceIdType(0)),
				SlotActivated: util.Ptr(false),
			},
			{
				SequenceId:    util.Ptr(PowerSequenceIdType(1)),
				SlotActivated: util.Ptr(false),
			},
		},
	}

	newData := PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []PowerTimeSlotScheduleDataType{
			{
				SequenceId:    util.Ptr(PowerSequenceIdType(1)),
				SlotActivated: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerTimeSlotScheduleData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, false, *item1.SlotActivated)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, true, *item2.SlotActivated)
}

func TestPowerTimeSlotValueListDataType_Update(t *testing.T) {
	sut := PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Value:      NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Value:      NewScaledNumberType(1),
			},
		},
	}

	newData := PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Value:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerTimeSlotValueData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestPowerTimeSlotScheduleConstraintsListDataType_Update(t *testing.T) {
	sut := PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:  util.Ptr(PowerSequenceIdType(0)),
				MinDuration: NewDurationType(1 * time.Second),
			},
			{
				SequenceId:  util.Ptr(PowerSequenceIdType(1)),
				MinDuration: NewDurationType(1 * time.Second),
			},
		},
	}

	newData := PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:  util.Ptr(PowerSequenceIdType(1)),
				MinDuration: NewDurationType(10 * time.Second),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerTimeSlotScheduleConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	duration, _ := item1.MinDuration.GetTimeDuration()
	assert.Equal(t, time.Duration(1*time.Second), duration)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	duration, _ = item2.MinDuration.GetTimeDuration()
	assert.Equal(t, time.Duration(10*time.Second), duration)
}

func TestPowerSequenceAlternativesRelationListDataType_Update(t *testing.T) {
	sut := PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []PowerSequenceAlternativesRelationDataType{
			{
				AlternativeId: util.Ptr(AlternativesIdType(0)),
				SequenceId:    []PowerSequenceIdType{0},
			},
			{
				AlternativeId: util.Ptr(AlternativesIdType(1)),
				SequenceId:    []PowerSequenceIdType{0},
			},
		},
	}

	newData := PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []PowerSequenceAlternativesRelationDataType{
			{
				AlternativeId: util.Ptr(AlternativesIdType(1)),
				SequenceId:    []PowerSequenceIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequenceAlternativesRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.AlternativeId))
	assert.Equal(t, 0, int(item1.SequenceId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.AlternativeId))
	assert.Equal(t, 1, int(item2.SequenceId[0]))
}

func TestPowerSequenceDescriptionListDataType_Update(t *testing.T) {
	sut := PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(0)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
			},
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
			},
		},
	}

	newData := PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequenceDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, EnergyDirectionTypeConsume, *item1.PositiveEnergyDirection)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, EnergyDirectionTypeProduce, *item2.PositiveEnergyDirection)
}

func TestPowerSequenceStateListDataType_Update(t *testing.T) {
	sut := PowerSequenceStateListDataType{
		PowerSequenceStateData: []PowerSequenceStateDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				State:      util.Ptr(PowerSequenceStateTypeRunning),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				State:      util.Ptr(PowerSequenceStateTypeRunning),
			},
		},
	}

	newData := PowerSequenceStateListDataType{
		PowerSequenceStateData: []PowerSequenceStateDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				State:      util.Ptr(PowerSequenceStateTypeCompleted),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequenceStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, PowerSequenceStateTypeRunning, *item1.State)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, PowerSequenceStateTypeCompleted, *item2.State)
}

func TestPowerSequenceScheduleListDataType_Update(t *testing.T) {
	sut := PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				EndTime:    NewAbsoluteOrRelativeTimeType("PT2H"),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				EndTime:    NewAbsoluteOrRelativeTimeType("PT2H"),
			},
		},
	}

	newData := PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				EndTime:    NewAbsoluteOrRelativeTimeType("PT4H"),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequenceScheduleData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, "PT2H", string(*item1.EndTime))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, "PT4H", string(*item2.EndTime))
}

func TestPowerSequenceScheduleConstraintsListDataType_Update(t *testing.T) {
	sut := PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(0)),
				EarliestEndTime: NewAbsoluteOrRelativeTimeType("PT2H"),
			},
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(1)),
				EarliestEndTime: NewAbsoluteOrRelativeTimeType("PT2H"),
			},
		},
	}

	newData := PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(1)),
				EarliestEndTime: NewAbsoluteOrRelativeTimeType("PT4H"),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequenceScheduleConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, "PT2H", string(*item1.EarliestEndTime))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, "PT4H", string(*item2.EarliestEndTime))
}

func TestPowerSequencePriceListDataType_Update(t *testing.T) {
	sut := PowerSequencePriceListDataType{
		PowerSequencePriceData: []PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Price:      NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Price:      NewScaledNumberType(1),
			},
		},
	}

	newData := PowerSequencePriceListDataType{
		PowerSequencePriceData: []PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Price:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequencePriceData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.Price.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.Price.GetValue())
}

func TestPowerSequenceSchedulePreferenceListDataType_Update(t *testing.T) {
	sut := PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Cheapest:   util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Cheapest:   util.Ptr(false),
			},
		},
	}

	newData := PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Cheapest:   util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.PowerSequenceSchedulePreferenceData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, false, *item1.Cheapest)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, true, *item2.Cheapest)
}
