package spine

import (
	"github.com/enbility/eebus-go/spine/model"
)

type EntityRemoteImpl struct {
	*EntityImpl
	device   *DeviceRemoteImpl
	features []*FeatureRemoteImpl
}

func NewEntityRemoteImpl(device *DeviceRemoteImpl, eType model.EntityTypeType, entityAddress []model.AddressEntityType) *EntityRemoteImpl {
	return &EntityRemoteImpl{
		EntityImpl: NewEntity(eType, device.Address(), entityAddress),
		device:     device,
	}
}

func (r *EntityRemoteImpl) Device() *DeviceRemoteImpl {
	return r.device
}

func (r *EntityRemoteImpl) AddFeature(f *FeatureRemoteImpl) {
	r.features = append(r.features, f)
}

func (r *EntityRemoteImpl) Features() []*FeatureRemoteImpl {
	return r.features
}

func (r *EntityRemoteImpl) Feature(addressFeature *model.AddressFeatureType) *FeatureRemoteImpl {
	for _, f := range r.features {
		if *f.Address().Feature == *addressFeature {
			return f
		}
	}
	return nil
}

func (r *EntityRemoteImpl) RemoveAllFeatures() {
	r.features = nil
}
