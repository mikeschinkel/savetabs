package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

const LeafMenuName = "leaf"

//goland:noinspection GoUnusedGlobalVariable
var (
	RootMenuType      = &MenuType{}
	GroupTypeMenuType = newStaticMenuType(RootMenuType, String{"gt"}, 2)

	BookmarkMenuType = newStaticMenuType(GroupTypeMenuType, GroupTypeBookmark, 0)
	TabGroupMenuType = newStaticMenuType(GroupTypeMenuType, GroupTypeTabGroup, 0)
	TagMenuType      = newStaticMenuType(GroupTypeMenuType, GroupTypeTag, 0)
	CategoryMenuType = newStaticMenuType(GroupTypeMenuType, GroupTypeCategory, 0)
	KeywordMenuType  = newStaticMenuType(GroupTypeMenuType, GroupTypeKeyword, 0)
	InvalidMenuType  = newStaticMenuType(GroupTypeMenuType, GroupTypeInvalid, 0)

	// LeafMenuArchetype should never be used directly but should be copies to make a LeadMenuType
	LeafMenuArchetype = newStaticMenuType(nil, String{LeafMenuName}, 0)
)

func CloneLeafMenuType() *MenuType {
	mt := &MenuType{}
	*mt = *LeafMenuArchetype
	return mt
}

type MenuType struct {
	Parent   *MenuType
	name     string
	IsLeaf   bool
	ParamCnt int
}

func (mt *MenuType) Id() (s string) {
	switch mt.Parent {
	case nil:
		s = "."
	case RootMenuType:
		s = strings.ToLower(mt.name)
	default:
		s = fmt.Sprintf("%s--%s", mt.Parent.Id(), strings.ToLower(mt.name))
	}
	return s
}

func (mt *MenuType) Level() (n int) {
	p := mt
	for p.Parent != nil {
		p = p.Parent
		n++
	}
	return n
}

func (mt *MenuType) Slice() []string {
	if mt.Parent == nil {
		return []string{}
	}
	return append(mt.Parent.Slice(), mt.name)
}

// Leaf returns the last substring separated by dashes ("-")
func (mt *MenuType) Leaf() string {
	if mt.Parent == nil {
		Panicf("Subtype called for a top level menu type: %v", mt)
	}
	return mt.name
}

func (mt *MenuType) SetName(name string) {
	mt.name = name
}

func (mt *MenuType) Name() string {
	return mt.name
}

var menuTypeMap = make(map[string]*MenuType, 0)
var MenuTypes []*MenuType
var menuTypeMutex sync.Mutex

func newStaticMenuType(parent *MenuType, name fmt.Stringer, paramCnt int) *MenuType {
	mt := NewMenuType(parent, name, paramCnt)
	menuTypeMutex.Lock()
	defer menuTypeMutex.Unlock()
	MenuTypes = append(MenuTypes, mt)
	menuTypeMap[mt.Id()] = mt
	return mt
}

type ParamsArgs struct {
	Equates  string
	Combines string
}

// Params provides string of parmeters with equates and joins.
// For URL query type equates +> '=' and combines => '&', e.g. x=1&y=2
// For URL path type equates => '--' and combines => '/', e.g. x--1/y--2
func (mt *MenuType) Params(args ParamsArgs) (p string) {
	s := mt.params([]string{})
	kvs := make([]string, len(s)/2)
	n := 0
	for i := 0; i < len(kvs)*2; i += 2 {
		kvs[n] = fmt.Sprintf("%s%s%s", s[i], args.Equates, s[i+1])
		n++
	}
	p = strings.Join(kvs, args.Combines)
	if len(s)%2 != 0 {
		p = fmt.Sprintf("%s:%s", p, s[n+1])
	}
	return p
}

func (mt *MenuType) params(s []string) []string {
	if mt.Parent != nil {
		s = mt.Parent.params(s)
	}
	if mt.Parent == nil {
		goto end
	}
	s = append(s, mt.name)

end:
	return s
}

func NewMenuType(parent *MenuType, name fmt.Stringer, paramCnt int) *MenuType {
	mt := &MenuType{
		Parent:   parent,
		name:     strings.ToLower(name.String()),
		IsLeaf:   name.String() == LeafMenuName,
		ParamCnt: paramCnt,
	}
	return mt
}

func MenuTypeByParentTypeAndMenuName(parent *MenuType, name string) (mt *MenuType, err error) {
	if nil == parent {
		err = errors.Join(ErrMenuTypeIsNil, fmt.Errorf("child_name=%s", name))
		goto end
	}
	mt, err = MenuTypeByValue(fmt.Sprintf("%s--%s", parent.Id(), name))
end:
	return mt, err
}

func MenuTypeByValue(value string) (mt *MenuType, err error) {
	mt, found := menuTypeMap[strings.ToLower(value)]
	if !found {
		err = errors.Join(ErrMenuTypeNotFound, fmt.Errorf("value=%s", value))
	}
	return mt, err
}

var menuTypesForRegexp string

// MenuTypesForRegexp contains a string for matching via regexp in the format
// `gt|gt-g|gt-c|gt-t|etc.`
func MenuTypesForRegexp() (s string) {
	var mts []string
	if menuTypesForRegexp == "" {
		menuTypeMutex.Lock()
		defer menuTypeMutex.Unlock()
		mts = ConvertSlice(MenuTypes, func(mt *MenuType) string {
			return mt.Id()
		})
		s = strings.Join(mts, "|")
	}
	return s
}
