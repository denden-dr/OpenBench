package events

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEvent struct {
	t EventType
}

func (e mockEvent) Type() EventType {
	return e.t
}

func TestAsyncEventBus_PublishAndSubscribe(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	bus := NewAsyncEventBus(10)
	defer bus.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	var received Event
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		received = event
		wg.Done()
		return nil
	})

	evt := mockEvent{t: "test.event"}
	err := bus.Publish(context.Background(), evt)
	must.NoError(err)

	// Wait with timeout
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event handler")
	}

	is.Equal(evt, received)
}

func TestAsyncEventBus_PanicRecovery(t *testing.T) {
	must := require.New(t)

	bus := NewAsyncEventBus(10)
	defer bus.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	bus.Subscribe("panic.event", func(ctx context.Context, event Event) error {
		defer wg.Done()
		panic("something went terribly wrong")
	})

	bus.Subscribe("panic.event", func(ctx context.Context, event Event) error {
		defer wg.Done()
		// This should still run even if the first one panicked
		return nil
	})

	evt := mockEvent{t: "panic.event"}
	err := bus.Publish(context.Background(), evt)
	must.NoError(err)

	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event handlers")
	}
}

func TestAsyncEventBus_HandlerError(t *testing.T) {
	must := require.New(t)

	bus := NewAsyncEventBus(10)
	defer bus.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	bus.Subscribe("error.event", func(ctx context.Context, event Event) error {
		defer wg.Done()
		return errors.New("handler error")
	})

	bus.Subscribe("error.event", func(ctx context.Context, event Event) error {
		defer wg.Done()
		// This should still run even if the first one returned an error
		return nil
	})

	evt := mockEvent{t: "error.event"}
	err := bus.Publish(context.Background(), evt)
	must.NoError(err)

	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event handlers")
	}
}
