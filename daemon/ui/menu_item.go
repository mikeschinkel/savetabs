package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/model"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

var _ shared.MenuItemParent = (*HTMLMenuItem)(nil)
var _ shared.MenuItemable = (*HTMLMenuItem)(nil)

type HTMLMenuItem struct {
	parent      shared.MenuItemParent
	localId     string
	Label       safehtml.HTML
	menuType    *shared.MenuType
	contextMenu *shared.ContextMenu
	dropItem    shared.DragDropItem
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

func (hmi HTMLMenuItem) LocalId() string {
	return hmi.localId
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
	DropItem shared.DragDropItem
}

func NewHTMLMenuItem(mi model.MenuItem, args *HTMLMenuItemArgs) HTMLMenuItem {
	return HTMLMenuItem{
		parent:   args.Parent,
		menuType: args.MenuType,
		dropItem: args.DropItem,
	}
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

func (hmi HTMLMenuItem) DropTarget() string {
	return hmi.dropItem.DropTarget()
}

func (hmi HTMLMenuItem) DropTypes() string {
	return hmi.dropItem.DropTypes()
}

func (hmi HTMLMenuItem) HTMLId() safehtml.Identifier {
	ft := shared.NewFilter(hmi.FilterType(), 0)
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
	ft := shared.NewFilter(hmi.FilterType(), 0)
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
	hmi.dropItem = args.DropItem
	//pmt := hmi.parent.MenuType()
	//hmi.filterType = pmt.FilterType
	hmi.menuType = args.MenuType
	return hmi
}

func (hmi HTMLMenuItem) IsLeaf() bool {
	return hmi.MenuType().IsLeaf()
}

func (hmi HTMLMenuItem) HasChildren() bool {
	return hmi.MenuType().HasChildren()
}

func (hmi HTMLMenuItem) IsTopLevelMenu() bool {
	return false
}

func (hmi HTMLMenuItem) NotTopLevelMenuItem() bool {
	return hmi.parent.Level() > 0
}

type SubmenuHTMLArgs struct {
	Menu     *HTMLMenu
	MenuType *shared.MenuType
}

// GetSubmenuHTML responds to HTTP GET request with an text/html response
// containing the HTMX=flavored HTML for a menu item, which also includes its
// children.
func GetSubmenuHTML(ctx Context, args SubmenuHTMLArgs) (hr HTMLResponse, err error) {
	var items model.MenuItems

	hr = NewHTMLResponse()

	mt := args.MenuType

	//if args.Menu == nil {
	//	panic("ERROR: A nil HTMLMenu was passed to ui.GetMenuItemHTML().")
	//}
	items, err = model.LoadMenuItems(ctx, model.LoadMenuItemParams{
		MenuType: mt,
	})
	if err != nil {
		goto end
	}
	args.Menu.MenuItems = shared.ConvertSlice(items.Items, func(item model.MenuItem) HTMLMenuItem {
		return newHTMLMenuItem(item, &HTMLMenuItemArgs{
			Parent:   args.Menu,
			MenuType: mt,
			DropItem: shared.NewDropItem(shared.LinkToGroupDragDrop, item.DBId), //TODO: Fix the 0 to a real ID
		})
	})
	hr.HTML, err = menuTemplate.ExecuteToHTML(args.Menu)
end:
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError

	}
	return hr, err
}
