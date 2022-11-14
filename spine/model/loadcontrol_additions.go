package model

// LoadControlLimitListDataType

var _ Updater = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []LoadControlLimitDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitListDataType).LoadControlLimitData
	}

	r.LoadControlLimitData = UpdateList(r.LoadControlLimitData, newData, filterPartial, filterDelete)
}

// LoadControlLimitDescriptionListDataType

var _ Updater = (*LoadControlLimitDescriptionListDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []LoadControlLimitDescriptionDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitDescriptionListDataType).LoadControlLimitDescriptionData
	}

	r.LoadControlLimitDescriptionData = UpdateList(r.LoadControlLimitDescriptionData, newData, filterPartial, filterDelete)
}
