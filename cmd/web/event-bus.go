package main

import (
	"fmt"
	"sync"
	"time"
)

type EventBus[t any] struct {
	subscribers map[EventType]map[Subscriber]SubscriberFunc
	mutex       sync.Mutex
}

type EventType string

type Subscriber string

type SubscriberFunc func(event string)

func NewEventBus[T any]() *EventBus[T] {
	return &EventBus[T]{
		subscribers: make(map[EventType]map[Subscriber]SubscriberFunc),
	}
}

func (eb *EventBus[T]) Subscribe(eventType EventType, subscriber Subscriber, fn SubscriberFunc) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = make(map[Subscriber]SubscriberFunc)
	}

	eb.subscribers[eventType][subscriber] = fn
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
		for _, fn := range subscribers {
			fn(event)
		}
	}
}

func x() {
	// Create a new event bus for string events.
	eventBus := NewEventBus[string]()

	// Define a subscriber function for string events.
	stringSubscriber := func(event string) {
		fmt.Println("Received string event:", event)
	}

	// Subscribe to a specific event type.
	eventBus.Subscribe("stringEvent", "mySubscriber", stringSubscriber)

	// Publish a string event.
	eventBus.Publish("stringEvent", "Hello, Event Bus!")

	// Unsubscribe the subscriber.
	eventBus.Unsubscribe("stringEvent", "mySubscriber")

	// The subscriber will not receive events after unsubscribing.
	eventBus.Publish("stringEvent", "This event won't be received.")

	// Sleep to allow time for the event bus to finish processing.
	time.Sleep(time.Second)

	eventBus.Publish("stringEvent", "Last event!")
}
