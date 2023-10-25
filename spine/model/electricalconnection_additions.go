package model

// ElectricalConnectionStateListDataType

var _ Updater = (*ElectricalConnectionStateListDataType)(nil)

func (r *ElectricalConnectionStateListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionStateDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionStateListDataType).ElectricalConnectionStateData
	}

	r.ElectricalConnectionStateData = UpdateList(r.ElectricalConnectionStateData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionPermittedValueSetListDataType

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionPermittedValueSetDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionPermittedValueSetListDataType).ElectricalConnectionPermittedValueSetData
	}

	r.ElectricalConnectionPermittedValueSetData = UpdateList(r.ElectricalConnectionPermittedValueSetData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionDescriptionListDataType

var _ Updater = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionDescriptionListDataType).ElectricalConnectionDescriptionData
	}

	r.ElectricalConnectionDescriptionData = UpdateList(r.ElectricalConnectionDescriptionData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionCharacteristicListDataType

var _ Updater = (*ElectricalConnectionCharacteristicListDataType)(nil)

func (r *ElectricalConnectionCharacteristicListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionCharacteristicDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionCharacteristicListDataType).ElectricalConnectionCharacteristicListData
	}

	r.ElectricalConnectionCharacteristicListData = UpdateList(r.ElectricalConnectionCharacteristicListData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionParameterDescriptionListDataType

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionParameterDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionParameterDescriptionListDataType).ElectricalConnectionParameterDescriptionData
	}

	r.ElectricalConnectionParameterDescriptionData = UpdateList(r.ElectricalConnectionParameterDescriptionData, newData, filterPartial, filterDelete)
}
