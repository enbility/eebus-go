package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestBillListDataType_Update(t *testing.T) {
	sut := BillListDataType{
		BillData: []BillDataType{
			{
				BillId:    util.Ptr(BillIdType(0)),
				ScopeType: util.Ptr(ScopeTypeTypeACCurrent),
			},
			{
				BillId:    util.Ptr(BillIdType(1)),
				ScopeType: util.Ptr(ScopeTypeTypeACCurrent),
			},
		},
	}

	newData := BillListDataType{
		BillData: []BillDataType{
			{
				BillId:    util.Ptr(BillIdType(1)),
				ScopeType: util.Ptr(ScopeTypeTypeACPower),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.BillData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, ScopeTypeTypeACCurrent, *item1.ScopeType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, ScopeTypeTypeACPower, *item2.ScopeType)
}

func TestBillConstraintsListDataType_Update(t *testing.T) {
	sut := BillConstraintsListDataType{
		BillConstraintsData: []BillConstraintsDataType{
			{
				BillId:           util.Ptr(BillIdType(0)),
				PositionCountMin: util.Ptr(BillPositionCountType(0)),
			},
			{
				BillId:           util.Ptr(BillIdType(1)),
				PositionCountMin: util.Ptr(BillPositionCountType(0)),
			},
		},
	}

	newData := BillConstraintsListDataType{
		BillConstraintsData: []BillConstraintsDataType{
			{
				BillId:           util.Ptr(BillIdType(1)),
				PositionCountMin: util.Ptr(BillPositionCountType(1)),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.BillConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, 0, int(*item1.PositionCountMin))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, 1, int(*item2.PositionCountMin))
}

func TestBillDescriptionListDataType_Update(t *testing.T) {
	sut := BillDescriptionListDataType{
		BillDescriptionData: []BillDescriptionDataType{
			{
				BillId:         util.Ptr(BillIdType(0)),
				UpdateRequired: util.Ptr(false),
			},
			{
				BillId:         util.Ptr(BillIdType(1)),
				UpdateRequired: util.Ptr(false),
			},
		},
	}

	newData := BillDescriptionListDataType{
		BillDescriptionData: []BillDescriptionDataType{
			{
				BillId:         util.Ptr(BillIdType(1)),
				UpdateRequired: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.BillDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, false, *item1.UpdateRequired)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, true, *item2.UpdateRequired)
}
