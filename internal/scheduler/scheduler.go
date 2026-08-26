package scheduler

import (
	"context"
	"time"
)

// Run executes operation repeatedly, retrying while it returns a non-nil error.
// The loop honours ctx: it stops as soon as ctx is cancelled, and the wait
// between attempts is interruptible so a cancelled request stops promptly
// instead of blocking on a sleep.
func Run(ctx context.Context, operation func(context.Context) error) {
	for {
		if err := ctx.Err(); err != nil {
			return
		}
		if operation(ctx) == nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond):
		}
	}
}
