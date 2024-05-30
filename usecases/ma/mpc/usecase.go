package mpc

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

type MPC struct {
	*usecase.UseCaseBase
}

var _ ucapi.MaMPCInterface = (*MPC)(nil)

func NewMPC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MPC {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCompressor,
		model.EntityTypeTypeElectricalImmersionHeater,
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeHeatPumpAppliance,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfPowerConsumption,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5},
		eventCB,
		validEntityTypes)

	uc := &MPC{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MPC) AddFeatures() {
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
func (e *MPC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeMonitoredUnit,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeElectricalConnection,
			model.FeatureTypeTypeMeasurement,
		},
	) {
		return false, nil
	}

	// check if measurement description contain data for the required scope
	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	filter := model.MeasurementDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	if data, err := measurement.GetDescriptionsForFilter(filter); data == nil || err != nil {
		return false, api.ErrDataNotAvailable
	}

	// check if electrical connection descriptions is provided
	electricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		return false, api.ErrFunctionNotSupported
	}

	if _, err = electricalConnection.GetDescriptionsForFilter(model.ElectricalConnectionDescriptionDataType{}); err != nil {
		return false, err
	}

	if _, err = electricalConnection.GetParameterDescriptionsForFilter(model.ElectricalConnectionParameterDescriptionDataType{}); err != nil {
		return false, err
	}

	return true, nil
}
