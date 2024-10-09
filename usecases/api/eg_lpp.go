package api

import (
	"time"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Energy Guard
// UseCase: Limitation of Power Production
type EgLPPInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current production limit data
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//
	// return values:
	//   - limit: load limit data
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	ProductionLimit(entity spineapi.EntityRemoteInterface) (limit LoadLimit, resultErr error)

	// send new LoadControlLimits
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//   - limit: load limit data
	//   - resultCB: callback function for handling the result response
	WriteProductionLimit(
		entity spineapi.EntityRemoteInterface,
		limit LoadLimit,
		resultCB func(result model.ResultDataType),
	) (*model.MsgCounterType, error)

	// Scenario 2

	// return Failsafe limit for the produced active (real) power of the
	// Controllable System. This limit becomes activated in "init" state or "failsafe state".
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//
	// return values:
	//   - positive values are used for production
	FailsafeProductionActivePowerLimit(entity spineapi.EntityRemoteInterface) (float64, error)

	// send new Failsafe Production Active Power Limit
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//   - value: the new limit in W
	WriteFailsafeProductionActivePowerLimit(entity spineapi.EntityRemoteInterface, value float64) (*model.MsgCounterType, error)

	// return minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state"
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//
	// return values:
	//   - negative values are used for production
	FailsafeDurationMinimum(entity spineapi.EntityRemoteInterface) (time.Duration, error)

	// send new Failsafe Duration Minimum
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//   - duration: the duration, between 2h and 24h
	WriteFailsafeDurationMinimum(entity spineapi.EntityRemoteInterface, duration time.Duration) (*model.MsgCounterType, error)

	// Scenario 3

	// start sending heartbeat from the local entity supporting this usecase
	//
	// the heartbeat is started by default when a non 0 timeout is set in the service configuration
	StartHeartbeat()

	// stop sending heartbeat from the local entity supporting this usecase
	StopHeartbeat()

	// check wether there was a heartbeat received within the last 2 minutes
	//
	// returns true, if the last heartbeat is within 2 minutes, otherwise false
	IsHeartbeatWithinDuration(entity spineapi.EntityRemoteInterface) bool

	// Scenario 4

	// return nominal maximum active (real) power the Controllable System is
	// able to produce according to the contract (EMS), device label or data sheet.
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	ProductionNominalMax(entity spineapi.EntityRemoteInterface) (float64, error)
}
