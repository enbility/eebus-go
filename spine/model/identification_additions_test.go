package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestIdentificationListDataType_Update(t *testing.T) {
	sut := IdentificationListDataType{
		IdentificationData: []IdentificationDataType{
			{
				IdentificationId:   util.Ptr(IdentificationIdType(0)),
				IdentificationType: util.Ptr(IdentificationTypeTypeEui48),
			},
			{
				IdentificationId:   util.Ptr(IdentificationIdType(1)),
				IdentificationType: util.Ptr(IdentificationTypeTypeEui48),
			},
		},
	}

	newData := IdentificationListDataType{
		IdentificationData: []IdentificationDataType{
			{
				IdentificationId:   util.Ptr(IdentificationIdType(1)),
				IdentificationType: util.Ptr(IdentificationTypeTypeEui64),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.IdentificationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IdentificationId))
	assert.Equal(t, IdentificationTypeTypeEui48, *item1.IdentificationType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IdentificationId))
	assert.Equal(t, IdentificationTypeTypeEui64, *item2.IdentificationType)
}
