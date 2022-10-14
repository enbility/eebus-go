package model

import "github.com/DerAndereAndi/eebus-go/util"

var _ Updater = (*MeasurementListDataType_Updater)(nil)
var _ UpdateDataProvider[MeasurementDataType] = (*MeasurementListDataType_Updater)(nil)

type MeasurementListDataType_Updater struct {
	*MeasurementListDataType
	newData       []MeasurementDataType
	filterPartial *FilterType
	filterDelete  *FilterType
}

func (r *MeasurementListDataType_Updater) DoUpdate() {
	r.MeasurementData = UpdateList[MeasurementDataType](r.MeasurementData, r.newData, r)
}

func (r *MeasurementListDataType_Updater) HasUpdateSelector() bool {
	return r.filterPartial != nil && r.filterPartial.MeasurementListDataSelectors != nil
}

func (r *MeasurementListDataType_Updater) UpdateSelectorMatch(item *MeasurementDataType) bool {
	return r.HasUpdateSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterPartial)
}

func (r *MeasurementListDataType_Updater) HasDeleteSelector() bool {
	return r.filterDelete != nil && r.filterDelete.MeasurementListDataSelectors != nil
}

func (r *MeasurementListDataType_Updater) DeleteSelectorMatch(item *MeasurementDataType) bool {
	return r.HasDeleteSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterDelete)
}

func (r *MeasurementListDataType_Updater) HasIdentifier(item *MeasurementDataType) bool {
	return item.MeasurementId != nil
}

func (r *MeasurementListDataType_Updater) CopyData(source *MeasurementDataType, dest *MeasurementDataType) {
	if source != nil && dest != nil {
		dest.Timestamp = source.Timestamp
		dest.EvaluationPeriod = source.EvaluationPeriod
		dest.ValueState = source.ValueState
		dest.Value = source.Value
		dest.ValueTendency = source.ValueTendency
	}
}

func (r *MeasurementListDataType_Updater) selectorHashKey(filter *FilterType) *string {
	var result *string = nil
	if filter != nil && filter.MeasurementListDataSelectors != nil {
		result = util.Ptr(measurementDataHashKey(
			filter.MeasurementListDataSelectors.MeasurementId,
		))
	}
	return result
}
