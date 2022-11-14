package model

// LoadControlEventListDataType

var _ Updater = (*LoadControlEventListDataType)(nil)

func (r *LoadControlEventListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []LoadControlEventDataType
	if newList != nil {
		newData = newList.(*LoadControlEventListDataType).LoadControlEventData
	}

	r.LoadControlEventData = UpdateList(r.LoadControlEventData, newData, filterPartial, filterDelete)
}

// LoadControlStateListDataType

var _ Updater = (*LoadControlStateListDataType)(nil)

func (r *LoadControlStateListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []LoadControlStateDataType
	if newList != nil {
		newData = newList.(*LoadControlStateListDataType).LoadControlStateData
	}

	r.LoadControlStateData = UpdateList(r.LoadControlStateData, newData, filterPartial, filterDelete)
}

// LoadControlLimitListDataType

var _ Updater = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []LoadControlLimitDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitListDataType).LoadControlLimitData
	}

	r.LoadControlLimitData = UpdateList(r.LoadControlLimitData, newData, filterPartial, filterDelete)
}

// LoadControlLimitConstraintsListDataType

var _ Updater = (*LoadControlLimitConstraintsListDataType)(nil)

func (r *LoadControlLimitConstraintsListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []LoadControlLimitConstraintsDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitConstraintsListDataType).LoadControlLimitConstraintsData
	}

	r.LoadControlLimitConstraintsData = UpdateList(r.LoadControlLimitConstraintsData, newData, filterPartial, filterDelete)
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
