package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type RequestDataResult struct {
	data        any
	errorResult *ErrorType
}

func NewRequestDataResult(data any, errorResult *ErrorType) *RequestDataResult {
	return &RequestDataResult{
		data:        data,
		errorResult: errorResult,
	}
}

func NewRequestDataResultError(err error) *RequestDataResult {
	return NewRequestDataResult(nil, NewErrorType(model.ErrorNumberTypeGeneralError, err.Error()))
}
