package model

// NodeManagementDestinationListDataType

var _ Updater = (*NodeManagementDestinationListDataType)(nil)

func (r *NodeManagementDestinationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []NodeManagementDestinationDataType
	if newList != nil {
		newData = newList.(*NodeManagementDestinationListDataType).NodeManagementDestinationData
	}

	r.NodeManagementDestinationData = UpdateList(r.NodeManagementDestinationData, newData, filterPartial, filterDelete)
}
