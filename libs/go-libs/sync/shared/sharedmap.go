package shared

type SharedMap[T any] struct {
	Shared[map[string]*Shared[T]]
}

func (s *SharedMap[T]) Add(index string, el *Shared[T]) {
	s.WLock()
	defer s.WUnlock()
	if _, ok := (*s.Val)[index]; !ok {
		(*s.Val)[index] = el
	}
}

func (s *SharedMap[T]) Remove(index string) {
	s.WLock()
	defer s.WUnlock()

	if _, ok := (*s.Val)[index]; ok {
		delete((*s.Val), index)

	}
}

func (s *SharedMap[T]) Get(index string) *Shared[T] {
	s.RLock()
	defer s.RUnlock()

	if sharedT, ok := (*s.Val)[index]; ok {
		return sharedT
	}

	return nil

}

func NewSharedMap[T any]() SharedMap[T] {
	m := make(map[string]*Shared[T])
	return SharedMap[T]{
		NewShared(&m),
	}
}
