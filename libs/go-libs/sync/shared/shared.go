package shared

import (
	"sync"
)

type Shared[T any] struct {
	Val *T
	mu  sync.RWMutex
}

func (s *Shared[T]) RLock() {
	s.mu.RLock()
}

func (s *Shared[T]) RUnlock() {
	s.mu.RUnlock()
}

func (s *Shared[T]) WLock() {
	s.mu.Lock()
}

func (s *Shared[T]) WUnlock() {
	s.mu.Unlock()
}

func (s *Shared[T]) RunWLock(f func()) {
	s.WLock()
	defer s.WUnlock()
	f()

}

func NewShared[T any](el *T) Shared[T] {
	return Shared[T]{
		Val: el,
	}
}
