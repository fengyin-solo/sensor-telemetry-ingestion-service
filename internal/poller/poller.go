package poller

import (
	"context"

	"telemetry/internal/backoff"
)

type Poller struct{ Attempts int }

func (p Poller) Poll(ctx context.Context, operation func(context.Context) error) error {
	attempts := p.Attempts
	if attempts < 1 {
		attempts = 1
	}
	return backoff.Run(context.Background(), attempts, operation)
}
