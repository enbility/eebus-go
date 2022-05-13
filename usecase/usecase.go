package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/ahmetb/go-linq/v3"
)

var entityTypeActorMap = map[model.EntityTypeType]model.UseCaseActorType{
	model.EntityTypeTypeEV:   model.UseCaseActorTypeEV,
	model.EntityTypeTypeEVSE: model.UseCaseActorTypeEVSE,
	model.EntityTypeTypeCEM:  model.UseCaseActorTypeCEM,
}

var useCaseValidActorsMap = map[model.UseCaseNameType][]model.UseCaseActorType{
	model.UseCaseNameTypeEVSECommissioningAndConfiguration: {model.UseCaseActorTypeEVSE, model.UseCaseActorTypeCEM},
	model.UseCaseNameTypeEVCommissioningAndConfiguration:   {model.UseCaseActorTypeEV, model.UseCaseActorTypeCEM},
}

type UseCaseImpl struct {
	Entity *spine.EntityLocalImpl
	Actor  model.UseCaseActorType

	name            model.UseCaseNameType
	scenarioSupport []model.UseCaseScenarioSupportType
}

func NewUseCase(entity *spine.EntityLocalImpl, ucEnumType model.UseCaseNameType, scenarioSupport []model.UseCaseScenarioSupportType) *UseCaseImpl {
	checkArguments(*entity.EntityImpl, ucEnumType)

	actor := entityTypeActorMap[entity.EntityType()]

	ucManager := entity.Device().UseCaseManager()
	ucManager.Add(actor, ucEnumType, scenarioSupport)

	return &UseCaseImpl{
		Entity:          entity,
		Actor:           actor,
		name:            model.UseCaseNameType(ucEnumType),
		scenarioSupport: scenarioSupport,
	}
}

func checkArguments(entity spine.EntityImpl, ucEnumType model.UseCaseNameType) {
	actor := entityTypeActorMap[entity.EntityType()]
	if actor == "" {
		panic(fmt.Errorf("cannot derive actor for entity type '%s'", entity.EntityType()))
	}

	if !linq.From(useCaseValidActorsMap[ucEnumType]).Contains(actor) {
		panic(fmt.Errorf("the actor '%s' is not valid for the use case '%s'", actor, ucEnumType))
	}
}

// either returns an existing feature or creates a new one
// for a given entity, featuretype and role
func (u *UseCaseImpl) EntityFeature(entity *spine.EntityLocalImpl, featureType model.FeatureTypeType, role model.RoleType, description string) *spine.FeatureLocalImpl {
	var f *spine.FeatureLocalImpl
	if t := entity.FeatureOfTypeAndRole(featureType, role); t != nil {
		var ok bool
		f, ok = t.(*spine.FeatureLocalImpl)
		if !ok {
			panic("found feature is not of type FeatureLocalImpl")
		}
	} else {
		f = spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, featureType, role)
		f.SetDescriptionString(description)
	}
	return f
}

// internal helper method for getting local and remote feature for a given featureType and a given remoteDevice
func (u *UseCaseImpl) GetLocalClientAndRemoteServerFeatures(featureType model.FeatureTypeType, remoteDevice *spine.DeviceRemoteImpl) (spine.FeatureLocal, *spine.FeatureRemoteImpl, error) {
	featureLocal := u.Entity.Device().FeatureByTypeAndRole(featureType, model.RoleTypeClient)
	featureRemote := remoteDevice.FeatureByTypeAndRole(featureType, model.RoleTypeServer)

	if featureLocal == nil || featureRemote == nil {
		return nil, nil, errors.New("local or remote feature not found")
	}

	return featureLocal, featureRemote, nil
}

func waitForRequest[T any](c chan T, maxDelay time.Duration) *T {
	timeout := time.After(maxDelay)

	select {
	case data := <-c:
		return &data
	case <-timeout:
		return nil
	}
}
