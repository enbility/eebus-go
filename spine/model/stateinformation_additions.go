package model

// StateInformationListDataType

var _ Updater = (*StateInformationListDataType)(nil)

func (r *StateInformationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []StateInformationDataType
	if newList != nil {
		newData = newList.(*StateInformationListDataType).StateInformationData
	}

	r.StateInformationData = UpdateList(r.StateInformationData, newData, filterPartial, filterDelete)
}
