package api

import "errors"

// ErrMetadataNotAvailable indicates that the meta data information is not available
// e.g. decsriptions, constraints, ...
var ErrMetadataNotAvailable = errors.New("meta data not available")

// ErrDataNotAvailable indicates that no data set is yet available
var ErrDataNotAvailable = errors.New("data not available")

// ErrDataInvalid indicates that the currently available data is not valid and should be ignored
var ErrDataInvalid = errors.New("data not valid")

// ErrDataForMetadataKeyNotFound indicates that no data item is found for the given key
var ErrDataForMetadataKeyNotFound = errors.New("data for key not found")

var ErrNotSupported = errors.New("not supported")

var ErrEntityNotFound = errors.New("entity not found")

var ErrUsecCaseNotSupported = errors.New("usecase is not supported")

var ErrFunctionNotSupported = errors.New("function is not supported")

var ErrOperationOnFunctionNotSupported = errors.New("operation is not supported on function")

var ErrMissingData = errors.New("missing data")

var ErrDeviceDisconnected = errors.New("device is disconnected")

var ErrNoCompatibleEntity = errors.New("no compatible entity")
