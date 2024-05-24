package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Filter struct {
	Type   FilterType
	Values []string
	Map    map[string]string
}

var filterTypeMap = make(map[string]*FilterType)

type FilterType struct {
	value  string
	Label  string
	Plural string
}

// FilterTypes is a convenience array to allow processing filter types.
// Note that it does not include `NoFilter`
var FilterTypes = make([]*FilterType, 0, 10)

var filterMutex sync.Mutex

func newFilterType(value, label, plural string) *FilterType {
	ft := &FilterType{
		value:  strings.ToLower(value),
		Label:  label,
		Plural: plural,
	}
	filterMutex.Lock()
	defer filterMutex.Unlock()
	dupFilterTypeCheck(ft)
	filterTypeMap[ft.value] = ft
	FilterTypes = append(FilterTypes, ft)
	return ft
}

func FilterTypeByValue(value string) (_ *FilterType, err error) {
	ft, ok := filterTypeMap[strings.ToLower(value)]
	if !ok {
		err = errors.Join(
			ErrInvalidFilterType,
			fmt.Errorf("filter_type=%s, valid_types=[%s]",
				ft,
				strings.Join(FilterTypeNames(), ","),
			),
		)
	}
	return ft, err
}

func FilterTypeNames() []string {
	return ConvertSlice(FilterTypes, func(ft *FilterType) string {
		return ft.value
	})
}

func (f FilterType) String() string {
	return f.value
}

var (
	GroupTypeFilterType = newFilterType("gt", "Group Type", "Group Types")
	GroupFilterType     = newFilterType("grp", "Group", "Groups")
	MetaFilterType      = newFilterType("m", "Meta Pair", "Meta Pairs")
)
