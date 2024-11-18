package lpc

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type LPC struct {
	*usecase.UseCaseBase
}

var _ ucapi.EgLPCInterface = (*LPC)(nil)

// Add support for the Limitation of Power Consumption (LPC) use case
// as a Energy Guard actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewLPC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *LPC {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeControllableSystem}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCEM,
		model.EntityTypeTypeCompressor,
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeHeatPumpAppliance,
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
		model.UseCaseNameTypeLimitationOfPowerConsumption,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &LPC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *LPC) AddFeatures() {
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
