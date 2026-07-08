package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"time"
)

type CemOHPCFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// Get all relevant info for an available power consumption
	//
	// return OptionalPowerConsumptionInfo struct with all available info if an optional power consumption has been announced
	OptionalPowerConsumption(entity spineapi.EntityRemoteInterface) (*OptionalPowerConsumptionInfo, error)

	// The availability of an optional consumption of power [OHPCF-011/1].
	//
	// return true if the optional consumption of power is possible
	OptionalPowerConsumptionAvailable(entity spineapi.EntityRemoteInterface) (bool, error)

	// The requested power estimate value [OHPCF-011/2/1].
	//
	// return the requested power estimate value
	RequestedPowerEstimate(entity spineapi.EntityRemoteInterface) (float64, error)

	// The maximal value for the requested power estimate [OHPCF-011/2/1].
	//
	// return the maximal value for the requested power estimate
	RequestedPowerMax(entity spineapi.EntityRemoteInterface) (float64, error)

	// Indication whether the consumption may be stopped by the CEM [OHPCF-011/5].
	//
	// return true if the consumption may be stopped
	ConsumptionIsStoppable(entity spineapi.EntityRemoteInterface) (bool, error)

	// Indication whether the consumption may be paused and resumed by the CEM [OHPCF-011/6]
	//
	// return true if the consumption may be paused
	ConsumptionIsPausable(entity spineapi.EntityRemoteInterface) (bool, error)

	// The start time of the process [OHPCF-012/1].
	//
	// return the start time of the process
	PowerConsumptionProcessStartTime(entity spineapi.EntityRemoteInterface) (time.Time, error)

	// The current state of this power consumption process [OHPCF-012/2].
	//
	// return the current state of this power consumption process
	PowerConsumptionProcessState(entity spineapi.EntityRemoteInterface) (CompressorPowerConsumptionStateType, error)

	// The minimal time a consumption process must last [OHPCF-008].
	//
	// return the minimal time a consumption process must last
	PowerConsumptionMinimalRunDuration(entity spineapi.EntityRemoteInterface) (time.Duration, error)

	// The minimal time a pause of a consumption process must last [OHPCF-009].
	//
	// return the minimal time a pause of a consumption process must last
	PowerConsumptionMinimalPauseDuration(entity spineapi.EntityRemoteInterface) (time.Duration, error)

	// Scenario 2

	// Schedule an optional power consumption process [OHPCF-004].
	//
	// note:
	// Rescheduling an already scheduled power consumption process is possible as long as the
	// scheduled process has not startet yet.
	//
	// parameters:
	//   - startIn: Delay from now until the power consumption starts (0 = start immediately)
	SchedulePowerConsumptionProcess(entity spineapi.EntityRemoteInterface, startIn time.Duration, resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType)) (*model.MsgCounterType, error)

	// stop (abort) the process [OHPCF-022/1].
	AbortPowerConsumptionProcess(entity spineapi.EntityRemoteInterface, resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType)) (*model.MsgCounterType, error)

	// pause the process [OHPCF-022/2].
	PausePowerConsumptionProcess(entity spineapi.EntityRemoteInterface, resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType)) (*model.MsgCounterType, error)

	// resume the process [OHPCF-022/3].
	ResumePowerConsumptionProcess(entity spineapi.EntityRemoteInterface, resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType)) (*model.MsgCounterType, error)
}
