package ohpcf

import (
	"errors"
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type OHPCF struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemOHPCFInterface = (*OHPCF)(nil)

func NewOHPCF(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *OHPCF {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeCompressor,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCompressor,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
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
		false,
	)

	uc := &OHPCF{
		UseCaseBase: usecase,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (o *OHPCF) AddFeatures() error {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeSmartEnergyManagementPs,
	}
	for _, feature := range clientFeatures {
		if f := o.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient); f == nil {
			return errors.New("failed to add feature: " + string(feature))
		}
	}

	return nil
}
