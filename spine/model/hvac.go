package model

type HvacSystemFunctionIdType uint

type HvacSystemFunctionTypeType string

const (
	HvacSystemFunctionTypeTypeHeating     HvacSystemFunctionTypeType = "heating"
	HvacSystemFunctionTypeTypeCooling     HvacSystemFunctionTypeType = "cooling"
	HvacSystemFunctionTypeTypeVentilation HvacSystemFunctionTypeType = "ventilation"
	HvacSystemFunctionTypeTypeDhw         HvacSystemFunctionTypeType = "dhw"
)

type HvacOperationModeIdType uint

type HvacOperationModeTypeType string

const (
	HvacOperationModeTypeTypeAuto HvacOperationModeTypeType = "auto"
	HvacOperationModeTypeTypeOn   HvacOperationModeTypeType = "on"
	HvacOperationModeTypeTypeOff  HvacOperationModeTypeType = "off"
	HvacOperationModeTypeTypeEco  HvacOperationModeTypeType = "eco"
)

type HvacOverrunIdType uint

type HvacOverrunTypeType string

const (
	HvacOverrunTypeTypeOneTimeDhw         HvacOverrunTypeType = "oneTimeDhw"
	HvacOverrunTypeTypeParty              HvacOverrunTypeType = "party"
	HvacOverrunTypeTypeSgReadyCondition1  HvacOverrunTypeType = "sgReadyCondition1"
	HvacOverrunTypeTypeSgReadyCondition3  HvacOverrunTypeType = "sgReadyCondition3"
	HvacOverrunTypeTypeSgReadyCondition4  HvacOverrunTypeType = "sgReadyCondition4"
	HvacOverrunTypeTypeOneDayAway         HvacOverrunTypeType = "oneDayAway"
	HvacOverrunTypeTypeOneDayAtHome       HvacOverrunTypeType = "oneDayAtHome"
	HvacOverrunTypeTypeOneTimeVentilation HvacOverrunTypeType = "oneTimeVentilation"
	HvacOverrunTypeTypeHvacSystemOff      HvacOverrunTypeType = "hvacSystemOff"
	HvacOverrunTypeTypeValveKick          HvacOverrunTypeType = "valveKick"
)

type HvacOverrunStatusType string

const (
	HvacOverrunStatusTypeActive   HvacOverrunStatusType = "active"
	HvacOverrunStatusTypeRunning  HvacOverrunStatusType = "running"
	HvacOverrunStatusTypeFinished HvacOverrunStatusType = "finished"
	HvacOverrunStatusTypeInactive HvacOverrunStatusType = "inactive"
)

type HvacSystemFunctionDataType struct {
	SystemFunctionId            *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty" eebus:"key"`
	CurrentOperationModeId      *HvacOperationModeIdType  `json:"currentOperationModeId,omitempty" eebus:"key"`
	IsOperationModeIdChangeable *bool                     `json:"isOperationModeIdChangeable,omitempty"`
	CurrentSetpointId           *SetpointIdType           `json:"currentSetpointId,omitempty" eebus:"key"`
	IsSetpointIdChangeable      *bool                     `json:"isSetpointIdChangeable,omitempty"`
	IsOverrunActive             *bool                     `json:"isOverrunActive,omitempty"`
}

type HvacSystemFunctionDataElementsType struct {
	SystemFunctionId            *ElementTagType `json:"systemFunctionId,omitempty"`
	CurrentOperationModeId      *ElementTagType `json:"currentOperationModeId,omitempty"`
	IsOperationModeIdChangeable *ElementTagType `json:"isOperationModeIdChangeable,omitempty"`
	CurrentSetpointId           *ElementTagType `json:"currentSetpointId,omitempty"`
	IsSetpointIdChangeable      *ElementTagType `json:"isSetpointIdChangeable,omitempty"`
	IsOverrunActive             *ElementTagType `json:"isOverrunActive,omitempty"`
}

