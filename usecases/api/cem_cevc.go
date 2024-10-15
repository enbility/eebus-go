package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Customer Energy Management
// UseCase: Coordinated EV Charging
type CemCEVCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// returns the current charging stratey
	//
	// parameters:
	//   - entity: the entity of the EV
	//
	// returns EVChargeStrategyTypeUnknown if it could not be determined, e.g.
	// if the vehicle communication is via IEC61851 or the EV doesn't provide
	// any information about its charging mode or plan
	ChargeStrategy(remoteEntity spineapi.EntityRemoteInterface) EVChargeStrategyType

	// returns the current energy demand
	//
	// parameters:
	//   - entity: the entity of the EV
	//
	// return values:
	//   - EVDemand: details about the actual demands from the EV
	//   - error: if no data is available
	//
	// if duration is 0, direct charging is active, otherwise timed charging is active
	EnergyDemand(remoteEntity spineapi.EntityRemoteInterface) (Demand, error)

	// Scenario 2

	TimeSlotConstraints(entity spineapi.EntityRemoteInterface) (TimeSlotConstraints, error)

	// send power limits to the EV
	//
	// parameters:
	//   - entity: the entity of the EV
	//   - data: the power limits
	//
	// if no data is provided, default power limits with the max possible value for 7 days will be sent
	WritePowerLimits(entity spineapi.EntityRemoteInterface, data []DurationSlotValue) error

	// Scenario 3

	// return the current incentive constraints
	//
	// parameters:
	//   - entity: the entity of the EV
	IncentiveConstraints(entity spineapi.EntityRemoteInterface) (IncentiveSlotConstraints, error)

	// send new incentives to the EV
	//
	// parameters:
	//   - entity: the entity of the EV
	//   - data: the incentive descriptions
	WriteIncentiveTableDescriptions(entity spineapi.EntityRemoteInterface, data []IncentiveTariffDescription) error

	// send incentives to the EV
	//
	// parameters:
	//   - entity: the entity of the EV
	//   - data: the incentives
	//
	// if no data is provided, default incentives with the same price for 7 days will be sent
	WriteIncentives(entity spineapi.EntityRemoteInterface, data []DurationSlotValue) error

	// Scenario 4

	// return the current charge plan constraints
	//
	// parameters:
	//   - entity: the entity of the EV
	ChargePlanConstraints(entity spineapi.EntityRemoteInterface) ([]DurationSlotValue, error)

	// return the current charge plan of the EV
	//
	// parameters:
	//   - entity: the entity of the EV
	ChargePlan(entity spineapi.EntityRemoteInterface) (ChargePlan, error)

	// Scenario 5 & 6

	// start sending heartbeat from the local CEM entity
	//
	// the heartbeat is started by default when a non 0 timeout is set in the service configuration
	StartHeartbeat()

	// stop sending heartbeat from the local CEM entity
	StopHeartbeat()

	// Scenario 7 & 8

	// set the local operating state of the local cem entity
	//
	// parameters:
	//   - failureState: if true, the operating state is set to failure, otherwise to normal
	SetOperatingState(failureState bool) error
}
