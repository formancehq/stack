package collectionutils

import (
	"reflect"
)

func Map[FROM any, TO any](input []FROM, mapper func(FROM) TO) []TO {
	ret := make([]TO, len(input))
	for i, input := range input {
		ret[i] = mapper(input)
	}
	return ret
}

func Filter[TYPE any](input []TYPE, filter func(TYPE) bool) []TYPE {
	ret := make([]TYPE, 0)
	for _, i := range input {
		if filter(i) {
			ret = append(ret, i)
		}
	}
	return ret
}

func Eq[T any](t T) func(T) bool {
	return func(t2 T) bool {
		return reflect.DeepEqual(t, t2)
	}
}

func Not[T any](t func(T) bool) func(T) bool {
	return func(t2 T) bool {
		return !t(t2)
	}
}
