package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type FunctionDataCmd interface {
	FunctionData
	ReadCmdType() model.CmdType
	ReplyCmdType() model.CmdType
	WriteCmdType() model.CmdType
	NotifyCmdType(partial bool) model.CmdType
}

var _ FunctionDataCmd = (*FunctionDataCmdImpl[int])(nil)

type FunctionDataCmdImpl[T any] struct {
	*FunctionDataImpl[T]
}

func NewFunctionDataCmd[T any](function model.FunctionType) *FunctionDataCmdImpl[T] {
	return &FunctionDataCmdImpl[T]{
		FunctionDataImpl: NewFunctionData[T](function),
	}
}

func (r *FunctionDataCmdImpl[T]) ReadCmdType() model.CmdType {
	cmd := createCmd[T](r.functionType, nil)
	return cmd
}

func (r *FunctionDataCmdImpl[T]) ReplyCmdType() model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	return cmd
}

func (r *FunctionDataCmdImpl[T]) WriteCmdType() model.CmdType {
	cmd := createCmd[T](r.functionType, r.data)
	return cmd
}


func (r *FunctionDataCmdImpl[T]) NotifyCmdType(partial bool) model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	cmd.Function = util.Ptr(model.FunctionType(r.functionType))
	cmd.Filter = filterType(partial)
	return cmd
}

func filterType(partial bool) []model.FilterType {
	if partial {
		return []model.FilterType{{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}}
	}
	return nil
}

func createCmd[T any](function model.FunctionType, data *T) model.CmdType {
	result := model.CmdType{}

	switch function {
	case model.FunctionTypeDeviceClassificationManufacturerData:
		result.DeviceClassificationManufacturerData = castData[model.DeviceClassificationManufacturerDataType](data)
	case model.FunctionTypeDeviceDiagnosisStateData:
		result.DeviceDiagnosisStateData = castData[model.DeviceDiagnosisStateDataType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData:
		result.DeviceConfigurationKeyValueDescriptionListData = castData[model.DeviceConfigurationKeyValueDescriptionListDataType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueListData:
		result.DeviceConfigurationKeyValueListData = castData[model.DeviceConfigurationKeyValueListDataType](data)
	case model.FunctionTypeIdentificationListData:
		result.IdentificationListData = castData[model.IdentificationListDataType](data)
	case model.FunctionTypeMeasurementConstraintsListData:
		result.MeasurementConstraintsListData = castData[model.MeasurementConstraintsListDataType](data)
	case model.FunctionTypeMeasurementDescriptionListData:
		result.MeasurementDescriptionListData = castData[model.MeasurementDescriptionListDataType](data)
	case model.FunctionTypeMeasurementListData:
		result.MeasurementListData = castData[model.MeasurementListDataType](data)
	case model.FunctionTypeElectricalConnectionParameterDescriptionListData:
		result.ElectricalConnectionParameterDescriptionListData = castData[model.ElectricalConnectionParameterDescriptionListDataType](data)
	case model.FunctionTypeElectricalConnectionDescriptionListData:
		result.ElectricalConnectionDescriptionListData = castData[model.ElectricalConnectionDescriptionListDataType](data)
	case model.FunctionTypeElectricalConnectionPermittedValueSetListData:
		result.ElectricalConnectionPermittedValueSetListData = castData[model.ElectricalConnectionPermittedValueSetListDataType](data)
		// add more model types here
	case model.FunctionTypeHvacOverrunListData:
		result.HvacOverrunListData = castData[model.HvacOverrunListDataType](data)
	case model.FunctionTypeHvacOverrunDescriptionListData:
		result.HvacOverrunDescriptionListData = castData[model.HvacOverrunDescriptionListDataType](data)
	}

	return result
}

func castData[D, S any](data *S) *D {
	if data == nil {
		return new(D)
	}
	return any(data).(*D)
}
