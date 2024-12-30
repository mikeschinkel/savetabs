package shared

import (
	"strconv"
	"strings"
)

// MapMapFunc maps a map (where `maps` is in map-reduce form) using a func that
// returns a mapped map element.
func MapMapFunc[M map[string]T, T, O any](m M, fn func(T) O) map[string]O {
	yy := make(map[string]O, len(m))
	for i, x := range m {
		yy[i] = fn(x)
	}
	return yy
}

type Int64Slice []int64

func (ii Int64Slice) Join(separator string) string {
	return strings.Join(ii.Strings(), separator)
}

func (ii Int64Slice) Strings() []string {
	return ConvertSlice([]int64(ii), func(i int64) string {
		return strconv.FormatInt(i, 10)
	})
}

func InvertSlice[S []T, M map[T]struct{}, T comparable](s S) M {
	m := make(M, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func ToUpperSlice(s []string) []string {
	for i := range s {
		s[i] = strings.ToUpper(s[i])
	}
	return s
}

// ConvertSlice converts a slice of one type to a slice of another type
func ConvertSlice[SF []F, ST []T, F, T any](from SF, fn func(F) T) ST {
	items := make(ST, len(from))
	for i, item := range from {
		items[i] = fn(item)
	}
	return items
}

// ConvertSliceWithFilter converts a slice of one type to a slice of another type with a filter func
func ConvertSliceWithFilter[SF []F, ST []T, F any, T any](from SF, fn func(F) (T, bool)) ST {
	items := make(ST, len(from))
	for i, item := range from {
		newItem, ok := fn(item)
		if !ok {
			continue
		}
		items[i] = newItem
	}
	return items
}

func ConvertSliceToMapFunc[M map[K]T, K comparable, S []T, T any](s S, fn func(M, T)) M {
	m := make(M, len(s))
	for _, t := range s {
		fn(m, t)
	}
	return m
}
