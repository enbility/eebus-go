package cdt

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type CDT struct {
	*usecase.UseCaseBase

	modes map[model.HvacOperationModeTypeType]model.SetpointIdType
}

var _ ucapi.CaCDTInterface = (*CDT)(nil)

// NewCDT creates a new CDT use case
func NewCDT(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *CDT {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeDHWCircuit,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeDHWCircuit,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeSetpoint,
				model.FeatureTypeTypeHvac,
			},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeConfigurationAppliance,
		model.UseCaseNameTypeConfigurationOfDhwTemperature,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &CDT{
		UseCaseBase: usecase,
		modes:       make(map[model.HvacOperationModeTypeType]model.SetpointIdType),
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

// AddFeatures adds the features required for the CDT use case
func (e *CDT) AddFeatures() {
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeSetpoint,
		model.FeatureTypeTypeHvac,
	}

	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}
