package outbox

import (
	"fmt"
	"sync"
)

type Repository struct {
	mu       sync.Mutex
	writes   int
	bySensor map[string]string
}

func NewRepository() *Repository {
	return &Repository{bySensor: make(map[string]string)}
}

func (r *Repository) Execute(sensorID string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.writes++
	eventID := fmt.Sprintf("event-%d", r.writes)
	r.bySensor[sensorID] = eventID
	return eventID
}

func (r *Repository) Writes() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.writes
}
