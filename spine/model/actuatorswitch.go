package model

type ActuatorSwitchFctType ActuatorSwitchFctEnumType

type ActuatorSwitchFctEnumType string

const (
	ActuatorSwitchFctEnumTypeOn     ActuatorSwitchFctEnumType = "on"
	ActuatorSwitchFctEnumTypeOff    ActuatorSwitchFctEnumType = "off"
	ActuatorSwitchFctEnumTypeToggle ActuatorSwitchFctEnumType = "toggle"
)

type ActuatorSwitchDataType struct {
	Function *ActuatorSwitchFctType `json:"function,omitempty"`
}

type ActuatorSwitchDataElementsType struct {
	Function *ElementTagType `json:"function,omitempty"`
}

type ActuatorSwitchDescriptionDataType struct {
	Label       *LabelType       `json:"label,omitempty"`
	Description *DescriptionType `json:"description,omitempty"`
}

type ActuatorSwitchDescriptionDataElementsType struct {
	Label       *ElementTagType `json:"label,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}
