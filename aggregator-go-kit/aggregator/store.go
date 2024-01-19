package aggregator

import (
	"context"
	"sync"
	"tolling/types"
)

type MemoryStorer struct {
	mx   sync.RWMutex
	data map[int]float64
}

func (s *MemoryStorer) Store(ctx context.Context, d types.Distance) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[d.OBUID] = d.Value

	return nil
}

func (s *MemoryStorer) Read(ctx context.Context, OBUID int) (float64, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	dist, ok := s.data[OBUID]
	return dist, ok
}

func newMemoryStorer() *MemoryStorer {
	return &MemoryStorer{data: make(map[int]float64), mx: sync.RWMutex{}}
}
