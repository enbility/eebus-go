package model

// SubscriptionManagementEntryListDataType

var _ Updater = (*SubscriptionManagementEntryListDataType)(nil)

func (r *SubscriptionManagementEntryListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SubscriptionManagementEntryDataType
	if newList != nil {
		newData = newList.(*SubscriptionManagementEntryListDataType).SubscriptionManagementEntryData
	}

	r.SubscriptionManagementEntryData = UpdateList(r.SubscriptionManagementEntryData, newData, filterPartial, filterDelete)
}
