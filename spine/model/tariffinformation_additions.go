package model

// TariffListDataType

var _ Updater = (*TariffListDataType)(nil)

func (r *TariffListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffDataType
	if newList != nil {
		newData = newList.(*TariffListDataType).TariffData
	}

	r.TariffData = UpdateList(r.TariffData, newData, filterPartial, filterDelete)
}

// TariffTierRelationListDataType

var _ Updater = (*TariffTierRelationListDataType)(nil)

func (r *TariffTierRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffTierRelationDataType
	if newList != nil {
		newData = newList.(*TariffTierRelationListDataType).TariffTierRelationData
	}

	r.TariffTierRelationData = UpdateList(r.TariffTierRelationData, newData, filterPartial, filterDelete)
}

// TariffBoundaryRelationListDataType

var _ Updater = (*TariffBoundaryRelationListDataType)(nil)

func (r *TariffBoundaryRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffBoundaryRelationDataType
	if newList != nil {
		newData = newList.(*TariffBoundaryRelationListDataType).TariffBoundaryRelationData
	}

	r.TariffBoundaryRelationData = UpdateList(r.TariffBoundaryRelationData, newData, filterPartial, filterDelete)
}

// TariffDescriptionListDataType

var _ Updater = (*TariffDescriptionListDataType)(nil)

func (r *TariffDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffDescriptionDataType
	if newList != nil {
		newData = newList.(*TariffDescriptionListDataType).TariffDescriptionData
	}

	r.TariffDescriptionData = UpdateList(r.TariffDescriptionData, newData, filterPartial, filterDelete)
}

// TierBoundaryListDataType

var _ Updater = (*TierBoundaryListDataType)(nil)

func (r *TierBoundaryListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierBoundaryDataType
	if newList != nil {
		newData = newList.(*TierBoundaryListDataType).TierBoundaryData
	}

	r.TierBoundaryData = UpdateList(r.TierBoundaryData, newData, filterPartial, filterDelete)
}

// TierBoundaryDescriptionListDataType

var _ Updater = (*TierBoundaryDescriptionListDataType)(nil)

func (r *TierBoundaryDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierBoundaryDescriptionDataType
	if newList != nil {
		newData = newList.(*TierBoundaryDescriptionListDataType).TierBoundaryDescriptionData
	}

	r.TierBoundaryDescriptionData = UpdateList(r.TierBoundaryDescriptionData, newData, filterPartial, filterDelete)
}

// CommodityListDataType

var _ Updater = (*CommodityListDataType)(nil)

func (r *CommodityListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []CommodityDataType
	if newList != nil {
		newData = newList.(*CommodityListDataType).CommodityData
	}

	r.CommodityData = UpdateList(r.CommodityData, newData, filterPartial, filterDelete)
}

// TierListDataType

var _ Updater = (*TierListDataType)(nil)

func (r *TierListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierDataType
	if newList != nil {
		newData = newList.(*TierListDataType).TierData
	}

	r.TierData = UpdateList(r.TierData, newData, filterPartial, filterDelete)
}

// TierIncentiveRelationListDataType

var _ Updater = (*TierIncentiveRelationListDataType)(nil)

func (r *TierIncentiveRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierIncentiveRelationDataType
	if newList != nil {
		newData = newList.(*TierIncentiveRelationListDataType).TierIncentiveRelationData
	}

	r.TierIncentiveRelationData = UpdateList(r.TierIncentiveRelationData, newData, filterPartial, filterDelete)
}

// TierDescriptionListDataType

var _ Updater = (*TierDescriptionListDataType)(nil)

func (r *TierDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierDescriptionDataType
	if newList != nil {
		newData = newList.(*TierDescriptionListDataType).TierDescriptionData
	}

	r.TierDescriptionData = UpdateList(r.TierDescriptionData, newData, filterPartial, filterDelete)
}

// IncentiveListDataType

var _ Updater = (*IncentiveListDataType)(nil)

func (r *IncentiveListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []IncentiveDataType
	if newList != nil {
		newData = newList.(*IncentiveListDataType).IncentiveData
	}

	r.IncentiveData = UpdateList(r.IncentiveData, newData, filterPartial, filterDelete)
}

// IncentiveDescriptionListDataType

var _ Updater = (*IncentiveDescriptionListDataType)(nil)

func (r *IncentiveDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []IncentiveDescriptionDataType
	if newList != nil {
		newData = newList.(*IncentiveDescriptionListDataType).IncentiveDescriptionData
	}

	r.IncentiveDescriptionData = UpdateList(r.IncentiveDescriptionData, newData, filterPartial, filterDelete)
}
