package spine

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

type DeviceLocalImpl struct {
	*DeviceImpl
	entities            []*EntityLocalImpl
	subscriptionManager SubscriptionManager
	bindingManager      BindingManager
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
func NewDeviceLocalImpl(brandName, deviceModel, serialNumber, deviceCode, deviceAddress string, deviceType model.DeviceTypeType, featureSet model.NetworkManagementFeatureSetType) *DeviceLocalImpl {
	address := model.AddressDeviceType(deviceAddress)

	var fSet *model.NetworkManagementFeatureSetType
	if len(featureSet) != 0 {
		fSet = &featureSet
	}

	res := &DeviceLocalImpl{
		DeviceImpl:          NewDeviceImpl(&address, &deviceType, fSet),
		subscriptionManager: NewSubscriptionManager(),
		bindingManager:      NewBindingManager(),
		remoteDevices:       make(map[string]*DeviceRemoteImpl),
		brandName:           brandName,
		deviceModel:         deviceModel,
		serialNumber:        serialNumber,
		deviceCode:          deviceCode,
	}

	res.addDeviceInformation()
	return res
}

func (r *DeviceLocalImpl) AddRemoteDevice(ski string, writeI SpineDataConnection) SpineDataProcessing {
	rDevice := NewDeviceRemoteImpl(r, ski, writeI)

	r.mux.Lock()
	r.remoteDevices[ski] = rDevice
	r.mux.Unlock()

	// Request Detailed Discovery Data
	_, _ = r.nodeManagement.RequestDetailedDiscovery(rDevice.ski, rDevice.address, rDevice.sender)

	// TODO: Add error handling
	// If the request returned an error, it should be retried until it does not

	// always add subscription, as it checks if it already exists
	Events.Subscribe(r)

	return rDevice
}

// React to some specific events
func (r *DeviceLocalImpl) HandleEvent(payload EventPayload) {
	// Subscribe to NodeManagment after DetailedDiscovery is received
	if payload.EventType != EventTypeDeviceChange {
		return
	}

	if payload.ChangeType != ElementChangeAdd {
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
		_, _ = r.nodeManagement.Subscribe(payload.Feature.Device(), payload.Feature.Address())

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

	r.remoteDevices[ski].CloseConnection()
	delete(r.remoteDevices, ski)

	// only unsubscribe if we don't have any remote devices left
	if len(r.remoteDevices) == 0 {
		Events.Unsubscribe(r)
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
		_ = remoteFeature.Sender().ResultError(message.RequestHeader, destAddr, NewErrorType(model.ErrorNumberTypeDestinationUnknown, errorMessage))

		return errors.New(errorMessage)
	}

	lfType := string(localFeature.Type())
	rfType := ""
	if remoteFeature != nil {
		remoteFeature.Type()
	}

	logging.Log.Debug(datagram.PrintMessageOverview(false, lfType, rfType))

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

func (r *DeviceLocalImpl) SubscriptionManager() SubscriptionManager {
	return r.subscriptionManager
}

func (r *DeviceLocalImpl) BindingManager() BindingManager {
	return r.bindingManager
}

func (r *DeviceLocalImpl) AddEntity(entity *EntityLocalImpl) {
	r.entities = append(r.entities, entity)

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeAdded)
}

func (r *DeviceLocalImpl) RemoveEntity(entity *EntityLocalImpl) {
	for i, e := range r.entities {
		if e == entity {
			r.entities = append(r.entities[:i], r.entities[i+1:]...)
			// TODO: delete subscriptions of removed entity (incl. delete call)
			break
		}
	}
}

func (r *DeviceLocalImpl) Entities() []*EntityLocalImpl {
	return r.entities
}

func (r *DeviceLocalImpl) Entity(id []model.AddressEntityType) *EntityLocalImpl {
	for _, e := range r.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
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

func (r *DeviceLocalImpl) FeatureByTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) FeatureLocal {
	if len(r.entities) < 1 {
		return nil
	}

	for _, entity := range r.entities {
		for _, feature := range entity.Features() {
			if feature.Type() == featureType && feature.Role() == role {
				return feature
			}
		}
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

func (r *DeviceLocalImpl) NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType) {
	subscriptions := r.SubscriptionManager().SubscriptionsOnFeature(*featureAddress)
	for _, subscription := range subscriptions {
		// TODO: error handling
		_, _ = subscription.clientFeature.Sender().Notify(subscription.serverFeature.Address(), subscription.clientFeature.Address(), cmd)
	}
}

func (r *DeviceLocalImpl) notifySubscribersOfEntity(entity *EntityLocalImpl, state model.NetworkManagementStateChangeType) {
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
