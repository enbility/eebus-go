package features

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type Feature interface {
	SubscribeForEntity() error
}

type FeatureImpl struct {
	featureType model.FeatureTypeType

	featureLocal  spine.FeatureLocal
	featureRemote *spine.FeatureRemoteImpl

	service *service.EEBUSService
	device  *spine.DeviceRemoteImpl
	entity  *spine.EntityRemoteImpl
}

var _ Feature = (*FeatureImpl)(nil)

func NewFeatureImpl(featureType model.FeatureTypeType, service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*FeatureImpl, error) {
	f := &FeatureImpl{
		featureType: featureType,
		service:     service,
		device:      entity.Device(),
		entity:      entity,
	}

	var err error
	f.featureLocal, f.featureRemote, err = service.GetLocalClientAndRemoteServerFeatures(f.featureType, entity)

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
func (f *FeatureImpl) requestData(function model.FunctionType) (*model.MsgCounterType, error) {
	fTypes := f.featureRemote.Operations()
	if _, exists := fTypes[function]; !exists {
		return nil, ErrFunctionNotSupported
	}

	if !fTypes[function].Read {
		return nil, ErrOperationOnFunctionNotSupported
	}

	msgCounter, fErr := f.featureLocal.RequestData(function, f.featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}
