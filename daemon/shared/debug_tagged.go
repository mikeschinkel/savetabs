//go:build debug

package shared

var filterMap = make(map[FilterType]struct{})

var dupFilterCheck = func(ft FilterType) {
	_, ok := filterMap[ft]
	if ok {
		Panicf("ERROR: Duplicate Filter Type declared: '%s' (Decalared filter types: %v)",
			ft,
			FilterTypes,
		)
	}
	filterMap[ft] = struct{}{}
}
