package model

// DirectControlActivityListDataType

var _ Updater = (*DirectControlActivityListDataType)(nil)

func (r *DirectControlActivityListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []DirectControlActivityDataType
	if newList != nil {
		newData = newList.(*DirectControlActivityListDataType).DirectControlActivityDataElements
	}

	r.DirectControlActivityDataElements = UpdateList(r.DirectControlActivityDataElements, newData, filterPartial, filterDelete)
}
