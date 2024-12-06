package cevc

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *CEVC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.evConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.TimeSeriesDescriptionListDataType:
		e.evTimeSeriesDescriptionDataUpdate(payload)

	case *model.TimeSeriesListDataType:
		e.evTimeSeriesDataUpdate(payload)

	case *model.IncentiveTableDescriptionDataType:
		e.evIncentiveTableDescriptionDataUpdate(payload)

	case *model.IncentiveTableConstraintsDataType:
		e.evIncentiveTableConstraintsDataUpdate(payload)

	case *model.IncentiveDataType:
		e.evIncentiveTableDataUpdate(payload)
	}
}

// an EV was connected
func (e *CEVC) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evDeviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		if !evDeviceConfiguration.HasSubscription() {
			if _, err := evDeviceConfiguration.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get device configuration descriptions
		if _, err := evDeviceConfiguration.RequestKeyValueDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity); err == nil {
		if !evTimeSeries.HasSubscription() {
			if _, err := evTimeSeries.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if !evTimeSeries.HasBinding() {
			if _, err := evTimeSeries.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get time series descriptions
		if _, err := evTimeSeries.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		// get time series constraints
		if _, err := evTimeSeries.RequestConstraints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, entity); err == nil {
		if !evIncentiveTable.HasSubscription() {
			if _, err := evIncentiveTable.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if !evIncentiveTable.HasBinding() {
			if _, err := evIncentiveTable.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get incentivetable descriptions
		if _, err := evIncentiveTable.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the time series description data of an EV was updated
func (e *CEVC) evTimeSeriesDescriptionDataUpdate(payload spineapi.EventPayload) {
	if evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, payload.Entity); err == nil {
		// get time series values
		if _, err := evTimeSeries.RequestData(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	// check if we are required to update the plan
	if !e.evCheckTimeSeriesDescriptionConstraintsUpdateRequired(payload.Entity) {
		return
	}

	_, err := e.EnergyDemand(payload.Entity)
	if err != nil {
		return
	}

	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyDemand)
	}

	_, err = e.TimeSlotConstraints(payload.Entity)
	if err != nil {
		logging.Log().Error("Error getting timeseries constraints:", err)
		return
	}

	_, err = e.IncentiveConstraints(payload.Entity)
	if err != nil {
		logging.Log().Error("Error getting incentive constraints:", err)
		return
	}

	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataRequestedPowerLimitsAndIncentives)
	}
}

// the load control limit data of an EV was updated
func (e *CEVC) evTimeSeriesDataUpdate(payload spineapi.EventPayload) {
	if _, err := e.ChargePlan(payload.Entity); err == nil && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateChargePlan)
	}

	if _, err := e.ChargePlanConstraints(payload.Entity); err == nil && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateTimeSlotConstraints)
	}
}

// the incentive table description data of an EV was updated
func (e *CEVC) evIncentiveTableDescriptionDataUpdate(payload spineapi.EventPayload) {
	if evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, payload.Entity); err == nil {
		// get time series values
		if _, err := evIncentiveTable.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// check if we are required to update the plan
	if e.evCheckIncentiveTableDescriptionUpdateRequired(payload.Entity) && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataRequestedIncentiveTableDescription)
	}
}

// the incentive table constraint data of an EV was updated
func (e *CEVC) evIncentiveTableConstraintsDataUpdate(payload spineapi.EventPayload) {
	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateIncentiveTable)
	}
}

// the incentive table data of an EV was updated
func (e *CEVC) evIncentiveTableDataUpdate(payload spineapi.EventPayload) {
	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateIncentiveTable)
	}
}

// check timeSeries descriptions if constraints element has updateRequired set to true
// as this triggers the CEM to send power tables within 20s
func (e *CEVC) evCheckTimeSeriesDescriptionConstraintsUpdateRequired(entity spineapi.EntityRemoteInterface) bool {
	evTimeSeries, err := client.NewTimeSeries(e.LocalEntity, entity)
	if err != nil {
		logging.Log().Error("timeseries feature not found")
		return false
	}

	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
	}
	data, err := evTimeSeries.GetDescriptionsForFilter(filter)
	if err != nil || len(data) == 0 {
		return false
	}

	if data[0].UpdateRequired != nil {
		return *data[0].UpdateRequired
	}

	return false
}

// check incentibeTable descriptions if the tariff description has updateRequired set to true
// as this triggers the CEM to send incentive tables within 20s
func (e *CEVC) evCheckIncentiveTableDescriptionUpdateRequired(entity spineapi.EntityRemoteInterface) bool {
	evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, entity)
	if err != nil {
		logging.Log().Error("incentivetable feature not found")
		return false
	}

	filter := model.TariffDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
	}
	data, err := evIncentiveTable.GetDescriptionsForFilter(filter)
	if err != nil || len(data) == 0 {
		return false
	}

	// only use the first description and therein the first tariff
	item := data[0].TariffDescription
	if item != nil && item.UpdateRequired != nil {
		return *item.UpdateRequired
	}

	return false
}
