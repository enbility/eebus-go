package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestLoadControlEventListDataType_Update(t *testing.T) {
	sut := model.LoadControlEventListDataType{
		LoadControlEventData: []model.LoadControlEventDataType{
			{
				EventId:            util.Ptr(model.LoadControlEventIdType(0)),
				EventActionConsume: util.Ptr(model.LoadControlEventActionTypeNormal),
			},
			{
				EventId:            util.Ptr(model.LoadControlEventIdType(1)),
				EventActionConsume: util.Ptr(model.LoadControlEventActionTypeNormal),
			},
		},
	}

	newData := model.LoadControlEventListDataType{
		LoadControlEventData: []model.LoadControlEventDataType{
			{
				EventId:            util.Ptr(model.LoadControlEventIdType(1)),
				EventActionConsume: util.Ptr(model.LoadControlEventActionTypeIncrease),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.LoadControlEventData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.EventId))
	assert.Equal(t, model.LoadControlEventActionTypeNormal, *item1.EventActionConsume)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.EventId))
	assert.Equal(t, model.LoadControlEventActionTypeIncrease, *item2.EventActionConsume)
}

func TestLoadControlStateListDataType_Update(t *testing.T) {
	sut := model.LoadControlStateListDataType{
		LoadControlStateData: []model.LoadControlStateDataType{
			{
				EventId:           util.Ptr(model.LoadControlEventIdType(0)),
				EventStateConsume: util.Ptr(model.LoadControlEventStateTypeEventAccepted),
			},
			{
				EventId:           util.Ptr(model.LoadControlEventIdType(1)),
				EventStateConsume: util.Ptr(model.LoadControlEventStateTypeEventAccepted),
			},
		},
	}

	newData := model.LoadControlStateListDataType{
		LoadControlStateData: []model.LoadControlStateDataType{
			{
				EventId:           util.Ptr(model.LoadControlEventIdType(1)),
				EventStateConsume: util.Ptr(model.LoadControlEventStateTypeEventStopped),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.LoadControlStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.EventId))
	assert.Equal(t, model.LoadControlEventStateTypeEventAccepted, *item1.EventStateConsume)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.EventId))
	assert.Equal(t, model.LoadControlEventStateTypeEventStopped, *item2.EventStateConsume)
}

func TestLoadControlLimitListDataType_Update(t *testing.T) {
	sut := model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(false),
			},
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(false),
			},
		},
	}

	newData := model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.LoadControlLimitData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.Equal(t, false, *item1.IsLimitChangeable)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.Equal(t, true, *item2.IsLimitChangeable)
}

func TestLoadControlLimitConstraintsListDataType_Update(t *testing.T) {
	sut := model.LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []model.LoadControlLimitConstraintsDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				ValueStepSize: model.NewScaledNumberType(1),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				ValueStepSize: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []model.LoadControlLimitConstraintsDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				ValueStepSize: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.LoadControlLimitConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.Equal(t, 1.0, float64(item1.ValueStepSize.GetValue()))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.Equal(t, 10.0, float64(item2.ValueStepSize.GetValue()))
}

func TestLoadControlLimitDescriptionListDataType_Update(t *testing.T) {
	sut := model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
			},
		},
	}

	newData := model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(model.LoadControlCategoryTypeOptimization),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.LoadControlLimitDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.Equal(t, model.LoadControlCategoryTypeObligation, *item1.LimitCategory)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.Equal(t, model.LoadControlCategoryTypeOptimization, *item2.LimitCategory)
}
