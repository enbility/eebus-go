package model

// ElectricalConnectionPermittedValueSetListDataType

var _ UpdaterFactory[ElectricalConnectionPermittedValueSetListDataType] = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) NewUpdater(
	newList *ElectricalConnectionPermittedValueSetListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	var newData []ElectricalConnectionPermittedValueSetDataType
	if newList != nil {
		newData = newList.ElectricalConnectionPermittedValueSetData
	}

	e := &ElectricalConnectionPermittedValueSetListDataType_Updater{
		ElectricalConnectionPermittedValueSetListDataType: r,
		newData: newData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}

	return e
}

type ElectricalConnectionPermittedValueSetListDataType_Updater struct {
	*ElectricalConnectionPermittedValueSetListDataType
	*FilterProvider
	newData []ElectricalConnectionPermittedValueSetDataType
}

func (r *ElectricalConnectionPermittedValueSetListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionPermittedValueSetData = UpdateList(r.ElectricalConnectionPermittedValueSetData, r.newData, r.filterPartial, r.filterDelete)
}

// ElectricalConnectionDescriptionListDataType

var _ UpdaterFactory[ElectricalConnectionDescriptionListDataType] = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) NewUpdater(
	newList *ElectricalConnectionDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionDescriptionListDataType_Updater{
		ElectricalConnectionDescriptionListDataType: r,
		newData: newList.ElectricalConnectionDescriptionData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

var _ Updater = (*ElectricalConnectionDescriptionListDataType_Updater)(nil)

type ElectricalConnectionDescriptionListDataType_Updater struct {
	*ElectricalConnectionDescriptionListDataType
	*FilterProvider
	newData []ElectricalConnectionDescriptionDataType
}

func (r *ElectricalConnectionDescriptionListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionDescriptionData = UpdateList(r.ElectricalConnectionDescriptionData, r.newData, r.filterPartial, r.filterDelete)
}

// ElectricalConnectionParameterDescriptionListDataType

var _ UpdaterFactory[ElectricalConnectionParameterDescriptionListDataType] = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) NewUpdater(
	newList *ElectricalConnectionParameterDescriptionListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &ElectricalConnectionParameterDescriptionListDataType_Updater{
		ElectricalConnectionParameterDescriptionListDataType: r,
		newData: newList.ElectricalConnectionParameterDescriptionData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType_Updater)(nil)

type ElectricalConnectionParameterDescriptionListDataType_Updater struct {
	*ElectricalConnectionParameterDescriptionListDataType
	*FilterProvider
	newData []ElectricalConnectionParameterDescriptionDataType
}

func (r *ElectricalConnectionParameterDescriptionListDataType_Updater) DoUpdate() {
	r.ElectricalConnectionParameterDescriptionData = UpdateList(r.ElectricalConnectionParameterDescriptionData, r.newData, r.filterPartial, r.filterDelete)
}
