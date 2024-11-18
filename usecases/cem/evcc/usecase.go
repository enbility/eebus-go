package evcc

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type EVCC struct {
	*usecase.UseCaseBase

	service api.ServiceInterface
}

var _ ucapi.CemEVCCInterface = (*EVCC)(nil)

// Add support for the EV Commissioning and Configuration (EVCEM) use case
// as a CEM actor
//
// Parameters:
//   - service: The service implementation
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewEVCC(
	service api.ServiceInterface,
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *EVCC {
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
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(2),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(3),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(4),
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeIdentification},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(5),
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceClassification},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(6),
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(7),
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(8),
			Mandatory: true,
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		"1.0.1",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &EVCC{
		UseCaseBase: usecase,
		service:     service,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *EVCC) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeIdentification,
		model.FeatureTypeTypeDeviceClassification,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeDeviceDiagnosis,
	}
	for _, feature := range clientFeatures {
		f := e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultCallback(e.HandleResponse)
	}
}
