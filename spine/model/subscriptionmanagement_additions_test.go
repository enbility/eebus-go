package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionManagementEntryListDataType_Update(t *testing.T) {
	sut := model.SubscriptionManagementEntryListDataType{
		SubscriptionManagementEntryData: []model.SubscriptionManagementEntryDataType{
			{
				SubscriptionId: util.Ptr(model.SubscriptionIdType(0)),
				Description:    util.Ptr(model.DescriptionType("old")),
			},
			{
				SubscriptionId: util.Ptr(model.SubscriptionIdType(1)),
				Description:    util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.SubscriptionManagementEntryListDataType{
		SubscriptionManagementEntryData: []model.SubscriptionManagementEntryDataType{
			{
				SubscriptionId: util.Ptr(model.SubscriptionIdType(1)),
				Description:    util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.SubscriptionManagementEntryData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SubscriptionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SubscriptionId))
	assert.Equal(t, "new", string(*item2.Description))
}
