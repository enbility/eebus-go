package model

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/util"
)

var _ UpdaterFactory[ElectricalConnectionPermittedValueSetListDataType] = (*ElectricalConnectionPermittedValueSetListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionPermittedValueSetDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) NewUpdater(
	newList *ElectricalConnectionPermittedValueSetListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionPermittedValueSetListDataType_Updater{
		ElectricalConnectionPermittedValueSetListDataType: r,
		newData:       newList.ElectricalConnectionPermittedValueSetData,
		filterPartial: filterPartial,
		filterDelete:  filterDelete,
	}
}

func (r ElectricalConnectionPermittedValueSetDataType) HashKey() string {
	return electricalConnectionPermittedValueSetDataHashKey(
		r.ElectricalConnectionId,
		r.ParameterId)
}

func electricalConnectionPermittedValueSetDataHashKey(electricalConnectionId *ElectricalConnectionIdType, parameterId *ElectricalConnectionParameterIdType) string {
	return fmt.Sprintf("%d|%d", *electricalConnectionId, *parameterId)
}
