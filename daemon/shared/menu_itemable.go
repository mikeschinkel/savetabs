package shared

import (
	"github.com/google/safehtml"
)

type Menu interface {
	HTMLId() safehtml.Identifier
	MenuType() *MenuType
	ItemURL() safehtml.URL
	SubmenuURL() safehtml.URL
	Level() int
	Parent() Menu
}
