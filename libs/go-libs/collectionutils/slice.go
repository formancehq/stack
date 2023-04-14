package collectionutils

func Map[FROM any, TO any](input []FROM, mapper func(FROM) TO) []TO {
	ret := make([]TO, len(input))
	for i, input := range input {
		ret[i] = mapper(input)
	}
	return ret
}
