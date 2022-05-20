package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

// TODO: move to separate file
func mapCmdToFunction(cmd model.CmdType) (*model.FunctionType, any, *ErrorType) {
	switch {
	case cmd.NodeManagementDetailedDiscoveryData != nil:
		return util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData), cmd.NodeManagementDetailedDiscoveryData, nil
	case cmd.DeviceClassificationManufacturerData != nil:
		return util.Ptr(model.FunctionTypeDeviceClassificationManufacturerData), cmd.DeviceClassificationManufacturerData, nil
	case cmd.DeviceDiagnosisStateData != nil:
		return util.Ptr(model.FunctionTypeDeviceDiagnosisStateData), cmd.DeviceDiagnosisStateData, nil
	}
	return nil, nil, NewErrorType(model.ErrorNumberTypeCommandNotSupported, "Function not found for cmd")
}

type FeatureLocal interface {
	Feature
	Data(function model.FunctionType) any
	SetData(function model.FunctionType, data any)
	Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType
	RequestData(
		function model.FunctionType,
		destination *FeatureRemoteImpl,
		requestChannel any) (*model.MsgCounterType, error)
	NotifyData(function model.FunctionType, destination *FeatureRemoteImpl) (*model.MsgCounterType, error)
	HandleMessage(message *Message) *ErrorType
}

var _ FeatureLocal = (*FeatureLocalImpl)(nil)

type FeatureLocalImpl struct {
	*FeatureImpl
	entity          *EntityLocalImpl
	functionDataMap map[model.FunctionType]FunctionDataCmd
}

func NewFeatureLocalImpl(id uint, entity *EntityLocalImpl, ftype model.FeatureTypeType, role model.RoleType) *FeatureLocalImpl {
	res := &FeatureLocalImpl{
		FeatureImpl: NewFeatureImpl(
			featureAddressType(id, entity.Address()),
			ftype,
			role),
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]FunctionDataCmd),
	}

	for _, fd := range CreateFunctionData[FunctionDataCmd](ftype) {
		res.functionDataMap[fd.Function()] = fd
	}
	res.operations = make(map[model.FunctionType]*Operations)

	return res
}

func (r *FeatureLocalImpl) Device() *DeviceLocalImpl {
	return r.entity.Device()
}

func (r *FeatureLocalImpl) Entity() *EntityLocalImpl {
	return r.entity
}

// Add supported function to the feature if its role is Server or Speical
func (r *FeatureLocalImpl) AddFunctionType(function model.FunctionType, read, write bool) {
	if r.role != model.RoleTypeServer && r.role != model.RoleTypeSpecial {
		return
	}
	if r.operations[function] != nil {
		return
	}
	r.operations[function] = NewOperations(read, write)
}

func (r *FeatureLocalImpl) Data(function model.FunctionType) any {
	return r.functionData(function).DataAny()
}

func (r *FeatureLocalImpl) SetData(function model.FunctionType, data any) {
	fd := r.functionData(function)
	fd.SetDataAny(data)

	r.Device().NotifySubscribers(r.Address(), []model.CmdType{fd.NotifyCmdType(false)})
}

func (r *FeatureLocalImpl) Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType {
	var funs []model.FunctionPropertyType
	for fun, operations := range r.operations {
		var functionType model.FunctionType = model.FunctionType(fun)
		sf := model.FunctionPropertyType{
			Function:           &functionType,
			PossibleOperations: operations.Information(),
		}

		funs = append(funs, sf)
	}

	res := model.NodeManagementDetailedDiscoveryFeatureInformationType{
		Description: &model.NetworkManagementFeatureDescriptionDataType{
			FeatureAddress:    r.Address(),
			FeatureType:       &r.ftype,
			Role:              &r.role,
			Description:       r.description,
			SupportedFunction: funs,
		},
	}

	return &res
}

func (r *FeatureLocalImpl) RequestData(
	function model.FunctionType,
	destination *FeatureRemoteImpl,
	requestChannel any) (*model.MsgCounterType, error) {

	fd := r.functionData(function)
	cmd := fd.ReadCmdType()

	msgCounter, err := destination.Sender().Request(model.CmdClassifierTypeRead, r.Address(), destination.Address(), false, []model.CmdType{cmd})
	if err == nil && requestChannel != nil {
		fd.AddPendingRequest(*msgCounter, requestChannel)
	}

	return msgCounter, err
}

func (r *FeatureLocalImpl) NotifyData(function model.FunctionType, destination *FeatureRemoteImpl) (*model.MsgCounterType, error) {
	fd := r.functionData(function)
	cmd := fd.NotifyCmdType(false)

	return destination.Sender().Request(model.CmdClassifierTypeRead, r.Address(), destination.Address(), false, []model.CmdType{cmd})
}

func (r *FeatureLocalImpl) HandleMessage(message *Message) *ErrorType {
	if message.Cmd.ResultData != nil {
		return r.processResult(message.CmdClassifier)
	}

	function, data, error := mapCmdToFunction(message.Cmd)
	if error != nil {
		return error
	}

	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		if err := r.processRead(*function, message.RequestHeader, message.featureRemote); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}
	case model.CmdClassifierTypeReply:
		if err := r.processReply(*function, data, message.RequestHeader, message.featureRemote); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}
	default:
		return NewErrorType(model.ErrorNumberTypeGeneralError, fmt.Sprintf("CmdClassifier not implemented: %s", message.CmdClassifier))
	}

	return nil
}

func (r *FeatureLocalImpl) processResult(cmdClassifier model.CmdClassifierType) *ErrorType {
	switch cmdClassifier {
	case model.CmdClassifierTypeResult:
		// TODO process the return result data for the message sent with the ID in msgCounterReference
		// error numbers explained in Resource Spec 3.11
		return nil

	default:
		return NewErrorType(model.ErrorNumberTypeGeneralError, fmt.Sprintf("ResultData CmdClassifierType %s not implemented", cmdClassifier))
	}
}

func (r *FeatureLocalImpl) processRead(function model.FunctionType, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) error {
	// is this a read request to a local server/special feature?
	if r.role == model.RoleTypeClient {
		// Read requests to a client feature are not allowed
		return featureRemote.Sender().Result(requestHeader, r.Address(), model.ErrorNumberTypeCommandRejected, nil)
	}

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
