package model

import (
	"fmt"
	"sort"

	"github.com/DerAndereAndi/eebus-go/util"
)

// LoadControlLimitListDataType

var _ UpdaterFactory[LoadControlLimitListDataType] = (*LoadControlLimitListDataType)(nil)
var _ util.HashKeyer = (*LoadControlLimitDataType)(nil)

func (r *LoadControlLimitListDataType) NewUpdater(
	newList *LoadControlLimitListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &LoadControlLimitListDataType_Updater{
		LoadControlLimitListDataType: r,
		newData:                      newList.LoadControlLimitData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

func (r LoadControlLimitDataType) HashKey() string {
	return fmt.Sprintf("%d", *r.LimitId)
}

var _ Updater = (*LoadControlLimitListDataType_Updater)(nil)
var _ UpdateDataProvider[LoadControlLimitDataType] = (*LoadControlLimitListDataType_Updater)(nil)

type LoadControlLimitListDataType_Updater struct {
	*LoadControlLimitListDataType
	*FilterProvider
	newData []LoadControlLimitDataType
}

func (r *LoadControlLimitListDataType_Updater) DoUpdate() {
	r.LoadControlLimitData = UpdateList[LoadControlLimitDataType](r.LoadControlLimitData, r.newData, r)
}

func (r *LoadControlLimitListDataType_Updater) HasSelector(filterType FilterEnumType) bool {
	filter := r.FilterForEnumType(filterType)

	return filter != nil && filter.LoadControlLimitListDataSelectors != nil
}

func (r *LoadControlLimitListDataType_Updater) SelectorMatch(filterType FilterEnumType, item *LoadControlLimitDataType) bool {
	filter := r.FilterForEnumType(filterType)

	if item == nil || filter == nil {
		return false
	}

	selector := filter.LoadControlLimitListDataSelectors
	if selector == nil {
		return false
	}

	if selector.LimitId != nil && *selector.LimitId != *item.LimitId {
		return false
	}

	return true
}

func (r *LoadControlLimitListDataType_Updater) Sort(data []LoadControlLimitDataType) []LoadControlLimitDataType {
	sort.Slice(data, func(i, j int) bool {
		item1 := data[i]
		item2 := data[j]
		if item1.LimitId != nil && item2.LimitId != nil && *item1.LimitId != *item2.LimitId {
			return *item1.LimitId < *item2.LimitId
		}

		return false
	})

	return data
}

func (r *LoadControlLimitListDataType_Updater) HasIdentifier(item *LoadControlLimitDataType) bool {
	return item.LimitId != nil
}

func (r *LoadControlLimitListDataType_Updater) CopyData(source *LoadControlLimitDataType, dest *LoadControlLimitDataType) {
	if source != nil && dest != nil {
		if source.IsLimitChangeable != nil {
			dest.IsLimitChangeable = source.IsLimitChangeable
		}
		if source.IsLimitActive != nil {
			dest.IsLimitActive = source.IsLimitActive
		}
		if source.TimePeriod != nil {
			dest.TimePeriod = source.TimePeriod
		}
		if source.Value != nil {
			dest.Value = source.Value
		}
	}
}

// LoadControlLimitDescriptionListDataType

var _ UpdaterFactory[LoadControlLimitDescriptionListDataType] = (*LoadControlLimitDescriptionListDataType)(nil)
var _ util.HashKeyer = (*LoadControlLimitDescriptionDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) NewUpdater(
	newList *LoadControlLimitDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &LoadControlLimitDescriptionListDataType_Updater{
		LoadControlLimitDescriptionListDataType: r,
		newData:                                 newList.LoadControlLimitDescriptionData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

func (r LoadControlLimitDescriptionDataType) HashKey() string {
	return fmt.Sprintf("%d|%d", *r.LimitId, *r.MeasurementId)
}

var _ Updater = (*LoadControlLimitDescriptionListDataType_Updater)(nil)
var _ UpdateDataProvider[LoadControlLimitDescriptionDataType] = (*LoadControlLimitDescriptionListDataType_Updater)(nil)

type LoadControlLimitDescriptionListDataType_Updater struct {
	*LoadControlLimitDescriptionListDataType
	*FilterProvider
	newData []LoadControlLimitDescriptionDataType
}

func (r *LoadControlLimitDescriptionListDataType_Updater) DoUpdate() {
	r.LoadControlLimitDescriptionData = UpdateList[LoadControlLimitDescriptionDataType](r.LoadControlLimitDescriptionData, r.newData, r)
}

func (r *LoadControlLimitDescriptionListDataType_Updater) HasSelector(filterType FilterEnumType) bool {
	filter := r.FilterForEnumType(filterType)

	return filter != nil && filter.LoadControlLimitDescriptionListDataSelectors != nil
}

func (r *LoadControlLimitDescriptionListDataType_Updater) SelectorMatch(filterType FilterEnumType, item *LoadControlLimitDescriptionDataType) bool {
	filter := r.FilterForEnumType(filterType)

	if item == nil || filter == nil {
		return false
	}

	selector := filter.LoadControlLimitDescriptionListDataSelectors
	if selector == nil {
		return false
	}

	if selector.LimitId != nil && *selector.LimitId != *item.LimitId {
		return false
	}

	if selector.LimitDirection != nil && *selector.LimitDirection != *item.LimitDirection {
		return false
	}

	if selector.LimitType != nil && *selector.LimitType != *item.LimitType {
		return false
	}

	if selector.MeasurementId != nil && *selector.MeasurementId != *item.MeasurementId {
		return false
	}

	if selector.ScopeType != nil && *selector.ScopeType != *item.ScopeType {
		return false
	}

	return true
}

func (r *LoadControlLimitDescriptionListDataType_Updater) Sort(data []LoadControlLimitDescriptionDataType) []LoadControlLimitDescriptionDataType {
	sort.Slice(data, func(i, j int) bool {
		item1 := data[i]
		item2 := data[j]
		if item1.LimitId != nil && item2.LimitId != nil && *item1.LimitId != *item2.LimitId {
			return *item1.LimitId < *item2.LimitId
		}

		if item1.MeasurementId != nil && item2.MeasurementId != nil {
			return *item1.MeasurementId < *item2.MeasurementId
		}

		return false
	})

	return data
}

func (r *LoadControlLimitDescriptionListDataType_Updater) HasIdentifier(item *LoadControlLimitDescriptionDataType) bool {
	return item.LimitId != nil
}

func (r *LoadControlLimitDescriptionListDataType_Updater) CopyData(source *LoadControlLimitDescriptionDataType, dest *LoadControlLimitDescriptionDataType) {
	if source != nil && dest != nil {
		if source.LimitType != nil {
			dest.LimitType = source.LimitType
		}
		if source.LimitCategory != nil {
			dest.LimitCategory = source.LimitCategory
		}
		if source.LimitDirection != nil {
			dest.LimitDirection = source.LimitDirection
		}
		if source.Unit != nil {
			dest.Unit = source.Unit
		}
		if source.ScopeType != nil {
			dest.ScopeType = source.ScopeType
		}
		if source.Label != nil {
			dest.Label = source.Label
		}
		if source.Description != nil {
			dest.Description = source.Description
		}
	}
}
