package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var _ FilterItem = (*GroupFilter)(nil)

type GroupFilter struct {
	Groups []FilterGroup
}

func (g GroupFilter) FilterType() *FilterType {
	return GroupFilterType
}

func (g GroupFilter) Label() string {
	return strings.Join(ConvertSlice(g.Groups, func(fg FilterGroup) string {
		return fg.GroupName
	}), ", ")
}

func (g GroupFilter) Filters() []any {
	return ConvertSlice(g.Groups, func(fg FilterGroup) any {
		return fg.String()
	})
}

var groupRegexp = regexp.MustCompile(`^([a-z]+):([a-z0-9-]+)$`)
var groupsRegexp = regexp.MustCompile(`^(([a-z]+):([a-z0-9-]+),?)+$`)

func ParseGroupFilter(value string) (gf GroupFilter, err error) {
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
		gf.Groups = append(gf.Groups, fg)
	}
	err = me.Err()
end:
	return gf, err
}
