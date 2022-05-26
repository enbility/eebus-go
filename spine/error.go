package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type ErrorType struct {
	ErrorNumber model.ErrorNumberType
	Description model.DescriptionType
}

func NewErrorType(errorNumber model.ErrorNumberType, description string) *ErrorType {
	return &ErrorType{
		ErrorNumber: errorNumber,
		Description: model.DescriptionType(description),
	}
}

func NewErrorTypeFromString(description string) *ErrorType {
	return NewErrorType(model.ErrorNumberTypeGeneralError, description)
}

func NewErrorTypeFromResult(result *model.ResultDataType) *ErrorType {
	return &ErrorType{
		ErrorNumber: *result.ErrorNumber,
		Description: *result.Description,
	}
}

func (e *ErrorType) Error() string {
	if len(e.Description) > 0 {
		return fmt.Sprintf("Error %d: %s", e.ErrorNumber, e.Description)
	}
	return fmt.Sprintf("Error %d", e.ErrorNumber)
}
