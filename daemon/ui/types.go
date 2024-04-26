package ui

import (
	"context"

	"github.com/google/safehtml"
)

type Context = context.Context

type MenuItemable interface {
	HTMLId() safehtml.Identifier
	MenuItemType() safehtml.Identifier
	LinksQueryParams() string
}

type FilterValueGetter interface {
	GetFilterLabel(string, string) string
	GetFilterValues(string) []string
	RawQuery() string
}

type Viewer interface {
	GetAlertHTML(Context, AlertType, string) ([]byte, int, error)
	GetErrorHTML(Context, error) ([]byte, int, error)
	GetMenuHTML(Context, string) ([]byte, int, error)
	GetLinkSetHTML(Context, string, FilterValueGetter, string) ([]byte, int, error)
	GetMenuItemHTML(Context, string, string) ([]byte, int, error)
}
