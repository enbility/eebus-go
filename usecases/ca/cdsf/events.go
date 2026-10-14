package cdsf

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *CDSF) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.dhwCircuitConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.HvacSystemFunctionDescriptionListDataType:
		e.hvacSystemFunctionDescriptionDataUpdate(payload.Entity)

	case *model.HvacSystemFunctionListDataType:
		e.hvacSystemFunctionDataUpdate(payload)

	case *model.HvacOverrunListDataType:
		e.hvacOverrunDataUpdate(payload)
	}
}

// process required steps when an DHW circuit device entity is connected
func (e *CDSF) dhwCircuitConnected(entity spineapi.EntityRemoteInterface) {
	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		if !hvac.HasSubscription() {
			if _, err := hvac.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		if _, err := hvac.RequestHvacSystemFunctionDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := hvac.RequestHvacOperationModeDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := hvac.RequestHvacOverrunDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the HVAC system function description data of a device was updated
func (e *CDSF) hvacSystemFunctionDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		// system function descriptions received, now get the relations and data
		if _, err := hvac.RequestHvacSystemFunctionOperationModeRelations(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := hvac.RequestHvacSystemFunctions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := hvac.RequestHvacOverruns(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the HVAC overrun data of a device was updated
func (e *CDSF) hvacOverrunDataUpdate(payload spineapi.EventPayload) {
	if data, ok := payload.Data.(*model.HvacOverrunListDataType); ok &&
		len(data.HvacOverrunData) > 0 && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateOverrun)
	}
}

// the HVAC system function data of a device was updated
func (e *CDSF) hvacSystemFunctionDataUpdate(payload spineapi.EventPayload) {
	if hvac, err := client.NewHvac(e.LocalEntity, payload.Entity); err == nil {
		filter := model.HvacSystemFunctionDataType{}
		if hvac.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateOperationMode)
		}
	}
}
