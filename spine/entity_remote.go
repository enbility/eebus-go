package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type EntityRemoteImpl struct {
	*EntityImpl
	device   *DeviceRemoteImpl
	features []*FeatureRemoteImpl
}

func NewEntityRemoteImpl(device *DeviceRemoteImpl, eType model.EntityTypeType, address []model.AddressEntityType) *EntityRemoteImpl {
	return &EntityRemoteImpl{
		EntityImpl: NewEntity(eType, address),
		device:     device,
	}
}

func (r *EntityRemoteImpl) Device() *DeviceRemoteImpl {
	return r.device
}

func (r *EntityRemoteImpl) AddFeature(f *FeatureRemoteImpl) {
	r.features = append(r.features, f)
}
