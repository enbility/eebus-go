package spine

import (
	"time"

	"github.com/enbility/eebus-go/spine/model"
)

type EventHandler interface {
	HandleEvent(EventPayload)
}

/* Device */

type Device interface {
	Address() *model.AddressDeviceType
	DeviceType() *model.DeviceTypeType
	FeatureSet() *model.NetworkManagementFeatureSetType
	DestinationData() model.NodeManagementDestinationDataType
}

type DeviceLocal interface {
	Device
	RemoveRemoteDeviceConnection(ski string)
	AddRemoteDeviceForSki(ski string, rDevice DeviceRemote)
	AddRemoteDevice(ski string, writeI SpineDataConnection) SpineDataProcessing
	RemoveRemoteDevice(ski string)
	RemoteDevices() []DeviceRemote
	RemoteDeviceForAddress(address model.AddressDeviceType) DeviceRemote
	RemoteDeviceForSki(ski string) DeviceRemote
	ProcessCmd(datagram model.DatagramType, remoteDevice DeviceRemote) error
	NodeManagement() *NodeManagementImpl
	SubscriptionManager() SubscriptionManager
	BindingManager() BindingManager
	HeartbeatManager() HeartbeatManager
	AddEntity(entity EntityLocal)
	RemoveEntity(entity EntityLocal)
	Entities() []EntityLocal
	Entity(id []model.AddressEntityType) EntityLocal
	EntityForType(entityType model.EntityTypeType) EntityLocal
	FeatureByAddress(address *model.FeatureAddressType) FeatureLocal
	NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType)
	Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType
}

type DeviceRemote interface {
	Device
	Ski() string
	SetAddress(address *model.AddressDeviceType)
	HandleIncomingSpineMesssage(message []byte) (*model.MsgCounterType, error)
	Sender() Sender
	Entity(id []model.AddressEntityType) EntityRemote
	Entities() []EntityRemote
	FeatureByAddress(address *model.FeatureAddressType) FeatureRemote
	RemoveByAddress(addr []model.AddressEntityType) EntityRemote
	FeatureByEntityTypeAndRole(entity EntityRemote, featureType model.FeatureTypeType, role model.RoleType) FeatureRemote
	UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType)
	AddEntityAndFeatures(initialData bool, data *model.NodeManagementDetailedDiscoveryDataType) ([]EntityRemote, error)
	AddEntity(entity EntityRemote) EntityRemote
	VerifyUseCaseScenariosAndFeaturesSupport(
		usecaseActor model.UseCaseActorType,
		usecaseName model.UseCaseNameType,
		scenarios []model.UseCaseScenarioSupportType,
		serverFeatures []model.FeatureTypeType,
	) bool
	CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error
}

// implemented by spine.DeviceLocalImpl and used by shipConnection
type DeviceLocalConnection interface {
	RemoveRemoteDeviceConnection(ski string)
	AddRemoteDevice(ski string, writeI SpineDataConnection) SpineDataProcessing
}

/* Entity */

type Entity interface {
	EntityType() model.EntityTypeType
	Address() *model.EntityAddressType
	Description() *model.DescriptionType
	SetDescription(d *model.DescriptionType)
	NextFeatureId() uint
}

type EntityLocal interface {
	Entity
	Device() DeviceLocal
	AddFeature(f FeatureLocal)
	GetOrAddFeature(featureType model.FeatureTypeType, role model.RoleType) FeatureLocal
	FeatureOfTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) FeatureLocal
	Features() []FeatureLocal
	Feature(addressFeature *model.AddressFeatureType) FeatureLocal
	Information() *model.NodeManagementDetailedDiscoveryEntityInformationType
	AddUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType,
		useCaseVersion model.SpecificationVersionType,
		useCaseDocumemtSubRevision string,
		useCaseAvailable bool,
		scenarios []model.UseCaseScenarioSupportType,
	)
	RemoveUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType,
	)
	RemoveAllUseCaseSupports()
	RemoveAllSubscriptions()
	RemoveAllBindings()
}

