package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type EntityLocalImpl struct {
	*EntityImpl
	device   *DeviceLocalImpl
	features []*FeatureLocalImpl
}

func NewEntityLocalImpl(device *DeviceLocalImpl, eType model.EntityTypeType, address []model.AddressEntityType) *EntityLocalImpl {
	return &EntityLocalImpl{
		EntityImpl: NewEntity(eType, address),
		device:     device,
	}
}

func (r *EntityLocalImpl) Device() *DeviceLocalImpl {
	return r.device
}

func (r *EntityLocalImpl) AddFeature(f *FeatureLocalImpl) {
	r.features = append(r.features, f)
}
