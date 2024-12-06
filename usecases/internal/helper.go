package internal

import (
	spineapi "github.com/enbility/spine-go/api"
)

func IsDeviceConnected(payload spineapi.EventPayload) bool {
	return payload.Device != nil &&
		payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeAdd
}

func IsDeviceDisconnected(payload spineapi.EventPayload) bool {
	return len(payload.Ski) > 0 &&
		payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeRemove
}

func IsEntityAdded(payload spineapi.EventPayload) bool {
	if payload.Entity != nil &&
		payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeAdd {
		return true
	}

	return false
}

func IsEntityRemoved(payload spineapi.EventPayload) bool {
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
