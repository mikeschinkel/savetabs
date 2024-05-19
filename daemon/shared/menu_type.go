package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type MenuType struct {
	Parent *MenuType
	name   string
	child  string
	Param  string
}

func (mt MenuType) Id() (s string) {
	switch mt.Parent {
	case nil:
		s = "."
	case &RootMenuType:
		s = strings.ToLower(mt.name)
	default:
		s = fmt.Sprintf("%s-%s", mt.Parent.Id(), strings.ToLower(mt.name))
	}
	return s
}

func (mt MenuType) Level() (n int) {
	p := &mt
	for p.Parent != nil {
		p = p.Parent
		n++
	}
	return n
}

func (mt MenuType) Slice() []string {
	if mt.Parent == nil {
		return []string{}
	}
	return append(mt.Parent.Slice(), mt.name)
}

// Leaf returns the last substring separated by dashes ("-")
func (mt MenuType) Leaf() string {
	if mt.Parent == nil {
		Panicf("Subtype called for a top level menu type: %v", mt)
	}
	return mt.name
}

func (mt MenuType) Name() string {
	return mt.name
}

var menuTypeMap = make(map[string]MenuType, 0)
var MenuTypes = []MenuType{RootMenuType}
var menuTypeMutex sync.Mutex

func newStaticMenuType(parent *MenuType, name, child fmt.Stringer) MenuType {
	mt := NewMenuType(parent, name, child)
	menuTypeMutex.Lock()
	defer menuTypeMutex.Unlock()
	MenuTypes = append(MenuTypes, mt)
	menuTypeMap[mt.Id()] = mt
	return mt
}

func (mt MenuType) Params() string {
	s := mt.params([]string{})
	kvs := make([]string, len(s)/2)
	n := 0
	for i := 0; i < len(kvs)*2; i += 2 {
		kvs[n] = fmt.Sprintf("%s=%s", s[i], s[i+1])
		n++
	}
	return strings.Join(kvs, "&")
}

func (mt MenuType) params(s []string) []string {
	if mt.Parent != nil {
		s = mt.Parent.params(s)
	}
	if mt.Parent == nil {
		goto end
	}
	s = append(s, mt.name)
	if mt.name != mt.Param {
		s = append(s, mt.Param)
	}

end:
	return s
}

func NewMenuType(parent *MenuType, name, child fmt.Stringer) MenuType {
	mt := MenuType{
		Parent: parent,
		name:   strings.ToLower(name.String()),
	}
	if child != nil {
		mt.child = strings.ToLower(child.String())
	}
	if parent != nil {
		mt.Param = parent.child
	}
	if mt.Param == "" {
		mt.Param = mt.name
	}
	return mt
}

func MenuTypeByParentTypeAndMenuName(parent *MenuType, name string) (mt MenuType, err error) {
	if nil == parent {
		err = errors.Join(ErrMenuTypeIsNil, fmt.Errorf("child_name=%s", name))
		goto end
	}
	mt, err = MenuTypeByValue(fmt.Sprintf("%s-%s", parent.Id(), name))
end:
	return mt, err
}

func MenuTypeByValue(value string) (mt MenuType, err error) {
	mt, found := menuTypeMap[strings.ToLower(value)]
	if !found {
		err = errors.Join(ErrMenuTypeNotFound, fmt.Errorf("value=%s", value))
	}
	return mt, err
}

//goland:noinspection GoUnusedGlobalVariable
var (
	RootMenuType      MenuType
	GroupTypeMenuType = newStaticMenuType(&RootMenuType, String{"gt"}, String{"grp"})

	BookmarkMenuType = newStaticMenuType(&GroupTypeMenuType, GroupTypeBookmark, nil)
	TabGroupMenuType = newStaticMenuType(&GroupTypeMenuType, GroupTypeTabGroup, nil)
	TagMenuType      = newStaticMenuType(&GroupTypeMenuType, GroupTypeTag, nil)
	CategoryMenuType = newStaticMenuType(&GroupTypeMenuType, GroupTypeCategory, nil)
	KeywordMenuType  = newStaticMenuType(&GroupTypeMenuType, GroupTypeKeyword, nil)
	InvalidMenuType  = newStaticMenuType(&GroupTypeMenuType, GroupTypeInvalid, nil)
)

var menuTypesForRegexp string

// MenuTypesForRegexp contains a string for matching via regexp in the format
// `gt|gt-g|gt-c|gt-t|etc.`
func MenuTypesForRegexp() (s string) {
	var mts []string
	if menuTypesForRegexp == "" {
		menuTypeMutex.Lock()
		defer menuTypeMutex.Unlock()
		mts = ConvertSlice(MenuTypes, func(mt MenuType) string {
			return mt.Id()
		})
		s = strings.Join(mts, "|")
	}
	return s
}
