package cevc

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type CEVC struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemCEVCInterface = (*CEVC)(nil)

// Add support for the Coordinated EV Charging (CEVC) use case
// as a CEM actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewCEVC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *CEVC {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeEV,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeTimeSeries},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(2),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeTimeSeries},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(3),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeIncentiveTable},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(4),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeTimeSeries},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(6),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(8),
			Mandatory: true,
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeEnergyBroker,
		model.UseCaseNameTypeCoordinatedEVCharging,
		"1.0.1",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &CEVC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *CEVC) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeTimeSeries,
		model.FeatureTypeTypeIncentiveTable,
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
