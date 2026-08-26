package backoff

import "context"

func Run(ctx context.Context, attempts int, operation func(context.Context) error) error {
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = operation(ctx)
		if lastErr == nil {
			return nil
		}
	}
	return lastErr
}
