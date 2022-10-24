package model

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/util"
)

var _ UpdaterFactory[MeasurementListDataType] = (*MeasurementListDataType)(nil)
var _ util.HashKeyer = (*MeasurementDataType)(nil)

func (r *MeasurementListDataType) NewUpdater(
	newList *MeasurementListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &MeasurementListDataType_Updater{
		MeasurementListDataType: r,
		newData:                 newList.MeasurementData,
		filterPartial:           filterPartial,
		filterDelete:            filterDelete,
	}
}

func (r MeasurementDataType) HashKey() string {
	return measurementDataHashKey(
		r.MeasurementId)
}

func measurementDataHashKey(measurementId *MeasurementIdType) string {
	return fmt.Sprintf("%d", *measurementId)
}
