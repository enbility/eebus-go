package spine

import (
	"fmt"
	"time"

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
	case cmd.DeviceConfigurationKeyValueDescriptionListData != nil:
		return util.Ptr(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData), cmd.DeviceConfigurationKeyValueDescriptionListData, nil
	case cmd.DeviceConfigurationKeyValueListData != nil:
		return util.Ptr(model.FunctionTypeDeviceConfigurationKeyValueListData), cmd.DeviceConfigurationKeyValueListData, nil
	case cmd.IdentificationListData != nil:
		return util.Ptr(model.FunctionTypeIdentificationListData), cmd.IdentificationListData, nil
	case cmd.MeasurementConstraintsListData != nil:
		return util.Ptr(model.FunctionTypeMeasurementConstraintsListData), cmd.MeasurementConstraintsListData, nil
	case cmd.MeasurementDescriptionListData != nil:
		return util.Ptr(model.FunctionTypeMeasurementDescriptionListData), cmd.MeasurementDescriptionListData, nil
	case cmd.MeasurementListData != nil:
		return util.Ptr(model.FunctionTypeMeasurementListData), cmd.MeasurementListData, nil
	case cmd.ElectricalConnectionParameterDescriptionListData != nil:
		return util.Ptr(model.FunctionTypeElectricalConnectionParameterDescriptionListData), cmd.ElectricalConnectionParameterDescriptionListData, nil
	case cmd.ElectricalConnectionDescriptionListData != nil:
		return util.Ptr(model.FunctionTypeElectricalConnectionDescriptionListData), cmd.ElectricalConnectionDescriptionListData, nil
	case cmd.ElectricalConnectionPermittedValueSetListData != nil:
		return util.Ptr(model.FunctionTypeElectricalConnectionPermittedValueSetListData), cmd.ElectricalConnectionPermittedValueSetListData, nil
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
		destination *FeatureRemoteImpl) (*model.MsgCounterType, *ErrorType)
	RequestDataBySenderAddress(
		function model.FunctionType,
		sender Sender,
		destinationAddress *model.FeatureAddressType,
		maxDelay time.Duration) (*model.MsgCounterType, *ErrorType)
	FetchRequestData(
		msgCounter model.MsgCounterType,
		destination *FeatureRemoteImpl) (any, *ErrorType)
	RequestAndFetchData(
		function model.FunctionType,
		destination *FeatureRemoteImpl) (any, *ErrorType)
	// Subscribes the local feature to the given destination feature; the go routine will block until the response is processed
	SubscribeAndWait(destination *FeatureRemoteImpl) *ErrorType
	NotifyData(function model.FunctionType, destination *FeatureRemoteImpl) (*model.MsgCounterType, *ErrorType)
	HandleMessage(message *Message) *ErrorType
}

var _ FeatureLocal = (*FeatureLocalImpl)(nil)

type FeatureLocalImpl struct {
	*FeatureImpl
	entity          *EntityLocalImpl
	functionDataMap map[model.FunctionType]FunctionDataCmd
	pendingRequests PendingRequests
}

func NewFeatureLocalImpl(id uint, entity *EntityLocalImpl, ftype model.FeatureTypeType, role model.RoleType) *FeatureLocalImpl {
	res := &FeatureLocalImpl{
		FeatureImpl: NewFeatureImpl(
			featureAddressType(id, entity.Address()),
			ftype,
			role),
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]FunctionDataCmd),
		pendingRequests: NewPendingRequest(),
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
	destination *FeatureRemoteImpl) (*model.MsgCounterType, *ErrorType) {
	return r.RequestDataBySenderAddress(function, destination.Sender(), destination.Address(), destination.MaxResponseDelayDuration())
}

func (r *FeatureLocalImpl) RequestDataBySenderAddress(
	function model.FunctionType,
	sender Sender,
	destinationAddress *model.FeatureAddressType,
	maxDelay time.Duration) (*model.MsgCounterType, *ErrorType) {

	fd := r.functionData(function)
	cmd := fd.ReadCmdType()

	msgCounter, err := sender.Request(model.CmdClassifierTypeRead, r.Address(), destinationAddress, false, []model.CmdType{cmd})
	if err == nil {
		r.pendingRequests.Add(*msgCounter, maxDelay)
		return msgCounter, nil
	}

	return msgCounter, NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
}

