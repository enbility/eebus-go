package model

type TariffIdType uint

type TariffCountType TariffIdType

type TierBoundaryIdType uint

type TierBoundaryCountType TierBoundaryIdType

type TierBoundaryTypeType string

const (
	TierBoundaryTypeTypePowerBoundary  TierBoundaryTypeType = "powerBoundary"
	TierBoundaryTypeTypeEnergyBoundary TierBoundaryTypeType = "energyBoundary"
	TierBoundaryTypeTypeCountBoundary  TierBoundaryTypeType = "countBoundary"
)

type CommodityIdType uint

type TierIdType uint

type TierCountType TierIdType

type TierTypeType string

const (
	TierTypeTypeFixedCost   TierTypeType = "fixedCost"
	TierTypeTypeDynamicCost TierTypeType = "dynamicCost"
)

type IncentiveIdType uint

type IncentiveCountType IncentiveIdType

type IncentiveTypeType string

const (
	IncentiveTypeTypeAbsoluteCost              IncentiveTypeType = "absoluteCost"
	IncentiveTypeTypeRelativeCost              IncentiveTypeType = "relativeCost"
	IncentiveTypeTypeRenewableEnergyPercentage IncentiveTypeType = "renewableEnergyPercentage"
	IncentiveTypeTypeCo2Emission               IncentiveTypeType = "co2Emission"
)

type IncentivePriorityType uint

type IncentiveValueTypeType string

const (
	IncentiveValueTypeTypeValue        IncentiveValueTypeType = "value"
	IncentiveValueTypeTypeAverageValue IncentiveValueTypeType = "averageValue"
	IncentiveValueTypeTypeMinvalue     IncentiveValueTypeType = "minValue"
	IncentiveValueTypeTypeMaxvalue     IncentiveValueTypeType = "maxValue"
)

type TariffOverallConstraintsDataType struct {
	MaxTariffCount         *TariffCountType       `json:"maxTariffCount,omitempty"`
	MaxBoundaryCount       *TierBoundaryCountType `json:"maxBoundaryCount,omitempty"`
	MaxTierCount           *TierCountType         `json:"maxTierCount,omitempty"`
	MaxIncentiveCount      *IncentiveCountType    `json:"maxIncentiveCount,omitempty"`
	MaxBoundariesPerTariff *TierBoundaryCountType `json:"maxBoundariesPerTariff,omitempty"`
	MaxTiersPerTariff      *TierCountType         `json:"maxTiersPerTariff,omitempty"`
	MaxBoundariesPerTier   *TierBoundaryCountType `json:"maxBoundariesPerTier,omitempty"`
	MaxIncentivesPerTier   *IncentiveCountType    `json:"maxIncentivesPerTier,omitempty"`
}

type TariffOverallConstraintsDataElementsType struct {
	MaxTariffCount         *ElementTagType `json:"maxTariffCount,omitempty"`
	MaxBoundaryCount       *ElementTagType `json:"maxBoundaryCount,omitempty"`
	MaxTierCount           *ElementTagType `json:"maxTierCount,omitempty"`
	MaxIncentiveCount      *ElementTagType `json:"maxIncentiveCount,omitempty"`
	MaxBoundariesPerTariff *ElementTagType `json:"maxBoundariesPerTariff,omitempty"`
	MaxTiersPerTariff      *ElementTagType `json:"maxTiersPerTariff,omitempty"`
	MaxBoundariesPerTier   *ElementTagType `json:"maxBoundariesPerTier,omitempty"`
	MaxIncentivesPerTier   *ElementTagType `json:"maxIncentivesPerTier,omitempty"`
}

type TariffDataType struct {
	TariffId     *TariffIdType `json:"tariffId,omitempty" eebus:"key"`
	ActiveTierId []TierIdType  `json:"activeTierId,omitempty"`
}

type TariffDataElementsType struct {
	TariffId     *ElementTagType `json:"tariffId,omitempty"`
	ActiveTierId *ElementTagType `json:"activeTierId,omitempty"`
}

type TariffListDataType struct {
	TariffData []TariffDataType `json:"tariffData,omitempty"`
}

type TariffListDataSelectorsType struct {
	TariffId     *TariffIdType `json:"tariffId,omitempty"`
	ActiveTierId *TierIdType   `json:"activeTierId,omitempty"`
}

type TariffTierRelationDataType struct {
	TariffId *TariffIdType `json:"tariffId,omitempty" eebus:"key"`
	TierId   []TierIdType  `json:"tierId,omitempty"`
}

