package ohpcf

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// Optimization of Heat Pump Compressor Function
type OHPCF struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemOHPCFInterface = (*OHPCF)(nil)

func NewOHPCF(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *OHPCF {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeCEM,
		model.UseCaseActorTypeCompressor,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCEM,
		model.EntityTypeTypeCompressor,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:       model.UseCaseScenarioSupportType(1),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeSmartEnergyManagementPs},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(2),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeSmartEnergyManagementPs},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionByHeatPumpCompressorFlexibility,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &OHPCF{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *OHPCF) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeSmartEnergyManagementPs,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}
