package model

// SetpointListDataType

var _ Updater = (*SetpointListDataType)(nil)

func (r *SetpointListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SetpointDataType
	if newList != nil {
		newData = newList.(*SetpointListDataType).SetpointData
	}

	r.SetpointData = UpdateList(r.SetpointData, newData, filterPartial, filterDelete)
}

// SetpointDescriptionListDataType

var _ Updater = (*SetpointDescriptionListDataType)(nil)

func (r *SetpointDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SetpointDescriptionDataType
	if newList != nil {
		newData = newList.(*SetpointDescriptionListDataType).SetpointDescriptionData
	}

	r.SetpointDescriptionData = UpdateList(r.SetpointDescriptionData, newData, filterPartial, filterDelete)
}
