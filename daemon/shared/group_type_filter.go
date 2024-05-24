package shared

import (
	"strings"
)

type GroupTypeFilter struct {
	GroupTypes []GroupType
}

var _ FilterItem = (*GroupFilter)(nil)

func (g GroupTypeFilter) FilterType() *FilterType {
	return GroupTypeFilterType
}

func (g GroupTypeFilter) Label() string {
	return strings.Join(ConvertSlice(g.GroupTypes, func(gt GroupType) string {
		return gt.Plural
	}), ", ")
}

func (g GroupTypeFilter) Filters() []any {
	return ConvertSlice(g.GroupTypes, func(gt GroupType) any {
		return gt.Lower()
	})
}

func ParseGroupTypeFilter(value string) (gf GroupTypeFilter, err error) {
	me := NewMultiErr()
	values := strings.Split(value, ",")
	gtf := GroupTypeFilter{
		GroupTypes: make([]GroupType, len(values)),
	}
	for i, groupType := range values {
		gt, err := ParseGroupTypeByLetter(groupType)
		if err != nil {
			me.Add(err)
			continue
		}
		gtf.GroupTypes[i] = gt
	}
	err = me.Err()
	return gtf, err
}
