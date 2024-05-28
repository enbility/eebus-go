package vapd

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

type VAPD struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemVAPDInterface = (*VAPD)(nil)

func NewVAPD(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *VAPD {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypePVSystem,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData,
		"1.0.1",
		"RC1",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		eventCB,
		validEntityTypes,
	)

	uc := &VAPD{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *VAPD) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
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
func (e *VAPD) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypePVSystem,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeDeviceConfiguration,
			model.FeatureTypeTypeElectricalConnection,
			model.FeatureTypeTypeMeasurement,
		},
	) {
		return false, nil
	}

	// check for required features
	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	// check if device configuration descriptions contains a required key name
	filter1 := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem),
	}
	if _, err = deviceConfiguration.GetKeyValueDescriptionsForFilter(filter1); err != nil {
		return false, err
	}

	electricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	// check if electrical connection descriptions and parameter descriptions are available name
	if _, err = electricalConnection.GetDescriptionsForFilter(model.ElectricalConnectionDescriptionDataType{}); err != nil {
		return false, err
	}
	if _, err = electricalConnection.GetParameterDescriptionsForFilter(model.ElectricalConnectionParameterDescriptionDataType{}); err != nil {
		return false, err
	}

	// check for required features
	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	// check if measurement descriptions contains a required scope
	filter2 := model.MeasurementDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	if data, err := measurement.GetDescriptionsForFilter(filter2); data == nil || err != nil {
		return false, api.ErrFunctionNotSupported
	}
	filter2.ScopeType = util.Ptr(model.ScopeTypeTypeACYieldTotal)
	if data, err := measurement.GetDescriptionsForFilter(filter2); data == nil || err != nil {
		return false, api.ErrFunctionNotSupported
	}

	return true, nil
}
