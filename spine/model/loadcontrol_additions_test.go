package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestLoadControlEventListDataType_Update(t *testing.T) {
	sut := LoadControlEventListDataType{
		LoadControlEventData: []LoadControlEventDataType{
			{
				EventId:            util.Ptr(LoadControlEventIdType(0)),
				EventActionConsume: util.Ptr(LoadControlEventActionTypeNormal),
			},
			{
				EventId:            util.Ptr(LoadControlEventIdType(1)),
				EventActionConsume: util.Ptr(LoadControlEventActionTypeNormal),
			},
		},
	}

	newData := LoadControlEventListDataType{
		LoadControlEventData: []LoadControlEventDataType{
			{
				EventId:            util.Ptr(LoadControlEventIdType(1)),
				EventActionConsume: util.Ptr(LoadControlEventActionTypeIncrease),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.LoadControlEventData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.EventId))
	assert.Equal(t, LoadControlEventActionTypeNormal, *item1.EventActionConsume)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.EventId))
	assert.Equal(t, LoadControlEventActionTypeIncrease, *item2.EventActionConsume)
}

func TestLoadControlStateListDataType_Update(t *testing.T) {
	sut := LoadControlStateListDataType{
		LoadControlStateData: []LoadControlStateDataType{
			{
				EventId:           util.Ptr(LoadControlEventIdType(0)),
				EventStateConsume: util.Ptr(LoadControlEventStateTypeEventAccepted),
			},
			{
				EventId:           util.Ptr(LoadControlEventIdType(1)),
				EventStateConsume: util.Ptr(LoadControlEventStateTypeEventAccepted),
			},
		},
	}

	newData := LoadControlStateListDataType{
		LoadControlStateData: []LoadControlStateDataType{
			{
				EventId:           util.Ptr(LoadControlEventIdType(1)),
				EventStateConsume: util.Ptr(LoadControlEventStateTypeEventStopped),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.LoadControlStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.EventId))
	assert.Equal(t, LoadControlEventStateTypeEventAccepted, *item1.EventStateConsume)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.EventId))
	assert.Equal(t, LoadControlEventStateTypeEventStopped, *item2.EventStateConsume)
}

func TestLoadControlLimitListDataType_Update(t *testing.T) {
	sut := LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(false),
			},
			{
				LimitId:           util.Ptr(LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(false),
			},
		},
	}

	newData := LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []LoadControlLimitConstraintsDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(0)),
				ValueStepSize: NewScaledNumberType(1),
			},
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				ValueStepSize: NewScaledNumberType(1),
			},
		},
	}

	newData := LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []LoadControlLimitConstraintsDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				ValueStepSize: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(0)),
				LimitCategory: util.Ptr(LoadControlCategoryTypeObligation),
			},
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(LoadControlCategoryTypeObligation),
			},
		},
	}

	newData := LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(LoadControlCategoryTypeOptimization),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.LoadControlLimitDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.Equal(t, LoadControlCategoryTypeObligation, *item1.LimitCategory)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.Equal(t, LoadControlCategoryTypeOptimization, *item2.LimitCategory)
}
