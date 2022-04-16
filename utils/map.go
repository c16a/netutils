package utils

func ToSingleValuedMap[K comparable, V any](input map[K][]V) map[K]V {
	out := make(map[K]V, len(input))
	for k, v := range input {
		out[k] = v[0]
	}
	return out
}

func ToMultiValuedMap[K comparable, V any](input map[K]V) map[K][]V {
	out := make(map[K][]V, len(input))
	for k, v := range input {
		out[k] = []V{v}
	}
	return out
}
