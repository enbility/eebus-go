package crcsf

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type CRCSF struct {
	*usecase.UseCaseBase
}

var _ ucapi.CaCRCSFInterface = (*CRCSF)(nil)

// Add support for the Configuration of Room Cooling System Function (CRCSF)
// use case as a Configuration Appliance actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewCRCSF(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *CRCSF {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeHVACRoom}
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeHvacRoom}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeHvac},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeConfigurationAppliance,
		model.UseCaseNameTypeConfigurationOfRoomCoolingSystemFunction,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
		false)

	uc := &CRCSF{
		UseCaseBase: usecase,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (e *CRCSF) AddFeatures() error {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeHvac,
	}
	for _, feature := range clientFeatures {
		if f := e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient); f == nil {
			return errors.New("could not add feature: " + string(feature))
		}
	}

	return nil
}
