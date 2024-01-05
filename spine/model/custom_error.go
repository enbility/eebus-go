package model

import (
	"fmt"

	"github.com/enbility/eebus-go/util"
)

type ErrorType struct {
	ErrorNumber ErrorNumberType
	Description *DescriptionType
}

func NewErrorType(errorNumber ErrorNumberType, description string) *ErrorType {
	return &ErrorType{
		ErrorNumber: errorNumber,
		Description: util.Ptr(DescriptionType(description)),
	}
}

func NewErrorTypeFromNumber(errorNumber ErrorNumberType) *ErrorType {
	return &ErrorType{
		ErrorNumber: errorNumber,
	}
}

func NewErrorTypeFromString(description string) *ErrorType {
	return NewErrorType(ErrorNumberTypeGeneralError, description)
}

func NewErrorTypeFromResult(result *ResultDataType) *ErrorType {
	if *result.ErrorNumber == ErrorNumberTypeNoError {
		return nil
	}

	return &ErrorType{
		ErrorNumber: *result.ErrorNumber,
		Description: result.Description,
	}
}

func (e *ErrorType) String() string {
	if e.Description != nil && len(*e.Description) > 0 {
		return fmt.Sprintf("Error %d: %s", e.ErrorNumber, *e.Description)
	}
	return fmt.Sprintf("Error %d", e.ErrorNumber)
}
