package oscev

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type OSCEV struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemOSCEVInterface = (*OSCEV)(nil)

// Add support for the Optimization of Self-Consumption during EV Charging (OSCEV) use case
// as a CEM actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewOSCEV(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *OSCEV {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeEV,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeLoadControl,
				model.FeatureTypeTypeElectricalConnection,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(3),
			Mandatory: true,
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
		"1.0.1",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &OSCEV{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *OSCEV) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeLoadControl,
		model.FeatureTypeTypeElectricalConnection,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}

	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
}
