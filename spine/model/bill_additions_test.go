package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestBillListDataType_Update(t *testing.T) {
	sut := model.BillListDataType{
		BillData: []model.BillDataType{
			{
				BillId:    util.Ptr(model.BillIdType(0)),
				ScopeType: util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				BillId:    util.Ptr(model.BillIdType(1)),
				ScopeType: util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	newData := model.BillListDataType{
		BillData: []model.BillDataType{
			{
				BillId:    util.Ptr(model.BillIdType(1)),
				ScopeType: util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.BillData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, model.ScopeTypeTypeACCurrent, *item1.ScopeType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, model.ScopeTypeTypeACPower, *item2.ScopeType)
}

func TestBillConstraintsListDataType_Update(t *testing.T) {
	sut := model.BillConstraintsListDataType{
		BillConstraintsData: []model.BillConstraintsDataType{
			{
				BillId:           util.Ptr(model.BillIdType(0)),
				PositionCountMin: util.Ptr(model.BillPositionCountType(0)),
			},
			{
				BillId:           util.Ptr(model.BillIdType(1)),
				PositionCountMin: util.Ptr(model.BillPositionCountType(0)),
			},
		},
	}

	newData := model.BillConstraintsListDataType{
		BillConstraintsData: []model.BillConstraintsDataType{
			{
				BillId:           util.Ptr(model.BillIdType(1)),
				PositionCountMin: util.Ptr(model.BillPositionCountType(1)),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.BillDescriptionListDataType{
		BillDescriptionData: []model.BillDescriptionDataType{
			{
				BillId:         util.Ptr(model.BillIdType(0)),
				UpdateRequired: util.Ptr(false),
			},
			{
				BillId:         util.Ptr(model.BillIdType(1)),
				UpdateRequired: util.Ptr(false),
			},
		},
	}

	newData := model.BillDescriptionListDataType{
		BillDescriptionData: []model.BillDescriptionDataType{
			{
				BillId:         util.Ptr(model.BillIdType(1)),
				UpdateRequired: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
