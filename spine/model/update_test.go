package model_test

import (
	"fmt"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type TestUpdateData struct {
	id       *int
	dataItem int
}

func (r TestUpdateData) HashKey() string {
	if r.id != nil {
		return fmt.Sprintf("%d", *r.id)
	} else {
		return ""
	}
}

var _ model.UpdateDataProvider[TestUpdateData] = (*TestUpdater)(nil)

type TestUpdater struct {
	updateSelectorHashKey *string
	deleteSelectorHashKey *string
}

func (r *TestUpdater) HasUpdateSelector() bool {
	return r.updateSelectorHashKey != nil
}

func (r *TestUpdater) UpdateSelectorMatch(item *TestUpdateData) bool {
	return r.updateSelectorHashKey != nil && item != nil &&
		item.HashKey() == *r.updateSelectorHashKey
}

func (r *TestUpdater) HasDeleteSelector() bool {
	return r.deleteSelectorHashKey != nil
}

func (r *TestUpdater) DeleteSelectorMatch(item *TestUpdateData) bool {
	return r.deleteSelectorHashKey != nil && item != nil &&
		item.HashKey() == *r.deleteSelectorHashKey
}

// determines if the identifiers of the passed item are set
func (r *TestUpdater) HasIdentifier(item *TestUpdateData) bool {
	return item.id != nil
}

// copies the data (not the identifiers) from the source to the destination item
func (r *TestUpdater) CopyData(source *TestUpdateData, dest *TestUpdateData) {
	dest.dataItem = source.dataItem
}

func TestUpdateList_NewItem(t *testing.T) {
	existingData := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}}
	newData := []TestUpdateData{{id: util.Ptr(2), dataItem: 2}}

	dataProvider := &TestUpdater{}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ChangedItem(t *testing.T) {
	existingData := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}}
	newData := []TestUpdateData{{id: util.Ptr(1), dataItem: 2}}

	dataProvider := &TestUpdater{}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_NewAndChangedItem(t *testing.T) {
	existingData := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}}
	newData := []TestUpdateData{{id: util.Ptr(1), dataItem: 2}, {id: util.Ptr(3), dataItem: 3}}

	dataProvider := &TestUpdater{}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 2}, {id: util.Ptr(3), dataItem: 3}}

	// Act
	result := model.UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ItemWithNoIdentifier(t *testing.T) {
	existingData := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}}
	newData := []TestUpdateData{{dataItem: 3}}

	dataProvider := &TestUpdater{}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 3}, {id: util.Ptr(2), dataItem: 3}}

	// Act
	result := model.UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_UpdateSelektor(t *testing.T) {
	existingData := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}}
	newData := []TestUpdateData{{dataItem: 3}}

	dataProvider := &TestUpdater{
		updateSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 3}, {id: util.Ptr(2), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_DeleteSelektor(t *testing.T) {
	existingData := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}}
	newData := []TestUpdateData{{id: util.Ptr(0), dataItem: 0}}

	dataProvider := &TestUpdater{
		deleteSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(2), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}
