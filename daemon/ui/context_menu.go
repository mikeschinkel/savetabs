package ui

import (
	"fmt"
	"net/http"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

var (
	contextMenuTemplate               = GetTemplate("context-menu")
	contextMenuItemRenameFormTemplate = GetTemplate("context-menu-item-rename-form")
)

type ContextMenuArgs struct {
	ContextMenu *ContextMenu
	Items       []ContextMenuItem
	DBId        int64
}

type ContextMenu struct {
	apiURL safehtml.URL
	Type   *shared.ContextMenuType
	DBId   int64
	Items  []ContextMenuItem
	Name   safehtml.HTML
}

func NewContextMenu(cmt *shared.ContextMenuType, host string) *ContextMenu {
	return &ContextMenu{
		apiURL: shared.MakeSafeAPIURL(host),
		Type:   cmt,
		Items:  make([]ContextMenuItem, 0),
	}
}
func (cm ContextMenu) String() string {
	return fmt.Sprintf("%s-%d", cm.Type.Name, cm.DBId)
}

func (cm ContextMenu) LoadName(ctx Context) (_ safehtml.HTML, err error) {
	var name string
	switch cm.Type {
	case shared.GroupContextMenuType:
		name, err = model.LoadGroupName(ctx, cm.DBId)
	}
	return shared.MakeSafeHTML(name), err
}

func (cm ContextMenu) RenameFormURL() safehtml.URL {
	return shared.MakeSafeURLf("%s/html/context-menu/%s/%d/rename-form",
		cm.apiURL,
		cm.Type.Name,
		cm.DBId,
	)
}

func (cm ContextMenu) NameURL() safehtml.URL {
	return shared.MakeSafeURLf("%s/context-menu/%s/%d/name",
		cm.apiURL,
		cm.Type.Name,
		cm.DBId,
	)
}

func GetContextMenuHTML(ctx Context, args ContextMenuArgs) (_ HTMLResponse, err error) {
	hr := NewHTMLResponse()
	cm := args.ContextMenu
	cm.Items = args.Items
	hr.HTML, err = contextMenuTemplate.ExecuteToHTML(cm)
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError
	}
	return hr, err
}

func GetContextMenuRenameFormHTML(ctx Context, args ContextMenuArgs) (_ HTMLResponse, err error) {
	cm := args.ContextMenu
	hr := NewHTMLResponse()

	cm.Name, err = cm.LoadName(ctx)
	if err != nil {
		goto end
	}

	hr.HTML, err = contextMenuItemRenameFormTemplate.ExecuteToHTML(cm)
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError
	}
end:
	return hr, err
}

func GetContextMenuItemNameHTML(ctx Context, name safehtml.HTML) (_ HTMLResponse, err error) {
	hr := NewHTMLResponse()
	hr.HTML = name
	return hr, err
}
