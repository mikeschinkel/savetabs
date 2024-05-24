package shared

import (
	"fmt"
	"net/url"
	"strings"
)

type FilterQuery struct {
	FilterItems []FilterItem
}

func (items FilterQuery) Label() string {
	sb := strings.Builder{}
	for _, ft := range items.FilterItems {
		sb.WriteString(fmt.Sprintf("%s: %s && ",
			ft.FilterType().Plural,
			ft.Label(),
		))
	}
	if sb.Len() == 0 {
		return ""
	}
	labels := sb.String()
	// Strip off trailing ' && ' with -4
	return labels[:len(labels)-4]
}

func ParseFilterQuery(u string) (items FilterQuery, err error) {
	var urlValues url.Values
	var me *MultiErr

	type parseOne func(string) (FilterItem, bool, error)
	var parse func([]string, parseOne) bool

	q, err := url.Parse(u)
	if err != nil {
		goto end
	}
	me = NewMultiErr()

	urlValues = q.Query()
	items.FilterItems = make([]FilterItem, 0)
	parse = func(vv []string, fn parseOne) (ok bool) {
		ok = true
		for _, value := range vv {
			item, found, err := fn(value)
			if err != nil {
				me.Add(err)
				ok = false
			}
			if item == nil {
				continue
			}
			if !found {
				continue
			}
			items.FilterItems = append(items.FilterItems, item)
		}
		return ok
	}
	for key, ftValues := range urlValues {
		ft, err := FilterTypeByValue(key)
		if err != nil {
			me.Add(err)
			continue
		}
		switch ft {
		case GroupTypeFilterType:
			fn := func(v string) (FilterItem, bool, error) {
				return ParseGroupFilter(v)
			}
			if !parse(ftValues, fn) {
				continue
			}
		case GroupFilterType:
			fn := func(v string) (FilterItem, bool, error) {
				return ParseGroupTypeFilter(v)
			}
			if !parse(ftValues, fn) {
				continue
			}
		case MetaFilterType:
			fn := func(v string) (FilterItem, bool, error) {
				return ParseMetaFilter(v)
			}
			if !parse(ftValues, fn) {
				continue
			}
		}
	}
	err = me.Err()
end:
	return items, err
}
