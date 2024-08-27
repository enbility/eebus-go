package cevc

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// returns the current charging strategy
func (e *CEVC) ChargeStrategy(entity spineapi.EntityRemoteInterface) ucapi.EVChargeStrategyType {
	if !e.IsCompatibleEntityType(entity) {
		return ucapi.EVChargeStrategyTypeUnknown
	}

	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return ucapi.EVChargeStrategyTypeUnknown
	}

	// only the time series data for singledemand is relevant for detecting the charging strategy
	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
	}
	data, err := evTimeSeries.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return ucapi.EVChargeStrategyTypeUnknown
	}

	// without time series slots, there is no known strategy
	if len(data[0].TimeSeriesSlot) == 0 {
		return ucapi.EVChargeStrategyTypeUnknown
	}

	// get the value for the first slot
	firstSlot := data[0].TimeSeriesSlot[0]

	switch {
	case firstSlot.Duration == nil:
		// if value is > 0 and duration does not exist, the EV is direct charging
		if firstSlot.Value != nil && firstSlot.Value.GetValue() > 0 {
			return ucapi.EVChargeStrategyTypeDirectCharging
		}

		// maxValue will show the maximum amount the battery could take
		return ucapi.EVChargeStrategyTypeNoDemand

	case firstSlot.Duration != nil:
		if _, err := firstSlot.Duration.GetTimeDuration(); err != nil {
			// we got an invalid duration
			return ucapi.EVChargeStrategyTypeUnknown
		}

		if firstSlot.MinValue != nil && firstSlot.MinValue.GetValue() > 0 {
			return ucapi.EVChargeStrategyTypeMinSoC
		}

		if firstSlot.Value != nil {
			if firstSlot.Value.GetValue() > 0 {
				// there is demand and a duration
				return ucapi.EVChargeStrategyTypeTimedCharging
			}

			return ucapi.EVChargeStrategyTypeNoDemand
		}
	}

	return ucapi.EVChargeStrategyTypeUnknown
}

// returns the current energy demand in Wh and the duration
func (e *CEVC) EnergyDemand(entity spineapi.EntityRemoteInterface) (ucapi.Demand, error) {
	demand := ucapi.Demand{}

	if !e.IsCompatibleEntityType(entity) {
		return demand, api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return demand, api.ErrDataNotAvailable
	}

	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
	}
	data, err := evTimeSeries.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return demand, api.ErrDataNotAvailable
	}

	// we need at least a time series slot
	if data[0].TimeSeriesSlot == nil {
		return demand, api.ErrDataNotAvailable
	}

	// get the value for the first slot, ignore all others, which
	// in the tests so far always have min/max/value 0
	firstSlot := data[0].TimeSeriesSlot[0]
	if firstSlot.MinValue != nil {
		demand.MinDemand = firstSlot.MinValue.GetValue()
	}
	if firstSlot.Value != nil {
		demand.OptDemand = firstSlot.Value.GetValue()
	}
	if firstSlot.MaxValue != nil {
		demand.MaxDemand = firstSlot.MaxValue.GetValue()
	}
	if firstSlot.Duration != nil {
		if tempDuration, err := firstSlot.Duration.GetTimeDuration(); err == nil {
			demand.DurationUntilEnd = tempDuration.Seconds()
		}
	}

	// start time has to be defined either in TimePeriod or the first slot
	relStartTime := time.Duration(0)

	startTimeSet := false
	if data[0].TimePeriod != nil && data[0].TimePeriod.StartTime != nil {
		if temp, err := data[0].TimePeriod.StartTime.GetTimeDuration(); err == nil {
			relStartTime = temp
			startTimeSet = true
		}
	}

	if !startTimeSet {
		if firstSlot.TimePeriod != nil && firstSlot.TimePeriod.StartTime != nil {
			if temp, err := firstSlot.TimePeriod.StartTime.GetTimeDuration(); err == nil {
				relStartTime = temp
			}
		}
	}

	demand.DurationUntilStart = relStartTime.Seconds()

	return demand, nil
}
