package evsoc

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type EVSOC struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemEVSOCInterface = (*EVSOC)(nil)

// Add support for the EV State of Charge (EVSOC) use case
// as a CEM actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewEVSOC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *EVSOC {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeEV,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeMeasurement},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVStateOfCharge,
		"1.0.0",
		"RC1",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &EVSOC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *EVSOC) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}

func (e *EVSOC) UpdateUseCaseAvailability(available bool) {
	e.LocalEntity.SetUseCaseAvailability(model.UseCaseFilterType{
		Actor:       model.UseCaseActorTypeCEM,
		UseCaseName: e.UseCaseName,
	}, available)
}
