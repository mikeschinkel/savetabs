package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Filter struct {
	Type FilterType
	Ids  []string
	Map  map[string]string
}

var filterTypeMap = make(map[string]*FilterType)

type FilterType struct {
	id      string
	Label   string
	Plural  string
	NewItem func(any) FilterItem
}

// FilterTypes is a convenience array to allow processing filter types.
// Note that it does not include `NoFilter`
var FilterTypes = make([]*FilterType, 0, 10)

var filterMutex sync.Mutex

func newFilterType(value, label, plural string, new func(any) FilterItem) *FilterType {
	ft := &FilterType{
		id:      strings.ToLower(value),
		Label:   label,
		Plural:  plural,
		NewItem: new,
	}
	filterMutex.Lock()
	defer filterMutex.Unlock()
	//dupFilterTypeCheck(ft)  // TODO: Search in git history for this func
	filterTypeMap[ft.id] = ft
	FilterTypes = append(FilterTypes, ft)
	return ft
}

func FilterTypeById(id string) (_ *FilterType, err error) {
	ft, ok := filterTypeMap[strings.ToLower(id)]
	if !ok {
		err = errors.Join(
			ErrInvalidFilterType,
			fmt.Errorf("filter_type=%s, valid_types=[%s]",
				ft,
				strings.Join(FilterTypeIds(), ","),
			),
		)
	}
	return ft, err
}

func FilterTypeIds() []string {
	return ConvertSlice(FilterTypes, func(ft *FilterType) string {
		return ft.id
	})
}

func (f FilterType) Id() string {
	return f.id
}
func (f FilterType) String() string {
	return f.id
}

var (
	GroupTypeFilterType = newFilterType("gt", "Group Type", "Group Types", newGroupTypeFilter)
	GroupFilterType     = newFilterType("grp", "Group", "Groups", newGroupFilter)
	MetaFilterType      = newFilterType("m", "Meta Pair", "Meta Pairs", newMetaFilter)
)
