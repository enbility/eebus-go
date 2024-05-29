package evcem

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type CemEVCEM struct {
	*usecase.UseCaseBase

	service api.ServiceInterface
}

var _ ucapi.CemEVCEMInterface = (*CemEVCEM)(nil)

func NewCemEVCEM(service api.ServiceInterface, localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *CemEVCEM {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging,
		"1.0.1",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		eventCB,
		validEntityTypes)

	uc := &CemEVCEM{
		UseCaseBase: usecase,
		service:     service,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *CemEVCEM) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
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
func (e *CemEVCEM) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName,
		nil,
		nil,
	) {
		return false, nil
	}

	return true, nil
}
