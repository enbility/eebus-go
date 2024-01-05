package spine

import (
	"fmt"
	"sync"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

var _ FeatureLocal = (*FeatureLocalImpl)(nil)

type FeatureLocalImpl struct {
	*FeatureImpl

	muxResultCB     sync.Mutex
	entity          *EntityLocalImpl
	functionDataMap map[model.FunctionType]FunctionDataCmd
	pendingRequests PendingRequests
	resultHandler   []FeatureResult
	resultCallback  map[model.MsgCounterType]func(result ResultMessage)

	mux sync.Mutex
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
		resultCallback:  make(map[model.MsgCounterType]func(result ResultMessage)),
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

// Add supported function to the feature if its role is Server or Special
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
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.functionData(function).DataAny()
}

func (r *FeatureLocalImpl) SetData(function model.FunctionType, data any) {
	r.mux.Lock()

	fd := r.functionData(function)
	fd.UpdateDataAny(data, nil, nil)

	r.mux.Unlock()

	r.Device().NotifySubscribers(r.Address(), fd.NotifyCmdType(nil, nil, false, nil))
}

func (r *FeatureLocalImpl) AddResultHandler(handler FeatureResult) {
	r.resultHandler = append(r.resultHandler, handler)
}

func (r *FeatureLocalImpl) AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg ResultMessage)) {
	r.muxResultCB.Lock()
	defer r.muxResultCB.Unlock()

	r.resultCallback[msgCounterReference] = function
}

func (r *FeatureLocalImpl) processResultCallbacks(msgCounterReference model.MsgCounterType, msg ResultMessage) {
	r.muxResultCB.Lock()
	defer r.muxResultCB.Unlock()

	cb, ok := r.resultCallback[msgCounterReference]
	if !ok {
		return
	}

	go cb(msg)

	delete(r.resultCallback, msgCounterReference)
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
	selector any,
	elements any,
	destination *FeatureRemoteImpl) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	cmd := fd.ReadCmdType(selector, elements)

	return r.RequestDataBySenderAddress(cmd, destination.Sender(), destination.Device().ski, destination.Address(), destination.MaxResponseDelayDuration())
}

func (r *FeatureLocalImpl) RequestDataBySenderAddress(
	cmd model.CmdType,
	sender Sender,
	deviceSki string,
	destinationAddress *model.FeatureAddressType,
	maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType) {

	msgCounter, err := sender.Request(model.CmdClassifierTypeRead, r.Address(), destinationAddress, false, []model.CmdType{cmd})
	if err == nil {
		r.pendingRequests.Add(deviceSki, *msgCounter, maxDelay)
		return msgCounter, nil
	}

	return msgCounter, model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
}

// Wait and return the response from destination for a message with the msgCounter ID
// this will block until the response is received
func (r *FeatureLocalImpl) FetchRequestData(
	msgCounter model.MsgCounterType,
	destination *FeatureRemoteImpl) (any, *model.ErrorType) {

	return r.pendingRequests.GetData(destination.Device().ski, msgCounter)
}

// Send a data request for function to destination and return the response
// this will block until the response is received
func (r *FeatureLocalImpl) RequestAndFetchData(
	function model.FunctionType,
	selector any,
	elements any,
	destination *FeatureRemoteImpl) (any, *model.ErrorType) {

	msgCounter, err := r.RequestData(function, selector, elements, destination)
	if err != nil {
		return nil, err
	}

	return r.FetchRequestData(*msgCounter, destination)
}

// Subscribe to a remote feature
func (r *FeatureLocalImpl) Subscribe(remoteDevice *DeviceRemoteImpl, remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if r.Role() == model.RoleTypeServer {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r))
	}

	msgCounter, err := remoteDevice.Sender().Subscribe(r.Address(), remoteAdress, r.ftype)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	return msgCounter, nil
}

