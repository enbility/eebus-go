package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type TestUpdateData struct {
	Id       *uint `eebus:"key"`
	DataItem *int
}

type TestUpdater struct {
	// updateSelectorHashKey *string
	// deleteSelectorHashKey *string
}

func TestUpdateList_NewItem(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}

	// Act
	result := UpdateList(existingData, newData, nil, nil)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ChangedItem(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}}

	// Act
	result := UpdateList(existingData, newData, nil, nil)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_NewAndChangedItem(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}, {Id: util.Ptr(uint(3)), DataItem: util.Ptr(int(3))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}, {Id: util.Ptr(uint(3)), DataItem: util.Ptr(int(3))}}

	// Act
	result := UpdateList(existingData, newData, nil, nil)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ItemWithNoIdentifier(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}
	newData := []TestUpdateData{{DataItem: util.Ptr(int(3))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(3))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(3))}}

	// Act
	result := UpdateList(existingData, newData, nil, nil)

	assert.Equal(t, expectedResult, result)
}

func TestRemoveFieldFromType(t *testing.T) {
	items := &LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId: util.Ptr(LoadControlLimitIdType(1)),
				Value:   NewScaledNumberType(16.0),
			},
		},
	}

	elements := &LoadControlLimitDataElementsType{
		Value: &ScaledNumberElementsType{},
	}

	RemoveElementFromItem(&items.LoadControlLimitData[0], elements)

	var nilValue *ScaledNumberType

	assert.Equal(t, nilValue, items.LoadControlLimitData[0].Value)
}

// TODO: Fix, as these tests won't work right now as TestUpdater doesn't use FilterProvider and its data structure
/*
func TestUpdateList_UpdateSelector(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(1), DataItem: 1}, {Id: util.Ptr(2), DataItem: 2}}
	newData := []TestUpdateData{{DataItem: 3}}

	dataProvider := &TestUpdater{
		updateSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{Id: util.Ptr(1), DataItem: 3}, {Id: util.Ptr(2), DataItem: 2}}

	// Act
	result := UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_DeleteSelector(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(1), DataItem: 1}, {Id: util.Ptr(2), DataItem: 2}}
	newData := []TestUpdateData{{Id: util.Ptr(0), DataItem: 0}}

	dataProvider := &TestUpdater{
		deleteSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{Id: util.Ptr(2), DataItem: 2}}

	// Act
	result := UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}
*/
