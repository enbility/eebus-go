package model

// SpecificationVersionListDataType

var _ Updater = (*SpecificationVersionListDataType)(nil)

func (r *SpecificationVersionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []SpecificationVersionDataType
	if newList != nil {
		newData = newList.(*SpecificationVersionListDataType).SpecificationVersionData
	}

	r.SpecificationVersionData = UpdateList(r.SpecificationVersionData, newData, filterPartial, filterDelete)
}
