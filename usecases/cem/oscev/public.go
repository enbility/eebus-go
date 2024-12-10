package oscev

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the min, max, default limits for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *OSCEV) CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, nil, nil, api.ErrNoCompatibleEntity
	}

	ec, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		return nil, nil, nil, err
	}

	lc, err := client.NewLoadControl(e.LocalEntity, entity)
	if err != nil {
		return nil, nil, nil, err
	}

	filter := model.LoadControlLimitDescriptionDataType {
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeRecommendation),
		Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:     util.Ptr(model.ScopeTypeTypeSelfConsumption),
	}
	
	limitDescs, err := lc.GetLimitDescriptionsForFilter(filter)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(limitDescs) == 0 {
		return nil, nil, nil, api.ErrDataNotAvailable
	}

	measDescs := make([]model.MeasurementDescriptionDataType, 0)
	for _, ld := range limitDescs {
		measId := ld.MeasurementId
		if measId == nil {
			return nil, nil, nil, api.ErrDataNotAvailable
		}
		md := model.MeasurementDescriptionDataType{
			MeasurementId: measId,
		}
		measDescs = append(measDescs, md)
	}
	return ec.GetPhaseCurrentLimits(measDescs)
}

// return the current loadcontrol recommendation limits
//
// parameters:
//   - entity: the entity of the EV
//
// return values:
//   - limits: per phase data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *OSCEV) LoadControlLimits(entity spineapi.EntityRemoteInterface) (
	limits []ucapi.LoadLimitsPhase, resultErr error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeRecommendation),
		Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:     util.Ptr(model.ScopeTypeTypeSelfConsumption),
	}
	return internal.LoadControlLimits(e.LocalEntity, entity, filter)
}

// send new LoadControlLimits to the remote EV
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - limits: a set of limits containing phase specific limit data
//   - resultCB: callback function for handling the result response
//
// recommendations:
// Sets a recommended charge power in A for each phase. This is mainly
// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
// The EV either needs to support the Optimization of Self Consumption usecase or
// the EVSE needs to be able map the recommendations into oligation limits which then
// works for all EVs communication either via IEC61851 or ISO15118.
func (e *OSCEV) WriteLoadControlLimits(
	entity spineapi.EntityRemoteInterface,
	limits []ucapi.LoadLimitsPhase,
	resultCB func(result model.ResultDataType),
) (*model.MsgCounterType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeRecommendation),
		Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:     util.Ptr(model.ScopeTypeTypeSelfConsumption),
	}
	return internal.WriteLoadControlPhaseLimits(e.LocalEntity, entity, filter, limits, resultCB)
}

// Scenario 2

// start sending heartbeat from the local CEM entity
//
// the heartbeat is started by default when a non 0 timeout is set in the service configuration
func (e *OSCEV) StartHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		_ = hm.StartHeartbeat()
	}
}

// stop sending heartbeat from the local CEM entity
func (e *OSCEV) StopHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		hm.StopHeartbeat()
	}
}

// Scenario 3

// set the local operating state of the local cem entity
//
// parameters:
//   - failureState: if true, the operating state is set to failure, otherwise to normal
func (e *OSCEV) SetOperatingState(failureState bool) error {
	lf, err := server.NewDeviceDiagnosis(e.LocalEntity)
	if err != nil {
		return err
	}

	state := model.DeviceDiagnosisOperatingStateTypeNormalOperation
	if failureState {
		state = model.DeviceDiagnosisOperatingStateTypeFailure
	}
	lf.SetLocalOperatingState(state)

	return nil
}
