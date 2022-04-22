package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

// TODO: move to separate file
func mapCmdToFunction(cmd model.CmdType) (*model.FunctionType, any, error) {
	switch {
	case cmd.DeviceClassificationManufacturerData != nil:
		return util.Ptr(model.FunctionTypeDeviceClassificationManufacturerData), cmd.DeviceClassificationManufacturerData, nil
	}
	return nil, nil, fmt.Errorf("Function not found for cmd")
}

type FeatureLocalImpl struct {
	*FeatureImpl
	entity          *EntityLocalImpl
	functionDataMap map[model.FunctionType]FunctionDataCmd
}

func NewFeatureLocalImpl(id uint, entity *EntityLocalImpl, ftype model.FeatureTypeType, role model.RoleType) *FeatureLocalImpl {
	res := &FeatureLocalImpl{
		FeatureImpl: NewFeatureImpl(
			featureAddressType(id, entity.Device().Address(), entity.Address()),
			ftype,
			role),
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]FunctionDataCmd),
	}

	for _, fd := range CreateFunctionData[FunctionDataCmd](ftype) {
		res.functionDataMap[fd.Function()] = fd
	}

	return res
}

func (r *FeatureLocalImpl) SetData(function model.FunctionType, data any) {
	fd := r.functionData(function)
	fd.SetDataAny(data)

	// TODO:
	//f.NotifySubscribers([]model.CmdType{fd.NotifyCmdType(false)})
}

func (r *FeatureLocalImpl) Data(function model.FunctionType) any {
	return r.functionData(function).DataAny()
}

func (r *FeatureLocalImpl) RequestData(
	function model.FunctionType,
	destination *FeatureRemoteImpl,
	requestChannel any) (*model.MsgCounterType, error) {

	fd := r.functionData(function)
	cmd := fd.ReadCmdType()

	msgCounter, err := destination.Sender().Request(model.CmdClassifierTypeRead, r.Address(), false, []model.CmdType{cmd})
	if err == nil && requestChannel != nil {
		fd.AddPendingRequest(*msgCounter, requestChannel)
	}

	return msgCounter, err
}

func (r *FeatureLocalImpl) HandleMessage(message *Message) error {
	if message.Cmd.ResultData != nil {
		return r.processResult(message.CmdClassifier)
	}

	function, data, error := mapCmdToFunction(message.Cmd)
	if error != nil {
		return error
	}

	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.processRead(*function, message.RequestHeader, message.featureRemote)
	case model.CmdClassifierTypeReply:
		return r.processReply(*function, data, message.RequestHeader, message.featureRemote)
	default:
		return fmt.Errorf("CmdClassifier not implemented: %s", message.CmdClassifier)
	}
}

func (r *FeatureLocalImpl) processResult(cmdClassifier model.CmdClassifierType) error {
	switch cmdClassifier {
	case model.CmdClassifierTypeResult:
		// TODO process the return result data for the message sent with the ID in msgCounterReference
		// error numbers explained in Resource Spec 3.11
		return nil

	default:
		return fmt.Errorf("ResultData CmdClassifierType %s not implemented", cmdClassifier)
	}
}

func (r *FeatureLocalImpl) processRead(function model.FunctionType, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) error {
	cmd := r.functionData(function).ReplyCmdType()
	err := featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)

	return err
}

func (r *FeatureLocalImpl) processReply(function model.FunctionType, data any, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) error {
	featureRemote.SetData(function, data)
	r.functionData(function).HandleReply(*requestHeader.MsgCounter, data)
	return nil
}

func (r *FeatureLocalImpl) functionData(function model.FunctionType) FunctionDataCmd {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
