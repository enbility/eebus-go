package cevc

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// returns the constraints for the time slots
func (e *CEVC) TimeSlotConstraints(entity spineapi.EntityRemoteInterface) (ucapi.TimeSlotConstraints, error) {
	result := ucapi.TimeSlotConstraints{}

	if !e.IsCompatibleEntityType(entity) {
		return result, api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return result, api.ErrDataNotAvailable
	}

	constraints, err := evTimeSeries.GetConstraints()
	if err != nil {
		return result, err
	}

	// only use the first constraint
	constraint := constraints[0]

	if constraint.SlotCountMin != nil {
		result.MinSlots = uint(*constraint.SlotCountMin)
	}
	if constraint.SlotCountMax != nil {
		result.MaxSlots = uint(*constraint.SlotCountMax)
	}
	if constraint.SlotDurationMin != nil {
		if duration, err := constraint.SlotDurationMin.GetTimeDuration(); err == nil {
			result.MinSlotDuration = duration
		}
	}
	if constraint.SlotDurationMax != nil {
		if duration, err := constraint.SlotDurationMax.GetTimeDuration(); err == nil {
			result.MaxSlotDuration = duration
		}
	}
	if constraint.SlotDurationStepSize != nil {
		if duration, err := constraint.SlotDurationStepSize.GetTimeDuration(); err == nil {
			result.SlotDurationStepSize = duration
		}
	}

	return result, nil
}

// send power limits to the EV
// if no data is provided, default power limits with the max possible value for 7 days will be sent
func (e *CEVC) WritePowerLimits(entity spineapi.EntityRemoteInterface, data []ucapi.DurationSlotValue) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		return api.ErrDataNotAvailable
	}

	if len(data) == 0 {
		data, err = e.defaultPowerLimits(entity)
		if err != nil {
			return err
		}
	}

	constraints, err := e.TimeSlotConstraints(entity)
	if err != nil {
		return err
	}

	if constraints.MinSlots != 0 && constraints.MinSlots > uint(len(data)) {
		return errors.New("too few charge slots provided")
	}

	if constraints.MaxSlots != 0 && constraints.MaxSlots < uint(len(data)) {
		return errors.New("too many charge slots provided")
	}

	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
	}
	desc, err := evTimeSeries.GetDescriptionsForFilter(filter)
	if err != nil || len(desc) == 0 {
		return api.ErrDataNotAvailable
	}

	timeSeriesSlots := []model.TimeSeriesSlotType{}
	var totalDuration time.Duration
	for index, slot := range data {
		relativeStart := totalDuration

		slotid := uint(index) //nolint:gosec // disable G115
		timeSeriesSlot := model.TimeSeriesSlotType{
			TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(slotid)),
			TimePeriod: &model.TimePeriodType{
				StartTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(relativeStart),
			},
			MaxValue: model.NewScaledNumberType(slot.Value),
		}

		// the last slot also needs an End Time
		if index == len(data)-1 {
			relativeEndTime := relativeStart + slot.Duration
			timeSeriesSlot.TimePeriod.EndTime = model.NewAbsoluteOrRelativeTimeTypeFromDuration(relativeEndTime)
		}
		timeSeriesSlots = append(timeSeriesSlots, timeSeriesSlot)

		totalDuration += slot.Duration
	}

	timeSeriesData := model.TimeSeriesDataType{
		TimeSeriesId: desc[0].TimeSeriesId,
		TimePeriod: &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromDuration(totalDuration),
		},
		TimeSeriesSlot: timeSeriesSlots,
	}

	_, err = evTimeSeries.WriteData([]model.TimeSeriesDataType{timeSeriesData})

	return err
}

func (e *CEVC) defaultPowerLimits(entity spineapi.EntityRemoteInterface) ([]ucapi.DurationSlotValue, error) {
	// send default power limits for the maximum timeframe
	// to fullfill spec, as there is no data provided
	logging.Log().Info("Fallback sending default power limits")

	evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		logging.Log().Error("electrical connection feature not found")
		return nil, err
	}

	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACPower),
	}
	paramDesc, err := evElectricalConnection.GetParameterDescriptionsForFilter(filter)
	if err != nil || len(paramDesc) == 0 || paramDesc[0].ParameterId == nil {
		logging.Log().Error("Error getting parameter descriptions:", err)
		return nil, err
	}

	filter2 := model.ElectricalConnectionPermittedValueSetDataType{
		ParameterId: paramDesc[0].ParameterId,
	}
	permitted, err := evElectricalConnection.GetPermittedValueSetForFilter(filter2)
	if err != nil || len(permitted) == 0 {
		logging.Log().Error("Error getting permitted values:", err)
		return nil, err
	}

	if len(permitted[0].PermittedValueSet) == 0 || len(permitted[0].PermittedValueSet[0].Range) == 0 {
		text := "No permitted value set available"
		logging.Log().Error(text)
		return nil, errors.New(text)
	}

	data := []ucapi.DurationSlotValue{
		{
			Duration: 7 * time.Hour * 24,
			Value:    permitted[0].PermittedValueSet[0].Range[0].Max.GetValue(),
		},
	}
	return data, nil
}
