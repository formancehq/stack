package shared

type SharedMapArr[T any] struct {
	Shared[map[string]*SharedArr[T]]
}

func (s *SharedMapArr[T]) Add(index string, el *Shared[T]) {
	s.WLock()
	defer s.WUnlock()

	if sharedArr, ok := (*s.Val)[index]; ok {
		sharedArr.Add(el)
	} else {
		newSharedArr := NewSharedArr[T]()
		newSharedArr.Add(el)
		(*s.Val)[index] = &newSharedArr
	}
}

func (s *SharedMapArr[T]) Adds(idxs []string, el *Shared[T]) {
	for _, str := range idxs {
		s.Add(str, el)
	}
}

func (s *SharedMapArr[T]) Remove(index string, el *Shared[T]) {
	s.WLock()
	defer s.WUnlock()

	if sharedArr, ok := (*s.Val)[index]; ok {
		ex := sharedArr.Remove(el)
		if ex != nil && sharedArr.Size() == 0 {
			delete((*s.Val), index)
		}
	}
}

func (s *SharedMapArr[T]) Removes(idxs []string, el *Shared[T]) {
	for _, str := range idxs {
		s.Remove(str, el)
	}
}

func (s *SharedMapArr[T]) Get(index string) *SharedArr[T] {
	s.RLock()
	defer s.RUnlock()
	if sharedArr, ok := (*s.Val)[index]; ok {
		return sharedArr
	}

	return nil

}

func NewSharedMapArr[T any]() SharedMapArr[T] {
	m := make(map[string]*SharedArr[T])
	return SharedMapArr[T]{
		NewShared(&m),
	}
}
