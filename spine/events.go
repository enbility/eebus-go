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
	EventTypeDeviceChange       EventType = iota // Sent after successful response of NodeManagementDetailedDiscovery
	EventTypeEntityChange                        // Sent after successful response of NodeManagementDetailedDiscovery
	EventTypeSubscriptionChange                  // Sent after successful subscription request from remote
	EventTypeDataChange                          // Sent after remote provided new data items for a function
)

type EventPayload struct {
	Ski        string
	EventType  EventType
	ChangeType ElementChangeType
	Device     *DeviceRemoteImpl // required for DetailedDiscovery Call
	Entity     *EntityRemoteImpl // required for DetailedDiscovery Call and Notify
	Feature    *FeatureRemoteImpl
	Data       any
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
