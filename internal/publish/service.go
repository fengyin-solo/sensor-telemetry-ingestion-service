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
	eventID := s.Repository.Execute(sensorID)
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		err = s.Client.Send(eventID)
		if err == nil {
			s.Audit.RecordSuccess(eventID)
			return nil
		}
	}
	return err
}
