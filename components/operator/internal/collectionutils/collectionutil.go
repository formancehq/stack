package collectionutils

// +kubebuilder:object:generate=false
type Array[T any] []T

func (f Array[T]) Filter(checkFn func(t T) bool) []T {
	ret := make([]T, 0)
	for _, item := range f {
		if checkFn(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

func (f Array[T]) First(checkFn func(t T) bool) *T {
	ret := f.Filter(checkFn)
	if len(ret) > 0 {
		return &ret[0]
	}
	return nil
}

func (f Array[T]) Append(t ...T) Array[T] {
	return append(f, t...)
}

func (f Array[T]) AppendIf(condition bool, fn func() []T) Array[T] {
	if !condition {
		return f
	}
	return f.Append(fn()...)
}

func NewArray[T any]() Array[T] {
	return Array[T]{}
}

func Filter[T any](items []T, checkFn func(t T) bool) []T {
	array := Array[T](items)
	return array.Filter(checkFn)
}

func First[T any](items []T, checkFn func(t T) bool) *T {
	array := Array[T](items)
	return array.First(checkFn)
}

func Equal[T comparable](value T) func(t T) bool {
	return func(t T) bool {
		return t == value
	}
}

func Map[T1 any, T2 any](v1 []T1, transformer func(T1) T2) []T2 {
	ret := make([]T2, 0)
	for _, v := range v1 {
		ret = append(ret, transformer(v))
	}
	return ret
}

func SliceFromMap[K comparable, V any](m map[K]V) []V {
	ret := make([]V, 0)
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}
