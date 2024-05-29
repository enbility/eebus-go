package lpp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type EgLPP struct {
	*usecase.UseCaseBase
}

var _ ucapi.EgLPPInterface = (*EgLPP)(nil)

func NewEgLPP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *EgLPP {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeEnergyGuard,
		model.UseCaseNameTypeLimitationOfPowerProduction,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4},
		eventCB,
		validEntityTypes)

	uc := &EgLPP{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *EgLPP) AddFeatures() {
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

func (e *EgLPP) UpdateUseCaseAvailability(available bool) {
	e.LocalEntity.SetUseCaseAvailability(model.UseCaseActorTypeEnergyGuard, e.UseCaseName, available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EgLPP) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEnergyGuard,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeDeviceDiagnosis,
			model.FeatureTypeTypeLoadControl,
			model.FeatureTypeTypeDeviceConfiguration,
		},
	) {
		return false, nil
	}

	if _, err := client.NewDeviceDiagnosis(e.LocalEntity, entity); err != nil {
		return false, api.ErrFunctionNotSupported
	}

	if _, err := client.NewLoadControl(e.LocalEntity, entity); err != nil {
		return false, api.ErrFunctionNotSupported
	}

	if _, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err != nil {
		return false, api.ErrFunctionNotSupported
	}

	return true, nil
}
