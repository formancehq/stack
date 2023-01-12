package fifo

import "sync"

type FIFO[ITEM any] struct {
	mu    sync.Mutex
	items []ITEM
}

func (s *FIFO[ITEM]) Pop() (ITEM, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		var i ITEM

		return i, false
	}

	ret := s.items[0]

	if len(s.items) == 1 {
		s.items = make([]ITEM, 0)

		return ret, true
	}

	s.items = s.items[1:]

	return ret, true
}

func (s *FIFO[ITEM]) Peek() (ITEM, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		var i ITEM

		return i, false
	}

	return s.items[0], true
}

func (s *FIFO[ITEM]) Push(i ITEM) *FIFO[ITEM] {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = append(s.items, i)

	return s
}

func (s *FIFO[ITEM]) Empty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.items) == 0
}
