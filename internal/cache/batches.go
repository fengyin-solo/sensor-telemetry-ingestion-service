package cache

import "sync"

type Batches struct {
	mu     sync.Mutex
	values map[string][]byte
}

func NewBatches() *Batches { return &Batches{values: make(map[string][]byte)} }

func (b *Batches) Put(id string, payload []byte) {
	b.mu.Lock()
	b.values[id] = append([]byte(nil), payload...)
	b.mu.Unlock()
}

func (b *Batches) Get(id string) []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]byte(nil), b.values[id]...)
}
