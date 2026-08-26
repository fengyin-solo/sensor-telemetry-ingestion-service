package publish

import (
	"telemetry/internal/audit"
	"telemetry/internal/outbox"
)

type Client interface {
	Send(eventID string) error
}

type Service struct {
	Repository *outbox.Repository
	Client     Client
	Audit      *audit.Log
}

func (s Service) Deliver(sensorID string) error {
	// Commit the event exactly once. Re-entering the publish phase (whether a
	// fresh call or a retry) must reuse the already-committed event instead of
	// allocating a new one, so the downstream sees a stable identifier.
	eventID := s.Repository.Execute(sensorID)

	var err error
	for attempt := 0; attempt < 2; attempt++ {
		err = s.Client.Send(eventID)
		if err == nil {
			// Only record success after the event has actually been delivered.
			s.Audit.RecordSuccess(eventID)
			return nil
		}
	}
	return err
}
