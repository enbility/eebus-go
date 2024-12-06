package mpc

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *MPC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.deviceConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.MeasurementDescriptionListDataType:
		e.deviceMeasurementDescriptionDataUpdate(payload.Entity)

	case *model.MeasurementListDataType:
		e.deviceMeasurementDataUpdate(payload)
	}
}

// process required steps when a device is connected
func (e *MPC) deviceConnected(entity spineapi.EntityRemoteInterface) {
	if electricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity); err == nil {
		if !electricalConnection.HasSubscription() {
			if _, err := electricalConnection.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get electrical connection parameter
		if _, err := electricalConnection.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := electricalConnection.RequestParameterDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}

	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		if !measurement.HasSubscription() {
			if _, err := measurement.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get measurement parameters
		if _, err := measurement.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := measurement.RequestConstraints(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the measurement descriptiondata of a device was updated
func (e *MPC) deviceMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestData(nil, nil); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of a device was updated
func (e *MPC) deviceMeasurementDataUpdate(payload spineapi.EventPayload) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 1
		filter := model.MeasurementDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
		}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePower)
		}

		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACPower)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePowerPerPhase)
		}

		// Scenario 2
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACEnergyConsumed)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyConsumed)
		}

		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACEnergyProduced)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyProduced)
		}

		// Scenario 3
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACCurrent)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentsPerPhase)
		}

		// Scenario 4
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACVoltage)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateVoltagePerPhase)
		}

		// Scenario 5
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACFrequency)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFrequency)
		}
	}
}
