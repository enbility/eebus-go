package model

import "github.com/DerAndereAndi/eebus-go/util"

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType_Updater)(nil)
var _ UpdateDataProvider[ElectricalConnectionPermittedValueSetDataType] = (*ElectricalConnectionPermittedValueSetListDataType_Updater)(nil)

type ElectricalConnectionPermittedValueSetListDataType_Updater struct {
	*ElectricalConnectionPermittedValueSetListDataType
	newData       []ElectricalConnectionPermittedValueSetDataType
	filterPartial *FilterType
	filterDelete  *FilterType
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionPermittedValueSetData = UpdateList[ElectricalConnectionPermittedValueSetDataType](r.ElectricalConnectionPermittedValueSetData, r.newData, r)
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) HasUpdateSelector() bool {
	return r.filterPartial != nil && r.filterPartial.ElectricalConnectionPermittedValueSetListDataSelectors != nil
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) UpdateSelectorMatch(item *ElectricalConnectionPermittedValueSetDataType) bool {
	return r.HasUpdateSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterPartial)
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) HasDeleteSelector() bool {
	return r.filterDelete != nil && r.filterDelete.ElectricalConnectionPermittedValueSetListDataSelectors != nil
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) DeleteSelectorMatch(item *ElectricalConnectionPermittedValueSetDataType) bool {
	return r.HasDeleteSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterDelete)
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) HasIdentifier(item *ElectricalConnectionPermittedValueSetDataType) bool {
	return item.ElectricalConnectionId != nil && item.ParameterId != nil
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) CopyData(source *ElectricalConnectionPermittedValueSetDataType, dest *ElectricalConnectionPermittedValueSetDataType) {
	if source != nil && dest != nil {
		dest.PermittedValueSet = source.PermittedValueSet
	}
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) selectorHashKey(filter *FilterType) *string {
	var result *string = nil
	if filter != nil && filter.ElectricalConnectionPermittedValueSetListDataSelectors != nil {
		result = util.Ptr(electricalConnectionPermittedValueSetDataHashKey(
			filter.ElectricalConnectionPermittedValueSetListDataSelectors.ElectricalConnectionId,
			filter.ElectricalConnectionPermittedValueSetListDataSelectors.ParameterId))
	}
	return result
}

var _ Updater = (*ElectricalConnectionDescriptionListDataType_Updater)(nil)
var _ UpdateDataProvider[ElectricalConnectionDescriptionDataType] = (*ElectricalConnectionDescriptionListDataType_Updater)(nil)

type ElectricalConnectionDescriptionListDataType_Updater struct {
	*ElectricalConnectionDescriptionListDataType
	newData       []ElectricalConnectionDescriptionDataType
	filterPartial *FilterType
	filterDelete  *FilterType
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionDescriptionData = UpdateList[ElectricalConnectionDescriptionDataType](r.ElectricalConnectionDescriptionData, r.newData, r)
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) HasUpdateSelector() bool {
	return r.filterPartial != nil && r.filterPartial.ElectricalConnectionDescriptionListDataSelectors != nil
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) UpdateSelectorMatch(item *ElectricalConnectionDescriptionDataType) bool {
	return r.HasUpdateSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterPartial)
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) HasDeleteSelector() bool {
	return r.filterDelete != nil && r.filterDelete.ElectricalConnectionDescriptionListDataSelectors != nil
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) DeleteSelectorMatch(item *ElectricalConnectionDescriptionDataType) bool {
	return r.HasDeleteSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterDelete)
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) HasIdentifier(item *ElectricalConnectionDescriptionDataType) bool {
	return item.ElectricalConnectionId != nil
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) CopyData(source *ElectricalConnectionDescriptionDataType, dest *ElectricalConnectionDescriptionDataType) {
	if source != nil && dest != nil {
		dest.AcConnectedPhases = source.AcConnectedPhases
		dest.AcRmsPeriodDuration = source.AcRmsPeriodDuration
		dest.Description = source.Description
		dest.Label = source.Label
		dest.PositiveEnergyDirection = source.PositiveEnergyDirection
		dest.PowerSupplyType = source.PowerSupplyType
		dest.ScopeType = source.ScopeType
	}
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) selectorHashKey(filter *FilterType) *string {
	var result *string = nil
	if filter != nil && filter.ElectricalConnectionDescriptionListDataSelectors != nil {
		result = util.Ptr(electricalConnectionDescriptionDataHashKey(
			filter.ElectricalConnectionDescriptionListDataSelectors.ElectricalConnectionId))
	}
	return result
}

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType_Updater)(nil)
var _ UpdateDataProvider[ElectricalConnectionParameterDescriptionDataType] = (*ElectricalConnectionParameterDescriptionListDataType_Updater)(nil)

type ElectricalConnectionParameterDescriptionListDataType_Updater struct {
	*ElectricalConnectionParameterDescriptionListDataType
	newData       []ElectricalConnectionParameterDescriptionDataType
	filterPartial *FilterType
	filterDelete  *FilterType
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionParameterDescriptionData = UpdateList[ElectricalConnectionParameterDescriptionDataType](r.ElectricalConnectionParameterDescriptionData, r.newData, r)
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) HasUpdateSelector() bool {
	return r.filterPartial != nil && r.filterPartial.ElectricalConnectionParameterDescriptionListDataSelectors != nil
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) UpdateSelectorMatch(item *ElectricalConnectionParameterDescriptionDataType) bool {
	return r.HasUpdateSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterPartial)
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) HasDeleteSelector() bool {
	return r.filterDelete != nil && r.filterDelete.ElectricalConnectionParameterDescriptionListDataSelectors != nil
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) DeleteSelectorMatch(item *ElectricalConnectionParameterDescriptionDataType) bool {
	return r.HasDeleteSelector() && item != nil &&
		item.HashKey() == *r.selectorHashKey(r.filterDelete)
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) HasIdentifier(item *ElectricalConnectionParameterDescriptionDataType) bool {
	return item.ElectricalConnectionId != nil
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) CopyData(source *ElectricalConnectionParameterDescriptionDataType, dest *ElectricalConnectionParameterDescriptionDataType) {
	if source != nil && dest != nil {
		dest.AcMeasuredHarmonic = source.AcMeasuredHarmonic
		dest.AcMeasuredInReferenceTo = source.AcMeasuredInReferenceTo
		dest.AcMeasuredPhases = source.AcMeasuredPhases
		dest.AcMeasurementType = source.AcMeasurementType
		dest.AcMeasurementVariant = source.AcMeasurementVariant
		dest.Description = source.Description
		dest.Label = source.Label
		dest.MeasurementId = source.MeasurementId
		dest.ParameterId = source.ParameterId
		dest.ScopeType = source.ScopeType
		dest.VoltageType = source.VoltageType
	}
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) selectorHashKey(filter *FilterType) *string {
	var result *string = nil
	if filter != nil && filter.ElectricalConnectionParameterDescriptionListDataSelectors != nil {
		result = util.Ptr(electricalConnectionParameterDescriptionDataHashKey(
			filter.ElectricalConnectionParameterDescriptionListDataSelectors.ElectricalConnectionId))
	}
	return result
}
