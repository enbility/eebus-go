package internal

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var PhaseNameMapping = []model.ElectricalConnectionPhaseNameType{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC}

func IsDeviceConnected(payload spineapi.EventPayload) bool {
	return payload.Device != nil &&
		payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeAdd
}

func IsDeviceDisconnected(payload spineapi.EventPayload) bool {
	return payload.Device != nil &&
		payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeRemove
}

func IsEntityConnected(payload spineapi.EventPayload) bool {
	if payload.Entity != nil &&
		payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeAdd {
		return true
	}

	return false
}

func IsEntityDisconnected(payload spineapi.EventPayload) bool {
	if payload.Entity != nil &&
		payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeRemove {
		return true
	}

	return false
}

func Deref(v *string) string {
	if v != nil {
		return string(*v)
	}
	return ""
}
