package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

var _ shared.MenuItemParent = (*HTMLMenuItem)(nil)
var _ shared.MenuItemable = (*HTMLMenuItem)(nil)

type HTMLMenuItem struct {
	parent      shared.MenuItemParent
	localId     string
	Label       safehtml.HTML
	menuType    *shared.MenuType
	contextMenu *shared.ContextMenu
}

func (hmi HTMLMenuItem) ContextMenuType() (id safehtml.Identifier) {
	return shared.MakeSafeId(hmi.contextMenu.Type.Name)
}
func (hmi HTMLMenuItem) ContextMenuDBId() int64 {
	return hmi.contextMenu.Id
}
func (hmi HTMLMenuItem) ContextMenuId() safehtml.Identifier {
	return shared.MakeSafeId(hmi.contextMenu.String())
}

func (hmi HTMLMenuItem) LocalId() safehtml.Identifier {
	return shared.MakeSafeId(hmi.localId)
}

func (hmi HTMLMenuItem) Parent() shared.MenuItemParent {
	return hmi.parent
}

func (hmi HTMLMenuItem) MenuType() *shared.MenuType {
	return hmi.menuType
}

func (hmi HTMLMenuItem) Level() int {
	return hmi.Parent().Level() + 1
}

type HTMLMenuItemArgs struct {
	Parent   shared.MenuItemParent
	MenuType *shared.MenuType
}

var zeroStateHTMLMenuItem HTMLMenuItem

func (hmi HTMLMenuItem) APIURL() safehtml.URL {
	return hmi.parent.APIURL()
}

func (hmi HTMLMenuItem) RenameEndpoint() safehtml.URL {
	return shared.MakeSafeURLf("%s/groups/99999999999/name", hmi.APIURL()) // TODO:  Update this to actual URL
}

func (hmi HTMLMenuItem) HTMLContextMenuURL() safehtml.URL {
	return shared.MakeSafeURLf("%s/html/context-menu", hmi.APIURL())
}

func (hmi HTMLMenuItem) HasContextMenu() bool {
	return hmi.contextMenu != nil
}

func (hmi HTMLMenuItem) HTMLId() safehtml.Identifier {
	ft := shared.NewFilterByFilterType(hmi.FilterType())
	id := ft.HTMLId(hmi)
	return shared.MakeSafeIdf("%s-%s", hmi.Parent().HTMLId(), id)
}

func (hmi HTMLMenuItem) ChildMenuURL() safehtml.URL {
	var u string
	if hmi.IsLeaf() {
		u = "#"
		goto end
	}
	// TODO: Change this to use URL encoding
	u = fmt.Sprintf("%s--%s", hmi.MenuType().FilterType.Id(), hmi.LocalId())
end:
	return shared.MakeSafeURL(u)
}

func (hmi HTMLMenuItem) FilterType() *shared.FilterType {
	return hmi.menuType.FilterType
}

func (hmi HTMLMenuItem) ContentQuery() safehtml.URL {
	var pcq string
	pmi, ok := hmi.Parent().(shared.MenuItemable)
	if ok {
		pcq = pmi.ContentQuery().String() + "&"
	}
	ft := shared.NewFilterByFilterType(hmi.FilterType())
	u := ft.ContentQuery(hmi)
	return shared.MakeSafeURL("?" + pcq + u)
}

func newHTMLMenuItem(mi model.MenuItem, args *HTMLMenuItemArgs) HTMLMenuItem {
	return HTMLMenuItem{}.Renew(mi, args)
}

func (hmi HTMLMenuItem) Renew(mi model.MenuItem, args *HTMLMenuItemArgs) HTMLMenuItem {
	if args == nil {
		shared.Panicf("RenewWithArgs: args must not be nil")
		// This next panic will never execute.
		// It is only  here so GoLand will stop saying "`args` may be nil."
		panic("")
	}
	hmi = zeroStateHTMLMenuItem
	hmi.Label = shared.MakeSafeHTML(mi.Label)
	hmi.localId = strings.ToLower(mi.LocalId)
	hmi.parent = args.Parent
	hmi.contextMenu = mi.ContextMenu
	//pmt := hmi.parent.MenuType()
	//hmi.filterType = pmt.FilterType
	hmi.menuType = args.MenuType
	return hmi
}

func (hmi HTMLMenuItem) IsLeaf() bool {
	return len(hmi.MenuType().Children) == 0
}

func (hmi HTMLMenuItem) IsTopLevelMenu() bool {
	return hmi.Level() == 1
}

func (hmi HTMLMenuItem) NotTopLevelMenu() bool {
	return !hmi.IsTopLevelMenu()
}

type MenuItemHTMLParams struct {
	Menu     *HTMLMenu
	MenuType *shared.MenuType
}

// GetMenuItemHTML responds to HTTP GET request with an text/html response
// containing the HTMX=flavored HTML for a menu item, which also includes its
// children.
func GetMenuItemHTML(ctx Context, p MenuItemHTMLParams) (hr HTMLResponse, err error) {
	var items model.MenuItems

	hr = NewHTMLResponse()

	if p.Menu == nil {
		panic("ERROR: A nil HTMLMenu was passed to ui.GetMenuItemHTML().")
	}

	items, err = model.LoadMenuItems(ctx, model.LoadMenuItemParams{
		MenuType: p.MenuType,
		Menu:     model.NewMenu(p.MenuType),
	})
	if err != nil {
		goto end
	}
	p.Menu.MenuItems = shared.ConvertSlice(items.Items, func(item model.MenuItem) HTMLMenuItem {
		return newHTMLMenuItem(item, &HTMLMenuItemArgs{
			Parent:   p.Menu,
			MenuType: p.MenuType,
		})
	})
	hr.HTML, err = menuTemplate.ExecuteToHTML(p.Menu)
end:
	if err != nil {
		hr.SetCode(http.StatusInternalServerError)
	}
	return hr, err
}
