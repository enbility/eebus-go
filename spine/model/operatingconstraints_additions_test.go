package model

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestOperatingConstraintsInterruptListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []OperatingConstraintsInterruptDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				IsPausable: util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				IsPausable: util.Ptr(false),
			},
		},
	}

	newData := OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []OperatingConstraintsInterruptDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				IsPausable: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsInterruptData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, false, *item1.IsPausable)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, true, *item2.IsPausable)
}

func TestOperatingConstraintsDurationListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []OperatingConstraintsDurationDataType{
			{
				SequenceId:        util.Ptr(PowerSequenceIdType(0)),
				ActiveDurationMin: NewDurationType(1 * time.Second),
			},
			{
				SequenceId:        util.Ptr(PowerSequenceIdType(1)),
				ActiveDurationMin: NewDurationType(1 * time.Second),
			},
		},
	}

	newData := OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []OperatingConstraintsDurationDataType{
			{
				SequenceId:        util.Ptr(PowerSequenceIdType(1)),
				ActiveDurationMin: NewDurationType(10 * time.Second),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsDurationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	duration, _ := item1.ActiveDurationMin.GetTimeDuration()
	assert.Equal(t, time.Duration(1*time.Second), duration)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	duration, _ = item2.ActiveDurationMin.GetTimeDuration()
	assert.Equal(t, time.Duration(10*time.Second), duration)
}

func TestOperatingConstraintsPowerDescriptionListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []OperatingConstraintsPowerDescriptionDataType{
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

	newData := OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []OperatingConstraintsPowerDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsPowerDescriptionData
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

func TestOperatingConstraintsPowerRangeListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				PowerMin:   NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				PowerMin:   NewScaledNumberType(1),
			},
		},
	}

	newData := OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				PowerMin:   NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsPowerRangeData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.PowerMin.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.PowerMin.GetValue())
}

func TestOperatingConstraintsPowerLevelListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Power:      NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Power:      NewScaledNumberType(1),
			},
		},
	}

	newData := OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Power:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsPowerLevelData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.Power.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.Power.GetValue())
}

func TestOperatingConstraintsResumeImplicationListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(0)),
				ResumeEnergyEstimated: NewScaledNumberType(1),
			},
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(1)),
				ResumeEnergyEstimated: NewScaledNumberType(1),
			},
		},
	}

	newData := OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(1)),
				ResumeEnergyEstimated: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsResumeImplicationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.ResumeEnergyEstimated.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.ResumeEnergyEstimated.GetValue())
}
