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

// Execute commits the event for sensorID exactly once and returns its ID.
// Re-entering the publish phase for an already-committed sensor must not
// allocate a new event or bump the write counter; subsequent calls reuse the
// existing event so retries send the same identifier the downstream already
// saw.
func (r *Repository) Execute(sensorID string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if eventID, ok := r.bySensor[sensorID]; ok {
		return eventID
	}
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
