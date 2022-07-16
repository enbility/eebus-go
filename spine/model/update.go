package model

import "github.com/DerAndereAndi/eebus-go/util"

type Updater interface {
	DoUpdate()
}

// interface which needs to be implemented by the model function type like 'ElectricalConnectionPermittedValueSetListDataType'
// so that partial updates of the function data is supported
type UpdaterFactory[T any] interface {
	NewUpdater(s *T, filterPartial *FilterType, filterDelete *FilterType) Updater
}

type UpdateDataProvider[T util.HashKeyer] interface {
	// is a partial selector given?
	HasUpdateSelector() bool
	// checks if the given item matches the partial selector
	UpdateSelectorMatch(*T) bool

	// is a delete selector given?
	HasDeleteSelector() bool
	// checks if the given item matches the delete selector
	DeleteSelectorMatch(*T) bool

	// determines if the identifiers of the passed item are set
	HasIdentifier(*T) bool
	// copies the data (not the identifiers) from the source to the destination item
	CopyData(source *T, dest *T)
}

// Generates a new list of function items by applying the rules mentioned in the spec
// (EEBus_SPINE_TS_ProtocolSpecification.pdf; chapter "5.3.4 Restricted function exchange with cmdOptions").
// The given data provider is used the get the current items and the items and the filters in the payload.
func UpdateList[T util.HashKeyer](existingData []T, newData []T, dataProvider UpdateDataProvider[T]) []T {
	if len(newData) == 0 {
		return existingData
	}

	// TODO: Check if only single fields should be considered here

	// process delete selector
	if dataProvider.HasDeleteSelector() {
		existingData = deleteSelectedData(existingData, dataProvider)
		if !dataProvider.HasUpdateSelector() {
			return existingData
		}
	}

	// process update selector
	if dataProvider.HasUpdateSelector() {
		return copyToSelectedData(existingData, dataProvider, &newData[0])
	}

	// check if items have no identifiers
	if !dataProvider.HasIdentifier(&newData[0]) {
		// no identifiers specified --> copy data to all existing items
		// (see EEBus_SPINE_TS_ProtocolSpecification.pdf, Table 7: Considered cmdOptions combinations for classifier "notify")
		return copyToAllData(existingData, dataProvider, &newData[0])
	}

	return util.Merge(existingData, newData)
}

func copyToSelectedData[T util.HashKeyer](existingData []T, dataProvider UpdateDataProvider[T], newData *T) []T {
	for i := range existingData {
		if dataProvider.UpdateSelectorMatch(util.Ptr(existingData[i])) {
			dataProvider.CopyData(newData, &existingData[i])
			break
		}
	}
	return existingData
}

func copyToAllData[T util.HashKeyer](existingData []T, dataProvider UpdateDataProvider[T], newData *T) []T {
	for i := range existingData {
		dataProvider.CopyData(newData, &existingData[i])
	}
	return existingData
}

func deleteSelectedData[T util.HashKeyer](existingData []T, dataProvider UpdateDataProvider[T]) []T {
	result := []T{}
	for i := range existingData {
		if !dataProvider.DeleteSelectorMatch(util.Ptr(existingData[i])) {
			result = append(result, existingData[i])
		}
	}
	return result
}
