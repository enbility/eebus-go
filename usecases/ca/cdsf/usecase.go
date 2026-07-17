package cdsf

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type CDSF struct {
	*usecase.UseCaseBase
}

var _ ucapi.CaCDSFInterface = (*CDSF)(nil)

// Add support for the Configuration of DHW System Function (CDSF)
// use case as a Configuration Appliance actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewCDSF(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *CDSF {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeDHWCircuit}
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeDHWCircuit}
	// the DHW circuit has to support at least one of the scenarios 1 and 2
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeHvac},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(2),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeHvac},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(3),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeHvac},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeConfigurationAppliance,
		model.UseCaseNameTypeConfigurationOfDhwSystemFunction,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
		false)

	uc := &CDSF{
		UseCaseBase: usecase,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (e *CDSF) AddFeatures() error {
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
