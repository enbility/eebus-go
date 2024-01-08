package spine

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

type DeviceLocalImpl struct {
	*DeviceImpl
	entities            []EntityLocal
	subscriptionManager SubscriptionManager
	bindingManager      BindingManager
	heartbeatManager    HeartbeatManager
	nodeManagement      *NodeManagementImpl

	remoteDevices map[string]*DeviceRemoteImpl

	brandName    string
	deviceModel  string
	deviceCode   string
	serialNumber string

	mux sync.Mutex
}

// BrandName is the brand
// DeviceModel is the model
// SerialNumber is the serial number
// DeviceCode is the SHIP id (accessMethods.id)
// DeviceAddress is the SPINE device address
func NewDeviceLocalImpl(brandName, deviceModel, serialNumber, deviceCode, deviceAddress string, deviceType model.DeviceTypeType, featureSet model.NetworkManagementFeatureSetType, heartbeatTimeout time.Duration) *DeviceLocalImpl {
	address := model.AddressDeviceType(deviceAddress)

	var fSet *model.NetworkManagementFeatureSetType
	if len(featureSet) != 0 {
		fSet = &featureSet
	}

	res := &DeviceLocalImpl{
		DeviceImpl:    NewDeviceImpl(&address, &deviceType, fSet),
		remoteDevices: make(map[string]*DeviceRemoteImpl),
		brandName:     brandName,
		deviceModel:   deviceModel,
		serialNumber:  serialNumber,
		deviceCode:    deviceCode,
	}

	res.subscriptionManager = NewSubscriptionManager(res)
	res.bindingManager = NewBindingManager(res)
	res.heartbeatManager = NewHeartbeatManager(res, res.subscriptionManager, heartbeatTimeout)

	res.addDeviceInformation()
	return res
}

var _ DeviceLocalConnection = (*DeviceLocalImpl)(nil)

func (r *DeviceLocalImpl) RemoveRemoteDeviceConnection(ski string) {
	remoteDevice := r.RemoteDeviceForSki(ski)

	r.RemoveRemoteDevice(ski)

	// inform about the disconnection
	payload := EventPayload{
		Ski:        ski,
		EventType:  EventTypeDeviceChange,
		ChangeType: ElementChangeRemove,
		Device:     remoteDevice,
	}
	Events.Publish(payload)
}

// Helper method used by tests and AddRemoteDevice
func (r *DeviceLocalImpl) AddRemoteDeviceForSki(ski string, rDevice *DeviceRemoteImpl) {
	r.mux.Lock()
	r.remoteDevices[ski] = rDevice
	r.mux.Unlock()
}

// Adds a new remote device with a given SKI and triggers SPINE requesting device details
func (r *DeviceLocalImpl) AddRemoteDevice(ski string, writeI SpineDataConnection) SpineDataProcessing {
	sender := NewSender(writeI)
	rDevice := NewDeviceRemoteImpl(r, ski, sender)

	r.AddRemoteDeviceForSki(ski, rDevice)

	// Request Detailed Discovery Data
	_, _ = r.nodeManagement.RequestDetailedDiscovery(rDevice.ski, rDevice.address, rDevice.sender)

	// TODO: Add error handling
	// If the request returned an error, it should be retried until it does not

	// always add subscription, as it checks if it already exists
	_ = Events.subscribe(EventHandlerLevelCore, r)

	return rDevice
}

// React to some specific events
func (r *DeviceLocalImpl) HandleEvent(payload EventPayload) {
	// Subscribe to NodeManagment after DetailedDiscovery is received
	if payload.EventType != EventTypeDeviceChange || payload.ChangeType != ElementChangeAdd {
		return
	}

	if payload.Data == nil {
		return
	}

	if len(payload.Ski) == 0 {
		return
	}

	if r.RemoteDeviceForSki(payload.Ski) == nil {
		return
	}

	switch payload.Data.(type) {
	case *model.NodeManagementDetailedDiscoveryDataType:
		_, _ = r.nodeManagement.Subscribe(payload.Feature.Address())

		// Request Use Case Data
		_, _ = r.nodeManagement.RequestUseCaseData(payload.Device.ski, payload.Device.Address(), payload.Device.Sender())
	}
}

