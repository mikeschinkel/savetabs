package ui

import (
	"github.com/google/safehtml"
	"savetabs/shared"
)

type ContextMenuItem struct {
	Label       safehtml.HTML
	ContextMenu *ContextMenu
}

func (cmi ContextMenuItem) Target() safehtml.URL {
	return shared.MakeSafeURLf("#%s", cmi.ContextMenu)
}

func (cmi ContextMenuItem) MethodURL() safehtml.URL {
	return cmi.ContextMenu.RenameFormURL() // TODO: Make this generic
}
