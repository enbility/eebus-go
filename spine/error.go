package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

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
