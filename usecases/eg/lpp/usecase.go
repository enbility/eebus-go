package lpp

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type LPP struct {
	*usecase.UseCaseBase
}

var _ ucapi.EgLPPInterface = (*LPP)(nil)

// Add support for the Limitation of Power Production (LPC) use case
// as a Energy Guard actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewLPP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *LPP {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeControllableSystem}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCEM,
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeLoadControl},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(2),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(3),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(4),
			Mandatory:      false,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeEnergyGuard,
		model.UseCaseNameTypeLimitationOfPowerProduction,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes)

	uc := &LPP{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *LPP) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceDiagnosis,
		model.FeatureTypeTypeLoadControl,
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeElectricalConnection,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}

	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
}

func (e *LPP) UpdateUseCaseAvailability(available bool) {
	e.LocalEntity.SetUseCaseAvailability(model.UseCaseFilterType{
		Actor:       model.UseCaseActorTypeEnergyGuard,
		UseCaseName: e.UseCaseName,
	}, available)
}