type TariffTierRelationDataElementsType struct {
	TariffId *ElementTagType `json:"tariffId,omitempty"`
	TierId   *ElementTagType `json:"tierId,omitempty"`
}

type TariffTierRelationListDataType struct {
	TariffTierRelationData []TariffTierRelationDataType `json:"tariffTierRelationData,omitempty"`
}

type TariffTierRelationListDataSelectorsType struct {
	TariffId *TariffIdType `json:"tariffId,omitempty"`
	TierId   *TierIdType   `json:"tierId,omitempty"`
}

type TariffBoundaryRelationDataType struct {
	TariffId   *TariffIdType        `json:"tariffId,omitempty" eebus:"key"`
	BoundaryId []TierBoundaryIdType `json:"boundaryId,omitempty"`
}

type TariffBoundaryRelationDataElementsType struct {
	TariffId   *ElementTagType `json:"tariffId,omitempty"`
	BoundaryId *ElementTagType `json:"boundaryId,omitempty"`
}

type TariffBoundaryRelationListDataType struct {
	TariffBoundaryRelationData []TariffBoundaryRelationDataType `json:"tariffBoundaryRelationData,omitempty"`
}

type TariffBoundaryRelationListDataSelectorsType struct {
	TariffId   *TariffIdType       `json:"tariffId,omitempty"`
	BoundaryId *TierBoundaryIdType `json:"boundaryId,omitempty"`
}

type TariffDescriptionDataType struct {
	TariffId        *TariffIdType      `json:"tariffId,omitempty" eebus:"key"`
	CommodityId     *CommodityIdType   `json:"commodityId,omitempty"`
	MeasurementId   *MeasurementIdType `json:"measurementId,omitempty"`
	TariffWriteable *bool              `json:"tariffWriteable,omitempty"`
	UpdateRequired  *bool              `json:"updateRequired,omitempty"`
	ScopeType       *ScopeTypeType     `json:"scopeType,omitempty"`
	Label           *LabelType         `json:"label,omitempty"`
	Description     *DescriptionType   `json:"description,omitempty"`
	SlotIdSupport   *bool              `json:"slotIdSupport,omitempty"`
}

type TariffDescriptionDataElementsType struct {
	TariffId        *ElementTagType `json:"tariffId,omitempty"`
	CommodityId     *ElementTagType `json:"commodityId,omitempty"`
	MeasurementId   *ElementTagType `json:"measurementId,omitempty"`
	TariffWriteable *ElementTagType `json:"tariffWriteable,omitempty"`
	UpdateRequired  *ElementTagType `json:"updateRequired,omitempty"`
	ScopeType       *ElementTagType `json:"scopeType,omitempty"`
	Label           *ElementTagType `json:"label,omitempty"`
	Description     *ElementTagType `json:"description,omitempty"`
	SlotIdSupport   *ElementTagType `json:"slotIdSupport,omitempty"`
}

type TariffDescriptionListDataType struct {
	TariffDescriptionData []TariffDescriptionDataType `json:"tariffDescriptionData,omitempty"`
}

type TariffDescriptionListDataSelectorsType struct {
	TariffId      *TariffIdType      `json:"tariffId,omitempty"`
	CommodityId   *CommodityIdType   `json:"commodityId,omitempty"`
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
	ScopeType     *ScopeTypeType     `json:"scopeType,omitempty"`
}

type TierBoundaryDataType struct {
	BoundaryId         *TierBoundaryIdType `json:"boundaryId,omitempty" eebus:"key"`
	TimePeriod         *TimePeriodType     `json:"timePeriod,omitempty"`
	TimeTableId        *TimeTableIdType    `json:"timeTableId,omitempty"`
	LowerBoundaryValue *ScaledNumberType   `json:"lowerBoundaryValue,omitempty"`
	UpperBoundaryValue *ScaledNumberType   `json:"upperBoundaryValue,omitempty"`
}

type TierBoundaryDataElementsType struct {
	BoundaryId         *ElementTagType           `json:"boundaryId,omitempty"`
	TimePeriod         *TimePeriodElementsType   `json:"timePeriod,omitempty"`
	TimeTableId        *ElementTagType           `json:"timeTableId,omitempty"`
	LowerBoundaryValue *ScaledNumberElementsType `json:"lowerBoundaryValue,omitempty"`
	UpperBoundaryValue *ScaledNumberElementsType `json:"upperBoundaryValue,omitempty"`
}

type TierBoundaryListDataType struct {
	TierBoundaryData []TierBoundaryDataType `json:"tierBoundaryData,omitempty"`
}

