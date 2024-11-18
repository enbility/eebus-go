package vapd

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type VAPD struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemVAPDInterface = (*VAPD)(nil)

// Add support for the Visualization of Aggregated Photovoltaic Data (VAPD) use case
// as a CEM actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewVAPD(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *VAPD {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypePVSystem,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypePVSystem,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(3),
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
		model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData,
		"1.0.1",
		"RC1",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &VAPD{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *VAPD) AddFeatures() {
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
