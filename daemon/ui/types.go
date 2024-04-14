package ui

import (
	"context"

	"github.com/google/safehtml"
)

type Context = context.Context

type MenuItemable interface {
	Identifier() safehtml.Identifier
}

type SlugsForGetter interface {
	GetSlugsFor(string) []string
}
