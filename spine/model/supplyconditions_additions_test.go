package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSupplyConditionListDataType_Update(t *testing.T) {
	sut := model.SupplyConditionListDataType{
		SupplyConditionData: []model.SupplyConditionDataType{
			{
				ConditionId: util.Ptr(model.ConditionIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				ConditionId: util.Ptr(model.ConditionIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.SupplyConditionListDataType{
		SupplyConditionData: []model.SupplyConditionDataType{
			{
				ConditionId: util.Ptr(model.ConditionIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.SupplyConditionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ConditionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ConditionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestSupplyConditionDescriptionListDataType_Update(t *testing.T) {
	sut := model.SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []model.SupplyConditionDescriptionDataType{
			{
				ConditionId: util.Ptr(model.ConditionIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				ConditionId: util.Ptr(model.ConditionIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []model.SupplyConditionDescriptionDataType{
			{
				ConditionId: util.Ptr(model.ConditionIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.SupplyConditionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ConditionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ConditionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestSupplyConditionThresholdRelationListDataType_Update(t *testing.T) {
	sut := model.SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []model.SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(model.ConditionIdType(0)),
				ThresholdId: []model.ThresholdIdType{0},
			},
			{
				ConditionId: util.Ptr(model.ConditionIdType(1)),
				ThresholdId: []model.ThresholdIdType{0},
			},
		},
	}

	newData := model.SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []model.SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(model.ConditionIdType(1)),
				ThresholdId: []model.ThresholdIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.SupplyConditionThresholdRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ConditionId))
	assert.Equal(t, 0, int(item1.ThresholdId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ConditionId))
	assert.Equal(t, 1, int(item2.ThresholdId[0]))
}
