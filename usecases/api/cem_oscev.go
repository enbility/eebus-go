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
	//
	// recommendations:
	// Sets a recommended charge power in A for each phase. This is mainly
	// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
	// The EV either needs to support the Optimization of Self Consumption usecase or
	// the EVSE needs to be able map the recommendations into oligation limits which then
	// works for all EVs communication either via IEC61851 or ISO15118.
	WriteLoadControlLimits(entity spineapi.EntityRemoteInterface, limits []LoadLimitsPhase) (*model.MsgCounterType, error)

	// Scenario 2

	// this is automatically covered by the SPINE implementation

	// Scenario 3

	// this is covered by the central CEM interface implementation
	// use that one to set the CEM's operation state which will inform all remote devices
}
