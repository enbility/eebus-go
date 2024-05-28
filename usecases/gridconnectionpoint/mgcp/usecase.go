package mgcp

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

type MGCP struct {
	*usecase.UseCaseBase
}

var _ ucapi.GcpMGCPInterface = (*MGCP)(nil)

func NewMGCP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MGCP {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCEM,
		model.EntityTypeTypeGridConnectionPointOfPremises,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfGridConnectionPoint,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7},
		eventCB,
		validEntityTypes)

	uc := &MGCP{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MGCP) AddFeatures() {
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
func (e *MGCP) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeGridConnectionPoint,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{2, 3, 4},
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
		ScopeType: util.Ptr(model.ScopeTypeTypeACPower),
	}
	data1, err1 := measurement.GetDescriptionsForFilter(filter)
	filter.ScopeType = util.Ptr(model.ScopeTypeTypeGridFeedIn)
	data2, err2 := measurement.GetDescriptionsForFilter(filter)
	filter.ScopeType = util.Ptr(model.ScopeTypeTypeGridConsumption)
	data3, err3 := measurement.GetDescriptionsForFilter(filter)
	if err1 != nil || err2 != nil || err3 != nil ||
		data1 == nil || data2 == nil || data3 == nil {
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

	return true, nil
}
