package evsoc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
)

type EVSOC struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemEVSOCInterface = (*EVSOC)(nil)

func NewEVSOC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *EVSOC {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVStateOfCharge,
		"1.0.0",
		"RC1",
		[]model.UseCaseScenarioSupportType{1},
		eventCB,
		validEntityTypes,
	)

	uc := &EVSOC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *EVSOC) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}

func (e *EVSOC) UpdateUseCaseAvailability(available bool) {
	e.LocalEntity.SetUseCaseAvailability(model.UseCaseActorTypeCEM, e.UseCaseName, available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EVSOC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeMeasurement},
	) {
		return false, nil
	}

	// check for required features
	evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil || evMeasurement == nil {
		return false, api.ErrFunctionNotSupported
	}

	// check if measurement description contains an element with scope SOC
	filter := model.MeasurementDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeStateOfCharge),
	}
	if data, err := evMeasurement.GetDescriptionsForFilter(filter); data == nil || err != nil {
		return false, api.ErrNoCompatibleEntity
	}

	return true, nil
}
