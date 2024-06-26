package oscev

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
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
