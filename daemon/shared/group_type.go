package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type GroupType struct {
	Type string
	Slug string
}

func (gt GroupType) String() string {
	return gt.Type
}

func (gt GroupType) Lower() string {
	return strings.ToLower(gt.Type)
}

func (gt GroupType) Upper() string {
	return strings.ToUpper(gt.Type)
}

func (gt GroupType) Empty() bool {
	return gt.Type == ""
}

var groupTypeBySlugMap = make(map[string]GroupType, 0)
var groupTypeByCodeMap = make(map[string]GroupType, 0)

var groupTypeMutex sync.Mutex

func newGroupType(typ, slug string) GroupType {
	groupTypeMutex.Lock()
	defer groupTypeMutex.Unlock()
	gt := GroupType{
		Type: typ,
		Slug: slug,
	}
	groupTypeBySlugMap[gt.Slug] = gt
	groupTypeByCodeMap[gt.Type] = gt
	return gt
}

var (
	GroupTypeBookmark = newGroupType("B", "bookmark")
	GroupTypeTabGroup = newGroupType("G", "tabgroup")
	GroupTypeTag      = newGroupType("T", "tag")
	GroupTypeCategory = newGroupType("C", "category")
	GroupTypeKeyword  = newGroupType("K", "keyword")
	GroupTypeInvalid  = newGroupType("I", "invalid")
)

func GroupTypeBySlug(slug string) (gt GroupType, err error) {
	var ok bool
	gt, ok = groupTypeBySlugMap[strings.ToLower(slug)]
	if ok {
		err = errors.Join(
			ErrGroupTypeNotFoundForSlug,
			fmt.Errorf("slug=%s", slug),
		)
	}
	return gt, err
}

func GroupTypeByType(typ string) (gt GroupType, err error) {
	var ok bool
	gt, ok = groupTypeByCodeMap[strings.ToUpper(typ)]
	if !ok {
		err = errors.Join(
			ErrGroupTypeNotFoundForType,
			fmt.Errorf("type=%s", typ),
		)
	}
	return gt, err
}
