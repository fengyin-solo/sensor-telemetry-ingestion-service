package store

import (
	"sync"

	"telemetry/internal/model"
)

type Registry struct {
	mu       sync.RWMutex
	readings map[string]model.Reading
}

func NewRegistry() *Registry {
	return &Registry{readings: make(map[string]model.Reading)}
}

func (r *Registry) Put(reading model.Reading) {
	r.mu.Lock()
	r.readings[reading.SensorID] = reading
	r.mu.Unlock()
}

func (r *Registry) Snapshot() map[string]model.Reading {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := make(map[string]model.Reading, len(r.readings))
	for id, reading := range r.readings {
		snapshot[id] = reading
	}
	return snapshot
}
