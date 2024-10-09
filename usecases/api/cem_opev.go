package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Customer Energy Management
// UseCase: Overload Protection by EV Charging Current Curtailment
type CemOPEVInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the min, max, default limits for each phase of the connected EV
	//
	// parameters:
	//   - entity: the entity of the EV
	CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error)

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
	LoadControlLimits(entity spineapi.EntityRemoteInterface) (limits []LoadLimitsPhase, resultErr error)

	// send new LoadControlLimits to the remote EV
	//
	// parameters:
	//   - entity: the entity of the EV
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
