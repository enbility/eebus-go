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

var _ UpdaterFactory[ElectricalConnectionDescriptionListDataType] = (*ElectricalConnectionDescriptionListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionDescriptionDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) NewUpdater(
	newList *ElectricalConnectionDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionDescriptionListDataType_Updater{
		ElectricalConnectionDescriptionListDataType: r,
		newData:       newList.ElectricalConnectionDescriptionData,
		filterPartial: filterPartial,
		filterDelete:  filterDelete,
	}
}

func (r ElectricalConnectionDescriptionDataType) HashKey() string {
	return electricalConnectionDescriptionDataHashKey(
		r.ElectricalConnectionId)
}

func electricalConnectionDescriptionDataHashKey(electricalConnectionId *ElectricalConnectionIdType) string {
	return fmt.Sprintf("%d", *electricalConnectionId)
}

var _ UpdaterFactory[ElectricalConnectionParameterDescriptionListDataType] = (*ElectricalConnectionParameterDescriptionListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionParameterDescriptionDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) NewUpdater(
	newList *ElectricalConnectionParameterDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionParameterDescriptionListDataType_Updater{
		ElectricalConnectionParameterDescriptionListDataType: r,
		newData:       newList.ElectricalConnectionParameterDescriptionData,
		filterPartial: filterPartial,
		filterDelete:  filterDelete,
	}
}

// TODO: selector should support any of electricalconnectionid, measurementid, parameterid
func (r ElectricalConnectionParameterDescriptionDataType) HashKey() string {
	return electricalConnectionDescriptionDataHashKey(
		r.ElectricalConnectionId)
}

func electricalConnectionParameterDescriptionDataHashKey(electricalConnectionId *ElectricalConnectionIdType) string {
	return fmt.Sprintf("%d", *electricalConnectionId)
}