type TierBoundaryListDataSelectorsType struct {
	BoundaryId *TierBoundaryIdType `json:"boundaryId,omitempty"`
}

type TierBoundaryDescriptionDataType struct {
	BoundaryId               *TierBoundaryIdType    `json:"boundaryId,omitempty" eebus:"key"`
	BoundaryType             *TierBoundaryTypeType  `json:"boundaryType,omitempty"`
	ValidForTierId           *TierIdType            `json:"validForTierId,omitempty"`
	SwitchToTierIdWhenLower  *TierIdType            `json:"switchToTierIdWhenLower,omitempty"`
	SwitchToTierIdWhenHigher *TierIdType            `json:"switchToTierIdWhenHigher,omitempty"`
	BoundaryUnit             *UnitOfMeasurementType `json:"boundaryUnit,omitempty"`
	Label                    *LabelType             `json:"label,omitempty"`
	Description              *DescriptionType       `json:"description,omitempty"`
}

type TierBoundaryDescriptionDataElementsType struct {
	BoundaryId               *ElementTagType `json:"boundaryId,omitempty"`
	BoundaryType             *ElementTagType `json:"boundaryType,omitempty"`
	ValidForTierId           *ElementTagType `json:"validForTierId,omitempty"`
	SwitchToTierIdWhenLower  *ElementTagType `json:"switchToTierIdWhenLower,omitempty"`
	SwitchToTierIdWhenHigher *ElementTagType `json:"switchToTierIdWhenHigher,omitempty"`
	BoundaryUnit             *ElementTagType `json:"boundaryUnit,omitempty"`
	Label                    *ElementTagType `json:"label,omitempty"`
	Description              *ElementTagType `json:"description,omitempty"`
}

type TierBoundaryDescriptionListDataType struct {
	TierBoundaryDescriptionData []TierBoundaryDescriptionDataType `json:"tierBoundaryDescriptionData,omitempty"`
}

type TierBoundaryDescriptionListDataSelectorsType struct {
	BoundaryId   *TierBoundaryIdType   `json:"boundaryId,omitempty"`
	BoundaryType *TierBoundaryTypeType `json:"boundaryType,omitempty"`
}

type CommodityDataType struct {
	CommodityId             *CommodityIdType     `json:"commodityId,omitempty" eebus:"key"`
	CommodityType           *CommodityTypeType   `json:"commodityType,omitempty"`
	PositiveEnergyDirection *EnergyDirectionType `json:"positiveEnergyDirection,omitempty"`
	Label                   *LabelType           `json:"label,omitempty"`
	Description             *DescriptionType     `json:"description,omitempty"`
}

