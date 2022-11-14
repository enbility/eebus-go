package model

// TaskManagementJobListDataType

var _ Updater = (*TaskManagementJobListDataType)(nil)

func (r *TaskManagementJobListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TaskManagementJobDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobListDataType).TaskManagementJobData
	}

	r.TaskManagementJobData = UpdateList(r.TaskManagementJobData, newData, filterPartial, filterDelete)
}

// TaskManagementJobRelationListDataType

var _ Updater = (*TaskManagementJobRelationListDataType)(nil)

func (r *TaskManagementJobRelationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TaskManagementJobRelationDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobRelationListDataType).TaskManagementJobRelationData
	}

	r.TaskManagementJobRelationData = UpdateList(r.TaskManagementJobRelationData, newData, filterPartial, filterDelete)
}

// TaskManagementJobDescriptionListDataType

var _ Updater = (*TaskManagementJobDescriptionListDataType)(nil)

func (r *TaskManagementJobDescriptionListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []TaskManagementJobDescriptionDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobDescriptionListDataType).TaskManagementJobDescriptionData
	}

	r.TaskManagementJobDescriptionData = UpdateList(r.TaskManagementJobDescriptionData, newData, filterPartial, filterDelete)
}
