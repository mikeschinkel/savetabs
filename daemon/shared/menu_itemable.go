package shared

import (
	"github.com/google/safehtml"
)

type MenuItemable interface {
	HTMLId() safehtml.Identifier
	MenuType() *MenuType
	LinksQuery() safehtml.URL
	SubmenuURL() safehtml.URL
	Level() int
	Parent() MenuItemable
	IsLeaf() bool
}
