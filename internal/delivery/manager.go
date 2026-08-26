package delivery

import (
	"errors"

	"telemetry/internal/state"
)

var ErrTemporary = errors.New("temporary delivery error")

type External interface {
	Do(idempotencyKey string) (late func(), err error)
}

type Manager struct{ Store *state.Store }

func (m Manager) Retry(id string, external External) error {
	m.Store.Update(state.Record{ID: id, Version: 1, Status: "running"})
	// Use a single idempotency key (the id itself) across all attempts so the
	// external system deduplicates: an earlier attempt that already took effect
	// is not repeated by a later attempt. Different deliveries use different ids,
	// so they never collide on the external side.
	idempotencyKey := id
	var winningCallback func()
	var winningErr error
	for attempt := 1; attempt <= 2; attempt++ {
		callback, err := external.Do(idempotencyKey)
		if err == nil {
			// This is the decisive attempt. Discard any callback captured by an
			// earlier attempt — it is stale and must not overwrite this result.
			winningCallback = callback
			winningErr = nil
			break
		}
		// On a temporary error, keep the latest callback in case this is the
		// final attempt, but never execute one from an earlier attempt.
		winningCallback = callback
		winningErr = err
		if !errors.Is(err, ErrTemporary) {
			break
		}
	}
	if winningErr == nil {
		m.Store.Update(state.Record{ID: id, Version: 2, Status: "success"})
		if winningCallback != nil {
			winningCallback()
		}
		return nil
	}
	return winningErr
}
