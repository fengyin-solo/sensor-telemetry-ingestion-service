package transmit

import (
	"errors"

	"telemetry/internal/adapter"
	"telemetry/internal/journal"
)

type Remote interface {
	Send(payload string) error
}

type Processor struct {
	Remote  Remote
	Journal *journal.Memory
}

func (p *Processor) Process(payload string) error {
	for attempt := 0; attempt < 3; attempt++ {
		err := adapter.Translate(p.Remote.Send(payload))
		if err == nil {
			return p.Journal.RecordOutcome(payload, nil)
		}
		if recordErr := p.Journal.RecordOutcome(payload, err); recordErr != nil {
			var rejected *adapter.RejectedError
			if errors.As(recordErr, &rejected) {
				return recordErr
			}
		}
		if attempt == 2 {
			return err
		}
	}
	return nil
}
