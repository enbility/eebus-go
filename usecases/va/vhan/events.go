package vhan

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *VHAN) HandleEvent(payload spineapi.EventPayload) {
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

	if _, ok := payload.Data.(*model.DeviceClassificationUserDataType); ok {
		e.deviceUserDataUpdate(payload)
	}
}

// process required steps when a device is connected
func (e *VHAN) deviceConnected(entity spineapi.EntityRemoteInterface) {
	if deviceClassification, err := client.NewDeviceClassification(e.LocalEntity, entity); err == nil {
		if !deviceClassification.HasSubscription() {
			if _, err := deviceClassification.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get the heating area name
		if _, err := deviceClassification.RequestUserData(); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the user data of a device was updated
func (e *VHAN) deviceUserDataUpdate(payload spineapi.EventPayload) {
	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateName)
	}
}
