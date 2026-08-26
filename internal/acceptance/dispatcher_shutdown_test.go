package acceptance

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"telemetry/internal/dispatcher"
)

func TestCancelledDispatchStopsRetriesBeforeShutdown(t *testing.T) {
	requestCtx, cancel := context.WithTimeout(context.Background(), 8*time.Millisecond)
	defer cancel()
	var calls atomic.Int32
	var worker dispatcher.Dispatcher
	worker.Start(requestCtx, func(ctx context.Context) error {
		calls.Add(1)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Millisecond):
			return errors.New("endpoint unavailable")
		}
	})
	<-requestCtx.Done()
	time.Sleep(8 * time.Millisecond)
	before := calls.Load()
	time.Sleep(8 * time.Millisecond)
	after := calls.Load()
	if after != before {
		t.Fatalf("dispatch retries continued after request cancellation: before=%d after=%d", before, after)
	}
	shutdownCtx, stop := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer stop()
	if err := worker.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("dispatcher did not shut down after cancellation: %v", err)
	}
}
