package shared

//goland:noinspection GoUnusedExportedFunction
func Ptr[T any](a T) *T {
	return &a
}
