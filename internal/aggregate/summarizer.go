package aggregate

import (
	"context"

	"telemetry/internal/model"
	"telemetry/internal/store"
)

type Summary struct {
	Count int
	Total int
}

type Summarizer struct {
	registry *store.Registry
}

func NewSummarizer(registry *store.Registry) *Summarizer {
	return &Summarizer{registry: registry}
}

func cloneReadings(input map[string]model.Reading) map[string]model.Reading {
	result := make(map[string]model.Reading, len(input))
	for id, reading := range input {
		result[id] = reading
	}
	return result
}

func (s *Summarizer) SnapshotThen(ctx context.Context, ready chan<- struct{}, proceed <-chan struct{}) (Summary, error) {
	snapshot := cloneReadings(s.registry.Snapshot())
	if ready != nil {
		close(ready)
	}
	if proceed != nil {
		select {
		case <-proceed:
		case <-ctx.Done():
			return Summary{}, ctx.Err()
		}
	}

	var summary Summary
	for _, reading := range snapshot {
		summary.Count++
		summary.Total += reading.Value
	}
	return summary, nil
}
