package scheduler

import (
	"context"
	"time"
)

func Run(ctx context.Context, operation func(context.Context) error) {
	for {
		if ctx.Err() != nil {
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
