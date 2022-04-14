package spine

type TaskManagementJobIdType uint

type TaskManagementJobStateType string

const (
// DirectControlActivityStateEnumType
// HvacOverrunStatusEnumType
// LoadControlEventStateEnumType
// PowerSequenceStateEnumType
)

type TaskManagementJobSourceType TaskManagementJobSourceEnumType

type TaskManagementJobSourceEnumType string

const (
	TaskManagementJobSourceEnumTypeInternalMechanism     TaskManagementJobSourceEnumType = "InternalMechanism"
	TaskManagementJobSourceEnumTypeUserInteraction       TaskManagementJobSourceEnumType = "UserInteraction"
	TaskManagementJobSourceEnumTypeExternalConfiguration TaskManagementJobSourceEnumType = "ExternalConfiguration"
)

type TaskManagementDirectControlRelatedType struct{}

type TaskManagementDirectControlRelatedElementsType struct{}

type TaskManagementHvacRelatedType struct {
	OverrunId *HvacOverrunIdType `json:"overrunId,omitempty"`
}

type TaskManagementHvacRelatedElementsType struct {
	OverrunId *ElementTagType `json:"overrunId,omitempty"`
}

type TaskManagementLoadControlReleatedType struct {
	EventId *LoadControlEventIdType `json:"eventId,omitempty"`
}

type TaskManagementLoadControlReleatedElementsType struct {
	EventId *ElementTagType `json:"eventId,omitempty"`
}

type TaskManagementPowerSequencesRelatedType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type TaskManagementPowerSequencesRelatedElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
}

type TaskManagementSmartEnergyManagementPsRelatedType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type TaskManagementSmartEnergyManagementPsRelatedElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
}

type TaskManagementJobDataType struct {
	JobId         *TaskManagementJobIdType    `json:"jobId,omitempty"`
	Timestamp     *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	JobState      *TaskManagementJobStateType `json:"jobState,omitempty"`
	ElapsedTime   *string                     `json:"elapsedTime,omitempty"`
	RemainingTime *string                     `json:"remainingTime,omitempty"`
}

type TaskManagementJobDataElementsType struct {
	JobId         *TaskManagementJobIdType    `json:"jobId,omitempty"`
	Timestamp     *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	JobState      *TaskManagementJobStateType `json:"jobState,omitempty"`
	ElapsedTime   *string                     `json:"elapsedTime,omitempty"`
	RemainingTime *string                     `json:"remainingTime,omitempty"`
}

type TaskManagementJobListDataType struct {
	TaskManagementJobData []TaskManagementJobDataType `json:"taskManagementJobData,omitempty"`
}

type TaskManagementJobListDataSelectorsType struct {
	JobId    *TaskManagementJobIdType    `json:"jobId,omitempty"`
	JobState *TaskManagementJobStateType `json:"jobState,omitempty"`
}

type TaskManagementJobRelationDataType struct {
	JobId                          *TaskManagementJobIdType                          `json:"jobId,omitempty"`
	DirectControlRelated           *TaskManagementDirectControlRelatedType           `json:"directControlRelated,omitempty"`
	HvacRelated                    *TaskManagementHvacRelatedType                    `json:"hvacRelated,omitempty"`
	LoadControlReleated            *TaskManagementLoadControlReleatedType            `json:"loadControlReleated,omitempty"`
	PowerSequencesRelated          *TaskManagementPowerSequencesRelatedType          `json:"powerSequencesRelated,omitempty"`
	SmartEnergyManagementPsRelated *TaskManagementSmartEnergyManagementPsRelatedType `json:"smartEnergyManagementPsRelated,omitempty"`
}

type TaskManagementJobRelationDataElementsType struct {
	JobId                          *ElementTagType                                           `json:"jobId,omitempty"`
	DirectControlRelated           *TaskManagementDirectControlRelatedElementsType           `json:"directControlRelated,omitempty"`
	HvacRelated                    *TaskManagementHvacRelatedElementsType                    `json:"hvacRelated,omitempty"`
	LoadControlReleated            *TaskManagementLoadControlReleatedElementsType            `json:"loadControlReleated,omitempty"`
	PowerSequencesRelated          *TaskManagementPowerSequencesRelatedElementsType          `json:"powerSequencesRelated,omitempty"`
	SmartEnergyManagementPsRelated *TaskManagementSmartEnergyManagementPsRelatedElementsType `json:"smartEnergyManagementPsRelated,omitempty"`
}

type TaskManagementJobRelationListDataType struct {
	TaskManagementJobRelationData []TaskManagementJobRelationDataType `json:"taskManagementJobRelationData,omitempty"`
}

type TaskManagementJobRelationListDataSelectorsType struct {
	JobId *TaskManagementJobIdType `json:"jobId,omitempty"`
}

type TaskManagementJobDescriptionDataType struct {
	JobId       *TaskManagementJobIdType     `json:"jobId,omitempty"`
	JobSource   *TaskManagementJobSourceType `json:"jobSource,omitempty"`
	Label       *LabelType                   `json:"label,omitempty"`
	Description *DescriptionType             `json:"description,omitempty"`
}

type TaskManagementJobDescriptionDataElementsType struct {
	JobId       *ElementTagType `json:"jobId,omitempty"`
	JobSource   *ElementTagType `json:"jobSource,omitempty"`
	Label       *ElementTagType `json:"label,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}

type TaskManagementJobDescriptionListDataType struct {
	TaskManagementJobDescriptionData []TaskManagementJobDescriptionDataType `json:"taskManagementJobDescriptionData,omitempty"`
}

type TaskManagementJobDescriptionListDataSelectorsType struct {
	JobId     *TaskManagementJobIdType     `json:"jobId,omitempty"`
	JobSource *TaskManagementJobSourceType `json:"jobSource,omitempty"`
}

type TaskManagementOverviewDataType struct {
	RemoteControllable *bool `json:"remoteControllable,omitempty"`
	JobsActive         *bool `json:"jobsActive,omitempty"`
}

type TaskManagementOverviewDataElementsType struct {
	RemoteControllable *ElementTagType `json:"remoteControllable,omitempty"`
	JobsActive         *ElementTagType `json:"jobsActive,omitempty"`
}
