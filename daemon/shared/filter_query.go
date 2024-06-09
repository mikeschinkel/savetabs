package shared

import (
	"fmt"
	"net/url"
	"strings"
)

type FilterQuery struct {
	ParentDBId  int64
	FilterItems []FilterItem
}

func (fq FilterQuery) String() (q string) {
	sb := strings.Builder{}
	sb.WriteByte('?')
	for _, item := range fq.FilterItems {
		sb.WriteString(item.String())
		sb.WriteByte('&')
	}
	return sb.String()[:sb.Len()-1]
}

func (fq FilterQuery) Label() string {
	sb := strings.Builder{}
	for _, ft := range fq.FilterItems {
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

func (fq FilterQuery) FilterItemByType(ft *FilterType) (fi FilterItem, index int) {
	for index, fi = range fq.FilterItems {
		if fi.FilterType() != ft {
			continue
		}
		break
	}
	return fi, index
}

func NewFilter(ft *FilterType, args any) (fi FilterItem) {
	return ft.NewItem(args)
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
	for ftId, ftValues := range urlValues {
		ft, err := FilterTypeById(ftId)
		if err != nil {
			me.Add(err)
			continue
		}
		switch ft {
		case GroupTypeFilterType:
			fn := func(v string) (FilterItem, bool, error) {
				return ParseGroupTypeFilter(v)
			}
			if !parse(ftValues, fn) {
				continue
			}
		case GroupFilterType:
			fn := func(v string) (FilterItem, bool, error) {
				return ParseGroupFilter(v)
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

func toQueryString[S []T, T fmt.Stringer](typeId string, items S) (q string) {
	var sb strings.Builder

	if len(items) == 0 {
		goto end
	}
	sb.WriteString(typeId)
	sb.WriteByte('=')
	for _, m := range items {
		sb.WriteString(m.String())
		sb.WriteByte(',')
	}
	q = sb.String()
	q = q[:len(q)-1]
end:
	return q
}
