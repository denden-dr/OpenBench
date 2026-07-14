package events

import (
	"context"
	"log/slog"
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

type eventWrapper struct {
	ctx   context.Context
	event Event
}

type AsyncEventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]EventHandler
	ch          chan eventWrapper
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewAsyncEventBus(bufferSize int) *AsyncEventBus {
	ctx, cancel := context.WithCancel(context.Background())
	bus := &AsyncEventBus{
		subscribers: make(map[EventType][]EventHandler),
		ch:          make(chan eventWrapper, bufferSize),
		ctx:         ctx,
		cancel:      cancel,
	}
	bus.wg.Add(1)
	go bus.worker()
	return bus
}

func (b *AsyncEventBus) worker() {
	defer b.wg.Done()
	for {
		select {
		case <-b.ctx.Done():
			// Process remaining events in the channel
			for wrap := range b.ch {
				b.dispatch(wrap.ctx, wrap.event)
			}
			return
		case wrap, ok := <-b.ch:
			if !ok {
				return
			}
			b.dispatch(wrap.ctx, wrap.event)
		}
	}
}

func (b *AsyncEventBus) dispatch(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers, exists := b.subscribers[event.Type()]
	b.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					slog.ErrorContext(ctx, "recovered from event handler panic",
						slog.Any("recover", r),
						slog.String("event_type", string(event.Type())),
					)
				}
			}()
			if err := handler(ctx, event); err != nil {
				slog.ErrorContext(ctx, "event handler failed",
					slog.Any("error", err),
					slog.String("event_type", string(event.Type())),
				)
			}
		}()
	}
}

func (b *AsyncEventBus) Publish(ctx context.Context, event Event) error {
	select {
	case <-b.ctx.Done():
		return context.Canceled
	case b.ch <- eventWrapper{ctx: ctx, event: event}:
		return nil
	}
}

func (b *AsyncEventBus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *AsyncEventBus) Close() {
	b.cancel()
	close(b.ch)
	b.wg.Wait()
}
