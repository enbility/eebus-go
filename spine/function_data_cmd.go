package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type FunctionDataCmdImpl[T any] struct {
	*FunctionDataImpl[T]
	pendingRequests PendingRequests[T]
}

func NewFunctionDataCmd[T any](function model.FunctionEnumType) *FunctionDataCmdImpl[T] {
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

func (r *FunctionDataCmdImpl[T]) NotifyCmdType(partial bool) model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	cmd.Function = util.Ptr(model.FunctionType(r.functionType))
	cmd.Filter = filterType(partial)
	return cmd
}

func (r *FunctionDataCmdImpl[T]) PendingRequest() PendingRequests[T] {
	return r.pendingRequests
}

func filterType(partial bool) []model.FilterType {
	if partial {
		return []model.FilterType{{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}}
	}
	return nil
}

func createCmd[T any](function model.FunctionEnumType, data *T) model.CmdType {
	result := model.CmdType{}

	switch function {
	case model.FunctionEnumTypeDeviceClassificationManufacturerData:
		result.DeviceClassificationManufacturerData = castData[model.DeviceClassificationManufacturerDataType](data)
	}

	return result
}

func castData[D, S any](data *S) *D {
	if data == nil {
		return new(D)
	}
	return any(data).(*D)
}
