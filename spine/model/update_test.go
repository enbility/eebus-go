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
	existingData          []TestUpdateData
	newData               []TestUpdateData
	updateSelectorHashKey *string
	deleteSelectorHashKey *string
}

func (r *TestUpdater) ExistingData() []TestUpdateData {
	return r.existingData
}

func (r *TestUpdater) NewData() []TestUpdateData {
	return r.newData
}

// the hash key of the update selector; nil if no selector was given
func (r *TestUpdater) UpdateSelectorHashKey() *string {
	return r.updateSelectorHashKey
}

// the hash key of the delete selector; nil if no selector was given
func (r *TestUpdater) DeleteSelectorHashKey() *string {
	return r.deleteSelectorHashKey
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
	dataProvider := &TestUpdater{
		existingData: []TestUpdateData{{id: util.Ptr(1), dataItem: 1}},
		newData:      []TestUpdateData{{id: util.Ptr(2), dataItem: 2}},
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ChangedItem(t *testing.T) {
	dataProvider := &TestUpdater{
		existingData: []TestUpdateData{{id: util.Ptr(1), dataItem: 1}},
		newData:      []TestUpdateData{{id: util.Ptr(1), dataItem: 2}},
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_NewAndChangedItem(t *testing.T) {
	dataProvider := &TestUpdater{
		existingData: []TestUpdateData{{id: util.Ptr(1), dataItem: 1}},
		newData:      []TestUpdateData{{id: util.Ptr(1), dataItem: 2}, {id: util.Ptr(3), dataItem: 3}},
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 2}, {id: util.Ptr(3), dataItem: 3}}

	// Act
	result := model.UpdateList[TestUpdateData](dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ItemWithNoIdentifier(t *testing.T) {
	dataProvider := &TestUpdater{
		existingData: []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}},
		newData:      []TestUpdateData{{dataItem: 3}},
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 3}, {id: util.Ptr(2), dataItem: 3}}

	// Act
	result := model.UpdateList[TestUpdateData](dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_UpdateSelektor(t *testing.T) {
	dataProvider := &TestUpdater{
		existingData:          []TestUpdateData{{id: util.Ptr(1), dataItem: 1}, {id: util.Ptr(2), dataItem: 2}},
		newData:               []TestUpdateData{{dataItem: 3}},
		updateSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{id: util.Ptr(1), dataItem: 3}, {id: util.Ptr(2), dataItem: 2}}

	// Act
	result := model.UpdateList[TestUpdateData](dataProvider)

	assert.Equal(t, expectedResult, result)
}
