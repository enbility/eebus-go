package opev

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
func (e *OPEV) CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, nil, nil, api.ErrNoCompatibleEntity
	}

	ec, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		return nil, nil, nil, err
	}

	meas, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return nil, nil, nil, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	measDesc, err := meas.GetDescriptionsForFilter(filter)
	if err != nil {
		return nil, nil, nil, err
	}

	return ec.GetPhaseCurrentLimits(measDesc)
}

// return the current loadcontrol obligation limits
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
func (e *OPEV) LoadControlLimits(entity spineapi.EntityRemoteInterface) (
	limits []ucapi.LoadLimitsPhase, resultErr error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
		ScopeType:     util.Ptr(model.ScopeTypeTypeOverloadProtection),
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
// Sets a maximum A limit for each phase that the EV may not exceed.
// Mainly used for implementing overload protection of the site or limiting the
// maximum charge power of EVs when the EV and EVSE communicate via IEC61851
// and with ISO15118 if the EV does not support the Optimization of Self Consumption
// usecase.
//
// note:
// For obligations to work for optimizing solar excess power, the EV needs to
// have an energy demand. Recommendations work even if the EV does not have an active
// energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
// In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific
// and needs to have specific EVSE support for the specific EV brand.
// In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
func (e *OPEV) WriteLoadControlLimits(
	entity spineapi.EntityRemoteInterface,
	limits []ucapi.LoadLimitsPhase,
	resultCB func(result model.ResultDataType),
) (*model.MsgCounterType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
		Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:     util.Ptr(model.ScopeTypeTypeOverloadProtection),
	}
	return internal.WriteLoadControlPhaseLimits(e.LocalEntity, entity, filter, limits, resultCB)
}

// Scenario 2

// start sending heartbeat from the local CEM entity
//
// the heartbeat is started by default when a non 0 timeout is set in the service configuration
func (e *OPEV) StartHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		_ = hm.StartHeartbeat()
	}
}

// stop sending heartbeat from the local CEM entity
func (e *OPEV) StopHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		hm.StopHeartbeat()
	}
}

// Scenario 3

// set the local operating state of the local cem entity
//
// parameters:
//   - failureState: if true, the operating state is set to failure, otherwise to normal
func (e *OPEV) SetOperatingState(failureState bool) error {
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
