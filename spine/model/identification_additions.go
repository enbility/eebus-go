package model

// IdentificationListDataType

var _ Updater = (*IdentificationListDataType)(nil)

func (r *IdentificationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []IdentificationDataType
	if newList != nil {
		newData = newList.(*IdentificationListDataType).IdentificationData
	}

	r.IdentificationData = UpdateList(r.IdentificationData, newData, filterPartial, filterDelete)
}
