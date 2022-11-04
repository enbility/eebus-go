package model

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/util"
)

var _ UpdaterFactory[LoadControlLimitListDataType] = (*LoadControlLimitListDataType)(nil)
var _ util.HashKeyer = (*LoadControlLimitDataType)(nil)

func (r *LoadControlLimitListDataType) NewUpdater(
	newList *LoadControlLimitListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &LoadControlLimitListDataType_Updater{
		LoadControlLimitListDataType: r,
		newData:                      newList.LoadControlLimitData,
		filterPartial:                filterPartial,
		filterDelete:                 filterDelete,
	}
}

func (r LoadControlLimitDataType) HashKey() string {
	return loadControlDataHashKey(
		r.LimitId)
}

var _ UpdaterFactory[LoadControlLimitDescriptionListDataType] = (*LoadControlLimitDescriptionListDataType)(nil)
var _ util.HashKeyer = (*LoadControlLimitDescriptionDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) NewUpdater(
	newList *LoadControlLimitDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &LoadControlLimitDescriptionListDataType_Updater{
		LoadControlLimitDescriptionListDataType: r,
		newData:                                 newList.LoadControlLimitDescriptionData,
		filterPartial:                           filterPartial,
		filterDelete:                            filterDelete,
	}
}

func (r LoadControlLimitDescriptionDataType) HashKey() string {
	return loadControlDataHashKey(
		r.LimitId)
}

func loadControlDataHashKey(limitId *LoadControlLimitIdType) string {
	return fmt.Sprintf("%d", *limitId)
}
