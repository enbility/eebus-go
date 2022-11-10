package model

import (
	"fmt"
	"sort"

	"github.com/DerAndereAndi/eebus-go/util"
)

// ElectricalConnectionPermittedValueSetListDataType

var _ UpdaterFactory[ElectricalConnectionPermittedValueSetListDataType] = (*ElectricalConnectionPermittedValueSetListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionPermittedValueSetDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) NewUpdater(
	newList *ElectricalConnectionPermittedValueSetListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionPermittedValueSetListDataType_Updater{
		ElectricalConnectionPermittedValueSetListDataType: r,
		newData: newList.ElectricalConnectionPermittedValueSetData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

func (r ElectricalConnectionPermittedValueSetDataType) HashKey() string {
	return fmt.Sprintf("%d|%d", *r.ElectricalConnectionId, *r.ParameterId)
}

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType_Updater)(nil)
var _ UpdateDataProvider[ElectricalConnectionPermittedValueSetDataType] = (*ElectricalConnectionPermittedValueSetListDataType_Updater)(nil)

type ElectricalConnectionPermittedValueSetListDataType_Updater struct {
	*ElectricalConnectionPermittedValueSetListDataType
	*FilterProvider
	newData []ElectricalConnectionPermittedValueSetDataType
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionPermittedValueSetData = UpdateList[ElectricalConnectionPermittedValueSetDataType](r.ElectricalConnectionPermittedValueSetData, r.newData, r)
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) HasSelector(filterType FilterEnumType) bool {
	filter := r.FilterForEnumType(filterType)

	return filter != nil && filter.ElectricalConnectionPermittedValueSetListDataSelectors != nil
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) SelectorMatch(filterType FilterEnumType, item *ElectricalConnectionPermittedValueSetDataType) bool {
	filter := r.FilterForEnumType(filterType)

	if item == nil || filter == nil {
		return false
	}

	selector := filter.ElectricalConnectionPermittedValueSetListDataSelectors
	if selector == nil {
		return false
	}

	if selector.ElectricalConnectionId != nil && *selector.ElectricalConnectionId != *item.ElectricalConnectionId {
		return false
	}

	if selector.ParameterId != nil && *selector.ParameterId != *item.ParameterId {
		return false
	}

	return true
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) Sort(data []ElectricalConnectionPermittedValueSetDataType) []ElectricalConnectionPermittedValueSetDataType {
	sort.Slice(data, func(i, j int) bool {
		item1 := data[i]
		item2 := data[j]
		if item1.ElectricalConnectionId != nil && item2.ElectricalConnectionId != nil && *item1.ElectricalConnectionId != *item2.ElectricalConnectionId {
			return *item1.ElectricalConnectionId < *item2.ElectricalConnectionId
		}

		if item1.ParameterId != nil && item2.ParameterId != nil && *item1.ParameterId != *item2.ParameterId {
			return *item1.ParameterId < *item2.ParameterId
		}

		return false
	})

	return data
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) HasIdentifier(item *ElectricalConnectionPermittedValueSetDataType) bool {
	return item.ElectricalConnectionId != nil && item.ParameterId != nil
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) CopyData(source *ElectricalConnectionPermittedValueSetDataType, dest *ElectricalConnectionPermittedValueSetDataType) {
	if source != nil && dest != nil {
		if source.PermittedValueSet != nil {
			dest.PermittedValueSet = source.PermittedValueSet
		}
	}
}

// ElectricalConnectionDescriptionListDataType

var _ UpdaterFactory[ElectricalConnectionDescriptionListDataType] = (*ElectricalConnectionDescriptionListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionDescriptionDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) NewUpdater(
	newList *ElectricalConnectionDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionDescriptionListDataType_Updater{
		ElectricalConnectionDescriptionListDataType: r,
		newData: newList.ElectricalConnectionDescriptionData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

func (r ElectricalConnectionDescriptionDataType) HashKey() string {
	return fmt.Sprintf("%d", *r.ElectricalConnectionId)
}

var _ Updater = (*ElectricalConnectionDescriptionListDataType_Updater)(nil)
var _ UpdateDataProvider[ElectricalConnectionDescriptionDataType] = (*ElectricalConnectionDescriptionListDataType_Updater)(nil)

type ElectricalConnectionDescriptionListDataType_Updater struct {
	*ElectricalConnectionDescriptionListDataType
	*FilterProvider
	newData []ElectricalConnectionDescriptionDataType
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionDescriptionData = UpdateList[ElectricalConnectionDescriptionDataType](r.ElectricalConnectionDescriptionData, r.newData, r)
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) HasSelector(filterType FilterEnumType) bool {
	filter := r.FilterForEnumType(filterType)

	return filter != nil && filter.ElectricalConnectionDescriptionListDataSelectors != nil
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) SelectorMatch(filterType FilterEnumType, item *ElectricalConnectionDescriptionDataType) bool {
	filter := r.FilterForEnumType(filterType)

	if item == nil || filter == nil {
		return false
	}

	selector := filter.ElectricalConnectionDescriptionListDataSelectors
	if selector == nil {
		return false
	}

	if selector.ElectricalConnectionId != nil && *selector.ElectricalConnectionId != *item.ElectricalConnectionId {
		return false
	}

	if selector.ScopeType != nil && *selector.ScopeType != *item.ScopeType {
		return false
	}

	return true
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) Sort(data []ElectricalConnectionDescriptionDataType) []ElectricalConnectionDescriptionDataType {
	sort.Slice(data, func(i, j int) bool {
		item1 := data[i]
		item2 := data[j]
		if item1.ElectricalConnectionId != nil && item2.ElectricalConnectionId != nil && *item1.ElectricalConnectionId != *item2.ElectricalConnectionId {
			return *item1.ElectricalConnectionId < *item2.ElectricalConnectionId
		}

		return false
	})

	return data
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) HasIdentifier(item *ElectricalConnectionDescriptionDataType) bool {
	return item.ElectricalConnectionId != nil
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) CopyData(source *ElectricalConnectionDescriptionDataType, dest *ElectricalConnectionDescriptionDataType) {
	if source != nil && dest != nil {
		if source.AcConnectedPhases != nil {
			dest.AcConnectedPhases = source.AcConnectedPhases
		}
		if source.AcRmsPeriodDuration != nil {
			dest.AcRmsPeriodDuration = source.AcRmsPeriodDuration
		}
		if source.Description != nil {
			dest.Description = source.Description
		}
		if source.Label != nil {
			dest.Label = source.Label
		}
		if source.PositiveEnergyDirection != nil {
			dest.PositiveEnergyDirection = source.PositiveEnergyDirection
		}
		if source.PowerSupplyType != nil {
			dest.PowerSupplyType = source.PowerSupplyType
		}
		if source.ScopeType != nil {
			dest.ScopeType = source.ScopeType
		}
	}
}

// ElectricalConnectionParameterDescriptionListDataType

var _ UpdaterFactory[ElectricalConnectionParameterDescriptionListDataType] = (*ElectricalConnectionParameterDescriptionListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionParameterDescriptionDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) NewUpdater(
	newList *ElectricalConnectionParameterDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionParameterDescriptionListDataType_Updater{
		ElectricalConnectionParameterDescriptionListDataType: r,
		newData: newList.ElectricalConnectionParameterDescriptionData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

// TODO: selector should support any of electricalconnectionid, measurementid, parameterid
func (r ElectricalConnectionParameterDescriptionDataType) HashKey() string {
	return fmt.Sprintf("%d|%d|%d", *r.ElectricalConnectionId, *r.ParameterId, *r.MeasurementId)
}

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType_Updater)(nil)
var _ UpdateDataProvider[ElectricalConnectionParameterDescriptionDataType] = (*ElectricalConnectionParameterDescriptionListDataType_Updater)(nil)

type ElectricalConnectionParameterDescriptionListDataType_Updater struct {
	*ElectricalConnectionParameterDescriptionListDataType
	*FilterProvider
	newData []ElectricalConnectionParameterDescriptionDataType
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionParameterDescriptionData = UpdateList[ElectricalConnectionParameterDescriptionDataType](r.ElectricalConnectionParameterDescriptionData, r.newData, r)
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) HasSelector(filterType FilterEnumType) bool {
	filter := r.FilterForEnumType(filterType)

	return filter != nil && filter.ElectricalConnectionParameterDescriptionListDataSelectors != nil
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) SelectorMatch(filterType FilterEnumType, item *ElectricalConnectionParameterDescriptionDataType) bool {
	filter := r.FilterForEnumType(filterType)

	if item == nil || filter == nil {
		return false
	}

	selector := filter.ElectricalConnectionParameterDescriptionListDataSelectors
	if selector == nil {
		return false
	}

	if selector.ElectricalConnectionId != nil && *selector.ElectricalConnectionId != *item.ElectricalConnectionId {
		return false
	}

	if selector.ParameterId != nil && *selector.ParameterId != *item.ParameterId {
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

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) Sort(data []ElectricalConnectionParameterDescriptionDataType) []ElectricalConnectionParameterDescriptionDataType {
	sort.Slice(data, func(i, j int) bool {
		item1 := data[i]
		item2 := data[j]
		if item1.ElectricalConnectionId != nil && item2.ElectricalConnectionId != nil && *item1.ElectricalConnectionId != *item2.ElectricalConnectionId {
			return *item1.ElectricalConnectionId < *item2.ElectricalConnectionId
		}

		if item1.ParameterId != nil && item2.ParameterId != nil && *item1.ParameterId != *item2.ParameterId {
			return *item1.ParameterId < *item2.ParameterId
		}

		if item1.MeasurementId != nil && item2.MeasurementId != nil {
			return *item1.MeasurementId < *item2.MeasurementId
		}

		return false
	})

	return data
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) HasIdentifier(item *ElectricalConnectionParameterDescriptionDataType) bool {
	return item.ElectricalConnectionId != nil
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) CopyData(source *ElectricalConnectionParameterDescriptionDataType, dest *ElectricalConnectionParameterDescriptionDataType) {
	if source != nil && dest != nil {
		if source.AcMeasuredHarmonic != nil {
			dest.AcMeasuredHarmonic = source.AcMeasuredHarmonic
		}
		if source.AcMeasuredInReferenceTo != nil {
			dest.AcMeasuredInReferenceTo = source.AcMeasuredInReferenceTo
		}
		if source.AcMeasuredPhases != nil {
			dest.AcMeasuredPhases = source.AcMeasuredPhases
		}
		if source.AcMeasurementType != nil {
			dest.AcMeasurementType = source.AcMeasurementType
		}
		if source.AcMeasurementVariant != nil {
			dest.AcMeasurementVariant = source.AcMeasurementVariant
		}
		if source.Description != nil {
			dest.Description = source.Description
		}
		if source.Label != nil {
			dest.Label = source.Label
		}
		if source.ScopeType != nil {
			dest.ScopeType = source.ScopeType
		}
		if source.VoltageType != nil {
			dest.VoltageType = source.VoltageType
		}
	}
}
