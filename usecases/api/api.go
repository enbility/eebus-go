package api

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

var PhaseNameMapping = []model.ElectricalConnectionPhaseNameType{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC}

// used to enable batch data updates for certain usecases
//
// a usecase that wants to provide batch update capabilities using this interface should
//
// 1. provide methods that return a type implementing this interface
// 2. provide an Update method that accepts a list of this interface
//
// The Update method can then iterate over the provided UpdateData, ensure all
// data points are supported, and then create a batched spine update request
type UpdateData interface {
	Supported() bool
	NotSupportedError() error
}

// used to enable batch data updates for MeaserumentData
//
// usecases can use this interface to provide batch update capabilities by
// implementing a method that takes a list of this interface and passes a list
// of MeasurementData to Measurement.UpdateDataForIds
type UpdateMeasurementData interface {
	UpdateData
	MeasurementData() api.MeasurementDataForID
}

// used to enable batch data updates for ConfigurationData
//
// usecases can use this interface to provide batch update capabilities by
// implementing a method that takes a list of this interface and updates the
// configuration data
type UpdateConfigurationData interface {
	UpdateData
	ConfigurationData() model.DeviceConfigurationKeyValueDataType
}