func (r *DeviceLocalImpl) RemoveRemoteDevice(ski string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.remoteDevices[ski] == nil {
		return
	}

	// remove all subscriptions for this device
	subscriptionMgr := r.SubscriptionManager()
	subscriptionMgr.RemoveSubscriptionsForDevice(r.remoteDevices[ski])

	// make sure Heartbeat Manager is up to date
	r.HeartbeatManager().UpdateHeartbeatOnSubscriptions()

	// remove all bindings for this device
	bindingMgr := r.BindingManager()
	bindingMgr.RemoveBindingsForDevice(r.remoteDevices[ski])

	delete(r.remoteDevices, ski)

	// only unsubscribe if we don't have any remote devices left
	if len(r.remoteDevices) == 0 {
		_ = Events.unsubscribe(EventHandlerLevelCore, r)
	}
}

func (r *DeviceLocalImpl) RemoteDevices() []*DeviceRemoteImpl {
	r.mux.Lock()
	defer r.mux.Unlock()

	res := make([]*DeviceRemoteImpl, 0)
	for _, rDevice := range r.remoteDevices {
		res = append(res, rDevice)
	}

	return res
}

func (r *DeviceLocalImpl) RemoteDeviceForAddress(address model.AddressDeviceType) *DeviceRemoteImpl {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.remoteDevices {
		if *item.address == address {
			return item
		}
	}

	return nil
}

func (r *DeviceLocalImpl) RemoteDeviceForSki(ski string) *DeviceRemoteImpl {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.remoteDevices[ski]
}

func (r *DeviceLocalImpl) ProcessCmd(datagram model.DatagramType, remoteDevice *DeviceRemoteImpl) error {
	destAddr := datagram.Header.AddressDestination
	localFeature := r.FeatureByAddress(destAddr)

	cmdClassifier := datagram.Header.CmdClassifier
	if len(datagram.Payload.Cmd) == 0 {
		return errors.New("no payload cmd content available")
	}
	cmd := datagram.Payload.Cmd[0]

	// TODO check if cmd.Function is the same as the provided cmd value
	filterPartial, filterDelete := cmd.ExtractFilter()

	remoteEntity := remoteDevice.Entity(datagram.Header.AddressSource.Entity)
	remoteFeature := remoteDevice.FeatureByAddress(datagram.Header.AddressSource)
	if remoteFeature == nil {
		return fmt.Errorf("invalid remote feature address: '%s'", datagram.Header.AddressSource)
	}

	message := &Message{
		RequestHeader: &datagram.Header,
		CmdClassifier: *cmdClassifier,
		Cmd:           cmd,
		FilterPartial: filterPartial,
		FilterDelete:  filterDelete,
		FeatureRemote: remoteFeature,
		EntityRemote:  remoteEntity,
		DeviceRemote:  remoteDevice,
	}

	ackRequest := datagram.Header.AckRequest

	if localFeature == nil {
		errorMessage := "invalid feature address"
		_ = remoteFeature.Sender().ResultError(message.RequestHeader, destAddr, model.NewErrorType(model.ErrorNumberTypeDestinationUnknown, errorMessage))

		return errors.New(errorMessage)
	}

	lfType := string(localFeature.Type())
	rfType := ""
	if remoteFeature != nil {
		remoteFeature.Type()
	}

	logging.Log().Debug(datagram.PrintMessageOverview(false, lfType, rfType))

	err := localFeature.HandleMessage(message)
	if err != nil {
		// TODO: add error description in a useful format

		// Don't send error responses for incoming resulterror messages
		if message.CmdClassifier != model.CmdClassifierTypeResult {
			_ = remoteFeature.Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
		}

		return errors.New(err.String())
	}
	if ackRequest != nil && *ackRequest {
		_ = remoteFeature.Sender().ResultSuccess(message.RequestHeader, localFeature.Address())
	}

	return nil
}

func (r *DeviceLocalImpl) NodeManagement() *NodeManagementImpl {
	return r.nodeManagement
}

func (r *DeviceLocalImpl) SubscriptionManager() SubscriptionManager {
	return r.subscriptionManager
}

func (r *DeviceLocalImpl) BindingManager() BindingManager {
	return r.bindingManager
}

func (r *DeviceLocalImpl) HeartbeatManager() HeartbeatManager {
	return r.heartbeatManager
}

func (r *DeviceLocalImpl) AddEntity(entity EntityLocal) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.entities = append(r.entities, entity)

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeAdded)
}

