package dispatcher

import (
	"context"
	"sync"

	"telemetry/internal/scheduler"
)

type Dispatcher struct{ workers sync.WaitGroup }

// Start launches a worker that runs operation under the supplied ctx.
// The ctx drives the worker's lifetime: when the request is cancelled (for
// example after a timeout) the worker stops retrying instead of looping
// forever, and Shutdown is not held up by background work that can no longer
// succeed. context.Background() is deliberately not used here.
func (d *Dispatcher) Start(ctx context.Context, operation func(context.Context) error) {
	d.workers.Add(1)
	go func() {
		defer d.workers.Done()
		scheduler.Run(ctx, operation)
	}()
}

func (d *Dispatcher) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		d.workers.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
