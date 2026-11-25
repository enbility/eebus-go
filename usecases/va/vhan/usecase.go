package vhan

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type VHAN struct {
	*usecase.UseCaseBase
}

var _ ucapi.VaVHANInterface = (*VHAN)(nil)

// Add support for the Visualization of Heating Area Name (VHAN)
// use case as a Visualization Appliance actor
func NewVHAN(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *VHAN {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeHeatingCircuit,
		model.UseCaseActorTypeHeatingZone,
		model.UseCaseActorTypeHVACRoom,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeHeatingCircuit,
		model.EntityTypeTypeHeatingZone,
		model.EntityTypeTypeHvacRoom,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceClassification},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(2),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceClassification},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(3),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceClassification},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeVisualizationAppliance,
		model.UseCaseNameTypeVisualizationOfHeatingAreaName,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
		false)

	uc := &VHAN{
		UseCaseBase: usecase,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (e *VHAN) AddFeatures() error {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceClassification,
	}
	for _, feature := range clientFeatures {
		if f := e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient); f == nil {
			return errors.New("could not add feature: " + string(feature))
		}
	}

	return nil
}
