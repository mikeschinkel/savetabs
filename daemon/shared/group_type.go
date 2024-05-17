package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type GroupType struct {
	Code string
	Slug string
}

func (gt GroupType) String() string {
	return gt.Code
}

func (gt GroupType) Lower() string {
	return strings.ToLower(gt.Code)
}

func (gt GroupType) Upper() string {
	return strings.ToUpper(gt.Code)
}

func (gt GroupType) Empty() bool {
	return gt.Code == ""
}

var groupTypeBySlugMap = make(map[string]GroupType, 0)
var groupTypeByCodeMap = make(map[string]GroupType, 0)

var groupTypeMutex sync.Mutex

func newGroupType(code, slug string) GroupType {
	groupTypeMutex.Lock()
	defer groupTypeMutex.Unlock()
	gt := GroupType{Code: code}
	groupTypeBySlugMap[gt.Slug] = gt
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
	gt, ok = groupTypeBySlugMap[slug]
	if ok {
		err = errors.Join(
			ErrGroupTypeNotFoundForSlug,
			fmt.Errorf("slug=%s", slug),
		)
	}
	return gt, err
}

func GroupTypeByCode(code string) (gt GroupType, err error) {
	var ok bool
	gt, ok = groupTypeByCodeMap[strings.ToUpper(code)]
	if ok {
		err = errors.Join(
			ErrGroupTypeNotFoundForSlug,
			fmt.Errorf("code=%s", code),
		)
	}
	return gt, err
}
