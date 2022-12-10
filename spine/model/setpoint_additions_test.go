package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSetpointListDataType_Update(t *testing.T) {
	sut := model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{
			{
				SetpointId: util.Ptr(model.SetpointIdType(0)),
				Value:      model.NewScaledNumberType(1),
			},
			{
				SetpointId: util.Ptr(model.SetpointIdType(1)),
				Value:      model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{
			{
				SetpointId: util.Ptr(model.SetpointIdType(1)),
				Value:      model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.SetpointData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SetpointId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SetpointId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestSetpointDescriptionListDataType_Update(t *testing.T) {
	sut := model.SetpointDescriptionListDataType{
		SetpointDescriptionData: []model.SetpointDescriptionDataType{
			{
				SetpointId:  util.Ptr(model.SetpointIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				SetpointId:  util.Ptr(model.SetpointIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.SetpointDescriptionListDataType{
		SetpointDescriptionData: []model.SetpointDescriptionDataType{
			{
				SetpointId:  util.Ptr(model.SetpointIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.SetpointDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SetpointId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SetpointId))
	assert.Equal(t, "new", string(*item2.Description))
}
