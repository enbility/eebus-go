package spine

import (
	"errors"
	"sync"

	"github.com/enbility/eebus-go/spine/model"
)

var Events events

type EventHandlerLevel uint

const (
	EventHandlerLevelCore        EventHandlerLevel = iota // Shall only be used by the core stack
	EventHandlerLevelMiddleware                           // Shall only be used by middleware implementations, e.g. CEMd
	EventHandlerLevelApplication                          // Shall only be used by applications
)

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
	EventTypeBindingChange                       // Sent after successful binding request from remote
	EventTypeDataChange                          // Sent after remote provided new data items for a function
)

type EventPayload struct {
	Ski           string            // required
	EventType     EventType         // required
	ChangeType    ElementChangeType // required
	Device        *DeviceRemoteImpl // required for DetailedDiscovery Call
	Entity        *EntityRemoteImpl // required for DetailedDiscovery Call and Notify
	Feature       *FeatureRemoteImpl
	CmdClassifier *model.CmdClassifierType // optional, used together with EventType EventTypeDataChange
	Data          any
}

type EventHandler interface {
	HandleEvent(EventPayload)
}

type eventHandlerItem struct {
	Level   EventHandlerLevel
	Handler EventHandler
}

type events struct {
	mu       sync.Mutex
	handlers []eventHandlerItem // event handling outside of the core stack
}

// will be used in EEBUS core directly to access the level EventHandlerLevelCore
func (r *events) subscribe(level EventHandlerLevel, handler EventHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	exists := false
	for _, item := range r.handlers {
		if item.Level == level && item.Handler == handler {
			exists = true
			break
		}
	}

	if !exists {
		newHandlerItem := eventHandlerItem{
			Level:   level,
			Handler: handler,
		}
		r.handlers = append(r.handlers, newHandlerItem)
	}

	return nil
}

// Subscribe to message events and handle them in
// the Eventhandler interface implementation
//
// returns an error if EventHandlerLevelCore is used as
// that is only allowed for internal use
func (r *events) Subscribe(level EventHandlerLevel, handler EventHandler) error {
	if level == EventHandlerLevelCore {
		return errors.New("This level is restricted to the EEBUS core implenentation!")
	}

	return r.subscribe(level, handler)
}

// will be used in EEBUS core directly to access the level EventHandlerLevelCore
func (r *events) unsubscribe(level EventHandlerLevel, handler EventHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var newHandlers []eventHandlerItem
	for _, item := range r.handlers {
		if item.Level != level && item.Handler != handler {
			newHandlers = append(newHandlers, item)
		}
	}

	r.handlers = newHandlers

	return nil
}

// Unsubscribe from getting events
func (r *events) Unsubscribe(level EventHandlerLevel, handler EventHandler) error {
	if level == EventHandlerLevelCore {
		return errors.New("This level is restricted to the EEBUS core implenentation!")
	}

	return r.unsubscribe(level, handler)
}

// Publish an event to all subscribers
func (r *events) Publish(payload EventPayload) {
	r.mu.Lock()

	// process subscribers by level
	handlerLevels := []EventHandlerLevel{
		EventHandlerLevelCore,
		EventHandlerLevelMiddleware,
		EventHandlerLevelApplication,
	}

	for _, level := range handlerLevels {
		for _, item := range r.handlers {
			if item.Level != level {
				continue
			}

			// do not run this asynchronously, to make sure all required
			// and expected actions are taken
			item.Handler.HandleEvent(payload)
		}
	}

	r.mu.Unlock()
}
