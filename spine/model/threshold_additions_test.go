package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestThresholdListDataType_Update(t *testing.T) {
	sut := model.ThresholdListDataType{
		ThresholdData: []model.ThresholdDataType{
			{
				ThresholdId:    util.Ptr(model.ThresholdIdType(0)),
				ThresholdValue: model.NewScaledNumberType(1),
			},
			{
				ThresholdId:    util.Ptr(model.ThresholdIdType(1)),
				ThresholdValue: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.ThresholdListDataType{
		ThresholdData: []model.ThresholdDataType{
			{
				ThresholdId:    util.Ptr(model.ThresholdIdType(1)),
				ThresholdValue: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []model.ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(model.ThresholdIdType(0)),
				ThresholdRangeMin: model.NewScaledNumberType(1),
			},
			{
				ThresholdId:       util.Ptr(model.ThresholdIdType(1)),
				ThresholdRangeMin: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []model.ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(model.ThresholdIdType(1)),
				ThresholdRangeMin: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []model.ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(model.ThresholdIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				ThresholdId: util.Ptr(model.ThresholdIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []model.ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(model.ThresholdIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
