package vabd

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *VABD) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.inverterConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.MeasurementDescriptionListDataType:
		e.inverterMeasurementDescriptionDataUpdate(payload.Entity)

	case *model.MeasurementListDataType:
		e.inverterMeasurementDataUpdate(payload)
	}
}

// process required steps when a grid device is connected
func (e *VABD) inverterConnected(entity spineapi.EntityRemoteInterface) {
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

// the measurement descriptiondata of an SMGW was updated
func (e *VABD) inverterMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestData(nil, nil); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *VABD) inverterMeasurementDataUpdate(payload spineapi.EventPayload) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 1
		filter := model.MeasurementDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
		}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePower)
		}

		// Scenario 2
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeCharge)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyCharged)
		}

		// Scenario 3
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeDischarge)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyDischarged)
		}

		// Scenario 4
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeStateOfCharge)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateStateOfCharge)
		}
	}
}
