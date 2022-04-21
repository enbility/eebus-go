package model

type ActuatorSwitchFctType string

const (
	ActuatorSwitchFctTypeOn     ActuatorSwitchFctType = "on"
	ActuatorSwitchFctTypeOff    ActuatorSwitchFctType = "off"
	ActuatorSwitchFctTypeToggle ActuatorSwitchFctType = "toggle"
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