type HvacSystemFunctionListDataType struct {
	HvacSystemFunctionData []HvacSystemFunctionDataType `json:"hvacSystemFunctionData,omitempty"`
}

type HvacSystemFunctionListDataSelectorsType struct {
	SystemFunctionId []HvacSystemFunctionIdType `json:"systemFunctionId,omitempty"`
}

type HvacSystemFunctionOperationModeRelationDataType struct {
	SystemFunctionId *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty" eebus:"key"`
	OperationModeId  *HvacOperationModeIdType  `json:"operationModeId,omitempty" eebus:"key"`
}

type HvacSystemFunctionOperationModeRelationDataElementsType struct {
	SystemFunctionId *ElementTagType `json:"systemFunctionId,omitempty"`
	OperationModeId  *ElementTagType `json:"operationModeId,omitempty"`
}

type HvacSystemFunctionOperationModeRelationListDataType struct {
	HvacSystemFunctionOperationModeRelationData []HvacSystemFunctionOperationModeRelationDataType `json:"hvacSystemFunctionOperationModeRelationData,omitempty"`
}

type HvacSystemFunctionOperationModeRelationListDataSelectorsType struct {
	SystemFunctionId []HvacSystemFunctionIdType `json:"systemFunctionId,omitempty"`
}

type HvacSystemFunctionSetpointRelationDataType struct {
	SystemFunctionId *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty" eebus:"key"`
	OperationModeId  *HvacOperationModeIdType  `json:"operationModeId,omitempty" eebus:"key"`
	SetpointId       *SetpointIdType           `json:"setpointId,omitempty" eebus:"key"`
}

type HvacSystemFunctionSetpointRelationDataElementsType struct {
	SystemFunctionId *ElementTagType `json:"systemFunctionId,omitempty"`
	OperationModeId  *ElementTagType `json:"operationModeId,omitempty"`
	SetpointId       *ElementTagType `json:"setpointId,omitempty"`
}

type HvacSystemFunctionSetpointRelationListDataType struct {
	HvacSystemFunctionSetpointRelationData []HvacSystemFunctionSetpointRelationDataType `json:"hvacSystemFunctionSetpointRelationData,omitempty"`
}

type HvacSystemFunctionSetpointRelationListDataSelectorsType struct {
	SystemFunctionId *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty"`
	OperationModeId  *HvacOperationModeIdType  `json:"operationModeId,omitempty"`
}

type HvacSystemFunctionPowerSequenceRelationDataType struct {
	SystemFunctionId *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty" eebus:"key"`
	SequenceId       []PowerSequenceIdType     `json:"sequenceId,omitempty"`
}

type HvacSystemFunctionPowerSequenceRelationDataElementsType struct {
	SystemFunctionId *ElementTagType `json:"systemFunctionId,omitempty"`
	SequenceId       *ElementTagType `json:"sequenceId,omitempty"`
}

type HvacSystemFunctionPowerSequenceRelationListDataType struct {
	HvacSystemFunctionPowerSequenceRelationData []HvacSystemFunctionPowerSequenceRelationDataType `json:"hvacSystemFunctionPowerSequenceRelationData,omitempty"`
}

type HvacSystemFunctionPowerSequenceRelationListDataSelectorsType struct {
	SystemFunctionId *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty"`
}

type HvacSystemFunctionDescriptionDataType struct {
	SystemFunctionId   *HvacSystemFunctionIdType   `json:"systemFunctionId,omitempty" eebus:"key"`
	SystemFunctionType *HvacSystemFunctionTypeType `json:"systemFunctionType,omitempty"`
	Label              *LabelType                  `json:"label,omitempty"`
	Description        *DescriptionType            `json:"description,omitempty"`
}

type HvacSystemFunctionDescriptionDataElementsType struct {
	SystemFunctionId   *ElementTagType `json:"systemFunctionId,omitempty"`
	SystemFunctionType *ElementTagType `json:"systemFunctionType,omitempty"`
	Label              *ElementTagType `json:"label,omitempty"`
	Description        *ElementTagType `json:"description,omitempty"`
}

