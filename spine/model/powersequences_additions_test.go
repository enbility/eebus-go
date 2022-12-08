package model_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestPowerTimeSlotScheduleListDataType_Update(t *testing.T) {
	sut := model.PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []model.PowerTimeSlotScheduleDataType{
			{
				SequenceId:    util.Ptr(model.PowerSequenceIdType(0)),
				SlotActivated: util.Ptr(false),
			},
			{
				SequenceId:    util.Ptr(model.PowerSequenceIdType(1)),
				SlotActivated: util.Ptr(false),
			},
		},
	}

	newData := model.PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []model.PowerTimeSlotScheduleDataType{
			{
				SequenceId:    util.Ptr(model.PowerSequenceIdType(1)),
				SlotActivated: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []model.PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				Value:      model.NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Value:      model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []model.PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Value:      model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []model.PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:  util.Ptr(model.PowerSequenceIdType(0)),
				MinDuration: model.NewDurationType(1 * time.Second),
			},
			{
				SequenceId:  util.Ptr(model.PowerSequenceIdType(1)),
				MinDuration: model.NewDurationType(1 * time.Second),
			},
		},
	}

	newData := model.PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []model.PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:  util.Ptr(model.PowerSequenceIdType(1)),
				MinDuration: model.NewDurationType(10 * time.Second),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []model.PowerSequenceAlternativesRelationDataType{
			{
				AlternativeId: util.Ptr(model.AlternativesIdType(0)),
				SequenceId:    []model.PowerSequenceIdType{0},
			},
			{
				AlternativeId: util.Ptr(model.AlternativesIdType(1)),
				SequenceId:    []model.PowerSequenceIdType{0},
			},
		},
	}

	newData := model.PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []model.PowerSequenceAlternativesRelationDataType{
			{
				AlternativeId: util.Ptr(model.AlternativesIdType(1)),
				SequenceId:    []model.PowerSequenceIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []model.PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(model.PowerSequenceIdType(0)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
			{
				SequenceId:              util.Ptr(model.PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	newData := model.PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []model.PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(model.PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeProduce),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.PowerSequenceDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, model.EnergyDirectionTypeConsume, *item1.PositiveEnergyDirection)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, model.EnergyDirectionTypeProduce, *item2.PositiveEnergyDirection)
}

func TestPowerSequenceStateListDataType_Update(t *testing.T) {
	sut := model.PowerSequenceStateListDataType{
		PowerSequenceStateData: []model.PowerSequenceStateDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				State:      util.Ptr(model.PowerSequenceStateTypeRunning),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				State:      util.Ptr(model.PowerSequenceStateTypeRunning),
			},
		},
	}

	newData := model.PowerSequenceStateListDataType{
		PowerSequenceStateData: []model.PowerSequenceStateDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				State:      util.Ptr(model.PowerSequenceStateTypeCompleted),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.PowerSequenceStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, model.PowerSequenceStateTypeRunning, *item1.State)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, model.PowerSequenceStateTypeCompleted, *item2.State)
}

func TestPowerSequenceScheduleListDataType_Update(t *testing.T) {
	sut := model.PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []model.PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				EndTime:    model.NewAbsoluteOrRelativeTimeType("PT2H"),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				EndTime:    model.NewAbsoluteOrRelativeTimeType("PT2H"),
			},
		},
	}

	newData := model.PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []model.PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				EndTime:    model.NewAbsoluteOrRelativeTimeType("PT4H"),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []model.PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:      util.Ptr(model.PowerSequenceIdType(0)),
				EarliestEndTime: model.NewAbsoluteOrRelativeTimeType("PT2H"),
			},
			{
				SequenceId:      util.Ptr(model.PowerSequenceIdType(1)),
				EarliestEndTime: model.NewAbsoluteOrRelativeTimeType("PT2H"),
			},
		},
	}

	newData := model.PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []model.PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:      util.Ptr(model.PowerSequenceIdType(1)),
				EarliestEndTime: model.NewAbsoluteOrRelativeTimeType("PT4H"),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerSequencePriceListDataType{
		PowerSequencePriceData: []model.PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				Price:      model.NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Price:      model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.PowerSequencePriceListDataType{
		PowerSequencePriceData: []model.PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Price:      model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []model.PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				Cheapest:   util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Cheapest:   util.Ptr(false),
			},
		},
	}

	newData := model.PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []model.PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Cheapest:   util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
