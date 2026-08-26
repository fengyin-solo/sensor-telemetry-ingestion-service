package audit

import "sync"

type Log struct {
	mu      sync.Mutex
	success []string
}

func (l *Log) RecordSuccess(eventID string) {
	l.mu.Lock()
	l.success = append(l.success, eventID)
	l.mu.Unlock()
}

func (l *Log) Successes() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.success...)
}
