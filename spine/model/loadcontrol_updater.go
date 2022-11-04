package model

import "github.com/DerAndereAndi/eebus-go/util"

var _ Updater = (*LoadControlLimitListDataType_Updater)(nil)
var _ UpdateDataProvider[LoadControlLimitDataType] = (*LoadControlLimitListDataType_Updater)(nil)

type LoadControlLimitListDataType_Updater struct {
	*LoadControlLimitListDataType
	newData       []LoadControlLimitDataType
	filterPartial *FilterType
	filterDelete  *FilterType
}

func (r *LoadControlLimitListDataType_Updater) DoUpdate() {
	r.LoadControlLimitData = UpdateList[LoadControlLimitDataType](r.LoadControlLimitData, r.newData, r)
}

func (r *LoadControlLimitListDataType_Updater) HasUpdateSelector() bool {
	return r.filterPartial != nil && r.filterPartial.LoadControlLimitListDataSelectors != nil
}

func (r *LoadControlLimitListDataType_Updater) UpdateSelectorMatch(item *LoadControlLimitDataType) bool {
	return r.HasUpdateSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterPartial)
}

func (r *LoadControlLimitListDataType_Updater) HasDeleteSelector() bool {
	return r.filterDelete != nil && r.filterDelete.LoadControlLimitListDataSelectors != nil
}

func (r *LoadControlLimitListDataType_Updater) DeleteSelectorMatch(item *LoadControlLimitDataType) bool {
	return r.HasDeleteSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterDelete)
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

var _ Updater = (*LoadControlLimitDescriptionListDataType_Updater)(nil)
var _ UpdateDataProvider[LoadControlLimitDescriptionDataType] = (*LoadControlLimitDescriptionListDataType_Updater)(nil)

type LoadControlLimitDescriptionListDataType_Updater struct {
	*LoadControlLimitDescriptionListDataType
	newData       []LoadControlLimitDescriptionDataType
	filterPartial *FilterType
	filterDelete  *FilterType
}

func (r *LoadControlLimitDescriptionListDataType_Updater) DoUpdate() {
	r.LoadControlLimitDescriptionData = UpdateList[LoadControlLimitDescriptionDataType](r.LoadControlLimitDescriptionData, r.newData, r)
}

func (r *LoadControlLimitDescriptionListDataType_Updater) HasUpdateSelector() bool {
	return r.filterPartial != nil && r.filterPartial.LoadControlLimitListDataSelectors != nil
}

func (r *LoadControlLimitDescriptionListDataType_Updater) UpdateSelectorMatch(item *LoadControlLimitDescriptionDataType) bool {
	return r.HasUpdateSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterPartial)
}

func (r *LoadControlLimitDescriptionListDataType_Updater) HasDeleteSelector() bool {
	return r.filterDelete != nil && r.filterDelete.LoadControlLimitDescriptionListDataSelectors != nil
}

func (r *LoadControlLimitDescriptionListDataType_Updater) DeleteSelectorMatch(item *LoadControlLimitDescriptionDataType) bool {
	return r.HasDeleteSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterDelete)
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
