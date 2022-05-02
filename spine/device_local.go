package spine

import (
	"errors"
	"reflect"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type DeviceLocalImpl struct {
	*DeviceImpl
	entities            []*EntityLocalImpl
	subscriptionManager SubscriptionManager
	nodeManagement      *NodeManagementImpl

	remoteDevices map[string]*DeviceRemoteImpl

	brandName   string
	deviceModel string
	deviceCode  string
}

func NewDeviceLocalImpl(brandName, deviceModel, deviceCode, deviceAddress string, deviceType model.DeviceTypeType) *DeviceLocalImpl {
	res := &DeviceLocalImpl{
		DeviceImpl:          NewDeviceImpl(model.AddressDeviceType(deviceAddress), deviceType),
		subscriptionManager: NewSubscriptionManager(),
		remoteDevices:       make(map[string]*DeviceRemoteImpl),
		brandName:           brandName,
		deviceModel:         deviceModel,
		deviceCode:          deviceCode,
	}

	res.addDeviceInformation()
	return res
}

func (r *DeviceLocalImpl) AddRemoteDevice(ski, deviceCode string, deviceType model.DeviceTypeType, readC <-chan []byte, writeC chan<- []byte) {
	rDevice := NewDeviceRemoteImpl(r, ski, deviceCode, deviceType, readC, writeC)
	r.remoteDevices[ski] = rDevice

	// Request Detailed Discovery Data
	r.nodeManagement.RequestDetailedDiscovery(&rDevice.address, rDevice.sender)

	Events.Subscribe(r)
}

// React to some specific events
func (r *DeviceLocalImpl) HandleEvent(payload EventPayload) {
	// Subscribe to NodeManagment after DetailedDiscovery is received
	if payload.EventType == EventTypeDeviceChange && payload.Data != nil {
		switch payload.Data.(type) {
		case *model.NodeManagementDetailedDiscoveryDataType:
			payload.Device.sender.Subscribe(r.nodeManagement.Address(), r.nodeManagement.Address(), model.FeatureTypeTypeNodeManagement)
		}
	}
}

func (r *DeviceLocalImpl) RemoveRemoteDevice(ski string) {
	if r.remoteDevices[ski] == nil {
		return
	}
	r.remoteDevices[ski].CloseConnection()
	delete(r.remoteDevices, ski)
}

func (r *DeviceLocalImpl) RemoteDevices() []*DeviceRemoteImpl {
	res := make([]*DeviceRemoteImpl, 0)
	for _, rDevice := range r.remoteDevices {
		res = append(res, rDevice)
	}
	return res
}

func (r *DeviceLocalImpl) ProcessCmd(datagram model.DatagramType, remoteDevice *DeviceRemoteImpl) error {
	destAddr := datagram.Header.AddressDestination
	localFeature := r.FeatureByAddress(destAddr)
	if localFeature == nil {
		return errors.New("invalid feature address")
	}

	cmdClassifier := datagram.Header.CmdClassifier
	cmd := datagram.Payload.Cmd[0]

	// isPartial
	isPartial := false
	functionCmd := cmd.Function
	filterCmd := cmd.Filter

	if functionCmd != nil && filterCmd != nil {
		// TODO check if the function is the same as the provided cmd value
		if len(filterCmd) > 0 {
			if filterCmd[0].CmdControl.Partial != nil {
				isPartial = true
			}
		}
	}

	remoteFeature := remoteDevice.FeatureByAddress(datagram.Header.AddressSource)

	message := &Message{
		RequestHeader: &datagram.Header,
		CmdClassifier: *cmdClassifier,
		Cmd:           cmd,
		IsPartial:     isPartial,
		featureRemote: remoteFeature,
		deviceRemote:  remoteDevice,
	}

	return localFeature.HandleMessage(message)
}

func (r *DeviceLocalImpl) SubscriptionManager() SubscriptionManager {
	return r.subscriptionManager
}

func (r *DeviceLocalImpl) AddEntity(entity *EntityLocalImpl) {
	r.entities = append(r.entities, entity)

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeAdded)
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
				Device: &r.address,
			},
			DeviceType: &r.dType,
			// TODO NetworkFeatureSet
			// NetworkFeatureSet: &smart,
		},
	}
	return &res
}

func (r *DeviceLocalImpl) NotifySubscribers(featureAddress *model.FeatureAddressType, cmd []model.CmdType) {
	subscriptions := r.SubscriptionManager().SubscriptionsOnFeature(*featureAddress)
	for _, subscription := range subscriptions {
		if err := subscription.clientFeature.Sender().Notify(
			subscription.serverFeature.Address(), subscription.clientFeature.Address(), cmd); err != nil {
			// TODO: error handling
		}
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

	cmd := []model.CmdType{{
		Function: util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData),
		Filter:   filterType(true),
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{
			SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
				SpecificationVersion: []model.SpecificationVersionDataType{model.SpecificationVersionDataType(SpecificationVersion)},
			},
			DeviceInformation:  deviceInformation,
			EntityInformation:  []model.NodeManagementDetailedDiscoveryEntityInformationType{entityInformation},
			FeatureInformation: featureInformation,
		},
	}}

	r.NotifySubscribers(r.nodeManagement.Address(), cmd)
}

func (r *DeviceLocalImpl) addDeviceInformation() {
	entityType := model.EntityTypeTypeDeviceInformation
	entity := NewEntityLocalImpl(r, entityType, []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)})

	{
		r.nodeManagement = NewNodeManagementImpl(entity.NextFeatureId(), entity)

		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementDetailedDiscoveryData, true, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementUseCaseData, true, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionData, true, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionRequestCall, false, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionDeleteCall, false, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementBindingData, true, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementBindingRequestCall, false, false)
		r.nodeManagement.AddFunctionType(model.FunctionTypeNodeManagementBindingDeleteCall, false, false)

		entity.AddFeature(r.nodeManagement)
	}
	{
		f := NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)

		f.AddFunctionType(model.FunctionTypeDeviceClassificationManufacturerData, true, false)

		manufacturerData := &model.DeviceClassificationManufacturerDataType{
			BrandName:  util.Ptr(model.DeviceClassificationStringType(r.brandName)),
			DeviceName: util.Ptr(model.DeviceClassificationStringType(r.deviceModel)),
			DeviceCode: util.Ptr(model.DeviceClassificationStringType(r.deviceCode)),
		}
		f.SetData(model.FunctionTypeDeviceClassificationManufacturerData, manufacturerData)

		entity.AddFeature(f)
	}

	r.entities = append(r.entities, entity)
}
