package spine

type ActuatorLevelFctType ActuatorLevelFctEnumType

type ActuatorLevelFctEnumType string

const (
	ActuatorLevelFctEnumTypeStart              ActuatorLevelFctEnumType = "start"
	ActuatorLevelFctEnumTypeUp                 ActuatorLevelFctEnumType = "up"
	ActuatorLevelFctEnumTypeDown               ActuatorLevelFctEnumType = "down"
	ActuatorLevelFctEnumTypeStop               ActuatorLevelFctEnumType = "stop"
	ActuatorLevelFctEnumTypePercentageAbsolute ActuatorLevelFctEnumType = "percentageAbsolute"
	ActuatorLevelFctEnumTypePercentageRelative ActuatorLevelFctEnumType = "percentageRelative"
	ActuatorLevelFctEnumTypeAbsolut            ActuatorLevelFctEnumType = "absolut"
	ActuatorLevelFctEnumTypeRelative           ActuatorLevelFctEnumType = "relative"
)

type ActuatorLevelDataType struct {
	Function *ActuatorLevelFctType `json:"function,omitempty"`
	Value    *ScaledNumberType     `json:"value,omitempty"`
}

type ActuatorLevelDataElementsType struct {
	Function *ElementTagType `json:"function,omitempty"`
	Value    *ElementTagType `json:"value,omitempty"`
}

type ActuatorLevelDescriptionDataType struct {
	Label            *LabelType             `json:"label,omitempty"`
	Description      *DescriptionType       `json:"description,omitempty"`
	LevelDefaultUnit *UnitOfMeasurementType `json:"levelDefaultUnit,omitempty"`
}

type ActuatorLevelDescriptionDataElementsType struct {
	Label            *ElementTagType `json:"label,omitempty"`
	Description      *ElementTagType `json:"description,omitempty"`
	LevelDefaultUnit *ElementTagType `json:"levelDefaultUnit,omitempty"`
}
