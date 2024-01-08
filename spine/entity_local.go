package spine

import (
	"github.com/enbility/eebus-go/spine/model"
)

var _ EntityLocal = (*EntityLocalImpl)(nil)

type EntityLocalImpl struct {
	*EntityImpl
	device   DeviceLocal
	features []FeatureLocal
}

func NewEntityLocalImpl(device DeviceLocal, eType model.EntityTypeType, entityAddress []model.AddressEntityType) *EntityLocalImpl {
	return &EntityLocalImpl{
		EntityImpl: NewEntity(eType, device.Address(), entityAddress),
		device:     device,
	}
}

func (r *EntityLocalImpl) Device() DeviceLocal {
	return r.device
}

// Add a feature to the entity if it is not already added
func (r *EntityLocalImpl) AddFeature(f FeatureLocal) {
	// check if this feature is already added
	for _, f2 := range r.features {
		if f2.Type() == f.Type() && f2.Role() == f.Role() {
			return
		}
	}
	r.features = append(r.features, f)
}

// either returns an existing feature or creates a new one
// for a given entity, featuretype and role
func (r *EntityLocalImpl) GetOrAddFeature(featureType model.FeatureTypeType, role model.RoleType) FeatureLocal {
	if f := r.FeatureOfTypeAndRole(featureType, role); f != nil {
		return f
	}
	f := NewFeatureLocalImpl(r.NextFeatureId(), r, featureType, role)

	description := string(featureType)
	switch role {
	case model.RoleTypeClient:
		description += " Client"
	case model.RoleTypeServer:
		description += " Server"
	}
	f.SetDescriptionString(description)
	r.features = append(r.features, f)

	if role == model.RoleTypeServer && featureType == model.FeatureTypeTypeDeviceDiagnosis {
		// Update HeartbeatManagerImpl
		r.device.HeartbeatManager().SetLocalFeature(r, f)
	}

	return f
}

func (r *EntityLocalImpl) FeatureOfTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) FeatureLocal {
	for _, f := range r.features {
		if f.Type() == featureType && f.Role() == role {
			return f
		}
	}
	return nil
}

func (r *EntityLocalImpl) Features() []FeatureLocal {
	return r.features
}

func (r *EntityLocalImpl) Feature(addressFeature *model.AddressFeatureType) FeatureLocal {
	if addressFeature == nil {
		return nil
	}
	for _, f := range r.features {
		if *f.Address().Feature == *addressFeature {
			return f
		}
	}
	return nil
}

func (r *EntityLocalImpl) Information() *model.NodeManagementDetailedDiscoveryEntityInformationType {
	res := &model.NodeManagementDetailedDiscoveryEntityInformationType{
		Description: &model.NetworkManagementEntityDescriptionDataType{
			EntityAddress: r.Address(),
			EntityType:    &r.eType,
		},
	}

	return res
}

// add a new usecase
func (r *EntityLocalImpl) AddUseCaseSupport(
	actor model.UseCaseActorType,
	useCaseName model.UseCaseNameType,
	useCaseVersion model.SpecificationVersionType,
	useCaseDocumemtSubRevision string,
	useCaseAvailable bool,
	scenarios []model.UseCaseScenarioSupportType,
) {
	nodeMgmt := r.device.NodeManagement()

	data := nodeMgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)
	if data == nil {
		data = &model.NodeManagementUseCaseDataType{}
	}

	address := model.FeatureAddressType{
		Device: r.address.Device,
		Entity: r.address.Entity,
	}

	data.AddUseCaseSupport(address, actor, useCaseName, useCaseVersion, useCaseDocumemtSubRevision, useCaseAvailable, scenarios)

	nodeMgmt.SetData(model.FunctionTypeNodeManagementUseCaseData, data)
}

// Remove a usecase with a given actor ans usecase name
func (r *EntityLocalImpl) RemoveUseCaseSupport(
	actor model.UseCaseActorType,
	useCaseName model.UseCaseNameType,
) {
	nodeMgmt := r.device.NodeManagement()

	data := nodeMgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)
	if data == nil {
		return
	}

	address := model.FeatureAddressType{
		Device: r.address.Device,
		Entity: r.address.Entity,
	}

	data.RemoveUseCaseSupport(address, actor, useCaseName)

	nodeMgmt.SetData(model.FunctionTypeNodeManagementUseCaseData, data)
}

// Remove all usecases
func (r *EntityLocalImpl) RemoveAllUseCaseSupports() {
	r.RemoveUseCaseSupport("", "")
}

// Remove all subscriptions
func (r *EntityLocalImpl) RemoveAllSubscriptions() {
	for _, item := range r.features {
		item.RemoveAllSubscriptions()
	}
}

// Remove all bindings
func (r *EntityLocalImpl) RemoveAllBindings() {
	for _, item := range r.features {
		item.RemoveAllBindings()
	}
}
