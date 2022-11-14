package model

// HvacSystemFunctionListDataType

var _ Updater = (*HvacSystemFunctionListDataType)(nil)

func (r *HvacSystemFunctionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionListDataType).HvacSystemFunctionData
	}

	r.HvacSystemFunctionData = UpdateList(r.HvacSystemFunctionData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionOperationModeRelationListDataType

var _ Updater = (*HvacSystemFunctionOperationModeRelationListDataType)(nil)

func (r *HvacSystemFunctionOperationModeRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionOperationModeRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionOperationModeRelationListDataType).HvacSystemFunctionOperationModeRelationData
	}

	r.HvacSystemFunctionOperationModeRelationData = UpdateList(r.HvacSystemFunctionOperationModeRelationData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionSetpointRelationListDataType

var _ Updater = (*HvacSystemFunctionSetpointRelationListDataType)(nil)

func (r *HvacSystemFunctionSetpointRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionSetpointRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionSetpointRelationListDataType).HvacSystemFunctionSetpointRelationData
	}

	r.HvacSystemFunctionSetpointRelationData = UpdateList(r.HvacSystemFunctionSetpointRelationData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionPowerSequenceRelationListDataType

var _ Updater = (*HvacSystemFunctionPowerSequenceRelationListDataType)(nil)

func (r *HvacSystemFunctionPowerSequenceRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionPowerSequenceRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionPowerSequenceRelationListDataType).HvacSystemFunctionPowerSequenceRelationData
	}

	r.HvacSystemFunctionPowerSequenceRelationData = UpdateList(r.HvacSystemFunctionPowerSequenceRelationData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionDescriptionListDataType

var _ Updater = (*HvacSystemFunctionDescriptionListDataType)(nil)

func (r *HvacSystemFunctionDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionDescriptionListDataType).HvacSystemFunctionDescriptionData
	}

	r.HvacSystemFunctionDescriptionData = UpdateList(r.HvacSystemFunctionDescriptionData, newData, filterPartial, filterDelete)
}

// HvacOperationModeDescriptionListDataType

var _ Updater = (*HvacOperationModeDescriptionListDataType)(nil)

func (r *HvacOperationModeDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacOperationModeDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOperationModeDescriptionListDataType).HvacOperationModeDescriptionData
	}

	r.HvacOperationModeDescriptionData = UpdateList(r.HvacOperationModeDescriptionData, newData, filterPartial, filterDelete)
}

// HvacOverrunListDataType

var _ Updater = (*HvacOverrunListDataType)(nil)

func (r *HvacOverrunListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacOverrunDataType
	if newList != nil {
		newData = newList.(*HvacOverrunListDataType).HvacOverrunData
	}

	r.HvacOverrunData = UpdateList(r.HvacOverrunData, newData, filterPartial, filterDelete)
}

// HvacOverrunDescriptionListDataType

var _ Updater = (*HvacOverrunDescriptionListDataType)(nil)

func (r *HvacOverrunDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacOverrunDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOverrunDescriptionListDataType).HvacOverrunDescriptionData
	}

	r.HvacOverrunDescriptionData = UpdateList(r.HvacOverrunDescriptionData, newData, filterPartial, filterDelete)
}
