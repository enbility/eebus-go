package vabd

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type VABD struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemVABDInterface = (*VABD)(nil)

// Add support for the Visualization of Aggregated Battery Data (VABD) use case
// as a CEM actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewVABD(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *VABD {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeBatterySystem,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeBatterySystem,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario: model.UseCaseScenarioSupportType(2),
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario: model.UseCaseScenarioSupportType(3),
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(4),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeVisualizationOfAggregatedBatteryData,
		"1.0.1",
		"RC1",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &VABD{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *VABD) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}
