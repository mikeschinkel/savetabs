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
	GetAlertHTML(Context, AlertType, Message) (safehtml.HTML, int, error)
	GetOOBAlertHTML(Context, AlertType, Message) (safehtml.HTML, int, error)
	GetErrorHTML(Context, error) (safehtml.HTML, int, error)
	GetMenuHTML(Context, string) (safehtml.HTML, int, error)
	GetLinkSetHTML(ctx Context, host, requestURI string, g FilterGetter) (safehtml.HTML, int, error)
	GetMenuItemHTML(Context, string, string) (safehtml.HTML, int, error)
}
