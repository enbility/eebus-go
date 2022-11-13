package model

// MeasurementListDataType

var _ UpdaterFactory[MeasurementListDataType] = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) NewUpdater(
	newList *MeasurementListDataType,
	filterPartial *FilterType,
	filterDelete *FilterType) Updater {

	return &MeasurementListDataType_Updater{
		MeasurementListDataType: r,
		newData:                 newList.MeasurementData,
		FilterProvider: &FilterProvider{
			filterPartial: filterPartial,
			filterDelete:  filterDelete,
		},
	}
}

var _ Updater = (*MeasurementListDataType_Updater)(nil)

type MeasurementListDataType_Updater struct {
	*MeasurementListDataType
	*FilterProvider
	newData []MeasurementDataType
}

func (r *MeasurementListDataType_Updater) DoUpdate() {
	r.MeasurementData = UpdateList(r.MeasurementData, r.newData, r.filterPartial, r.filterDelete)
}

func (r *MeasurementListDataType_Updater) HasIdentifier(item *MeasurementDataType) bool {
	return item.MeasurementId != nil
}
