package model

type SmartEnergyManagementPsAlternativesRelationType PowerSequenceAlternativesRelationDataType // ignoring the custom changes

type SmartEnergyManagementPsAlternativesRelationElementsType PowerSequenceAlternativesRelationDataElementsType // ignoring changes

type SmartEnergyManagementPsAlternativesType struct {
	Relation      *SmartEnergyManagementPsAlternativesRelationType `json:"relation,omitempty"`
	PowerSequence []SmartEnergyManagementPsPowerSequenceType       `json:"powerSequence,omitempty"`
}

type SmartEnergyManagementPsAlternativesElementsType struct {
	Relation      *SmartEnergyManagementPsAlternativesRelationElementsType `json:"relation,omitempty"`
	PowerSequence *SmartEnergyManagementPsPowerSequenceElementsType        `json:"powerSequence,omitempty"`
}

type SmartEnergyManagementPsPowerSequenceType struct {
	Description                           *PowerSequenceDescriptionDataType              `json:"description,omitempty"`                           // ignoring changes
	State                                 *PowerSequenceStateDataType                    `json:"state,omitempty"`                                 // ignoring changes
	Schedule                              *PowerSequenceScheduleDataType                 `json:"schedule,omitempty"`                              // ignoring changes
	ScheduleConstraints                   *PowerSequenceScheduleConstraintsDataType      `json:"scheduleConstraints,omitempty"`                   // ignoring changes
	SchedulePreference                    *PowerSequenceSchedulePreferenceDataType       `json:"schedulePreference,omitempty"`                    // ignoring changes
	OperatingConstraintsInterrupt         *OperatingConstraintsInterruptDataType         `json:"operatingConstraintsInterrupt,omitempty"`         // ignoring changes
	OperatingConstraintsDuration          *OperatingConstraintsDurationDataType          `json:"operatingConstraintsDuration,omitempty"`          // ignoring changes
	OperatingConstraintsResumeImplication *OperatingConstraintsResumeImplicationDataType `json:"operatingConstraintsResumeImplication,omitempty"` // ignoring changes
	PowerTimeSlot                         []SmartEnergyManagementPsPowerTimeSlotType     `json:"powerTimeSlot,omitempty"`                         // ignoring changes
}

type SmartEnergyManagementPsPowerSequenceElementsType struct {
	Description                           *PowerSequenceDescriptionDataElementsType              `json:"description,omitempty"`
	State                                 *PowerSequenceStateDataElementsType                    `json:"state,omitempty"`
	Schedule                              *PowerSequenceScheduleDataElementsType                 `json:"schedule,omitempty"`
	ScheduleConstraints                   *PowerSequenceScheduleConstraintsDataElementsType      `json:"scheduleConstraints,omitempty"`
	SchedulePreference                    *PowerSequenceSchedulePreferenceDataElementsType       `json:"schedulePreference,omitempty"`
	OperatingConstraintsInterrupt         *OperatingConstraintsInterruptDataElementsType         `json:"operatingConstraintsInterrupt,omitempty"`
	OperatingConstraintsDuration          *OperatingConstraintsDurationDataElementsType          `json:"operatingConstraintsDuration,omitempty"`
	OperatingConstraintsResumeImplication *OperatingConstraintsResumeImplicationDataElementsType `json:"operatingConstraintsResumeImplication,omitempty"`
	PowerTimeSlot                         *SmartEnergyManagementPsPowerTimeSlotElementsType      `json:"powerTimeSlot,omitempty"`
}

type SmartEnergyManagementPsPowerTimeSlotType struct {
	Schedule            *PowerTimeSlotScheduleDataType                     `json:"schedule,omitempty"` // ignoring changes
	ValueList           *SmartEnergyManagementPsPowerTimeSlotValueListType `json:"valueList,omitempty"`
	ScheduleConstraints *PowerTimeSlotScheduleConstraintsDataType          `json:"scheduleConstraints,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPowerTimeSlotElementsType struct {
	Schedule            *PowerTimeSlotScheduleDataElementsType                     `json:"schedule,omitempty"` // ignoring changes
	ValueList           *SmartEnergyManagementPsPowerTimeSlotValueListElementsType `json:"valueList,omitempty"`
	ScheduleConstraints *PowerTimeSlotScheduleConstraintsDataElementsType          `json:"scheduleConstraints,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPowerTimeSlotValueListType struct {
	Value []PowerTimeSlotValueDataType `json:"value,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPowerTimeSlotValueListElementsType struct {
	Value *PowerTimeSlotValueDataElementsType `json:"value,omitempty"`
}

type SmartEnergyManagementPsDataType struct {
	NodeScheduleInformation *PowerSequenceNodeScheduleInformationDataType `json:"nodeScheduleInformation,omitempty"` // ignoring changes
	Alternatives            []SmartEnergyManagementPsAlternativesType     `json:"alternatives,omitempty"`
}

type SmartEnergyManagementPsDataElementsType struct {
	NodeScheduleInformation *PowerSequenceNodeScheduleInformationDataElementsType `json:"nodeScheduleInformation,omitempty"`
	Alternatives            *SmartEnergyManagementPsAlternativesElementsType      `json:"alternatives,omitempty"`
}

type SmartEnergyManagementPsDataSelectorsType struct {
	AlternativesRelation     *PowerSequenceAlternativesRelationListDataSelectorsType `json:"alternativesRelation,omitempty"`     // ignoring changes
	PowerSequenceDescription *PowerSequenceDescriptionListDataSelectorsType          `json:"powerSequenceDescription,omitempty"` // ignoring changes
	PowerTimeSlotSchedule    *PowerTimeSlotScheduleListDataSelectorsType             `json:"powerTimeSlotSchedule,omitempty"`    // ignoring changes
	PowerTimeSlotValue       *PowerTimeSlotValueListDataSelectorsType                `json:"powerTimeSlotValue,omitempty"`       // ignoring changes
}

type SmartEnergyManagementPsPriceDataType struct {
	Price *PowerSequencePriceDataType `json:"price,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPriceDataElementsType struct {
	Price *PowerSequencePriceDataElementsType `json:"price,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPriceDataSelectorsType struct {
	Price *PowerSequencePriceListDataSelectorsType `json:"price,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsConfigurationRequestCallType struct {
	ScheduleConfigurationRequest *PowerSequenceScheduleConfigurationRequestCallType `json:"scheduleConfigurationRequest,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsConfigurationRequestCallElementsType struct {
	ScheduleConfigurationRequest *PowerSequenceScheduleConfigurationRequestCallElementsType `json:"scheduleConfigurationRequest,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPriceCalculationRequestCallType struct {
	PriceCalculationRequest *PowerSequencePriceCalculationRequestCallType `json:"priceCalculationRequest,omitempty"` // ignoring changes
}

type SmartEnergyManagementPsPriceCalculationRequestCallElementsType struct {
	PriceCalculationRequest *PowerSequencePriceCalculationRequestCallElementsType `json:"priceCalculationRequest,omitempty"` // ignoring changes
}
