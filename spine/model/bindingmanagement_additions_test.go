package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestBindingManagementEntryListDataType_Update(t *testing.T) {
	sut := BindingManagementEntryListDataType{
		BindingManagementEntryData: []BindingManagementEntryDataType{
			{
				BindingId:   util.Ptr(BindingIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				BindingId:   util.Ptr(BindingIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := BindingManagementEntryListDataType{
		BindingManagementEntryData: []BindingManagementEntryDataType{
			{
				BindingId:   util.Ptr(BindingIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.BindingManagementEntryData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BindingId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BindingId))
	assert.Equal(t, "new", string(*item2.Description))
}
