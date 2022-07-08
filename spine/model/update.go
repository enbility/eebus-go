package model

import "github.com/DerAndereAndi/eebus-go/util"

type Updater interface {
	DoUpdate()
}

type UpdaterFactory[T any] interface {
	NewUpdater(s *T, filterPartial *FilterType, filterDelete *FilterType) Updater
}

type UpdateDataProvider[T util.HashKeyer] interface {
	ExistingData() []T
	NewData() []T
	UpdateSelektorHashKey() *string
	DeleteSelektorHashKey() *string
	HasIdentifier(*T) bool
	CopyData(source *T, dest *T)
}

func UpdateList[T util.HashKeyer](dataProvider UpdateDataProvider[T]) []T {
	newData := dataProvider.NewData()
	existingData := dataProvider.ExistingData()
	if len(newData) == 0 {
		return existingData
	}

	// TODO: consider filterDelete
	// TODO: Check if only single fields should be considered here

	// check if selector is used
	updateSelectorHashKey := dataProvider.UpdateSelektorHashKey()
	if updateSelectorHashKey != nil {
		return copyToSelectedData(dataProvider, &newData[0], *updateSelectorHashKey)
	}

	// check if items have no identifiers
	if !dataProvider.HasIdentifier(&newData[0]) {
		// no identifiers specified --> copy data to all existing items
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
