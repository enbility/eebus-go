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
