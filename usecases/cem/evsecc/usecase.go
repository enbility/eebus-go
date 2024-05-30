package evsecc

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type EVSECC struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemEVSECCInterface = (*EVSECC)(nil)

func NewEVSECC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *EVSECC {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEVSE,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		"1.0.1",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2},
		eventCB,
		validEntityTypes)

	uc := &EVSECC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *EVSECC) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceClassification,
		model.FeatureTypeTypeDeviceDiagnosis,
	}

	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EVSECC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{2},
		[]model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
	) {
		// Workaround for the Porsche Mobile Charger Connect that falsely reports
		// the usecase to be on the EV actor
		if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
			model.UseCaseActorTypeEV,
			e.UseCaseName,
			[]model.UseCaseScenarioSupportType{2},
			[]model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
		) {
			return false, nil
		}
	}

	return true, nil
}
