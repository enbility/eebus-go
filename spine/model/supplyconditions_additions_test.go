package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSupplyConditionListDataType_Update(t *testing.T) {
	sut := SupplyConditionListDataType{
		SupplyConditionData: []SupplyConditionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := SupplyConditionListDataType{
		SupplyConditionData: []SupplyConditionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []SupplyConditionDescriptionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []SupplyConditionDescriptionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(0)),
				ThresholdId: []ThresholdIdType{0},
			},
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				ThresholdId: []ThresholdIdType{0},
			},
		},
	}

	newData := SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				ThresholdId: []ThresholdIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
