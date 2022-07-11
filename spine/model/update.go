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
	// current items in the function data
	ExistingData() []T
	// items in the payload
	NewData() []T
	// the hash key of the update selector; nil if no selector was given
	UpdateSelectorHashKey() *string
	// the hash key of the delete selector; nil if no selector was given
	DeleteSelectorHashKey() *string
	// determines if the identifiers of the passed item are set
	HasIdentifier(*T) bool
	// copies the data (not the identifiers) from the source to the destination item
	CopyData(source *T, dest *T)
}

// Generates a new list of function items by applying the rules mentioned in the spec
// (EEBus_SPINE_TS_ProtocolSpecification.pdf; chapter "5.3.4 Restricted function exchange with cmdOptions").
// The given data provider is used the get the current items and the items and the filters in the payload.
func UpdateList[T util.HashKeyer](dataProvider UpdateDataProvider[T]) []T {
	newData := dataProvider.NewData()
	existingData := dataProvider.ExistingData()
	if len(newData) == 0 {
		return existingData
	}

	// TODO: consider filterDelete
	// TODO: Check if only single fields should be considered here

	// check if selector is used
	updateSelectorHashKey := dataProvider.UpdateSelectorHashKey()
	if updateSelectorHashKey != nil {
		return copyToSelectedData(dataProvider, &newData[0], *updateSelectorHashKey)
	}

	// check if items have no identifiers
	if !dataProvider.HasIdentifier(&newData[0]) {
		// no identifiers specified --> copy data to all existing items
		// (see EEBus_SPINE_TS_ProtocolSpecification.pdf, Table 7: Considered cmdOptions combinations for classifier "notify")
		return copyToAllData(dataProvider, &newData[0])
	}

	return util.Merge(existingData, newData)
}

func copyToSelectedData[T util.HashKeyer](dataProvider UpdateDataProvider[T], newData *T, selectorHashKey string) []T {
	existingData := dataProvider.ExistingData()
	for i := range existingData {
		if existingData[i].HashKey() == selectorHashKey {
			dataProvider.CopyData(newData, &existingData[i])
			break
		}
	}
	return existingData
}

func copyToAllData[T util.HashKeyer](dataProvider UpdateDataProvider[T], newData *T) []T {
	existingData := dataProvider.ExistingData()
	for i := range existingData {
		dataProvider.CopyData(newData, &existingData[i])
	}
	return existingData
}
