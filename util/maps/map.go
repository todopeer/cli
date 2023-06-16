package maps

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func GetKeys[K comparable, V any](m map[K]V) []K {
	res := make([]K, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func SortedKByV[K comparable, V constraints.Ordered](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool {
		return m[res[i]] < m[res[j]]
	})
	return res
}
