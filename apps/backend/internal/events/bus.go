package events

import (
	"context"
	"sync"
)

type EventType string

type Event interface {
	Type() EventType
}

type EventHandler func(ctx context.Context, event Event) error

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType EventType, handler EventHandler)
}

type SyncEventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]EventHandler
}

func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{
		subscribers: make(map[EventType][]EventHandler),
	}
}

func (b *SyncEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, exists := b.subscribers[event.Type()]
	b.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (b *SyncEventBus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}
