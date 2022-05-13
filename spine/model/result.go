package model

type ErrorNumberType uint

const (
	ErrorNumberTypeNoError ErrorNumberType = iota
	ErrorNumberTypeGeneralError
	ErrorNumberTypeTimeout
	ErrorNumberTypeOverload
	ErrorNumberTypeDestinationUnknown
	ErrorNumberTypeDestinationUnreachable
	ErrorNumberTypeCommandNotSupported
	ErrorNumberTypeCommandRejected
	ErrorNumberTypeRestrictedFunctionExchangeCombinationNotSupported
	ErrorNumberBindingIsNecessaryForThisCommand
)

type ResultDataType struct {
	ErrorNumber *ErrorNumberType `json:"errorNumber,omitempty"`
	Description *DescriptionType `json:"description,omitempty"`
}
