package model

// LoadControlLimitListDataType

var _ Updater = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementDataType
	if newList != nil {
		newData = newList.(*MeasurementListDataType).MeasurementData
	}

	r.MeasurementData = UpdateList(r.MeasurementData, newData, filterPartial, filterDelete)
}
