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
	var late func()
	for attempt := 1; attempt <= 2; attempt++ {
		callback, err := external.Do(id)
		if callback != nil {
			late = callback
		}
		if err == nil {
			m.Store.Update(state.Record{ID: id, Version: 2, Status: "success"})
			if late != nil {
				late()
			}
			return nil
		}
		if !errors.Is(err, ErrTemporary) {
			return err
		}
	}
	return ErrTemporary
}
