package shared

import (
	"fmt"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func Ptr[T any](a T) *T {
	return &a
}

func EnsureNotNil[T any](v *T, def T) T {
	if v == nil {
		v = &def
	}
	return *v
}

func ToUpperSlice(s []string) []string {
	for i := range s {
		s[i] = strings.ToUpper(s[i])
	}
	return s
}

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func ConvertSlice[SF []F, ST []T, F any, T any](from SF, fn func(F) T) ST {
	items := make(ST, len(from))
	for i, item := range from {
		items[i] = fn(item)
	}
	return items
}

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
