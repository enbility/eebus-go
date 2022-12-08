package model

// ThresholdListDataType

var _ Updater = (*ThresholdListDataType)(nil)

func (r *ThresholdListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ThresholdDataType
	if newList != nil {
		newData = newList.(*ThresholdListDataType).ThresholdData
	}

	r.ThresholdData = UpdateList(r.ThresholdData, newData, filterPartial, filterDelete)
}

// ThresholdConstraintsListDataType

var _ Updater = (*ThresholdConstraintsListDataType)(nil)

func (r *ThresholdConstraintsListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ThresholdConstraintsDataType
	if newList != nil {
		newData = newList.(*ThresholdConstraintsListDataType).ThresholdConstraintsData
	}

	r.ThresholdConstraintsData = UpdateList(r.ThresholdConstraintsData, newData, filterPartial, filterDelete)
}

// ThresholdDescriptionListDataType

var _ Updater = (*ThresholdDescriptionListDataType)(nil)

func (r *ThresholdDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []ThresholdDescriptionDataType
	if newList != nil {
		newData = newList.(*ThresholdDescriptionListDataType).ThresholdDescriptionData
	}

	r.ThresholdDescriptionData = UpdateList(r.ThresholdDescriptionData, newData, filterPartial, filterDelete)
}
