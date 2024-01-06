package spine

import (
	"time"

	"github.com/enbility/eebus-go/spine/model"
)

type EventHandler interface {
	HandleEvent(EventPayload)
}

/* Device */

// implemented by spine.DeviceLocalImpl and used by shipConnection
type DeviceLocalConnection interface {
	RemoveRemoteDeviceConnection(ski string)
	AddRemoteDevice(ski string, writeI SpineDataConnection) SpineDataProcessing
}

/* Feature */

type Feature interface {
	Address() *model.FeatureAddressType
	Type() model.FeatureTypeType
	Role() model.RoleType
}

type FeatureLocal interface {
	Feature
	Data(function model.FunctionType) any
	SetData(function model.FunctionType, data any)
	AddResultHandler(handler FeatureResult)
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg ResultMessage))
	Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType
	AddFunctionType(function model.FunctionType, read, write bool)
	RequestData(
		function model.FunctionType,
		selector any,
		elements any,
		destination *FeatureRemoteImpl) (*model.MsgCounterType, *model.ErrorType)
	RequestDataBySenderAddress(
		cmd model.CmdType,
		sender Sender,
		destinationSki string,
		destinationAddress *model.FeatureAddressType,
		maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType)
	FetchRequestData(
		msgCounter model.MsgCounterType,
		destination *FeatureRemoteImpl) (any, *model.ErrorType)
	RequestAndFetchData(
		function model.FunctionType,
		selector any,
		elements any,
		destination *FeatureRemoteImpl) (any, *model.ErrorType)
	Subscribe(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// SubscribeAndWait(remoteDevice *DeviceRemoteImpl, remoteAdress *model.FeatureAddressType) *ErrorType // Subscribes the local feature to the given destination feature; the go routine will block until the response is processed
	RemoveSubscription(remoteAddress *model.FeatureAddressType)
	RemoveAllSubscriptions()
	Bind(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// BindAndWait(remoteDevice *DeviceRemoteImpl, remoteAddress *model.FeatureAddressType) *ErrorType
	RemoveBinding(remoteAddress *model.FeatureAddressType)
	RemoveAllBindings()
	NotifyData(
		function model.FunctionType,
		deleteSelector, partialSelector any,
		partialWithoutSelector bool,
		deleteElements any,
		destination *FeatureRemoteImpl) (*model.MsgCounterType, *model.ErrorType)
	WriteData(
		function model.FunctionType,
		deleteSelector, partialSelector any,
		deleteElements any,
		destination *FeatureRemoteImpl) (*model.MsgCounterType, *model.ErrorType)
	HandleMessage(message *Message) *model.ErrorType
}

type FeatureResult interface {
	HandleResult(ResultMessage)
}

/* Functions */

type FunctionDataCmd interface {
	FunctionData
	ReadCmdType(partialSelector any, elements any) model.CmdType
	ReplyCmdType(partial bool) model.CmdType
	NotifyCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType
	WriteCmdType(deleteSelector, partialSelector any, deleteElements any) model.CmdType
}

type FunctionData interface {
	Function() model.FunctionType
	DataAny() any
	UpdateDataAny(data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
}

/* Sender */

type ComControl interface {
	// This must be connected to the correct remote device !!
	SendSpineMessage(datagram model.DatagramType) error
}

//go:generate mockery --name=Sender

type Sender interface {
	// Sends a read cmd to request some data
	Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error)
	// Sends a result cmd with no error to indicate that a message was processed successfully
	ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error
	// Sends a result cmd with error information to indicate that a message processing failed
	ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error
	// Sends a reply cmd to response to a read cmd
	Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error
	// Sends a call cmd with a subscription request
	Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a subscription delete request
	Unsubscribe(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding request
	Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding delte request
	Unbind(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error)
	// Sends a notify cmd to indicate that a subscribed feature changed
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// Sends a write cmd, setting properties of remote features
	Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// return the datagram for a given msgCounter (only availbe for Notify messasges!), error if not found
	DatagramForMsgCounter(msgCounter model.MsgCounterType) (model.DatagramType, error)
}

/* PendingRequests */

type PendingRequests interface {
	Add(ski string, counter model.MsgCounterType, maxDelay time.Duration)
	SetData(ski string, counter model.MsgCounterType, data any) *model.ErrorType
	SetResult(ski string, counter model.MsgCounterType, errorResult *model.ErrorType) *model.ErrorType
	GetData(ski string, counter model.MsgCounterType) (any, *model.ErrorType)
	Remove(ski string, counter model.MsgCounterType) *model.ErrorType
}

/* Bindings */

// implemented by spine.BindingManagerImpl
type BindingManager interface {
	AddBinding(remoteDevice *DeviceRemoteImpl, data model.BindingManagementRequestCallType) error
	RemoveBinding(data model.BindingManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error
	RemoveBindingsForDevice(remoteDevice *DeviceRemoteImpl)
	RemoveBindingsForEntity(remoteEntity *EntityRemoteImpl)
	Bindings(remoteDevice *DeviceRemoteImpl) []*BindingEntry
	BindingsOnFeature(featureAddress model.FeatureAddressType) []*BindingEntry
}

/* Subscription Manager */

type SubscriptionManager interface {
	AddSubscription(remoteDevice *DeviceRemoteImpl, data model.SubscriptionManagementRequestCallType) error
	RemoveSubscription(data model.SubscriptionManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error
	RemoveSubscriptionsForDevice(remoteDevice *DeviceRemoteImpl)
	RemoveSubscriptionsForEntity(remoteEntity *EntityRemoteImpl)
	Subscriptions(remoteDevice *DeviceRemoteImpl) []*SubscriptionEntry
	SubscriptionsOnFeature(featureAddress model.FeatureAddressType) []*SubscriptionEntry
}

/* Heartbeats */

type HeartbeatManager interface {
	IsHeartbeatRunning() bool
	UpdateHeartbeatOnSubscriptions()
	StartHeartbeat() error
	StopHeartbeat()
}

/* Processing */

//go:generate mockery --name=SpineDataProcessing

// Used to pass an incoming SPINE message from a SHIP connection to the proper DeviceRemoteImpl
//
// Implemented by DeviceRemoteImpl, used by ShipConnection
type SpineDataProcessing interface {
	HandleIncomingSpineMesssage(message []byte) (*model.MsgCounterType, error)
}

//go:generate mockery --name=SpineDataConnection

// Used to pass an outgoing SPINE message from a DeviceLocalImpl to the SHIP connection
//
// Implemented by ShipConnection, used by DeviceLocalImpl
type SpineDataConnection interface {
	WriteSpineMessage(message []byte)
}
