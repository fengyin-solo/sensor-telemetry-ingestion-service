package journal

import "sync"

type Memory struct {
	mu      sync.Mutex
	entries []string
}

func (m *Memory) RecordOutcome(payload string, operationErr error) error {
	m.mu.Lock()
	m.entries = append(m.entries, payload)
	m.mu.Unlock()
	return operationErr
}

func (m *Memory) Entries() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.entries...)
}
