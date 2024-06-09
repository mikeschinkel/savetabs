package shared

import (
	"errors"
	"fmt"
)

type FilterGroup struct {
	GroupType GroupType
	GroupName string
}

var filterGroupRegex = groupRegexp

func ParseFilterGroup(group string) (fg FilterGroup, err error) {
	match := filterGroupRegex.FindStringSubmatch(group)
	if match == nil {
		err = ErrInvalidGroupFilterFormat
		goto end
	}
	fg.GroupType, err = ParseGroupTypeByLetter(match[1])
	if err != nil {
		goto end
	}
	fg.GroupName = match[2]
end:
	if err != nil {
		err = errors.Join(err, fmt.Errorf("group=%s", group))
	}
	return fg, err
}

func (fg FilterGroup) String() string {
	return fg.Slug()
}

func (fg FilterGroup) Slug() string {
	return fmt.Sprintf("%s:%s",
		fg.GroupType.Lower(),
		fg.GroupName,
	)
}
