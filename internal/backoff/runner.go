package backoff

import "context"

func Run(ctx context.Context, attempts int, operation func(context.Context) error) error {
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		lastErr = operation(context.Background())
		if lastErr == nil {
			return nil
		}
	}
	return lastErr
}
