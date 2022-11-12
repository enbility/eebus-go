package features

import (
	"errors"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type Feature interface {
	SubscribeForEntity() error
}

type FeatureImpl struct {
	featureType model.FeatureTypeType

	localRole  model.RoleType
	remoteRole model.RoleType

	spineLocalDevice *spine.DeviceLocalImpl

	featureLocal  spine.FeatureLocal
	featureRemote *spine.FeatureRemoteImpl

	device *spine.DeviceRemoteImpl
	entity *spine.EntityRemoteImpl
}

var _ Feature = (*FeatureImpl)(nil)

func NewFeatureImpl(featureType model.FeatureTypeType, localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*FeatureImpl, error) {
	f := &FeatureImpl{
		featureType:      featureType,
		localRole:        localRole,
		remoteRole:       remoteRole,
		spineLocalDevice: spineLocalDevice,
		device:           entity.Device(),
		entity:           entity,
	}

	var err error
	f.featureLocal, f.featureRemote, err = f.getLocalClientAndRemoteServerFeatures()

	return f, err
}

// subscribe to the feature for a the entity
func (f *FeatureImpl) SubscribeForEntity() error {
	if _, fErr := f.featureLocal.Subscribe(f.featureRemote.Device(), f.featureRemote.Address()); fErr != nil {
		return errors.New(fErr.String())
	}

	return nil
}

// bind to the feature of a the entity
func (f *FeatureImpl) Bind() error {
	if _, fErr := f.featureLocal.Bind(f.featureRemote.Device(), f.featureRemote.Address()); fErr != nil {
		return errors.New(fErr.String())
	}

	return nil
}

// helper method which adds checking if the feature is available and the operation is allowed
// selectors and elements are used if specific data should be requested by using
// model.FilterType DataSelectors (selectors) and/or DataElements (elements)
// both should use the proper data types for the used function
func (f *FeatureImpl) requestData(function model.FunctionType, selectors any, elements any) (*model.MsgCounterType, error) {
	if f.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	fTypes := f.featureRemote.Operations()
	if _, exists := fTypes[function]; !exists {
		return nil, ErrFunctionNotSupported
	}

	if !fTypes[function].Read {
		return nil, ErrOperationOnFunctionNotSupported
	}

	msgCounter, fErr := f.featureLocal.RequestData(function, selectors, elements, f.featureRemote)
	if fErr != nil {
		logging.Log.Error(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// internal helper method for getting local and remote feature for a given featureType and a given remoteDevice
func (f *FeatureImpl) getLocalClientAndRemoteServerFeatures() (spine.FeatureLocal, *spine.FeatureRemoteImpl, error) {
	if f.entity == nil {
		return nil, nil, errors.New("invalid remote entity provided")
	}

	featureLocal := f.spineLocalDevice.FeatureByTypeAndRole(f.featureType, f.localRole)
	featureRemote := f.entity.Device().FeatureByEntityTypeAndRole(f.entity, f.featureType, f.remoteRole)

	if featureLocal == nil {
		return nil, nil, errors.New("local feature not found")
	}

	if featureRemote == nil {
		return nil, nil, errors.New("remote feature not found")
	}

	return featureLocal, featureRemote, nil
}
