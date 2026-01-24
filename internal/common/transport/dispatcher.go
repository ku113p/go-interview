package transport

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrQueueFull = errors.New("transport dispatcher queue full")
	ErrClosed    = errors.New("transport dispatcher closed")
)

type Dispatcher struct {
	handler  HandlerFunc
	requests chan Envelope
	results  chan<- Result

	wg   sync.WaitGroup
	once sync.Once

	closed atomic.Bool

	workers int
}

func NewDispatcher(
	handler HandlerFunc,
	queueSize int,
	workers int,
	results chan<- Result,
) (*Dispatcher, error) {
	if handler == nil {
		return nil, fmt.Errorf("transport dispatcher: handler is required")
	}

	if queueSize <= 0 {
		return nil, fmt.Errorf("transport dispatcher: queue size must be positive")
	}

	if workers <= 0 {
		return nil, fmt.Errorf("transport dispatcher: workers must be positive")
	}

	if results == nil {
		return nil, fmt.Errorf("transport dispatcher: results channel is required")
	}

	d := &Dispatcher{
		handler:  handler,
		requests: make(chan Envelope, queueSize),
		results:  results,
		workers:  workers,
	}

	d.startWorkers()
	return d, nil
}

func (d *Dispatcher) Send(ctx context.Context, env Envelope) error {
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if d.closed.Load() {
		return ErrClosed
	}

	env.Ctx = ctx

	select {
	case d.requests <- env:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ErrQueueFull
	}
}

func (d *Dispatcher) Shutdown(ctx context.Context) error {
	d.once.Do(func() {
		d.closed.Store(true)
		close(d.requests)
	})

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done // ensure goroutines complete
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (d *Dispatcher) startWorkers() {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.runWorker()
	}
}
