package model

// AlarmListDataType

var _ Updater = (*AlarmListDataType)(nil)

func (r *AlarmListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []AlarmDataType
	if newList != nil {
		newData = newList.(*AlarmListDataType).AlarmListData
	}

	r.AlarmListData = UpdateList(r.AlarmListData, newData, filterPartial, filterDelete)
}
