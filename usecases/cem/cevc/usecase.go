package cevc

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

type CEVC struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemCEVCInterface = (*CEVC)(nil)

func NewCEVC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *CEVC {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeCoordinatedEVCharging,
		"1.0.1",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		eventCB,
		validEntityTypes,
	)

	uc := &CEVC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *CEVC) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeTimeSeries,
		model.FeatureTypeTypeIncentiveTable,
		model.FeatureTypeTypeElectricalConnection,
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
func (e *CEVC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{2, 3, 4, 5, 6, 7, 8},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeTimeSeries,
			model.FeatureTypeTypeIncentiveTable,
		},
	) {
		return false, nil
	}

	// check for required features
	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}
	evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	// check if timeseries descriptions contains constraints data
	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
	}
	if _, err = evTimeSeries.GetDescriptionsForFilter(filter); err != nil {
		return false, err
	}

	// check if incentive table descriptions contains data for the required scope
	filter2 := model.TariffDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
	}
	if _, err = evIncentiveTable.GetDescriptionsForFilter(filter2); err != nil {
		return false, err
	}

	return true, nil
}
