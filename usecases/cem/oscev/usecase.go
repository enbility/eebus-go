package oscev

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
)

type OSCEV struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemOSCEVInterface = (*OSCEV)(nil)

func NewOSCEV(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *OSCEV {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCompressor,
		model.EntityTypeTypeElectricalImmersionHeater,
		model.EntityTypeTypeEV,
		model.EntityTypeTypeHeatPumpAppliance,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
		"1.0.1",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		eventCB,
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

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *OSCEV) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		[]model.FeatureTypeType{model.FeatureTypeTypeLoadControl},
	) {
		return false, nil
	}

	// check for required features
	evLoadControl, err := client.NewLoadControl(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	// check if loadcontrol limit descriptions contains a recommendation category
	filter := model.LoadControlLimitDescriptionDataType{
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeRecommendation),
	}
	if data, err := evLoadControl.GetLimitDescriptionsForFilter(filter); err != nil || len(data) == 0 {
		return false, api.ErrFunctionNotSupported
	}

	return true, nil
}
