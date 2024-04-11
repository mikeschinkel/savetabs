package ui

import (
	"github.com/google/safehtml"
)

type MenuItemable interface {
	Identifier() safehtml.Identifier
}
