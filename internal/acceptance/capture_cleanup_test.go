package acceptance

import (
	"errors"
	"fmt"
	"testing"

	"telemetry/internal/capture"
)

var errCaptureRead = errors.New("capture c is corrupt")

type trackedOpener struct {
	open    int
	maxOpen int
	limit   int
}

func (o *trackedOpener) Open(name string) (capture.Resource, error) {
	if o.open >= o.limit {
		return nil, fmt.Errorf("capture handle limit reached")
	}
	o.open++
	if o.open > o.maxOpen {
		o.maxOpen = o.open
	}
	return &trackedResource{name: name, owner: o}, nil
}

type trackedResource struct {
	name  string
	owner *trackedOpener
}

func (r *trackedResource) Read() error {
	if r.name == "c" {
		return errCaptureRead
	}
	return nil
}

func (r *trackedResource) Close() error {
	r.owner.open--
	return nil
}

type transactionProbe struct{ commits, rollbacks, audits int }

func (p *transactionProbe) Commit() error   { p.commits++; return nil }
func (p *transactionProbe) Rollback() error { p.rollbacks++; return nil }
func (p *transactionProbe) AuditSuccess()   { p.audits++ }

func TestCaptureFailureClosesHandlesAndRollsBack(t *testing.T) {
	opener := &trackedOpener{limit: 2}
	backend := &transactionProbe{}
	err := capture.Process(opener, backend, []string{"a", "b", "c", "d"})
	if !errors.Is(err, errCaptureRead) {
		t.Fatalf("capture processing did not preserve the read failure: %v", err)
	}
	if opener.maxOpen != 1 || opener.open != 0 {
		t.Fatalf("capture handles were not released per item: max=%d remaining=%d", opener.maxOpen, opener.open)
	}
	if backend.rollbacks != 1 || backend.commits != 0 || backend.audits != 0 {
		t.Fatalf("failed capture used the wrong transaction outcome: commits=%d rollbacks=%d audits=%d", backend.commits, backend.rollbacks, backend.audits)
	}
}
