package model

// MessagingListDataType

var _ Updater = (*MessagingListDataType)(nil)

func (r *MessagingListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MessagingDataType
	if newList != nil {
		newData = newList.(*MessagingListDataType).MessagingData
	}

	r.MessagingData = UpdateList(r.MessagingData, newData, filterPartial, filterDelete)
}
