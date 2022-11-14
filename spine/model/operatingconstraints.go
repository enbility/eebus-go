package model

type OperatingConstraintsInterruptDataType struct {
	SequenceId                  *PowerSequenceIdType `json:"sequenceId,omitempty"`
	IsPausable                  *bool                `json:"isPausable,omitempty"`
	IsStoppable                 *bool                `json:"isStoppable,omitempty"`
	NotInterruptibleAtHighPower *bool                `json:"notInterruptibleAtHighPower,omitempty"`
	MaxCyclesPerDay             *uint                `json:"maxCyclesPerDay,omitempty"`
}

type OperatingConstraintsInterruptDataElementsType struct {
	SequenceId                  *ElementTagType `json:"sequenceId,omitempty"`
	IsPausable                  *ElementTagType `json:"isPausable,omitempty"`
	IsStoppable                 *ElementTagType `json:"isStoppable,omitempty"`
	NotInterruptibleAtHighPower *ElementTagType `json:"notInterruptibleAtHighPower,omitempty"`
	MaxCyclesPerDay             *ElementTagType `json:"maxCyclesPerDay,omitempty"`
}

type OperatingConstraintsInterruptListDataType struct {
	OperatingConstraintsInterruptData []OperatingConstraintsInterruptDataType `json:"operatingConstraintsInterruptData,omitempty"`
}

type OperatingConstraintsInterruptListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type OperatingConstraintsDurationDataType struct {
	SequenceId           *PowerSequenceIdType `json:"sequenceId,omitempty" eebus:"key"`
	ActiveDurationMin    *DurationType        `json:"activeDurationMin,omitempty"`
	ActiveDurationMax    *DurationType        `json:"activeDurationMax,omitempty"`
	PauseDurationMin     *DurationType        `json:"pauseDurationMin,omitempty"`
	PauseDurationMax     *DurationType        `json:"pauseDurationMax,omitempty"`
	ActiveDurationSumMin *DurationType        `json:"activeDurationSumMin,omitempty"`
	ActiveDurationSumMax *DurationType        `json:"activeDurationSumMax,omitempty"`
}

type OperatingConstraintsDurationDataElementsType struct {
	SequenceId           *ElementTagType `json:"sequenceId,omitempty"`
	ActiveDurationMin    *ElementTagType `json:"activeDurationMin,omitempty"`
	ActiveDurationMax    *ElementTagType `json:"activeDurationMax,omitempty"`
	PauseDurationMin     *ElementTagType `json:"pauseDurationMin,omitempty"`
	PauseDurationMax     *ElementTagType `json:"pauseDurationMax,omitempty"`
	ActiveDurationSumMin *ElementTagType `json:"activeDurationSumMin,omitempty"`
	ActiveDurationSumMax *ElementTagType `json:"activeDurationSumMax,omitempty"`
}

type OperatingConstraintsDurationListDataType struct {
	OperatingConstraintsDurationData []OperatingConstraintsDurationDataType `json:"operatingConstraintsDurationData,omitempty"`
}

type OperatingConstraintsDurationListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type OperatingConstraintsPowerDescriptionDataType struct {
	SequenceId              *PowerSequenceIdType   `json:"sequenceId,omitempty" eebus:"key"`
	PositiveEnergyDirection *EnergyDirectionType   `json:"positiveEnergyDirection,omitempty"`
	PowerUnit               *UnitOfMeasurementType `json:"powerUnit,omitempty"`
	EnergyUnit              *UnitOfMeasurementType `json:"energyUnit,omitempty"`
	Description             *DescriptionType       `json:"description,omitempty"`
}

type OperatingConstraintsPowerDescriptionDataElementsType struct {
	SequenceId              *ElementTagType `json:"sequenceId,omitempty"`
	PositiveEnergyDirection *ElementTagType `json:"positiveEnergyDirection,omitempty"`
	PowerUnit               *ElementTagType `json:"powerUnit,omitempty"`
	EnergyUnit              *ElementTagType `json:"energyUnit,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
}

type OperatingConstraintsPowerDescriptionListDataType struct {
	OperatingConstraintsPowerDescriptionData []OperatingConstraintsPowerDescriptionDataType `json:"operatingConstraintsPowerDescriptionData,omitempty"`
}

type OperatingConstraintsPowerDescriptionListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type OperatingConstraintsPowerRangeDataType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty" eebus:"key"`
	PowerMin   *ScaledNumberType    `json:"powerMin,omitempty"`
	PowerMax   *ScaledNumberType    `json:"powerMax,omitempty"`
	EnergyMin  *ScaledNumberType    `json:"energyMin,omitempty"`
	EnergyMax  *ScaledNumberType    `json:"energyMax,omitempty"`
}

type OperatingConstraintsPowerRangeDataElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
	PowerMin   *ElementTagType `json:"powerMin,omitempty"`
	PowerMax   *ElementTagType `json:"powerMax,omitempty"`
	EnergyMin  *ElementTagType `json:"energyMin,omitempty"`
	EnergyMax  *ElementTagType `json:"energyMax,omitempty"`
}

type OperatingConstraintsPowerRangeListDataType struct {
	OperatingConstraintsPowerRangeData []OperatingConstraintsPowerRangeDataType `json:"operatingConstraintsPowerRangeData,omitempty"`
}

type OperatingConstraintsPowerRangeListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type OperatingConstraintsPowerLevelDataType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty" eebus:"key"`
	Power      *ScaledNumberType    `json:"power,omitempty"`
}

type OperatingConstraintsPowerLevelDataElementsType struct {
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
	Power      *ElementTagType `json:"power,omitempty"`
}

type OperatingConstraintsPowerLevelListDataType struct {
	OperatingConstraintsPowerLevelData []OperatingConstraintsPowerLevelDataType `json:"operatingConstraintsPowerLevelData,omitempty"`
}

type OperatingConstraintsPowerLevelListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}

type OperatingConstraintsResumeImplicationDataType struct {
	SequenceId            *PowerSequenceIdType   `json:"sequenceId,omitempty" eebus:"key"`
	ResumeEnergyEstimated *ScaledNumberType      `json:"resumeEnergyEstimated,omitempty"`
	EnergyUnit            *UnitOfMeasurementType `json:"energyUnit,omitempty"`
	ResumeCostEstimated   *ScaledNumberType      `json:"resumeCostEstimated,omitempty"`
	Currency              *CurrencyType          `json:"currency,omitempty"`
}

type OperatingConstraintsResumeImplicationDataElementsType struct {
	SequenceId            *ElementTagType           `json:"sequenceId,omitempty"`
	ResumeEnergyEstimated *ScaledNumberElementsType `json:"resumeEnergyEstimated,omitempty"`
	EnergyUnit            *ElementTagType           `json:"energyUnit,omitempty"`
	ResumeCostEstimated   *ScaledNumberElementsType `json:"resumeCostEstimated,omitempty"`
	Currency              *ElementTagType           `json:"currency,omitempty"`
}

type OperatingConstraintsResumeImplicationListDataType struct {
	OperatingConstraintsResumeImplicationData []OperatingConstraintsResumeImplicationDataType `json:"operatingConstraintsResumeImplicationData,omitempty"`
}

type OperatingConstraintsResumeImplicationListDataSelectorsType struct {
	SequenceId *PowerSequenceIdType `json:"sequenceId,omitempty"`
}
