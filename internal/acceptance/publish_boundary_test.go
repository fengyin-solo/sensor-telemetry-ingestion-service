package acceptance

import (
	"errors"
	"testing"

	"telemetry/internal/audit"
	"telemetry/internal/outbox"
	"telemetry/internal/publish"
)

type failOnceClient struct {
	calls []string
}

func (c *failOnceClient) Send(eventID string) error {
	c.calls = append(c.calls, eventID)
	if len(c.calls) == 1 {
		return errors.New("publisher unavailable")
	}
	return nil
}

func TestPublishRetryReusesCommittedEventAndAuditsOnce(t *testing.T) {
	repository := outbox.NewRepository()
	repository.Execute("sensor-a")
	client := &failOnceClient{}
	auditLog := &audit.Log{}
	err := (publish.Service{Repository: repository, Client: client, Audit: auditLog}).Deliver("sensor-a")
	if err != nil {
		t.Fatalf("event publication did not recover on retry: %v", err)
	}
	if repository.Writes() != 1 {
		t.Fatalf("publication retry repeated the telemetry transaction: writes=%d", repository.Writes())
	}
	if len(client.calls) != 2 || client.calls[0] != client.calls[1] {
		t.Fatalf("publication retry did not reuse one event: calls=%v", client.calls)
	}
	if successes := auditLog.Successes(); len(successes) != 1 || successes[0] != client.calls[1] {
		t.Fatalf("success audit was emitted before final publication: successes=%v", successes)
	}
}
