package model

// SupplyConditionListDataType

var _ Updater = (*SupplyConditionListDataType)(nil)

func (r *SupplyConditionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SupplyConditionDataType
	if newList != nil {
		newData = newList.(*SupplyConditionListDataType).SupplyConditionData
	}

	r.SupplyConditionData = UpdateList(r.SupplyConditionData, newData, filterPartial, filterDelete)
}

// SupplyConditionDescriptionListDataType

var _ Updater = (*SupplyConditionDescriptionListDataType)(nil)

func (r *SupplyConditionDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SupplyConditionDescriptionDataType
	if newList != nil {
		newData = newList.(*SupplyConditionDescriptionListDataType).SupplyConditionDescriptionData
	}

	r.SupplyConditionDescriptionData = UpdateList(r.SupplyConditionDescriptionData, newData, filterPartial, filterDelete)
}

// SupplyConditionThresholdRelationListDataType

var _ Updater = (*SupplyConditionThresholdRelationListDataType)(nil)

func (r *SupplyConditionThresholdRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SupplyConditionThresholdRelationDataType
	if newList != nil {
		newData = newList.(*SupplyConditionThresholdRelationListDataType).SupplyConditionThresholdRelationData
	}

	r.SupplyConditionThresholdRelationData = UpdateList(r.SupplyConditionThresholdRelationData, newData, filterPartial, filterDelete)
}
