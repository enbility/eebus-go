package model

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/util"
)

// MeasurementListDataType

var _ UpdaterFactory[MeasurementListDataType] = (*MeasurementListDataType)(nil)
var _ util.HashKeyer = (*MeasurementDataType)(nil)

func (r *MeasurementListDataType) NewUpdater(
	newList *MeasurementListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &MeasurementListDataType_Updater{
		MeasurementListDataType: r,
		newData:                 newList.MeasurementData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

func (r MeasurementDataType) HashKey() string {
	return measurementDataHashKey(
		r.MeasurementId)
}

func measurementDataHashKey(measurementId *MeasurementIdType) string {
	return fmt.Sprintf("%d", *measurementId)
}

var _ Updater = (*MeasurementListDataType_Updater)(nil)
var _ UpdateDataProvider[MeasurementDataType] = (*MeasurementListDataType_Updater)(nil)

type MeasurementListDataType_Updater struct {
	*MeasurementListDataType
	*FilterProvider
	newData []MeasurementDataType
}

func (r *MeasurementListDataType_Updater) DoUpdate() {
	r.MeasurementData = UpdateList[MeasurementDataType](r.MeasurementData, r.newData, r)
}

func (r *MeasurementListDataType_Updater) HasSelector(filterType FilterEnumType) bool {
	filter := r.FilterForEnumType(filterType)

	return filter != nil && filter.MeasurementListDataSelectors != nil
}

func (r *MeasurementListDataType_Updater) SelectorMatch(filterType FilterEnumType, item *MeasurementDataType) bool {
	filter := r.FilterForEnumType(filterType)

	return r.HasSelector(filterType) && item != nil && filter != nil &&
		item.HashKey() == *r.selectorHashKey(filter)
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
