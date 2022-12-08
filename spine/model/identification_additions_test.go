package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestIdentificationListDataType_Update(t *testing.T) {
	sut := model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:   util.Ptr(model.IdentificationIdType(0)),
				IdentificationType: util.Ptr(model.IdentificationTypeTypeEui48),
			},
			{
				IdentificationId:   util.Ptr(model.IdentificationIdType(1)),
				IdentificationType: util.Ptr(model.IdentificationTypeTypeEui48),
			},
		},
	}

	newData := model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:   util.Ptr(model.IdentificationIdType(1)),
				IdentificationType: util.Ptr(model.IdentificationTypeTypeEui64),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.IdentificationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IdentificationId))
	assert.Equal(t, model.IdentificationTypeTypeEui48, *item1.IdentificationType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IdentificationId))
	assert.Equal(t, model.IdentificationTypeTypeEui64, *item2.IdentificationType)
}
