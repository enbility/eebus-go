package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type EntityLocalImpl struct {
	*EntityImpl
	device   *DeviceLocalImpl
	features []FeatureLocal
}

func NewEntityLocalImpl(device *DeviceLocalImpl, eType model.EntityTypeType, entityAddress []model.AddressEntityType) *EntityLocalImpl {
	return &EntityLocalImpl{
		EntityImpl: NewEntity(eType, device.Address(), entityAddress),
		device:     device,
	}
}

func (r *EntityLocalImpl) Device() *DeviceLocalImpl {
	return r.device
}

func (r *EntityLocalImpl) AddFeature(f FeatureLocal) {
	r.features = append(r.features, f)
}

func (r *EntityLocalImpl) Features() []FeatureLocal {
	return r.features
}

func (r *EntityLocalImpl) Feature(addressFeature *model.AddressFeatureType) FeatureLocal {
	for _, f := range r.features {
		if *f.Address().Feature == *addressFeature {
			return f
		}
	}
	return nil
}

func (r *EntityLocalImpl) Information() *model.NodeManagementDetailedDiscoveryEntityInformationType {
	res := model.NodeManagementDetailedDiscoveryEntityInformationType{
		Description: &model.NetworkManagementEntityDescriptionDataType{
			EntityAddress: r.Address(),
			EntityType:    &r.eType,
		},
	}

	return &res
}
