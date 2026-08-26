package batch

import (
	"context"

	"telemetry/internal/model"
	"telemetry/internal/source"
)

func Collect(ctx context.Context, sensorIDs []string) ([]model.Reading, error) {
	readings, errs := source.Stream(ctx, sensorIDs)
	result := make([]model.Reading, 0, len(sensorIDs))
	for readings != nil || errs != nil {
		select {
		case reading, ok := <-readings:
			if !ok {
				readings = nil
				continue
			}
			result = append(result, reading)
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return result, nil
}
