package mdt

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type MDT struct {
	*usecase.UseCaseBase
}

var _ ucapi.MaMDTInterface = (*MDT)(nil)

// Add support for the Monitoring of Domestic Hot Water (DHW) Temperature (MDT)
// use case as a Monitoring Appliance actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewMDT(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MDT {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeDHWCircuit}
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeDHWCircuit}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeMeasurement},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfDhwTemperature,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
		false)

	uc := &MDT{
		UseCaseBase: usecase,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (e *MDT) AddFeatures() error {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		if f := e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient); f == nil {
			return errors.New("could not add feature: " + string(feature))
		}
	}

	return nil
}
