package transport

import (
	"context"
	"runtime"
	"time"
)

func (d *Dispatcher) runWorker() {
	defer d.wg.Done()

	for {
		env, ok := <-d.requests
		if !ok {
			return
		}
		d.handle(env)
	}
}

func (d *Dispatcher) handle(env Envelope) {
	ctx := env.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	started := time.Now()
	value, err := d.handler(ctx, env)
	select {
	case d.results <- NewResult(env.TaskID, env.Module, env.Action, value, err, started):
	default:
		runtime.Gosched()
		d.results <- NewResult(env.TaskID, env.Module, env.Action, value, err, started)
	}
}
