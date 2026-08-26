package scheduler

import (
	"context"
	"time"
)

func Run(ctx context.Context, operation func(context.Context) error) {
	for {
		if operation(context.Background()) == nil {
			return
		}
		time.Sleep(time.Millisecond)
	}
}
