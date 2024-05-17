package shared

import (
	"strings"
	"sync"
)

type Filter struct {
	Type   FilterType
	Values []string
}

func NewFilter(ft FilterType, values []string) Filter {
	return Filter{
		Type:   ft,
		Values: values,
	}
}

type FilterQuery struct {
	FilterLabel Label
	Filters     FilterMap
	FilterTypes []FilterType
}

type FilterMap map[FilterType]Filter

type FilterType struct {
	value string
}

// FilterTypes is a convenience array to allow processing filter types.
// Note that it does not include `NoFilter`
var FilterTypes = make([]FilterType, 0, 10)

var filterMutex sync.Mutex

func newFilterType(value string) FilterType {
	ft := FilterType{value: strings.ToUpper(value)}
	filterMutex.Lock()
	defer filterMutex.Unlock()
	dupFilterCheck(ft)
	FilterTypes = append(FilterTypes, ft)
	return ft
}

func (f FilterType) String() string {
	return f.value
}

var (
	GroupTypeFilter = newFilterType("GT")
	MetaFilter      = newFilterType("M")
	BookmarkFilter  = newFilterType(GroupTypeBookmark.Upper())
	TabGroupFilter  = newFilterType(GroupTypeTabGroup.Upper())
	TagFilter       = newFilterType(GroupTypeTag.Upper())
	CategoryFilter  = newFilterType(GroupTypeCategory.Upper())
	KeywordFilter   = newFilterType(GroupTypeKeyword.Upper())
	InvalidFilter   = newFilterType(GroupTypeInvalid.Upper())
	NoFilter        = newFilterType("_")
)