type HvacSystemFunctionDescriptionListDataType struct {
	HvacSystemFunctionDescriptionData []HvacSystemFunctionDescriptionDataType `json:"hvacSystemFunctionDescriptionData,omitempty"`
}

type HvacSystemFunctionDescriptionListDataSelectorsType struct {
	SystemFunctionId *HvacSystemFunctionIdType `json:"systemFunctionId,omitempty"`
}

type HvacOperationModeDescriptionDataType struct {
	OperationModeId   *HvacOperationModeIdType   `json:"operationModeId,omitempty" eebus:"key"`
	OperationModeType *HvacOperationModeTypeType `json:"operationModeType,omitempty"`
	Label             *LabelType                 `json:"label,omitempty"`
	Description       *DescriptionType           `json:"description,omitempty"`
}

type HvacOperationModeDescriptionDataElementsType struct {
	OperationModeId   *ElementTagType `json:"operationModeId,omitempty"`
	OperationModeType *ElementTagType `json:"operationModeType,omitempty"`
	Label             *ElementTagType `json:"label,omitempty"`
	Description       *ElementTagType `json:"description,omitempty"`
}

type HvacOperationModeDescriptionListDataType struct {
	HvacOperationModeDescriptionData []HvacOperationModeDescriptionDataType `json:"hvacOperationModeDescriptionData,omitempty"`
}

type HvacOperationModeDescriptionListDataSelectorsType struct {
	OperationModeId *HvacOperationModeIdType `json:"operationModeId,omitempty"`
}

type HvacOverrunDataType struct {
	OverrunId                 *HvacOverrunIdType     `json:"overrunId,omitempty" eebus:"key"`
	OverrunStatus             *HvacOverrunStatusType `json:"overrunStatus,omitempty"`
	TimeTableId               *TimeTableIdType       `json:"timeTableId,omitempty" eebus:"key"`
	IsOverrunStatusChangeable *bool                  `json:"isOverrunStatusChangeable,omitempty"`
}

type HvacOverrunDataElementsType struct {
	OverrunId                 *ElementTagType `json:"overrunId,omitempty"`
	OverrunStatus             *ElementTagType `json:"overrunStatus,omitempty"`
	TimeTableId               *ElementTagType `json:"timeTableId,omitempty"`
	IsOverrunStatusChangeable *ElementTagType `json:"isOverrunStatusChangeable,omitempty"`
}

type HvacOverrunListDataType struct {
	HvacOverrunData []HvacOverrunDataType `json:"hvacOverrunData,omitempty"`
}

type HvacOverrunListDataSelectorsType struct {
	OverrunId *HvacOverrunIdType `json:"overrunId,omitempty"`
}

type HvacOverrunDescriptionDataType struct {
	OverrunId                *HvacOverrunIdType         `json:"overrunId,omitempty" eebus:"key"`
	OverrunType              *HvacOverrunTypeType       `json:"overrunType,omitempty"`
	AffectedSystemFunctionId []HvacSystemFunctionIdType `json:"affectedSystemFunctionId,omitempty"`
	Label                    *LabelType                 `json:"label,omitempty"`
	Description              *DescriptionType           `json:"description,omitempty"`
}

type HvacOverrunDescriptionDataElementsType struct {
	OverrunId                *ElementTagType `json:"overrunId,omitempty"`
	OverrunType              *ElementTagType `json:"overrunType,omitempty"`
	AffectedSystemFunctionId *ElementTagType `json:"affectedSystemFunctionId,omitempty"`
	Label                    *ElementTagType `json:"label,omitempty"`
	Description              *ElementTagType `json:"description,omitempty"`
}

type HvacOverrunDescriptionListDataType struct {
	HvacOverrunDescriptionData []HvacOverrunDescriptionDataType `json:"hvacOverrunDescriptionData,omitempty"`
}

type HvacOverrunDescriptionListDataSelectorsType struct {
	OverrunId *HvacOverrunIdType `json:"overrunId,omitempty"`
}
