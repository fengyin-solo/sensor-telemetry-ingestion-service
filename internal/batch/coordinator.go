package batch

import (
	"context"

	"telemetry/internal/model"
	"telemetry/internal/source"
)

func Collect(ctx context.Context, sensorIDs []string) ([]model.Reading, error) {
	readings, errs := source.Stream(ctx, sensorIDs)
	result := make([]model.Reading, 0, len(sensorIDs))
	for reading := range readings {
		result = append(result, reading)
	}
	if err := <-errs; err != nil {
		return nil, err
	}
	return result, nil
}
