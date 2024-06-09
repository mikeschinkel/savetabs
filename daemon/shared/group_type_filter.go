package shared

import (
	"fmt"
	"strings"
)

type GroupTypeFilter struct {
	GroupTypes []GroupType
}

func newGroupTypeFilter(args any) FilterItem {
	return GroupTypeFilter{
		GroupTypes: make([]GroupType, 0),
	}
}

func (gtf GroupTypeFilter) String() string {
	return toQueryString(GroupTypeFilterType.Id(), gtf.GroupTypes)
}

var _ FilterItem = (*GroupTypeFilter)(nil)

func (gtf GroupTypeFilter) HTMLId(mi MenuItemable) string {
	mt := mi.Parent().MenuType()
	return fmt.Sprintf("%s-%s", mt.Name(), mi.LocalId())
}
func (gtf GroupTypeFilter) ContentQuery(mi MenuItemable) (u string) {
	mt := mi.Parent().MenuType()
	return fmt.Sprintf("%s=%s", mt.Name(), mi.LocalId())
}

func (gtf GroupTypeFilter) FilterType() *FilterType {
	return GroupTypeFilterType
}

func (gtf GroupTypeFilter) Label() string {
	return strings.Join(ConvertSlice(gtf.GroupTypes, func(gt GroupType) string {
		return gt.Plural
	}), ", ")
}

func (gtf GroupTypeFilter) Filters() []any {
	return ConvertSlice(gtf.GroupTypes, func(gt GroupType) any {
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
