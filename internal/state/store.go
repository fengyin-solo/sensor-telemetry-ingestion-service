package state

import "sync"

type Record struct {
	ID      string
	Version int
	Status  string
}

type Store struct {
	mu      sync.Mutex
	current map[string]Record
	cache   map[string]Record
}

func NewStore() *Store {
	return &Store{current: make(map[string]Record), cache: make(map[string]Record)}
}

// Update applies record as the new state for its ID, but only if it is not
// older than what is already stored. A late callback from an earlier attempt
// carries a stale version and must not overwrite a fresher result, so the
// update is dropped when its version is behind the current one.
func (s *Store) Update(record Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.current[record.ID]; ok && record.Version < existing.Version {
		return
	}
	s.current[record.ID] = record
	s.cache[record.ID] = record
}

// Replace unconditionally writes record, used to reset to a known baseline.
func (s *Store) Replace(record Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current[record.ID] = record
	s.cache[record.ID] = record
}

func (s *Store) Get(id string) Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current[id]
}

func (s *Store) Cached(id string) Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cache[id]
}
