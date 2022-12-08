package model

type ThresholdIdType uint

type ThresholdTypeType string

const (
	ThresholdTypeTypeGoodAbove                ThresholdTypeType = "goodAbove"
	ThresholdTypeTypeBadAbove                 ThresholdTypeType = "badAbove"
	ThresholdTypeTypeGoodBelow                ThresholdTypeType = "goodBelow"
	ThresholdTypeTypeBadBelow                 ThresholdTypeType = "badBelow"
	ThresholdTypeTypeMinValueThreshold        ThresholdTypeType = "minValueThreshold"
	ThresholdTypeTypeMaxValueThreshold        ThresholdTypeType = "maxValueThreshold"
	ThresholdTypeTypeMinValueThresholdExtreme ThresholdTypeType = "minValueThresholdExtreme"
	ThresholdTypeTypeMaxValueThresholdExtreme ThresholdTypeType = "maxValueThresholdExtreme"
	ThresholdTypeTypeSagThreshold             ThresholdTypeType = "sagThreshold"
	ThresholdTypeTypeSwellThreshold           ThresholdTypeType = "swellThreshold"
)

type ThresholdDataType struct {
	ThresholdId    *ThresholdIdType  `json:"thresholdId,omitempty" eebus:"key"`
	ThresholdValue *ScaledNumberType `json:"thresholdValue,omitempty"`
}

type ThresholdDataElementsType struct {
	ThresholdId    *ElementTagType           `json:"thresholdId,omitempty"`
	ThresholdValue *ScaledNumberElementsType `json:"thresholdValue,omitempty"`
}

type ThresholdListDataType struct {
	ThresholdData []ThresholdDataType `json:"thresholdData,omitempty"`
}

type ThresholdListDataSelectorsType struct {
	ThresholdId *ThresholdIdType `json:"thresholdId,omitempty"`
}

type ThresholdConstraintsDataType struct {
	ThresholdId       *ThresholdIdType  `json:"thresholdId,omitempty" eebus:"key"`
	ThresholdRangeMin *ScaledNumberType `json:"thresholdRangeMin,omitempty"`
	ThresholdRangeMax *ScaledNumberType `json:"thresholdRangeMax,omitempty"`
	ThresholdStepSize *ScaledNumberType `json:"thresholdStepSize,omitempty"`
}

type ThresholdConstraintsDataElementsType struct {
	ThresholdId       *ElementTagType           `json:"thresholdId,omitempty"`
	ThresholdRangeMin *ScaledNumberElementsType `json:"thresholdRangeMin,omitempty"`
	ThresholdRangeMax *ScaledNumberElementsType `json:"thresholdRangeMax,omitempty"`
	ThresholdStepSize *ScaledNumberElementsType `json:"thresholdStepSize,omitempty"`
}

type ThresholdConstraintsListDataType struct {
	ThresholdConstraintsData []ThresholdConstraintsDataType `json:"thresholdConstraintsData,omitempty"`
}

type ThresholdConstraintsListDataSelectorsType struct {
	ThresholdId *ThresholdIdType `json:"thresholdId,omitempty"`
}

type ThresholdDescriptionDataType struct {
	ThresholdId   *ThresholdIdType       `json:"thresholdId,omitempty" eebus:"key"`
	ThresholdType *ThresholdTypeType     `json:"thresholdType,omitempty"`
	Unit          *UnitOfMeasurementType `json:"unit,omitempty"`
	ScopeType     *ScopeTypeType         `json:"scopeType,omitempty"`
	Label         *LabelType             `json:"label,omitempty"`
	Description   *DescriptionType       `json:"description,omitempty"`
}

type ThresholdDescriptionDataElementsType struct {
	ThresholdId   *ElementTagType `json:"thresholdId,omitempty"`
	ThresholdType *ElementTagType `json:"thresholdType,omitempty"`
	Unit          *ElementTagType `json:"unit,omitempty"`
	ScopeType     *ElementTagType `json:"scopeType,omitempty"`
	Label         *ElementTagType `json:"label,omitempty"`
	Description   *ElementTagType `json:"description,omitempty"`
}

type ThresholdDescriptionListDataType struct {
	ThresholdDescriptionData []ThresholdDescriptionDataType `json:"thresholdDescriptionData,omitempty"`
}

type ThresholdDescriptionListDataSelectorsType struct {
	ThresholdId *ThresholdIdType `json:"thresholdId,omitempty"`
	ScopeType   *ScopeTypeType   `json:"scopeType,omitempty"`
}
