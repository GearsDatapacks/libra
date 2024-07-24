package printer

import (
	"cmp"
	"slices"
)

type KeyValue[K, V any] struct {
	Key   K
	Value V
}

func SortMap[K cmp.Ordered, V any](m map[K]V) []KeyValue[K, V] {
	// We do this to ensure consistent order for our tests
	values := make([]KeyValue[K, V], 0, len(m))
	for key, value := range m {
		values = append(values, KeyValue[K, V]{Key: key, Value: value})
	}
	slices.SortFunc(values, func(a, b KeyValue[K, V]) int {
		return cmp.Compare(a.Key, b.Key)
	})
	return values
}
