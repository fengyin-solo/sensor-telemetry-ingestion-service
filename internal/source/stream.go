package source

import (
	"context"
	"errors"
	"strings"

	"telemetry/internal/model"
)

var ErrMalformed = errors.New("malformed telemetry frame")

func Stream(ctx context.Context, sensorIDs []string) (<-chan model.Reading, <-chan error) {
	readings := make(chan model.Reading)
	errs := make(chan error, 1)
	go func() {
		defer close(readings)
		defer close(errs)
		for index, sensorID := range sensorIDs {
			if strings.HasPrefix(sensorID, "bad:") {
				errs <- ErrMalformed
				return
			}
			select {
			case readings <- model.Reading{SensorID: sensorID, Value: index + 1}:
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			}
		}
	}()
	return readings, errs
}