type CommodityDataElementsType struct {
	CommodityId             *ElementTagType `json:"commodityId,omitempty"`
	CommodityType           *ElementTagType `json:"commodityType,omitempty"`
	PositiveEnergyDirection *ElementTagType `json:"positiveEnergyDirection,omitempty"`
	Label                   *ElementTagType `json:"label,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
}

type CommodityListDataType struct {
	CommodityData []CommodityDataType `json:"commodityData,omitempty"`
}

type CommodityListDataSelectorsType struct {
	CommodityId   *CommodityIdType   `json:"commodityId,omitempty"`
	CommodityType *CommodityTypeType `json:"commodityType,omitempty"`
}

type TierDataType struct {
	TierId            *TierIdType       `json:"tierId,omitempty" eebus:"key"`
	TimePeriod        *TimePeriodType   `json:"timePeriod,omitempty"`
	TimeTableId       *TimeTableIdType  `json:"timeTableId,omitempty"`
	ActiveIncentiveId []IncentiveIdType `json:"activeIncentiveId,omitempty"`
}

type TierDataElementsType struct {
	TierId            *ElementTagType `json:"tierId,omitempty"`
	TimePeriod        *ElementTagType `json:"timePeriod,omitempty"`
	TimeTableId       *ElementTagType `json:"timeTableId,omitempty"`
	ActiveIncentiveId *ElementTagType `json:"activeIncentiveId,omitempty"`
}

type TierListDataType struct {
	TierData []TierDataType `json:"tierData,omitempty"`
}

type TierListDataSelectorsType struct {
	TierId            *TierIdType      `json:"tierId,omitempty"`
	ActiveIncentiveId *IncentiveIdType `json:"activeIncentiveId,omitempty"`
}

type TierIncentiveRelationDataType struct {
	TierId      *TierIdType       `json:"tierId,omitempty" eebus:"key"`
	IncentiveId []IncentiveIdType `json:"incentiveId,omitempty"`
}

type TierIncentiveRelationDataElementsType struct {
	TierId      *ElementTagType `json:"tierId,omitempty"`
	IncentiveId *ElementTagType `json:"incentiveId,omitempty"`
}

type TierIncentiveRelationListDataType struct {
	TierIncentiveRelationData []TierIncentiveRelationDataType `json:"tierIncentiveRelationData,omitempty"`
}

type TierIncentiveRelationListDataSelectorsType struct {
	TierId      *TierIdType      `json:"tierId,omitempty"`
	IncentiveId *IncentiveIdType `json:"incentiveId,omitempty"`
}

type TierDescriptionDataType struct {
	TierId      *TierIdType      `json:"tierId,omitempty" eebus:"key"`
	TierType    *TierTypeType    `json:"tierType,omitempty"`
	Label       *LabelType       `json:"label,omitempty"`
	Description *DescriptionType `json:"description,omitempty"`
}

type TierDescriptionDataElementsType struct {
	TierId      *ElementTagType `json:"tierId,omitempty"`
	TierType    *ElementTagType `json:"tierType,omitempty"`
	Label       *ElementTagType `json:"label,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}

type TierDescriptionListDataType struct {
	TierDescriptionData []TierDescriptionDataType `json:"tierDescriptionData,omitempty"`
}

type TierDescriptionListDataSelectorsType struct {
	TierId   *TierIdType   `json:"tierId,omitempty"`
	TierType *TierTypeType `json:"tierType,omitempty"`
}

type IncentiveDataType struct {
	IncentiveId *IncentiveIdType            `json:"incentiveId,omitempty" eebus:"key"`
	ValueType   *IncentiveValueTypeType     `json:"valueType,omitempty"`
	Timestamp   *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	TimePeriod  *TimePeriodType             `json:"timePeriod,omitempty"`
	TimeTableId *TimeTableIdType            `json:"timeTableId,omitempty"`
	Value       *ScaledNumberType           `json:"value,omitempty"`
}

type IncentiveDataElementsType struct {
	IncentiveId *ElementTagType           `json:"incentiveId,omitempty"`
	ValueType   *ElementTagType           `json:"valueType,omitempty"`
	Timestamp   *ElementTagType           `json:"timestamp,omitempty"`
	TimePeriod  *TimePeriodElementsType   `json:"timePeriod,omitempty"`
	TimeTableId *ElementTagType           `json:"timeTableId,omitempty"`
	Value       *ScaledNumberElementsType `json:"value,omitempty"`
}

type IncentiveListDataType struct {
	IncentiveData []IncentiveDataType `json:"incentiveData,omitempty"`
}

type IncentiveListDataSelectorsType struct {
	IncentiveId       *IncentiveIdType        `json:"incentiveId,omitempty"`
	ValueType         *IncentiveValueTypeType `json:"valueType,omitempty"`
	TimestampInterval *TimestampIntervalType  `json:"timestampInterval,omitempty"`
}

type IncentiveDescriptionDataType struct {
	IncentiveId       *IncentiveIdType       `json:"incentiveId,omitempty" eebus:"key"`
	IncentiveType     *IncentiveTypeType     `json:"incentiveType,omitempty"`
	IncentivePriority *IncentivePriorityType `json:"incentivePriority,omitempty"`
	Currency          *CurrencyType          `json:"currency,omitempty"`
	Unit              *UnitOfMeasurementType `json:"unit,omitempty"`
	Label             *LabelType             `json:"label,omitempty"`
	Description       *DescriptionType       `json:"description,omitempty"`
}

type IncentiveDescriptionDataElementsType struct {
	IncentiveId       *ElementTagType `json:"incentiveId,omitempty"`
	IncentiveType     *ElementTagType `json:"incentiveType,omitempty"`
	IncentivePriority *ElementTagType `json:"incentivePriority,omitempty"`
	Currency          *ElementTagType `json:"currency,omitempty"`
	Unit              *ElementTagType `json:"unit,omitempty"`
	Label             *ElementTagType `json:"label,omitempty"`
	Description       *ElementTagType `json:"description,omitempty"`
}

type IncentiveDescriptionListDataType struct {
	IncentiveDescriptionData []IncentiveDescriptionDataType `json:"incentiveDescriptionData,omitempty"`
}

type IncentiveDescriptionListDataSelectorsType struct {
	IncentiveId   *IncentiveIdType   `json:"incentiveId,omitempty"`
	IncentiveType *IncentiveTypeType `json:"incentiveType,omitempty"`
}
