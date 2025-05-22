package event_bus

import (
	"sync"
)

type EventBus[T any] struct {
	subscribers map[eventType]map[SubscriberID]Subscriber
	mutex       sync.RWMutex
}

type eventType string
type SubscriberID string
type Subscriber func(string)

func NewEventBus[T any]() *EventBus[T] {
	return &EventBus[T]{
		subscribers: make(map[eventType]map[SubscriberID]Subscriber),
	}
}

func (eb *EventBus[T]) Subscribe(eventType eventType, id SubscriberID, subscriber Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = make(map[SubscriberID]Subscriber)
	}

	eb.subscribers[eventType][id] = subscriber
}

func (eb *EventBus[T]) Unsubscribe(eventType eventType, id SubscriberID) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if subscribers, ok := eb.subscribers[eventType]; ok {
		delete(subscribers, id)
		if len(subscribers) == 0 {
			delete(eb.subscribers, eventType)
		}
	}
}

func (eb *EventBus[T]) Publish(eventType eventType, event string) {
	eb.mutex.RLock()
	subscribers := make(map[SubscriberID]Subscriber)
	for id, fn := range eb.subscribers[eventType] {
		subscribers[id] = fn
	}
	eb.mutex.RUnlock()

	for _, fn := range subscribers {
		fn(event)
	}
}
