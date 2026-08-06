package crht

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *CRHT) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.hvacRoomConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.HvacSystemFunctionDescriptionListDataType:
		e.hvacSystemFunctionDescriptionDataUpdate(payload.Entity)

	case *model.SetpointDescriptionListDataType:
		e.setpointDescriptionDataUpdate(payload.Entity)

	case *model.SetpointListDataType:
		e.setpointDataUpdate(payload)

	case *model.SetpointConstraintsListDataType:
		e.setpointConstraintsDataUpdate(payload)
	}
}

// process required steps when an HVAC room device entity is connected
func (e *CRHT) hvacRoomConnected(entity spineapi.EntityRemoteInterface) {
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
	}

	if setpoint, err := client.NewSetpoint(e.LocalEntity, entity); err == nil {
		if !setpoint.HasSubscription() {
			if _, err := setpoint.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		selector := &model.SetpointDescriptionListDataSelectorsType{
			ScopeType: util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
		}
		if _, err := setpoint.RequestSetpointDescriptions(selector, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the HVAC system function description data of a device was updated
func (e *CRHT) hvacSystemFunctionDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		// system function descriptions received, now get the setpoint relations
		if _, err := hvac.RequestHvacSystemFunctionSetpointRelations(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the setpoint description data of a device was updated
func (e *CRHT) setpointDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if setpoint, err := client.NewSetpoint(e.LocalEntity, entity); err == nil {
		// setpoint descriptions received, now get the constraints and data
		if _, err := setpoint.RequestSetpointConstraints(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := setpoint.RequestSetpoints(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the setpoint data of a device was updated
func (e *CRHT) setpointDataUpdate(payload spineapi.EventPayload) {
	if setpoint, err := client.NewSetpoint(e.LocalEntity, payload.Entity); err == nil {
		filter := model.SetpointDataType{}
		if setpoint.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSetpoints)
		}
	}
}

// the setpoint constraints data of a device was updated
func (e *CRHT) setpointConstraintsDataUpdate(payload spineapi.EventPayload) {
	if data, ok := payload.Data.(*model.SetpointConstraintsListDataType); ok &&
		len(data.SetpointConstraintsData) > 0 && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSetpointConstraints)
	}
}
