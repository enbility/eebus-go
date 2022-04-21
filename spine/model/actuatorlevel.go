package model

type ActuatorLevelFctType string

const (
	ActuatorLevelFctTypeStart              ActuatorLevelFctType = "start"
	ActuatorLevelFctTypeUp                 ActuatorLevelFctType = "up"
	ActuatorLevelFctTypeDown               ActuatorLevelFctType = "down"
	ActuatorLevelFctTypeStop               ActuatorLevelFctType = "stop"
	ActuatorLevelFctTypePercentageAbsolute ActuatorLevelFctType = "percentageAbsolute"
	ActuatorLevelFctTypePercentageRelative ActuatorLevelFctType = "percentageRelative"
	ActuatorLevelFctTypeAbsolut            ActuatorLevelFctType = "absolut"
	ActuatorLevelFctTypeRelative           ActuatorLevelFctType = "relative"
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
