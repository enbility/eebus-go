package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMessagingListDataType_Update(t *testing.T) {
	sut := MessagingListDataType{
		MessagingData: []MessagingDataType{
			{
				MessagingNumber: util.Ptr(MessagingNumberType(0)),
				Text:            util.Ptr(MessagingDataTextType("old")),
			},
			{
				MessagingNumber: util.Ptr(MessagingNumberType(1)),
				Text:            util.Ptr(MessagingDataTextType("old")),
			},
		},
	}

	newData := MessagingListDataType{
		MessagingData: []MessagingDataType{
			{
				MessagingNumber: util.Ptr(MessagingNumberType(1)),
				Text:            util.Ptr(MessagingDataTextType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.MessagingData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MessagingNumber))
	assert.Equal(t, "old", string(*item1.Text))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MessagingNumber))
	assert.Equal(t, "new", string(*item2.Text))
}
