package evcem

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type EVCEM struct {
	*usecase.UseCaseBase

	service api.ServiceInterface
}

var _ ucapi.CemEVCEMInterface = (*EVCEM)(nil)

// Add support for the Measurement of Electricity during EV Charging (EVCEM) use case
// as a CEM actor
//
// Parameters:
//   - service: The service implementation
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewEVCEM(
	service api.ServiceInterface,
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *EVCEM {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeEV,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario: model.UseCaseScenarioSupportType(1),
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
		{
			Scenario: model.UseCaseScenarioSupportType(2),
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
		{
			Scenario: model.UseCaseScenarioSupportType(3),
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeElectricalConnection,
				model.FeatureTypeTypeMeasurement,
			},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging,
		"1.0.1",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes)

	uc := &EVCEM{
		UseCaseBase: usecase,
		service:     service,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *EVCEM) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}