/*
TODO: check if this function is needed and can be fixed, see https://github.com/enbility/eebus-go/issues/31
// Subscribe to a remote feature and wait for the result
func (r *FeatureLocalImpl) SubscribeAndWait(remoteDevice *DeviceRemoteImpl, remoteAdress *model.FeatureAddressType) *ErrorType {
	if r.Role() == model.RoleTypeServer {
		return NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r))
	}

	msgCounter, err := remoteDevice.Sender().Subscribe(r.Address(), remoteAdress, r.ftype)
	if err != nil {
		return NewErrorTypeFromString(err.Error())
	}

	maxDelay := defaultMaxResponseDelay
	rf := remoteDevice.FeatureByAddress(NodeManagementAddress(remoteDevice.Address()))
	if rf != nil {
		maxDelay = rf.MaxResponseDelayDuration()
	}

	r.pendingRequests.Add(*msgCounter, maxDelay)
	// this will block the go routine until the response is procedded
	_, result := r.pendingRequests.GetData(*msgCounter)
	// TODO: activate polling when subscription failed

	return result
}
*/

// Bind to a remote feature
func (r *FeatureLocalImpl) Bind(remoteDevice *DeviceRemoteImpl, remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if r.Role() == model.RoleTypeServer {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r))
	}

	msgCounter, err := remoteDevice.Sender().Bind(r.Address(), remoteAddress, r.ftype)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	return msgCounter, nil
}

/*
TODO: check if this function is needed and can be fixed, see https://github.com/enbility/eebus-go/issues/31
// Bind to a remote feature and wait for the result
func (r *FeatureLocalImpl) BindAndWait(remoteDevice *DeviceRemoteImpl, remoteAddress *model.FeatureAddressType) *ErrorType {
	if r.Role() == model.RoleTypeServer {
		return NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r))
	}

	msgCounter, err := remoteDevice.Sender().Bind(r.Address(), remoteAddress, r.ftype)
	if err != nil {
		return NewErrorTypeFromString(err.Error())
	}

	maxDelay := defaultMaxResponseDelay
	rf := remoteDevice.FeatureByAddress(remoteAddress)
	if rf != nil {
		maxDelay = rf.MaxResponseDelayDuration()
	}

	r.pendingRequests.Add(*msgCounter, maxDelay)
	// this will block the go routine until the response is procedded
	_, result := r.pendingRequests.GetData(*msgCounter)
	// TODO: activate polling when binding failed

	return result
}
*/

// Send a notification message with the current data of function to the destination
func (r *FeatureLocalImpl) NotifyData(
	function model.FunctionType,
	deleteSelector, partialSelector any,
	partialWithoutSelector bool,
	deleteElements any,
	destination *FeatureRemoteImpl) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	cmd := fd.NotifyCmdType(deleteSelector, partialSelector, partialWithoutSelector, deleteElements)

	msgCounter, err := destination.Sender().Request(model.CmdClassifierTypeRead, r.Address(), destination.Address(), false, []model.CmdType{cmd})
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}
	return msgCounter, nil
}

// Send a write message with provided data of function to the destination
func (r *FeatureLocalImpl) WriteData(
	function model.FunctionType,
	deleteSelector, partialSelector any,
	deleteElements any,
	destination *FeatureRemoteImpl) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	cmd := fd.WriteCmdType(deleteSelector, partialSelector, deleteElements)

	msgCounter, err := destination.Sender().Write(r.Address(), destination.Address(), cmd)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	return msgCounter, nil
}

func (r *FeatureLocalImpl) HandleMessage(message *Message) *model.ErrorType {
	if message.Cmd.ResultData != nil {
		return r.processResult(message)
	}

	cmdData, err := message.Cmd.Data()
	if err != nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, err.Error())
	}
	if cmdData.Function == nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "No function found for cmd data")
	}

	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		if err := r.processRead(*cmdData.Function, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeReply:
		if err := r.processReply(*cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeNotify:
		if err := r.processNotify(*cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete, message.FeatureRemote); err != nil {
			return err
		}
	default:
		return model.NewErrorTypeFromString(fmt.Sprintf("CmdClassifier not implemented: %s", message.CmdClassifier))
	}

	return nil
}

