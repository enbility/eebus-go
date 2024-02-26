package features

import (
	"errors"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Feature struct {
	featureType model.FeatureTypeType

	localRole  model.RoleType
	remoteRole model.RoleType

	spineLocalDevice api.DeviceLocalInterface
	localEntity      api.EntityLocalInterface

	featureLocal  api.FeatureLocalInterface
	featureRemote api.FeatureRemoteInterface

	remoteDevice api.DeviceRemoteInterface
	remoteEntity api.EntityRemoteInterface
}

var _ FeatureInterface = (*Feature)(nil)

func NewFeature(featureType model.FeatureTypeType, localRole, remoteRole model.RoleType, localEntity api.EntityLocalInterface, remoteEntity api.EntityRemoteInterface) (*Feature, error) {
	if localEntity == nil {
		return nil, errors.New("local entity is nil")
	}

	if remoteEntity == nil {
		return nil, errors.New("remote entity is nil")
	}

	f := &Feature{
		featureType:      featureType,
		localRole:        localRole,
		remoteRole:       remoteRole,
		spineLocalDevice: localEntity.Device(),
		localEntity:      localEntity,
		remoteDevice:     remoteEntity.Device(),
		remoteEntity:     remoteEntity,
	}

	var err error
	f.featureLocal, f.featureRemote, err = f.getLocalClientAndRemoteServerFeatures()

	return f, err
}

// check if there is a subscription to the remote feature
func (f *Feature) HasSubscription() bool {
	subscription := f.featureLocal.HasSubscriptionToRemote(f.featureRemote.Address())
	return subscription
}

// subscribe to the feature of the entity
func (f *Feature) Subscribe() (*model.MsgCounterType, error) {
	msgCounter, fErr := f.featureLocal.SubscribeToRemote(f.featureRemote.Address())

	if fErr != nil {
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// unssubscribe to the feature of the entity
func (f *Feature) Unsubscribe() (*model.MsgCounterType, error) {
	msgCounter, fErr := f.featureLocal.RemoveRemoteSubscription(f.featureRemote.Address())

	if fErr != nil {
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// check if there is a binding to the remote feature
func (f *Feature) HasBinding() bool {
	binding := f.featureLocal.HasBindingToRemote(f.featureRemote.Address())
	return binding
}

// bind to the feature of the entity
func (f *Feature) Bind() (*model.MsgCounterType, error) {
	msgCounter, fErr := f.featureLocal.BindToRemote(f.featureRemote.Address())
	if fErr != nil {
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// remove a binding to the feature of the entity
func (f *Feature) Unbind() (*model.MsgCounterType, error) {
	msgCounter, fErr := f.featureLocal.RemoveRemoteBinding(f.featureRemote.Address())

	if fErr != nil {
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// add a callback function to be invoked once a result to a msgCounter came in
func (f *Feature) AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResultMessage)) {
	f.featureLocal.AddResultCallback(msgCounterReference, function)
}

// helper method which adds checking if the feature is available and the operation is allowed
// selectors and elements are used if specific data should be requested by using
// model.FilterType DataSelectors (selectors) and/or DataElements (elements)
// both should use the proper data types for the used function
func (f *Feature) requestData(function model.FunctionType, selectors any, elements any) (*model.MsgCounterType, error) {
	if f.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	fTypes := f.featureRemote.Operations()
	if _, exists := fTypes[function]; !exists {
		return nil, ErrFunctionNotSupported
	}

	if !fTypes[function].Read() {
		return nil, ErrOperationOnFunctionNotSupported
	}

	msgCounter, fErr := f.featureLocal.RequestRemoteData(function, selectors, elements, f.featureRemote)
	if fErr != nil {
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// internal helper method for getting local and remote feature for a given featureType and a given remoteDevice
func (f *Feature) getLocalClientAndRemoteServerFeatures() (api.FeatureLocalInterface, api.FeatureRemoteInterface, error) {
	featureLocal := f.localEntity.FeatureOfTypeAndRole(f.featureType, f.localRole)
	if featureLocal == nil && f.localRole == model.RoleTypeClient {
		featureLocal = f.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeGeneric, f.localRole)
	}
	if featureLocal == nil {
		return nil, nil, errors.New("local feature not found")
	}

	featureRemote := f.remoteEntity.Device().FeatureByEntityTypeAndRole(f.remoteEntity, f.featureType, f.remoteRole)
	if featureRemote == nil && f.localRole == model.RoleTypeClient {
		featureRemote = f.remoteEntity.Device().FeatureByEntityTypeAndRole(f.remoteEntity, model.FeatureTypeTypeGeneric, f.localRole)
	}

	if featureRemote == nil {
		return nil, nil, errors.New("remote feature not found")
	}

	return featureLocal, featureRemote, nil
}
