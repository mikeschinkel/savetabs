package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var _ FilterItem = (*GroupFilter)(nil)

type GroupFilter struct {
	FilterGroups []FilterGroup
}

func (gf GroupFilter) HTMLId(mi MenuItemable) string {
	mt := mi.MenuType()
	id := fmt.Sprintf("%s-%s-%s",
		mt.FilterType.Id(),
		mt.Name(),
		mi.LocalId(),
	)
	return id
}

func (gf GroupFilter) ContentQuery(mi MenuItemable) (u string) {
	mt := mi.MenuType()
	u = fmt.Sprintf("%s=%s:%s",
		mt.FilterType.Id(),
		mt.Name(),
		mi.LocalId(),
	)
	return u
}

func newGroupFilter() GroupFilter {
	return GroupFilter{
		FilterGroups: make([]FilterGroup, 0),
	}
}

func (gf GroupFilter) String() string {
	return toQueryString(GroupFilterType.Id(), gf.FilterGroups)
}

func (g GroupFilter) FilterType() *FilterType {
	return GroupFilterType
}

func (g GroupFilter) Label() string {
	return strings.Join(ConvertSlice(g.FilterGroups, func(fg FilterGroup) string {
		return fg.GroupName
	}), ", ")
}

func (g GroupFilter) Filters() []any {
	return ConvertSlice(g.FilterGroups, func(fg FilterGroup) any {
		return fg.String()
	})
}

var groupRegexp = regexp.MustCompile(`^([a-z]+):([a-z0-9-]+)$`)
var groupsRegexp = regexp.MustCompile(`^(([a-z]+):([a-z0-9-]+),?)+$`)

func ParseGroupFilter(value string) (gf GroupFilter, found bool, err error) {
	var fg FilterGroup

	me := NewMultiErr()

	matches := groupsRegexp.FindAllStringSubmatch(value, -1)
	if matches == nil {
		err = errors.Join(
			ErrInvalidGroupFilterFormat,
			fmt.Errorf("filter_values=%s, format_expected='<group_type>:<group_name>'", value),
		)
		goto end
	}
	for _, match := range matches {
		fg, err = ParseFilterGroup(match[0])
		if err != nil {
			me.Add(err)
		}
		gf.FilterGroups = append(gf.FilterGroups, fg)
	}
	err = me.Err()
	found = true
end:
	return gf, found, err
}
