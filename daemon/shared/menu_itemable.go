package shared

import (
	"github.com/google/safehtml"
)

type MenuItemParent interface {
	APIURL() safehtml.URL
	HTMLId() safehtml.Identifier
	MenuType() *MenuType
	Level() int
}
type MenuItemable interface {
	HTMLId() safehtml.Identifier
	ContentQuery() safehtml.URL
	ChildMenuURL() safehtml.URL
	Level() int
	MenuType() *MenuType
	Parent() MenuItemParent
	LocalId() string
	IsLeaf() bool
}
