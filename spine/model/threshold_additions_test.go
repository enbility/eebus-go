package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestThresholdListDataType_Update(t *testing.T) {
	sut := ThresholdListDataType{
		ThresholdData: []ThresholdDataType{
			{
				ThresholdId:    util.Ptr(ThresholdIdType(0)),
				ThresholdValue: NewScaledNumberType(1),
			},
			{
				ThresholdId:    util.Ptr(ThresholdIdType(1)),
				ThresholdValue: NewScaledNumberType(1),
			},
		},
	}

	newData := ThresholdListDataType{
		ThresholdData: []ThresholdDataType{
			{
				ThresholdId:    util.Ptr(ThresholdIdType(1)),
				ThresholdValue: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.ThresholdData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ThresholdId))
	assert.Equal(t, 1.0, item1.ThresholdValue.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ThresholdId))
	assert.Equal(t, 10.0, item2.ThresholdValue.GetValue())
}

func TestThresholdConstraintsListDataType_Update(t *testing.T) {
	sut := ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(ThresholdIdType(0)),
				ThresholdRangeMin: NewScaledNumberType(1),
			},
			{
				ThresholdId:       util.Ptr(ThresholdIdType(1)),
				ThresholdRangeMin: NewScaledNumberType(1),
			},
		},
	}

	newData := ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(ThresholdIdType(1)),
				ThresholdRangeMin: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.ThresholdConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ThresholdId))
	assert.Equal(t, 1.0, item1.ThresholdRangeMin.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ThresholdId))
	assert.Equal(t, 10.0, item2.ThresholdRangeMin.GetValue())
}

func TestThresholdDescriptionListDataType_Update(t *testing.T) {
	sut := ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(ThresholdIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				ThresholdId: util.Ptr(ThresholdIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(ThresholdIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.ThresholdDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ThresholdId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ThresholdId))
	assert.Equal(t, "new", string(*item2.Description))
}
