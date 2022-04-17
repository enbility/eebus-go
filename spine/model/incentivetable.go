package model

type IncentiveTableType struct {
	Tariff        *TariffDataType                   `json:"tariff,omitempty"` // ignoring changes
	IncentiveSlot []IncentiveTableIncentiveSlotType `json:"incentiveSlot,omitempty"`
}

type IncentiveTableElementsType struct {
	Tariff        *TariffDataElementsType                  `json:"tariff,omitempty"` // ignoring changes
	IncentiveSlot *IncentiveTableIncentiveSlotElementsType `json:"incentiveSlot,omitempty"`
}

type IncentiveTableIncentiveSlotType struct {
	TimeInterval *TimeTableDataType       `json:"timeInterval,omitempty"` // ignoring changes
	Tier         []IncentiveTableTierType `json:"tier,omitempty"`
}

type IncentiveTableIncentiveSlotElementsType struct {
	TimeInterval *TimeTableDataElementsType      `json:"timeInterval,omitempty"` // ignoring changes
	Tier         *IncentiveTableTierElementsType `json:"tier,omitempty"`
}

type IncentiveTableTierType struct {
	Tier      *TierDataType          `json:"tier,omitempty"`      // ignoring changes
	Boundary  []TierBoundaryDataType `json:"boundary,omitempty"`  // ignoring changes
	Incentive []IncentiveDataType    `json:"incentive,omitempty"` // ignoring changes
}

type IncentiveTableTierElementsType struct {
	Tier      *TierDataElementsType         `json:"tier,omitempty"`      // ignoring changes
	Boundary  *TierBoundaryDataElementsType `json:"boundary,omitempty"`  // ignoring changes
	Incentive *IncentiveDataElementsType    `json:"incentive,omitempty"` // ignoring changes
}

type IncentiveTableDataType struct {
	IncentiveTable []IncentiveTableType `json:"incentiveTable,omitempty"`
}

type IncentiveTableDataElementsType struct {
	IncentiveTable *IncentiveTableElementsType `json:"incentiveTable,omitempty"`
}

type IncentiveTableDataSelectorsType struct {
	Tariff *TariffListDataSelectorsType `json:"tariff,omitempty"`
}

type IncentiveTableDescriptionType struct {
	TariffDescription *TariffDescriptionDataType          `json:"tariffDescription,omitempty"`
	Tier              []IncentiveTableDescriptionTierType `json:"tier,omitempty"`
}

type IncentiveTableDescriptionElementsType struct {
	TariffDescription *TariffDescriptionDataElementsType `json:"tariffDescription,omitempty"`
	Tier              *IncentiveTableDescriptionTierType `json:"tier,omitempty"`
}

type IncentiveTableDescriptionTierType struct {
	TierDescription      *TierDescriptionDataType          `json:"tierDescription,omitempty"`
	BoundaryDescription  []TierBoundaryDescriptionDataType `json:"boundaryDescription,omitempty"`
	IncentiveDescription []IncentiveDescriptionDataType    `json:"incentiveDescription,omitempty"`
}

type IncentiveTableDescriptionTierElementsType struct {
	TierDescription      *TierDescriptionDataElementsType         `json:"tierDescription,omitempty"`
	BoundaryDescription  *TierBoundaryDescriptionDataElementsType `json:"boundaryDescription,omitempty"`
	IncentiveDescription *IncentiveDescriptionDataElementsType    `json:"incentiveDescription,omitempty"`
}

type IncentiveTableDescriptionDataType struct {
	IncentiveTableDescription []IncentiveTableDescriptionType `json:"incentiveTableDescription,omitempty"`
}

type IncentiveTableDescriptionDataElementsType struct {
	IncentiveTableDescription *IncentiveTableDescriptionElementsType `json:"incentiveTableDescription,omitempty"`
}

type IncentiveTableDescriptionDataSelectorsType struct {
	TariffDescription *TariffDescriptionListDataSelectorsType `json:"tariffDescription,omitempty"`
}

type IncentiveTableConstraintsType struct {
	Tariff                   *TariffDataType                   `json:"tariff,omitempty"`
	TariffConstraints        *TariffOverallConstraintsDataType `json:"tariffConstraints,omitempty"`
	IncentiveSlotConstraints *TimeTableConstraintsDataType     `json:"incentiveSlotConstraints,omitempty"`
}

type IncentiveTableConstraintsElementsType struct {
	Tariff                   *TariffDataElementsType                   `json:"tariff,omitempty"`
	TariffConstraints        *TariffOverallConstraintsDataElementsType `json:"tariffConstraints,omitempty"`
	IncentiveSlotConstraints *TimeTableConstraintsDataElementsType     `json:"incentiveSlotConstraints,omitempty"`
}

type IncentiveTableConstraintsDataType struct {
	IncentiveTableConstraints []IncentiveTableConstraintsType `json:"incentiveTableConstraints,omitempty"`
}

type IncentiveTableConstraintsDataElementsType struct {
	IncentiveTableConstraints *IncentiveTableConstraintsElementsType `json:"incentiveTableConstraints,omitempty"`
}

type IncentiveTableConstraintsDataSelectorsType struct {
	Tariff *TariffListDataSelectorsType `json:"tariff,omitempty"`
}
