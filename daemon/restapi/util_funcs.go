package restapi

import (
	"context"
)

func newContext() context.Context {
	return context.Background()
}

//goland:noinspection GoUnusedFunction
func deleteElement[T any](slice []T, index int) []T {
	// Copy the elements following the index one position to the left.
	copy(slice[index:], slice[index+1:])
	// Return the slice without the last element.
	return slice[:len(slice)-1]
}
