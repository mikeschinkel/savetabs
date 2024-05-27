package ui

import (
	"net/http"

	"github.com/google/safehtml"
	"savetabs/shared"
)

var contextMenuTemplate = GetTemplate("context-menu")

type ContextMenuArgs struct {
	APIURL      safehtml.URL
	ContextMenu ContextMenu
	Items       []ContextMenuItem
}

type ContextMenuItems struct {
	Items []ContextMenuItem
}

func newContextMenuItems(items []ContextMenuItem) ContextMenuItems {
	return ContextMenuItems{
		Items: items,
	}
}

type ContextMenu struct {
	Type  *shared.ContextMenuType
	Items []ContextMenuItem
}

func NewContextMenu(cmt *shared.ContextMenuType) ContextMenu {
	return ContextMenu{
		Type:  cmt,
		Items: make([]ContextMenuItem, 0),
	}
}

type ContextMenuItem struct {
	Label           safehtml.HTML
	ContextMenuType *shared.ContextMenuType
}

func GetContextMenuHTML(ctx Context, args ContextMenuArgs) (_ HTMLResponse, err error) {
	var items ContextMenuItems

	hr := NewHTMLResponse()

	items = newContextMenuItems(args.Items)
	hr.HTML, err = menuTemplate.ExecuteToHTML(items)

	if err != nil {
		hr.SetCode(http.StatusInternalServerError)
	}
	return hr, err
}
