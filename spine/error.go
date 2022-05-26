package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type ErrorType struct {
	ErrorNumber model.ErrorNumberType
	Description *model.DescriptionType
}

func NewErrorType(errorNumber model.ErrorNumberType, description string) *ErrorType {
	return &ErrorType{
		ErrorNumber: errorNumber,
		Description: util.Ptr(model.DescriptionType(description)),
	}
}

func NewErrorTypeFromNumber(errorNumber model.ErrorNumberType) *ErrorType {
	return &ErrorType{
		ErrorNumber: errorNumber,
	}
}

func NewErrorTypeFromString(description string) *ErrorType {
	return NewErrorType(model.ErrorNumberTypeGeneralError, description)
}

func NewErrorTypeFromResult(result *model.ResultDataType) *ErrorType {
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
