package ui

import (
	"context"

	"github.com/google/safehtml"
)

type Context = context.Context

type MenuItemable interface {
	Identifier() safehtml.Identifier
	MenuItemType() safehtml.Identifier
	LinksQueryParams() string
}

type FilterValueGetter interface {
	GetFilterValues(string) []string
}

type Viewer interface {
	GetMenuHTML(Context, string) ([]byte, int, error)
	GetLinksHTML(Context, string, FilterValueGetter) ([]byte, int, error)
	GetMenuItemHTML(Context, string, string) ([]byte, int, error)
}
