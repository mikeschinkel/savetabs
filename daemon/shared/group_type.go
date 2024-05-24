package shared

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

type GroupType struct {
	Letter string
	Slug   string
	Label  string
	Plural string
}

func (gt GroupType) String() string {
	return gt.Letter
}

func (gt GroupType) Lower() string {
	return strings.ToLower(gt.Letter)
}

func (gt GroupType) Upper() string {
	return strings.ToUpper(gt.Letter)
}

func (gt GroupType) Empty() bool {
	return gt.Letter == ""
}

var groupTypeBySlugMap = make(map[string]GroupType, 0)
var groupTypeByLetterMap = make(map[string]GroupType, 0)

var groupTypeMutex sync.Mutex

func newGroupType(ltr, slug, label, plural string) GroupType {
	groupTypeMutex.Lock()
	defer groupTypeMutex.Unlock()
	gt := GroupType{
		Letter: ltr,
		Slug:   slug,
		Label:  label,
		Plural: plural,
	}
	groupTypeBySlugMap[gt.Slug] = gt
	groupTypeByLetterMap[gt.Letter] = gt
	return gt
}

var (
	GroupTypeBookmark = newGroupType("B", "bookmark", "Bookmark", "Bookmarks")
	GroupTypeTabGroup = newGroupType("G", "tabgroup", "Tag Groups", "Tag Groups")
	GroupTypeTag      = newGroupType("T", "tag", "Tag", "Tags")
	GroupTypeCategory = newGroupType("C", "category", "Category", "Categories")
	GroupTypeKeyword  = newGroupType("K", "keyword", "Keyword", "Keywords")
	GroupTypeInvalid  = newGroupType("I", "invalid", "Invalid", "Invalids")
)

func ParseGroupTypeBySlug(slug string) (gt GroupType, err error) {
	var ok bool
	if strings.Contains(slug, "-") {
		slog.Warn("Group type slugs with dashes exist", "slug", slug)
		slug = strings.Replace(slug, "-", "", -1)
	}
	gt, ok = groupTypeBySlugMap[strings.ToLower(slug)]
	if !ok {
		err = errors.Join(
			ErrGroupTypeNotFoundForSlug,
			fmt.Errorf("slug=%s", slug),
		)
	}
	return gt, err
}

func ParseGroupTypeByLetter(ltr string) (gt GroupType, err error) {
	var ok bool
	gt, ok = groupTypeByLetterMap[strings.ToUpper(ltr)]
	if !ok {
		err = errors.Join(
			ErrGroupTypeNotFoundForLetter,
			fmt.Errorf("letter=%s", ltr),
		)
	}
	return gt, err
}
