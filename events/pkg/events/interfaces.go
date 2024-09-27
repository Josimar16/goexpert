package events

import (
	"sync"
	"time"
)

type Event interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
}

type EventHandler interface {
	Handle(event Event, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	RegisterHandler(eventName string, handler EventHandler) error
	Dispatch(event Event) error
	Remove(eventName string, handler EventHandler) error
	Has(eventName string, handler EventHandler) bool
	Clear() error
}
