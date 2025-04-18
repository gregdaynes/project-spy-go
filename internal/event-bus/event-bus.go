package event_bus

import (
	"sync"
)

type EventBus[t any] struct {
	subscribers map[eventType]map[Subscriber]struct{}
	mutex       sync.Mutex
}

type eventType string

type Subscriber *func(string)

func NewEventBus[T any]() *EventBus[T] {

	return &EventBus[T]{
		subscribers: make(map[eventType]map[Subscriber]struct{}),
	}
}

func (eb *EventBus[T]) Subscribe(eventType eventType, subscriber Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = make(map[Subscriber]struct{})
	}

	eb.subscribers[eventType][subscriber] = struct{}{}
}

func (eb *EventBus[T]) Unsubscribe(eventType eventType, subscriber Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if subscribers, ok := eb.subscribers[eventType]; ok {
		delete(subscribers, subscriber)

		if len(subscribers) == 0 {
			delete(eb.subscribers, eventType)
		}
	}
}

func (eb *EventBus[T]) Publish(eventType eventType, event string) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if subscribers, ok := eb.subscribers[eventType]; ok {
		for fn := range subscribers {
			(*fn)(event)
		}
	}
}
