package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type FunctionDataCmd interface {
	FunctionData
	ReadCmdType() model.CmdType
	ReplyCmdType() model.CmdType
	NotifyCmdType(partial bool) model.CmdType
	AddPendingRequest(counter model.MsgCounterType, requestChannel any)
	HandleReply(remoteDevice *DeviceRemoteImpl, counter model.MsgCounterType, data any)
}

var _ FunctionDataCmd = (*FunctionDataCmdImpl[int])(nil)

type FunctionDataCmdImpl[T any] struct {
	*FunctionDataImpl[T]
	pendingRequests PendingRequests[*T]
}

func NewFunctionDataCmd[T any](function model.FunctionType) *FunctionDataCmdImpl[T] {
	return &FunctionDataCmdImpl[T]{
		FunctionDataImpl: NewFunctionData[T](function),
		pendingRequests:  make(PendingRequests[*T]),
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

func (r *FunctionDataCmdImpl[T]) AddPendingRequest(counter model.MsgCounterType, requestChannel any) {
	r.pendingRequests.Add(counter, requestChannel.(chan *T))
}

func (r *FunctionDataCmdImpl[T]) HandleReply(remoteDevice *DeviceRemoteImpl, counter model.MsgCounterType, data any) {
	if err := r.pendingRequests.Handle(counter, data.(*T)); err != nil {
		payload := EventPayload{
			Ski:        remoteDevice.ski,
			EventType:  EventTypeDataChange,
			ChangeType: ElementChangeUpdate,
			Device:     remoteDevice,
			Data:       data,
		}
		Events.Publish(payload)
	}
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
		// add more model types here
	}

	return result
}

func castData[D, S any](data *S) *D {
	if data == nil {
		return new(D)
	}
	return any(data).(*D)
}
