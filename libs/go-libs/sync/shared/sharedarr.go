package shared

import (
	"sync"
)

type SharedArr[T any] struct {
	Shared[[]*Shared[T]]
}

func (s *SharedArr[T]) Add(new *Shared[T]) {
	s.WLock()
	defer s.WUnlock()

	s.UnsafeAdd(new)
}

func (s *SharedArr[T]) UnsafeAdd(new *Shared[T]) {
	idx := s.UnsafeFind(new)
	if idx < 0 {
		*s.Val = append(*s.Val, new)
	}
}

func (s *SharedArr[T]) Remove(new *Shared[T]) *Shared[T] {
	s.WLock()
	defer s.WUnlock()

	return s.UnsafeRemove(new)
}

func (s *SharedArr[T]) UnsafeRemove(new *Shared[T]) *Shared[T] {
	idx := s.UnsafeFind(new)

	if idx == -1 {
		return nil
	}

	if idx == 0 {
		*s.Val = (*s.Val)[1:]
	} else if idx == len(*s.Val) {
		*s.Val = (*s.Val)[:len(*s.Val)-1]
	} else {
		*s.Val = append((*s.Val)[:idx], (*s.Val)[idx+1:]...)
	}
	return new
}

func (s *SharedArr[T]) Find(el *Shared[T]) int {
	s.RLock()
	defer s.RUnlock()

	return s.UnsafeFind(el)
}

func (s *SharedArr[T]) UnsafeFind(el *Shared[T]) int {
	for idx, ptr := range *s.Val {
		if ptr == el {
			return idx
		}
	}
	return -1
}

func (s *SharedArr[T]) FindElement(f func(*Shared[T]) bool) *Shared[T] {
	s.RLock()
	defer s.RUnlock()

	return s.UnsafeFindElement(f)
}

func (s *SharedArr[T]) UnsafeFindElement(f func(*Shared[T]) bool) *Shared[T] {

	for _, shared := range *s.Val {
		if f(shared) {
			return shared
		}
	}

	return nil
}

func (s *SharedArr[T]) Apply(f func(*Shared[T])) {
	s.RLock()
	defer s.RUnlock()

	s.UnsafeApply(f)
}

func (s *SharedArr[T]) UnsafeApply(f func(*Shared[T])) {
	for _, shared := range *s.Val {
		shared.WLock()
		f(shared)
		shared.WUnlock()
	}
}

func (s *SharedArr[T]) AsyncApply(f func(_ *Shared[T], wg *sync.WaitGroup)) {
	s.RLock()
	defer s.RUnlock()

	s.UnsafeAsyncApply(f)

}

func (s *SharedArr[T]) UnsafeAsyncApply(f func(_ *Shared[T], wg *sync.WaitGroup)) {
	var wg sync.WaitGroup

	for _, shared := range *s.Val {
		s := shared

		wg.Add(1)
		go func() {
			s.WLock()
			f(s, &wg)
			s.WUnlock()
		}()
	}

	wg.Wait()
}

func (s *SharedArr[T]) Filter(f func(*Shared[T]) bool) *SharedArr[T] {

	s.RLock()
	defer s.RUnlock()

	return s.UnsafeFilter(f)

}

func (s *SharedArr[T]) UnsafeFilter(f func(*Shared[T]) bool) *SharedArr[T] {
	newSharedArr := NewSharedArr[T]()

	s.Apply(func(s *Shared[T]) {
		if f(s) {
			newSharedArr.Add(s)
		}
	})

	return &newSharedArr

}

func (s *SharedArr[T]) Size() int {
	s.RLock()
	defer s.RUnlock()

	return s.UnsafeSize()
}

func (s *SharedArr[T]) UnsafeSize() int {
	return len((*s.Val))
}

func (s *SharedArr[T]) Empty() SharedArr[T] {
	s.WLock()
	defer s.WUnlock()

	return s.UnsafeEmpty()
}

func (s *SharedArr[T]) UnsafeEmpty() SharedArr[T] {

	copy := (*s.Val)
	(*s.Val) = make([]*Shared[T], 0)
	return SharedArr[T]{
		NewShared(&copy),
	}
}

func (s *SharedArr[T]) Merge(s2 *SharedArr[T]) {
	s.WLock()
	defer s.WUnlock()

	s.UnsafeMerge(s2)
}

func (s *SharedArr[T]) UnsafeMerge(s2 *SharedArr[T]) {
	s2.RLock()
	defer s2.RUnlock()

	for _, el := range *s2.Val {
		s.Add(el)
	}
}

func (s *SharedArr[T]) Extract() *[]*T {
	arr := make([]*T, 0)
	s.Apply(func(s *Shared[T]) {
		arr = append(arr, s.Val)
	})
	return &arr
}

func (s *SharedArr[T]) ExtractCopy() *[]T {
	arr := make([]T, 0)
	s.Apply(func(s *Shared[T]) {
		arr = append(arr, *s.Val)
	})
	return &arr
}

func (s *SharedArr[T]) From(arr *[]*T) *SharedArr[T] {
	sharedArr := make([]*Shared[T], 0)
	for _, el := range *arr {
		shared := NewShared(el)
		sharedArr = append(sharedArr, &shared)
	}
	s.Val = &sharedArr
	return s
}

func NewSharedArr[T any]() SharedArr[T] {
	arr := make([]*Shared[T], 0)
	return SharedArr[T]{
		NewShared(&arr),
	}
}

func UnsafeExtract[T any](sharedArr *SharedArr[T]) []*T {

	arrT := make([]*T, 0)
	for _, el := range *sharedArr.Val {
		arrT = append(arrT, el.Val)
	}
	return arrT

}