// this will block until the response is received
func (r *FeatureLocalImpl) FetchRequestData(
	msgCounter model.MsgCounterType,
	destination *FeatureRemoteImpl) (any, *ErrorType) {

	return r.pendingRequests.GetData(msgCounter)
}

// this will block until the response is received
func (r *FeatureLocalImpl) RequestAndFetchData(
	function model.FunctionType,
	destination *FeatureRemoteImpl) (any, *ErrorType) {

	msgCounter, err := r.RequestData(function, destination)
	if err != nil {
		return nil, err
	}

	return r.FetchRequestData(*msgCounter, destination)
}

func (r *FeatureLocalImpl) SubscribeAndWait(destination *FeatureRemoteImpl) *ErrorType {
	if r.Role() == model.RoleTypeServer {
		return NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r))
	}

	msgCounter, err := destination.Sender().Subscribe(r.Address(), destination.Address(), r.ftype)
	if err != nil {
		return NewErrorTypeFromString(err.Error())
	}

	maxDelay := defaultMaxResponseDelay
	rf := destination.Device().FeatureByAddress(NodeManagementAddress(destination.Device().Address()))
	if rf != nil {
		maxDelay = rf.MaxResponseDelayDuration()
	}

	r.pendingRequests.Add(*msgCounter, maxDelay)
	// this will block the go routine until the response is procedded
	_, result := r.pendingRequests.GetData(*msgCounter)
	// TODO: activate polling when subscription failed

	return result
}

func (r *FeatureLocalImpl) NotifyData(function model.FunctionType, destination *FeatureRemoteImpl) (*model.MsgCounterType, *ErrorType) {
	fd := r.functionData(function)
	cmd := fd.NotifyCmdType(false)

	msgCounter, err := destination.Sender().Request(model.CmdClassifierTypeRead, r.Address(), destination.Address(), false, []model.CmdType{cmd})
	if err != nil {
		return nil, NewErrorTypeFromString(err.Error())
	}
	return msgCounter, nil
}

func (r *FeatureLocalImpl) HandleMessage(message *Message) *ErrorType {
	if message.Cmd.ResultData != nil {
		return r.processResult(message)
	}

	function, data, error := mapCmdToFunction(message.Cmd)
	if error != nil {
		return error
	}

	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		if err := r.processRead(*function, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeReply:
		if err := r.processReply(*function, data, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	default:
		return NewErrorTypeFromString(fmt.Sprintf("CmdClassifier not implemented: %s", message.CmdClassifier))
	}

	return nil
}

func (r *FeatureLocalImpl) processResult(message *Message) *ErrorType {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeResult:
		if *message.Cmd.ResultData.ErrorNumber != model.ErrorNumberTypeNoError {
			// TODO process the return result data for the message sent with the ID in msgCounterReference
			// error numbers explained in Resource Spec 3.11
			fmt.Printf("Error Result received: %s", string(*message.Cmd.ResultData.Description))
		}
		return r.pendingRequests.SetResult(*message.RequestHeader.MsgCounterReference, NewErrorTypeFromResult(message.Cmd.ResultData))

	default:
		return NewErrorType(
			model.ErrorNumberTypeGeneralError,
			fmt.Sprintf("ResultData CmdClassifierType %s not implemented", message.CmdClassifier))
	}
}

func (r *FeatureLocalImpl) processRead(function model.FunctionType, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) *ErrorType {
	// is this a read request to a local server/special feature?
	if r.role == model.RoleTypeClient {
		// Read requests to a client feature are not allowed
		return NewErrorTypeFromNumber(model.ErrorNumberTypeCommandRejected)
	}

	cmd := r.functionData(function).ReplyCmdType()
	err := featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)

	return NewErrorTypeFromString(err.Error())
}

func (r *FeatureLocalImpl) processReply(function model.FunctionType, data any, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) *ErrorType {
	featureRemote.SetData(function, data)
	if err := r.pendingRequests.SetData(*requestHeader.MsgCounterReference, data); err != nil {
		payload := EventPayload{
			Ski:        featureRemote.Device().ski,
			EventType:  EventTypeDataChange,
			ChangeType: ElementChangeUpdate,
			Feature:    featureRemote,
			Data:       data,
		}
		Events.Publish(payload)
	}

	return nil
}

func (r *FeatureLocalImpl) functionData(function model.FunctionType) FunctionDataCmd {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
