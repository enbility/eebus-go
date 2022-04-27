package spine

import (
	"reflect"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type DeviceLocalImpl struct {
	*DeviceImpl
	entities            []*EntityLocalImpl
	subscriptionManager SubscriptionManager
	nodeManagement      *NodeManagementImpl
}

func NewDeviceLocalImpl(address model.AddressDeviceType) *DeviceLocalImpl {
	res := &DeviceLocalImpl{
		DeviceImpl:          NewDeviceImpl(address),
		subscriptionManager: NewSubscriptionManager(),
	}

	res.addDeviceInformation()
	return res
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
		entity.AddFeature(r.nodeManagement.FeatureLocalImpl)
	}
	{
		f := NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
		entity.AddFeature(f)
	}

	r.entities = append(r.entities, entity)
}
