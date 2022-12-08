package model_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestOperatingConstraintsInterruptListDataType_Update(t *testing.T) {
	sut := model.OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []model.OperatingConstraintsInterruptDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				IsPausable: util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				IsPausable: util.Ptr(false),
			},
		},
	}

	newData := model.OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []model.OperatingConstraintsInterruptDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				IsPausable: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []model.OperatingConstraintsDurationDataType{
			{
				SequenceId:        util.Ptr(model.PowerSequenceIdType(0)),
				ActiveDurationMin: model.NewDurationType(1 * time.Second),
			},
			{
				SequenceId:        util.Ptr(model.PowerSequenceIdType(1)),
				ActiveDurationMin: model.NewDurationType(1 * time.Second),
			},
		},
	}

	newData := model.OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []model.OperatingConstraintsDurationDataType{
			{
				SequenceId:        util.Ptr(model.PowerSequenceIdType(1)),
				ActiveDurationMin: model.NewDurationType(10 * time.Second),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []model.OperatingConstraintsPowerDescriptionDataType{
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

	newData := model.OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []model.OperatingConstraintsPowerDescriptionDataType{
			{
				SequenceId:              util.Ptr(model.PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeProduce),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.OperatingConstraintsPowerDescriptionData
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

func TestOperatingConstraintsPowerRangeListDataType_Update(t *testing.T) {
	sut := model.OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []model.OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				PowerMin:   model.NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				PowerMin:   model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []model.OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				PowerMin:   model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []model.OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(0)),
				Power:      model.NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Power:      model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []model.OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				Power:      model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []model.OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(model.PowerSequenceIdType(0)),
				ResumeEnergyEstimated: model.NewScaledNumberType(1),
			},
			{
				SequenceId:            util.Ptr(model.PowerSequenceIdType(1)),
				ResumeEnergyEstimated: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []model.OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(model.PowerSequenceIdType(1)),
				ResumeEnergyEstimated: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
