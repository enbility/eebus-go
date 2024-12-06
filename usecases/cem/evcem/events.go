package evcem

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *EVCEM) HandleEvent(payload spineapi.EventPayload) {
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
	case *model.ElectricalConnectionDescriptionListDataType:
		e.evElectricalConnectionDescriptionDataUpdate(payload)

	case *model.MeasurementDescriptionListDataType:
		e.evMeasurementDescriptionDataUpdate(payload.Entity)

	case *model.MeasurementListDataType:
		e.evMeasurementDataUpdate(payload)
	}
}

// an EV was connected
func (e *EVCEM) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions

	if evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity); err == nil {
		if !evElectricalConnection.HasSubscription() {
			if _, err := evElectricalConnection.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get electrical connection descriptions
		if _, err := evElectricalConnection.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical connection parameter descriptions
		if _, err := evElectricalConnection.RequestParameterDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		if !evMeasurement.HasSubscription() {
			if _, err := evMeasurement.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get measurement descriptions
		if _, err := evMeasurement.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement constraints
		if _, err := evMeasurement.RequestConstraints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the electrical connection description data of an EV was updated
func (e *EVCEM) evElectricalConnectionDescriptionDataUpdate(payload spineapi.EventPayload) {
	if payload.Data == nil {
		return
	}

	data, ok := payload.Data.(*model.ElectricalConnectionDescriptionListDataType)
	if !ok {
		return
	}

	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId != nil && item.AcConnectedPhases != nil && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePhasesConnected)
			return
		}
	}
}

// the measurement description data of an EV was updated
func (e *EVCEM) evMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		// get measurement values
		if _, err := evMeasurement.RequestData(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the measurement data of an EV was updated
func (e *EVCEM) evMeasurementDataUpdate(payload spineapi.EventPayload) {
	if evMeasurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 1
		filter := model.MeasurementDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeACCurrent),
		}
		if evMeasurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentPerPhase)
		}

		// Scenario 2
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACPower)
		if evMeasurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePowerPerPhase)
		}

		// Scenario 3
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeCharge)
		if evMeasurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyCharged)
		}
	}
}
