package spine

import "sync"

var Events events

type ElementChangeType uint16

const (
	ElementChangeAdd ElementChangeType = iota
	ElementChangeUpdate
	ElementChangeRemove
)

type EventType uint16

const (
	EventTypeDeviceChange EventType = iota
	EventTypeEntityChange
	EventTypeSubscriptionChange
	EventTypeDataChange
)

type EventPayload struct {
	Ski        string
	EventType  EventType
	ChangeType ElementChangeType
	Device     *DeviceRemoteImpl // required for DetailedDiscovery Call
	Entity     *EntityRemoteImpl // required for DetailedDiscovery Call and Notify
	Feature    *FeatureRemoteImpl
	Data       interface{}
}

type EventHandler interface {
	HandleEvent(EventPayload)
}

type events struct {
	mu       sync.Mutex
	handlers []EventHandler
}

func (r *events) Subscribe(handler EventHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers = append(r.handlers, handler)
}

func (r *events) Publish(payload EventPayload) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, handler := range r.handlers {
		go handler.HandleEvent(payload)
	}
}
