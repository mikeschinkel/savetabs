package shared

import (
	"github.com/google/safehtml"
)

type MenuItemable interface {
	HTMLId() safehtml.Identifier
	MenuType() *MenuType
	ItemURL() safehtml.URL
	SubmenuURL() safehtml.URL
	Level() int
	Parent() MenuItemable
}
