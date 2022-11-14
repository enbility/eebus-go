package model

// NetworkManagementDeviceDescriptionListDataType

var _ Updater = (*NetworkManagementDeviceDescriptionListDataType)(nil)

func (r *NetworkManagementDeviceDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []NetworkManagementDeviceDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementDeviceDescriptionListDataType).NetworkManagementDeviceDescriptionData
	}

	r.NetworkManagementDeviceDescriptionData = UpdateList(r.NetworkManagementDeviceDescriptionData, newData, filterPartial, filterDelete)
}

// NetworkManagementEntityDescriptionListDataType

var _ Updater = (*NetworkManagementEntityDescriptionListDataType)(nil)

func (r *NetworkManagementEntityDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []NetworkManagementEntityDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementEntityDescriptionListDataType).NetworkManagementEntityDescriptionData
	}

	r.NetworkManagementEntityDescriptionData = UpdateList(r.NetworkManagementEntityDescriptionData, newData, filterPartial, filterDelete)
}

// NetworkManagementFeatureDescriptionListDataType

var _ Updater = (*NetworkManagementFeatureDescriptionListDataType)(nil)

func (r *NetworkManagementFeatureDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []NetworkManagementFeatureDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementFeatureDescriptionListDataType).NetworkManagementFeatureDescriptionData
	}

	r.NetworkManagementFeatureDescriptionData = UpdateList(r.NetworkManagementFeatureDescriptionData, newData, filterPartial, filterDelete)
}
