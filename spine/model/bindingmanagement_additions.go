package model

// BindingManagementEntryListDataType

var _ Updater = (*BindingManagementEntryListDataType)(nil)

func (r *BindingManagementEntryListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []BindingManagementEntryDataType
	if newList != nil {
		newData = newList.(*BindingManagementEntryListDataType).BindingManagementEntryData
	}

	r.BindingManagementEntryData = UpdateList(r.BindingManagementEntryData, newData, filterPartial, filterDelete)
}