type EntityRemote interface {
	Entity
	Device() DeviceRemote
	AddFeature(f FeatureRemote)
	Features() []FeatureRemote
	Feature(addressFeature *model.AddressFeatureType) FeatureRemote
	RemoveAllFeatures()
}

/* Feature */

type Feature interface {
	Address() *model.FeatureAddressType
	Type() model.FeatureTypeType
	Role() model.RoleType
	Operations() map[model.FunctionType]*Operations
}

type FeatureRemote interface {
	Feature
	DataCopy(function model.FunctionType) any
	SetData(function model.FunctionType, data any)
	UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
	SetDescription(desc *model.DescriptionType)
	Sender() Sender
	Device() DeviceRemote
	Entity() EntityRemote
	SetOperations(functions []model.FunctionPropertyType)
	SetMaxResponseDelay(delay *model.MaxResponseDelayType)
	MaxResponseDelayDuration() time.Duration
}

type FeatureLocal interface {
	Feature
	DataCopy(function model.FunctionType) any
	SetData(function model.FunctionType, data any)
	AddResultHandler(handler FeatureResult)
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg ResultMessage))
	Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType
	AddFunctionType(function model.FunctionType, read, write bool)
	RequestData(
		function model.FunctionType,
		selector any,
		elements any,
		destination FeatureRemote) (*model.MsgCounterType, *model.ErrorType)
	RequestDataBySenderAddress(
		cmd model.CmdType,
		sender Sender,
		destinationSki string,
		destinationAddress *model.FeatureAddressType,
		maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType)
	FetchRequestData(
		msgCounter model.MsgCounterType,
		destination FeatureRemote) (any, *model.ErrorType)
	RequestAndFetchData(
		function model.FunctionType,
		selector any,
		elements any,
		destination FeatureRemote) (any, *model.ErrorType)
	Subscribe(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// SubscribeAndWait(remoteDevice DeviceRemote, remoteAdress *model.FeatureAddressType) *ErrorType // Subscribes the local feature to the given destination feature; the go routine will block until the response is processed
	RemoveSubscription(remoteAddress *model.FeatureAddressType)
	RemoveAllSubscriptions()
	Bind(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// BindAndWait(remoteDevice DeviceRemote, remoteAddress *model.FeatureAddressType) *ErrorType
	RemoveBinding(remoteAddress *model.FeatureAddressType)
	RemoveAllBindings()
	NotifyData(
		function model.FunctionType,
		deleteSelector, partialSelector any,
		partialWithoutSelector bool,
		deleteElements any,
		destination FeatureRemote) (*model.MsgCounterType, *model.ErrorType)
	WriteData(
		function model.FunctionType,
		deleteSelector, partialSelector any,
		deleteElements any,
		destination FeatureRemote) (*model.MsgCounterType, *model.ErrorType)
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
	DataCopyAny() any
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
	AddBinding(remoteDevice DeviceRemote, data model.BindingManagementRequestCallType) error
	RemoveBinding(data model.BindingManagementDeleteCallType, remoteDevice DeviceRemote) error
	RemoveBindingsForDevice(remoteDevice DeviceRemote)
	RemoveBindingsForEntity(remoteEntity EntityRemote)
	Bindings(remoteDevice DeviceRemote) []*BindingEntry
	BindingsOnFeature(featureAddress model.FeatureAddressType) []*BindingEntry
}

/* Subscription Manager */

type SubscriptionManager interface {
	AddSubscription(remoteDevice DeviceRemote, data model.SubscriptionManagementRequestCallType) error
	RemoveSubscription(data model.SubscriptionManagementDeleteCallType, remoteDevice DeviceRemote) error
	RemoveSubscriptionsForDevice(remoteDevice DeviceRemote)
	RemoveSubscriptionsForEntity(remoteEntity EntityRemote)
	Subscriptions(remoteDevice DeviceRemote) []*SubscriptionEntry
	SubscriptionsOnFeature(featureAddress model.FeatureAddressType) []*SubscriptionEntry
}

/* Heartbeats */

type HeartbeatManager interface {
	IsHeartbeatRunning() bool
	UpdateHeartbeatOnSubscriptions()
	SetLocalFeature(entity *EntityLocalImpl, feature FeatureLocal)
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
