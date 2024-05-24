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

func ParseGroupTypeFilter(value string) (gf GroupTypeFilter, found bool, err error) {
	var gtf GroupTypeFilter

	me := NewMultiErr()
	values := strings.Split(value, ",")
	if len(values) == 0 {
		goto end
	}
	gtf = GroupTypeFilter{
		GroupTypes: make([]GroupType, 0, len(values)),
	}
	for _, groupType := range values {
		gt, err := ParseGroupTypeByLetter(groupType)
		if err != nil {
			me.Add(err)
			continue
		}
		found = true
		gtf.GroupTypes = append(gtf.GroupTypes, gt)
	}
	err = me.Err()
end:
	return gtf, found, err
}
