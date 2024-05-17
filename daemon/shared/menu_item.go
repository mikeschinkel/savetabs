package shared

import (
	"strings"
	"sync"
)

type MenuType struct {
	value string
}

func (mt MenuType) String() string {
	return strings.ToLower(mt.value)
}

var MenuTypes []MenuType
var menuTypeMutex sync.Mutex

func NewMenuType(value string) MenuType {
	menuTypeMutex.Lock()
	defer menuTypeMutex.Unlock()
	mt := MenuType{value: value}
	MenuTypes = append(MenuTypes, mt)
	return mt
}

var (
	GroupTypeMenuType = NewMenuType("gt")
	GroupMenuType     = NewMenuType("grp")
)

var menuTypesForRegexp string

// MenuTypesForRegexp contains a string for matching via regexp in the format
// `gt|grp`.
func MenuTypesForRegexp() (s string) {
	var mts []string
	if menuTypesForRegexp == "" {
		menuTypeMutex.Lock()
		defer menuTypeMutex.Unlock()
		mts = ConvertSlice(MenuTypes, func(mt MenuType) string {
			return mt.String()
		})
		s = strings.Join(mts, "|")
	}
	return s
}
