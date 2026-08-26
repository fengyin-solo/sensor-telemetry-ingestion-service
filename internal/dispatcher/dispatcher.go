package dispatcher

import (
	"context"
	"sync"

	"telemetry/internal/scheduler"
)

type Dispatcher struct{ workers sync.WaitGroup }

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
