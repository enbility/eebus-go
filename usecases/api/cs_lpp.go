package api

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Controllable System
// UseCase: Limitation of Power Production
type CsLPPInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current loadcontrol limit data
	//
	// return values:
	//   - limit: load limit data
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	ProductionLimit() (LoadLimit, error)

	// set the current loadcontrol limit data
	SetProductionLimit(limit LoadLimit) (resultErr error)

	// return the currently pending incoming consumption write limits
	PendingProductionLimits() map[model.MsgCounterType]LoadLimit

	// accept or deny an incoming consumption write limit
	//
	// parameters:
	//  - msg: the incoming write message
	//  - approve: if the write limit for msg should be approved or not
	//  - reason: the reason why the approval is denied, otherwise an empty string
	ApproveOrDenyProductionLimit(msgCounter model.MsgCounterType, approve bool, reason string)

	// Scenario 2

	// return Failsafe limit for the produced active (real) power of the
	// Controllable System. This limit becomes activated in "init" state or "failsafe state".
	//
	// return values:
	//   - value: the power limit in W
	//   - changeable: boolean if the client service can change the limit
	FailsafeProductionActivePowerLimit() (value float64, isChangeable bool, resultErr error)

	// set Failsafe limit for the produced active (real) power of the
	// Controllable System. This limit becomes activated in "init" state or "failsafe state".
	//
	// parameters:
	//   - value: the power limit in W
	//   - changeable: boolean if the client service can change the limit
	SetFailsafeProductionActivePowerLimit(value float64, changeable bool) (resultErr error)

	// return minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state"
	//
	// return values:
	//   - value: the power limit in W
	//   - changeable: boolean if the client service can change the limit
	FailsafeDurationMinimum() (duration time.Duration, isChangeable bool, resultErr error)

	// set minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state"
	//
	// parameters:
	//   - duration: has to be >= 2h and <= 24h
	//   - changeable: boolean if the client service can change this value
	SetFailsafeDurationMinimum(duration time.Duration, changeable bool) (resultErr error)

	// Scenario 3

	// this is automatically covered by the SPINE implementation
	//
	// returns true, if the last heartbeat is within 2 minutes, otherwise false
	IsHeartbeatWithinDuration() bool

	// Scenario 4

	// return nominal maximum active (real) power the Controllable System is
	// allowed to produce due to the customer's contract.
	ContractualProductionNominalMax() (float64, error)

	// set nominal maximum active (real) power the Controllable System is
	// allowed to produce due to the customer's contract.
	//
	// parameters:
	//   - value: contractual nominal max power production in W
	SetContractualProductionNominalMax(value float64) (resultErr error)
}
