package main

import (
	"sync"
)

type EventBus[t any] struct {
	subscribers map[EventType]map[Subscriber]struct{}
	mutex       sync.Mutex
}

type EventType string

type Subscriber *func(string)

func NewEventBus[T any]() *EventBus[T] {

	return &EventBus[T]{
		subscribers: make(map[EventType]map[Subscriber]struct{}),
	}
}

func (eb *EventBus[T]) Subscribe(eventType EventType, subscriber Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = make(map[Subscriber]struct{})
	}

	eb.subscribers[eventType][subscriber] = struct{}{}
}

func (eb *EventBus[T]) Unsubscribe(eventType EventType, subscriber Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if subscribers, ok := eb.subscribers[eventType]; ok {
		delete(subscribers, subscriber)

		if len(subscribers) == 0 {
			delete(eb.subscribers, eventType)
		}
	}
}

func (eb *EventBus[T]) Publish(eventType EventType, event string) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if subscribers, ok := eb.subscribers[eventType]; ok {
		for fn := range subscribers {
			(*fn)(event)
		}
	}
}

// func x() {
// 	// Create a new event bus for string events.
// 	eventBus := NewEventBus[string]()
//
// 	// Define a subscriber function for string events.
// 	stringSubscriber := func(event string) {
// 		fmt.Println("Received string event:", event)
// 	}
//
// 	// Subscribe to a specific event type.
// 	eventBus.Subscribe("stringEvent", &stringSubscriber)
//
// 	// Publish a string event.
// 	eventBus.Publish("stringEvent", "Hello, Event Bus!")
//
// 	// Unsubscribe the subscriber.
// 	// eventBus.Unsubscribe("stringEvent", &stringSubscriber)
//
// 	// The subscriber will not receive events after unsubscribing.
// 	eventBus.Publish("stringEvent", "This event won't be received.")
//
// 	// Sleep to allow time for the event bus to finish processing.
// 	time.Sleep(time.Second)
//
// 	eventBus.Publish("stringEvent", "Last event!")
// }
