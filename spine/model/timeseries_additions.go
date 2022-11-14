package model

// TimeSeriesListDataType

var _ Updater = (*TimeSeriesListDataType)(nil)

func (r *TimeSeriesListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeSeriesDataType
	if newList != nil {
		newData = newList.(*TimeSeriesListDataType).TimeSeriesData
	}

	r.TimeSeriesData = UpdateList(r.TimeSeriesData, newData, filterPartial, filterDelete)
}

// TimeSeriesDescriptionListDataType

var _ Updater = (*TimeSeriesDescriptionListDataType)(nil)

func (r *TimeSeriesDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeSeriesDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeSeriesDescriptionListDataType).TimeSeriesDescriptionData
	}

	r.TimeSeriesDescriptionData = UpdateList(r.TimeSeriesDescriptionData, newData, filterPartial, filterDelete)
}

// TimeSeriesConstraintsListDataType

var _ Updater = (*TimeSeriesConstraintsListDataType)(nil)

func (r *TimeSeriesConstraintsListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeSeriesConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeSeriesConstraintsListDataType).TimeSeriesConstraintsData
	}

	r.TimeSeriesConstraintsData = UpdateList(r.TimeSeriesConstraintsData, newData, filterPartial, filterDelete)
}
