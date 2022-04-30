package spine

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
	EventTypeDataChange
)

type EventPayload struct {
	EventType  EventType
	ChangeType ElementChangeType
	Device     *DeviceRemoteImpl
	Entity     *EntityRemoteImpl
	Data       interface{}
}

type EventHandler interface {
	HandleEvent(EventPayload)
}

type events struct {
	handlers []EventHandler
}

func (r *events) Subscribe(handler EventHandler) {
	r.handlers = append(r.handlers, handler)
}

func (r events) Publish(payload EventPayload) {
	for _, handler := range r.handlers {
		go handler.HandleEvent(payload)
	}
}