func (r *FeatureLocalImpl) processResult(message *Message) *model.ErrorType {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeResult:
		if *message.Cmd.ResultData.ErrorNumber != model.ErrorNumberTypeNoError {
			// error numbers explained in Resource Spec 3.11
			errorString := fmt.Sprintf("Error Result received %d", *message.Cmd.ResultData.ErrorNumber)
			if message.Cmd.ResultData.Description != nil {
				errorString += fmt.Sprintf(": %s", *message.Cmd.ResultData.Description)
			}
			logging.Log().Debug(errorString)
		}

		// we don't need to populate this error as requests don't require a pendingRequest entry
		_ = r.pendingRequests.SetResult(message.DeviceRemote.ski, *message.RequestHeader.MsgCounterReference, model.NewErrorTypeFromResult(message.Cmd.ResultData))

		if message.RequestHeader.MsgCounterReference == nil {
			return nil
		}

		// call the Features Error Handler
		errorMsg := ResultMessage{
			MsgCounterReference: *message.RequestHeader.MsgCounterReference,
			Result:              message.Cmd.ResultData,
			FeatureLocal:        r,
			FeatureRemote:       message.FeatureRemote,
			DeviceRemote:        message.DeviceRemote,
		}

		if r.resultHandler != nil {
			for _, item := range r.resultHandler {
				go item.HandleResult(errorMsg)
			}
		}

		r.processResultCallbacks(*message.RequestHeader.MsgCounterReference, errorMsg)

		return nil

	default:
		return model.NewErrorType(
			model.ErrorNumberTypeGeneralError,
			fmt.Sprintf("ResultData CmdClassifierType %s not implemented", message.CmdClassifier))
	}
}

func (r *FeatureLocalImpl) processRead(function model.FunctionType, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) *model.ErrorType {
	// is this a read request to a local server/special feature?
	if r.role == model.RoleTypeClient {
		// Read requests to a client feature are not allowed
		return model.NewErrorTypeFromNumber(model.ErrorNumberTypeCommandRejected)
	}

	cmd := r.functionData(function).ReplyCmdType(false)
	if err := featureRemote.Sender().Reply(requestHeader, r.Address(), cmd); err != nil {
		return model.NewErrorTypeFromString(err.Error())
	}

	return nil
}

func (r *FeatureLocalImpl) processReply(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, requestHeader *model.HeaderType, featureRemote *FeatureRemoteImpl) *model.ErrorType {
	featureRemote.UpdateData(function, data, filterPartial, filterDelete)
	_ = r.pendingRequests.SetData(featureRemote.Device().ski, *requestHeader.MsgCounterReference, data)
	// an error in SetData only means that there is no pendingRequest waiting for this dataset
	// so this is nothing to consider as an error to return

	// the data was updated, so send an event, other event handlers may watch out for this as well
	payload := EventPayload{
		Ski:           featureRemote.Device().ski,
		EventType:     EventTypeDataChange,
		ChangeType:    ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		CmdClassifier: util.Ptr(model.CmdClassifierTypeReply),
		Data:          data,
	}
	Events.Publish(payload)

	return nil
}

func (r *FeatureLocalImpl) processNotify(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, featureRemote *FeatureRemoteImpl) *model.ErrorType {
	featureRemote.UpdateData(function, data, filterPartial, filterDelete)

	payload := EventPayload{
		Ski:           featureRemote.Device().ski,
		EventType:     EventTypeDataChange,
		ChangeType:    ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		CmdClassifier: util.Ptr(model.CmdClassifierTypeNotify),
		Data:          data,
	}
	Events.Publish(payload)

	return nil
}

func (r *FeatureLocalImpl) functionData(function model.FunctionType) FunctionDataCmd {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
