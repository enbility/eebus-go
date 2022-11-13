package model

// LoadControlLimitListDataType

var _ UpdaterFactory[LoadControlLimitListDataType] = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) NewUpdater(
	newList *LoadControlLimitListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &LoadControlLimitListDataType_Updater{
		LoadControlLimitListDataType: r,
		newData:                      newList.LoadControlLimitData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

var _ Updater = (*LoadControlLimitListDataType_Updater)(nil)

type LoadControlLimitListDataType_Updater struct {
	*LoadControlLimitListDataType
	*FilterProvider
	newData []LoadControlLimitDataType
}

func (r *LoadControlLimitListDataType_Updater) DoUpdate() {
	r.LoadControlLimitData = UpdateList(r.LoadControlLimitData, r.newData, r.filterPartial, r.filterDelete)
}

func (r *LoadControlLimitListDataType_Updater) HasIdentifier(item *LoadControlLimitDataType) bool {
	return item.LimitId != nil
}

// LoadControlLimitDescriptionListDataType

var _ UpdaterFactory[LoadControlLimitDescriptionListDataType] = (*LoadControlLimitDescriptionListDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) NewUpdater(
	newList *LoadControlLimitDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &LoadControlLimitDescriptionListDataType_Updater{
		LoadControlLimitDescriptionListDataType: r,
		newData:                                 newList.LoadControlLimitDescriptionData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

var _ Updater = (*LoadControlLimitDescriptionListDataType_Updater)(nil)

type LoadControlLimitDescriptionListDataType_Updater struct {
	*LoadControlLimitDescriptionListDataType
	*FilterProvider
	newData []LoadControlLimitDescriptionDataType
}

func (r *LoadControlLimitDescriptionListDataType_Updater) DoUpdate() {
	r.LoadControlLimitDescriptionData = UpdateList(r.LoadControlLimitDescriptionData, r.newData, r.filterPartial, r.filterDelete)
}

func (r *LoadControlLimitDescriptionListDataType_Updater) HasIdentifier(item *LoadControlLimitDescriptionDataType) bool {
	return item.LimitId != nil
}
