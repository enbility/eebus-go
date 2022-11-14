package model

// UseCaseInformationListDataType

var _ Updater = (*UseCaseInformationListDataType)(nil)

func (r *UseCaseInformationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []UseCaseInformationDataType
	if newList != nil {
		newData = newList.(*UseCaseInformationListDataType).UseCaseInformationData
	}

	r.UseCaseInformationData = UpdateList(r.UseCaseInformationData, newData, filterPartial, filterDelete)
}
