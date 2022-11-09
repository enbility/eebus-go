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
	return fmt.Sprintf("%d", r.MeasurementId)
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

	if item == nil || filter == nil {
		return false
	}

	selector := filter.MeasurementListDataSelectors
	if selector == nil {
		return false
	}

	if selector.MeasurementId != nil && *selector.MeasurementId != *item.MeasurementId {
		return false
	}

	if selector.ValueType != nil && *selector.ValueType != *item.ValueType {
		return false
	}

	// TODO: Add selector.TimestampInterval

	return true
}

func (r *MeasurementListDataType_Updater) HasIdentifier(item *MeasurementDataType) bool {
	return item.MeasurementId != nil
}

func (r *MeasurementListDataType_Updater) CopyData(source *MeasurementDataType, dest *MeasurementDataType) {
	if source != nil && dest != nil {
		if source.Timestamp != nil {
			dest.Timestamp = source.Timestamp
		}
		if source.EvaluationPeriod != nil {
			dest.EvaluationPeriod = source.EvaluationPeriod
		}
		if source.ValueState != nil {
			dest.ValueState = source.ValueState
		}
		if source.Value != nil {
			dest.Value = source.Value
		}
		if source.ValueTendency != nil {
			dest.ValueTendency = source.ValueTendency
		}
	}
}
