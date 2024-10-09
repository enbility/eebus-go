package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Customer Energy Management
// UseCase: Optimization of Self-Consumption During EV Charging
type CemOSCEVInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the min, max, default limits for each phase of the connected EV
	//
	// parameters:
	//   - entity: the entity of the EV
	CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error)

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
	LoadControlLimits(entity spineapi.EntityRemoteInterface) (limits []LoadLimitsPhase, resultErr error)

	// send new LoadControlLimits to the remote EV
	//
	// parameters:
	//   - entity: the entity of the EV
	//   - limits: a set of limits containing phase specific limit data
	//   - resultCB: callback function for handling the result response
	//
	// recommendations:
	// Sets a recommended charge power in A for each phase. This is mainly
	// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
	// The EV either needs to support the Optimization of Self Consumption usecase or
	// the EVSE needs to be able map the recommendations into oligation limits which then
	// works for all EVs communication either via IEC61851 or ISO15118.
	WriteLoadControlLimits(
		entity spineapi.EntityRemoteInterface,
		limits []LoadLimitsPhase,
		resultCB func(result model.ResultDataType),
	) (*model.MsgCounterType, error)

	// Scenario 2

	// start sending heartbeat from the local CEM entity
	//
	// the heartbeat is started by default when a non 0 timeout is set in the service configuration
	StartHeartbeat()

	// stop sending heartbeat from the local CEM entity
	StopHeartbeat()

	// Scenario 3

	// set the local operating state of the local cem entity
	//
	// parameters:
	//   - failureState: if true, the operating state is set to failure, otherwise to normal
	SetOperatingState(failureState bool) error
}
