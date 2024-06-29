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

func (e *CEVC) ChargePlanConstraints(entity spineapi.EntityRemoteInterface) ([]ucapi.DurationSlotValue, error) {
	constraints := []ucapi.DurationSlotValue{}

	if !e.IsCompatibleEntityType(entity) {
		return constraints, api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return constraints, api.ErrDataNotAvailable
	}

	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
	}
	data, err := evTimeSeries.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return constraints, api.ErrDataNotAvailable
	}

	// we need at least a time series slot
	if data[0].TimeSeriesSlot == nil {
		return constraints, api.ErrDataNotAvailable
	}

	// get the values for all slots
	for _, slot := range data[0].TimeSeriesSlot {
		newSlot := ucapi.DurationSlotValue{}

		if slot.Duration != nil {
			if duration, err := slot.Duration.GetTimeDuration(); err == nil {
				newSlot.Duration = duration
			}
		} else if slot.TimePeriod != nil {
			var slotStart, slotEnd time.Time
			if slot.TimePeriod.StartTime != nil {
				if time, err := slot.TimePeriod.StartTime.GetTime(); err == nil {
					slotStart = time
				}
			}
			if slot.TimePeriod.EndTime != nil {
				if time, err := slot.TimePeriod.EndTime.GetTime(); err == nil {
					slotEnd = time
				}
			}
			newSlot.Duration = slotEnd.Sub(slotStart)
		}

		if slot.MaxValue != nil {
			newSlot.Value = slot.MaxValue.GetValue()
		}

		constraints = append(constraints, newSlot)
	}

	return constraints, nil
}

func (e *CEVC) ChargePlan(entity spineapi.EntityRemoteInterface) (ucapi.ChargePlan, error) {
	plan := ucapi.ChargePlan{}

	if !e.IsCompatibleEntityType(entity) {
		return plan, api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return plan, api.ErrDataNotAvailable
	}

	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypePlan),
	}
	data, err := evTimeSeries.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return plan, api.ErrDataNotAvailable
	}

	// we need at least a time series slot
	if data[0].TimeSeriesSlot == nil {
		return plan, api.ErrDataNotAvailable
	}

	startAvailable := false
	// check the start time relative to now of the plan, default is now
	currentStart := time.Now()
	currentEnd := currentStart
	if data[0].TimePeriod != nil && data[0].TimePeriod.StartTime != nil {
		if start, err := data[0].TimePeriod.StartTime.GetTimeDuration(); err == nil {
			currentStart = currentStart.Add(start)
			startAvailable = true
		}
	}

	// get the values for all slots
	for index, slot := range data[0].TimeSeriesSlot {
		newSlot := ucapi.ChargePlanSlotValue{}

		slotStartDefined := false
		if index == 0 && startAvailable && (slot.TimePeriod == nil || slot.TimePeriod.StartTime == nil) {
			newSlot.Start = currentStart
			slotStartDefined = true
		}
		if slot.TimePeriod != nil && slot.TimePeriod.StartTime != nil {
			if time, err := slot.TimePeriod.StartTime.GetTime(); err == nil {
				newSlot.Start = time
				slotStartDefined = true
			}
		}
		if !slotStartDefined {
			newSlot.Start = currentEnd
		}

		if slot.Duration != nil {
			if duration, err := slot.Duration.GetTimeDuration(); err == nil {
				newSlot.End = newSlot.Start.Add(duration)
				currentEnd = newSlot.End
			}
		} else if slot.TimePeriod != nil && slot.TimePeriod.EndTime != nil {
			if time, err := slot.TimePeriod.StartTime.GetTime(); err == nil {
				newSlot.End = time
				currentEnd = newSlot.End
			}
		}

		if slot.Value != nil {
			newSlot.Value = slot.Value.GetValue()
		}
		if slot.MinValue != nil {
			newSlot.MinValue = slot.MinValue.GetValue()
		}
		if slot.MaxValue != nil {
			newSlot.MaxValue = slot.MaxValue.GetValue()
		}

		plan.Slots = append(plan.Slots, newSlot)
	}

	return plan, nil
}
