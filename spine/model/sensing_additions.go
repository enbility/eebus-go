package model

// SensingListDataType

var _ Updater = (*SensingListDataType)(nil)

func (r *SensingListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SensingDataType
	if newList != nil {
		newData = newList.(*SensingListDataType).SensingData
	}

	r.SensingData = UpdateList(r.SensingData, newData, filterPartial, filterDelete)
}
