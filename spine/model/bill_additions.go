package model

// BillListDataType

var _ Updater = (*BillListDataType)(nil)

func (r *BillListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []BillDataType
	if newList != nil {
		newData = newList.(*BillListDataType).BillData
	}

	r.BillData = UpdateList(r.BillData, newData, filterPartial, filterDelete)
}

// BillConstraintsListDataType

var _ Updater = (*BillConstraintsListDataType)(nil)

func (r *BillConstraintsListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []BillConstraintsDataType
	if newList != nil {
		newData = newList.(*BillConstraintsListDataType).BillConstraintsData
	}

	r.BillConstraintsData = UpdateList(r.BillConstraintsData, newData, filterPartial, filterDelete)
}

// BillDescriptionListDataType

var _ Updater = (*BillDescriptionListDataType)(nil)

func (r *BillDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []BillDescriptionDataType
	if newList != nil {
		newData = newList.(*BillDescriptionListDataType).BillDescriptionData
	}

	r.BillDescriptionData = UpdateList(r.BillDescriptionData, newData, filterPartial, filterDelete)
}
