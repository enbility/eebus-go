package opev

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

type OPEV struct {
	*usecase.UseCaseBase
}

var _ ucapi.CemOPEVInterface = (*OPEV)(nil)

func NewOPEV(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *OPEV {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
		"1.0.1",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		eventCB,
		validEntityTypes,
	)

	uc := &OPEV{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *OPEV) AddFeatures() {
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
func (e *OPEV) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
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

	// check if loadcontrol limit descriptions contains a obligation category
	filter := model.LoadControlLimitDescriptionDataType{
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
	}
	if data, err := evLoadControl.GetLimitDescriptionsForFilter(filter); err != nil || len(data) == 0 {
		return false, api.ErrFunctionNotSupported
	}

	return true, nil
}
