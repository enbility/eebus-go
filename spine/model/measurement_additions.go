package model

// MeasurementListDataType

var _ Updater = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementDataType
	if newList != nil {
		newData = newList.(*MeasurementListDataType).MeasurementData
	}

	r.MeasurementData = UpdateList(r.MeasurementData, newData, filterPartial, filterDelete)
}

// MeasurementSeriesListDataType

var _ Updater = (*MeasurementSeriesListDataType)(nil)

func (r *MeasurementSeriesListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementSeriesDataType
	if newList != nil {
		newData = newList.(*MeasurementSeriesListDataType).MeasurementSeriesData
	}

	r.MeasurementSeriesData = UpdateList(r.MeasurementSeriesData, newData, filterPartial, filterDelete)
}

// MeasurementConstraintsListDataType

var _ Updater = (*MeasurementConstraintsListDataType)(nil)

func (r *MeasurementConstraintsListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementConstraintsDataType
	if newList != nil {
		newData = newList.(*MeasurementConstraintsListDataType).MeasurementConstraintsData
	}

	r.MeasurementConstraintsData = UpdateList(r.MeasurementConstraintsData, newData, filterPartial, filterDelete)
}

// MeasurementDescriptionListDataType

var _ Updater = (*MeasurementDescriptionListDataType)(nil)

func (r *MeasurementDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementDescriptionDataType
	if newList != nil {
		newData = newList.(*MeasurementDescriptionListDataType).MeasurementDescriptionData
	}

	r.MeasurementDescriptionData = UpdateList(r.MeasurementDescriptionData, newData, filterPartial, filterDelete)
}

// MeasurementThresholdRelationListDataType

var _ Updater = (*MeasurementThresholdRelationListDataType)(nil)

func (r *MeasurementThresholdRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementThresholdRelationDataType
	if newList != nil {
		newData = newList.(*MeasurementThresholdRelationListDataType).MeasurementThresholdRelationData
	}

	r.MeasurementThresholdRelationData = UpdateList(r.MeasurementThresholdRelationData, newData, filterPartial, filterDelete)
}
