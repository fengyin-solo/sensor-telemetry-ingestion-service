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

func (s *Store) Update(record Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.current[record.ID]; ok && record.Version <= existing.Version {
		return
	}
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
