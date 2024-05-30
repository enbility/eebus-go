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

func NewEVCC(
	service api.ServiceInterface,
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *EVCC {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		"1.0.1",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
		eventCB,
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

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EVCC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 8},
		[]model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
	) {
		return false, nil
	}

	return true, nil
}
