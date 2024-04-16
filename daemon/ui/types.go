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
