package model

type ErrorNumberType uint

type ResultDataType struct {
	ErrorNumber *ErrorNumberType `json:"errorNumber,omitempty"`
	Description *DescriptionType `json:"description,omitempty"`
}
