package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var _ FilterItem = (*GroupFilter)(nil)

type GroupFilter struct {
	dbId         int64
	FilterGroups []FilterGroup
}

func (gf GroupFilter) DBId() int64 {
	return gf.dbId
}

func (gf GroupFilter) WithDBId(id int64) GroupFilter {
	// Damn that needless "ineffective assignment" warning by GoLand
	gf.dbId = id
	return gf
}

func (gf GroupFilter) FilterGroupByType(groupType GroupType) (fg *FilterGroup) {
	for _, g := range gf.FilterGroups {
		if g.GroupType != groupType {
			continue
		}
		fg = &g
		break
	}
	return fg
}

func (gf GroupFilter) HTMLId(mi MenuItemable) string {
	mt := mi.MenuType()
	id := fmt.Sprintf("%s-%s-%s",
		mt.FilterType.Id(),
		mt.Name(),
		mi.LocalId(),
	)
	return Slugify(id)
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

func newGroupFilter(args any) FilterItem {
	if args == nil {
		panic("No GroupId not provided for GroupFilterType")
	}
	var groupId int64
	var _groupId int
	var ok bool
	groupId, ok = args.(int64)
	if !ok {
		// If a literal 0 is passed it will be `int`, not `int64` so we need to capa
		_groupId, ok = args.(int)
		groupId = int64(_groupId)
	}
	if !ok {
		Panicf("Invalid type for GroupId provided for GroupFilterType; %#v provided instead", args)
	}
	return GroupFilter{
		dbId:         groupId,
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
