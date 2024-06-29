package server

import (
	"errors"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Feature struct {
	featureType model.FeatureTypeType

	localRole model.RoleType

	spineLocalDevice spineapi.DeviceLocalInterface
	localEntity      spineapi.EntityLocalInterface

	featureLocal spineapi.FeatureLocalInterface
}

// var _ api.FeatureServerInterface = (*Feature)(nil)

func NewFeature(
	featureType model.FeatureTypeType,
	localEntity spineapi.EntityLocalInterface) (*Feature, error) {
	if localEntity == nil {
		return nil, errors.New("local entity is nil")
	}

	f := &Feature{
		featureType:      featureType,
		localRole:        model.RoleTypeServer,
		spineLocalDevice: localEntity.Device(),
		localEntity:      localEntity,
	}

	f.featureLocal = f.localEntity.FeatureOfTypeAndRole(f.featureType, f.localRole)
	if f.featureLocal == nil {
		return nil, errors.New("local feature not found")
	}

	return f, nil
}
