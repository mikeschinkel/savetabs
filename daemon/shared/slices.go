package shared

// MapSliceFunc maps a slice (where `map` is in map-reduce form) using a func that
// returns a mapped slice element.
func MapSliceFunc[S []T, T, O any](slice S, fn func(T) O) []O {
	yy := make([]O, len(slice))
	for i, x := range slice {
		yy[i] = fn(x)
	}
	return yy
}

// MapMapFunc maps a map (where `maps` is in map-reduce form) using a func that
// returns a mapped map element.
func MapMapFunc[M map[string]T, T, O any](m M, fn func(T) O) map[string]O {
	yy := make(map[string]O, len(m))
	for i, x := range m {
		yy[i] = fn(x)
	}
	return yy
}
