package shared

import (
	"fmt"
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

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
