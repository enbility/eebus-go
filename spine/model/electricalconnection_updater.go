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
	r.ElectricalConnectionPermittedValueSetData = UpdateList[ElectricalConnectionPermittedValueSetDataType](r)
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) ExistingData() []ElectricalConnectionPermittedValueSetDataType {
	return r.ElectricalConnectionPermittedValueSetData
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) NewData() []ElectricalConnectionPermittedValueSetDataType {
	return r.newData
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) UpdateSelectorHashKey() *string {
	return r.selectorHashKey(r.filterPartial)
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) DeleteSelectorHashKey() *string {
	return r.selectorHashKey(r.filterDelete)
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
