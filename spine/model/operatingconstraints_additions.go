package model

// OperatingConstraintsInterruptListDataType

var _ Updater = (*OperatingConstraintsInterruptListDataType)(nil)

func (r *OperatingConstraintsInterruptListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsInterruptDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsInterruptListDataType).OperatingConstraintsInterruptData
	}

	r.OperatingConstraintsInterruptData = UpdateList(r.OperatingConstraintsInterruptData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsDurationListDataType

var _ Updater = (*OperatingConstraintsDurationListDataType)(nil)

func (r *OperatingConstraintsDurationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsDurationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsDurationListDataType).OperatingConstraintsDurationData
	}

	r.OperatingConstraintsDurationData = UpdateList(r.OperatingConstraintsDurationData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsPowerDescriptionListDataType

var _ Updater = (*OperatingConstraintsPowerDescriptionListDataType)(nil)

func (r *OperatingConstraintsPowerDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsPowerDescriptionDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerDescriptionListDataType).OperatingConstraintsPowerDescriptionData
	}

	r.OperatingConstraintsPowerDescriptionData = UpdateList(r.OperatingConstraintsPowerDescriptionData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsPowerRangeListDataType

var _ Updater = (*OperatingConstraintsPowerRangeListDataType)(nil)

func (r *OperatingConstraintsPowerRangeListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsPowerRangeDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerRangeListDataType).OperatingConstraintsPowerRangeData
	}

	r.OperatingConstraintsPowerRangeData = UpdateList(r.OperatingConstraintsPowerRangeData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsPowerLevelListDataType

var _ Updater = (*OperatingConstraintsPowerLevelListDataType)(nil)

func (r *OperatingConstraintsPowerLevelListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsPowerLevelDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerLevelListDataType).OperatingConstraintsPowerLevelData
	}

	r.OperatingConstraintsPowerLevelData = UpdateList(r.OperatingConstraintsPowerLevelData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsResumeImplicationListDataType

var _ Updater = (*OperatingConstraintsResumeImplicationListDataType)(nil)

func (r *OperatingConstraintsResumeImplicationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsResumeImplicationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsResumeImplicationListDataType).OperatingConstraintsResumeImplicationData
	}

	r.OperatingConstraintsResumeImplicationData = UpdateList(r.OperatingConstraintsResumeImplicationData, newData, filterPartial, filterDelete)
}
