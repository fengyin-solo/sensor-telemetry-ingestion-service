package acceptance

import (
	"testing"

	"telemetry/internal/delivery"
	"telemetry/internal/state"
)

type idempotentExternal struct {
	seen    map[string]bool
	effects int
	calls   int
	store   *state.Store
	id      string
}

func (e *idempotentExternal) Do(key string) (func(), error) {
	e.calls++
	if e.seen[key] {
		return nil, nil
	}
	e.seen[key] = true
	e.effects++
	if e.calls == 1 {
		return func() {
			e.store.Update(state.Record{ID: e.id, Version: 1, Status: "running"})
		}, delivery.ErrTemporary
	}
	return nil, nil
}

func TestRetryKeepsOneEffectAndRejectsLateState(t *testing.T) {
	storeInstance := state.NewStore()
	external := &idempotentExternal{seen: make(map[string]bool), store: storeInstance, id: "delivery-a"}
	err := (delivery.Manager{Store: storeInstance}).Retry("delivery-a", external)
	if err != nil {
		t.Fatalf("delivery retry did not finish successfully: %v", err)
	}
	if external.effects != 1 {
		t.Fatalf("delivery retry repeated the external effect: effects=%d", external.effects)
	}
	current := storeInstance.Get("delivery-a")
	cached := storeInstance.Cached("delivery-a")
	if current.Status != "success" || current.Version != 2 {
		t.Fatalf("late callback moved delivery out of success: %+v", current)
	}
	if cached != current {
		t.Fatalf("delivery detail and cached status disagree: current=%+v cached=%+v", current, cached)
	}
}
