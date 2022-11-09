package model

import (
	"fmt"

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
	return loadControlDataHashKey(
		r.LimitId)
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

	return r.HasSelector(filterType) && item != nil && filter != nil &&
		item.HashKey() == *r.selectorHashKey(filter)
}

func (r *LoadControlLimitListDataType_Updater) HasIdentifier(item *LoadControlLimitDataType) bool {
	return item.LimitId != nil
}

func (r *LoadControlLimitListDataType_Updater) CopyData(source *LoadControlLimitDataType, dest *LoadControlLimitDataType) {
	if source != nil && dest != nil {
		dest.IsLimitChangeable = source.IsLimitChangeable
		dest.IsLimitActive = source.IsLimitActive
		dest.TimePeriod = source.TimePeriod
		dest.Value = source.Value
	}
}

func (r *LoadControlLimitListDataType_Updater) selectorHashKey(filter *FilterType) *string {
	var result *string = nil
	if filter != nil && filter.LoadControlLimitListDataSelectors != nil {
		result = util.Ptr(loadControlDataHashKey(
			filter.LoadControlLimitListDataSelectors.LimitId,
		))
	}
	return result
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
	return loadControlDataHashKey(
		r.LimitId)
}

func loadControlDataHashKey(limitId *LoadControlLimitIdType) string {
	return fmt.Sprintf("%d", *limitId)
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

	return r.HasSelector(filterType) && item != nil && filter != nil &&
		item.HashKey() == *r.selectorHashKey(filter)
}

func (r *LoadControlLimitDescriptionListDataType_Updater) HasIdentifier(item *LoadControlLimitDescriptionDataType) bool {
	return item.LimitId != nil
}

func (r *LoadControlLimitDescriptionListDataType_Updater) CopyData(source *LoadControlLimitDescriptionDataType, dest *LoadControlLimitDescriptionDataType) {
	if source != nil && dest != nil {
		dest.LimitType = source.LimitType
		dest.LimitCategory = source.LimitCategory
		dest.LimitDirection = source.LimitDirection
		dest.MeasurementId = source.MeasurementId
		dest.Unit = source.Unit
		dest.ScopeType = source.ScopeType
		dest.Label = source.Label
		dest.Description = source.Description
	}
}

func (r *LoadControlLimitDescriptionListDataType_Updater) selectorHashKey(filter *FilterType) *string {
	var result *string = nil
	if filter != nil && filter.LoadControlLimitDescriptionListDataSelectors != nil {
		result = util.Ptr(loadControlDataHashKey(
			filter.LoadControlLimitDescriptionListDataSelectors.LimitId,
		))
	}
	return result
}