func (r *DeviceLocalImpl) RemoveEntity(entity EntityLocal) {
	entity.RemoveAllUseCaseSupports()
	entity.RemoveAllSubscriptions()
	entity.RemoveAllBindings()

	r.mux.Lock()
	defer r.mux.Unlock()

	var entities []EntityLocal
	for _, e := range r.entities {
		if e != entity {
			entities = append(entities, e)
		}
	}

	r.entities = entities
}

func (r *DeviceLocalImpl) Entities() []EntityLocal {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.entities
}

func (r *DeviceLocalImpl) Entity(id []model.AddressEntityType) EntityLocal {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, e := range r.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

func (r *DeviceLocalImpl) EntityForType(entityType model.EntityTypeType) EntityLocal {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, e := range r.entities {
		if e.EntityType() == entityType {
			return e
		}
	}
	return nil
}

func (r *DeviceLocalImpl) FeatureByAddress(address *model.FeatureAddressType) FeatureLocal {
	entity := r.Entity(address.Entity)
	if entity != nil {
		return entity.Feature(address.Feature)
	}
	return nil
}

func (r *DeviceLocalImpl) Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType {
	res := model.NodeManagementDetailedDiscoveryDeviceInformationType{
		Description: &model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{
				Device: r.address,
			},
			DeviceType:        r.dType,
			NetworkFeatureSet: r.featureSet,
		},
	}
	return &res
}

// send a notify message to all remote devices
func (r *DeviceLocalImpl) NotifyUseCaseData() {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, remoteDevice := range r.remoteDevices {
		// TODO: add error management
		_, _ = r.nodeManagement.NotifyUseCaseData(remoteDevice)
	}
}

func (r *DeviceLocalImpl) NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType) {
	subscriptions := r.SubscriptionManager().SubscriptionsOnFeature(*featureAddress)
	for _, subscription := range subscriptions {
		// TODO: error handling
		_, _ = subscription.clientFeature.Sender().Notify(subscription.serverFeature.Address(), subscription.clientFeature.Address(), cmd)
	}
}

func (r *DeviceLocalImpl) notifySubscribersOfEntity(entity EntityLocal, state model.NetworkManagementStateChangeType) {
	deviceInformation := r.Information()
	entityInformation := *entity.Information()
	entityInformation.Description.LastStateChange = &state

	var featureInformation []model.NodeManagementDetailedDiscoveryFeatureInformationType
	if state == model.NetworkManagementStateChangeTypeAdded {
		for _, f := range entity.Features() {
			featureInformation = append(featureInformation, *f.Information())
		}
	}

	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData),
		Filter:   filterEmptyPartial(),
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{
			SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
				SpecificationVersion: []model.SpecificationVersionDataType{model.SpecificationVersionDataType(SpecificationVersion)},
			},
			DeviceInformation:  deviceInformation,
			EntityInformation:  []model.NodeManagementDetailedDiscoveryEntityInformationType{entityInformation},
			FeatureInformation: featureInformation,
		},
	}

	r.NotifySubscribers(r.nodeManagement.Address(), cmd)
}

func (r *DeviceLocalImpl) addDeviceInformation() {
	entityType := model.EntityTypeTypeDeviceInformation
	entity := NewEntityLocalImpl(r, entityType, []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)})

	{
		r.nodeManagement = NewNodeManagementImpl(entity.NextFeatureId(), entity)
		entity.AddFeature(r.nodeManagement)
	}
	{
		f := NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)

		f.AddFunctionType(model.FunctionTypeDeviceClassificationManufacturerData, true, false)

		manufacturerData := &model.DeviceClassificationManufacturerDataType{
			BrandName:    util.Ptr(model.DeviceClassificationStringType(r.brandName)),
			VendorName:   util.Ptr(model.DeviceClassificationStringType(r.brandName)),
			DeviceName:   util.Ptr(model.DeviceClassificationStringType(r.deviceModel)),
			DeviceCode:   util.Ptr(model.DeviceClassificationStringType(r.deviceCode)),
			SerialNumber: util.Ptr(model.DeviceClassificationStringType(r.serialNumber)),
		}
		f.SetData(model.FunctionTypeDeviceClassificationManufacturerData, manufacturerData)

		entity.AddFeature(f)
	}

	r.entities = append(r.entities, entity)
}
