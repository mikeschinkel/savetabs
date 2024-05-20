package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

var _ shared.Menu = (*HTMLMenuItem)(nil)

type HTMLMenuItem struct {
	shared.Menu
	LocalId   string
	Label     safehtml.HTML
	Type      *shared.MenuType
	IconState IconState
}

type HTMLMenuItemArgs struct {
	IconState IconState
	shared.Menu
}

var zeroStateHTMLMenuItem HTMLMenuItem

func (hmi HTMLMenuItem) MenuType() *shared.MenuType {
	return hmi.Type
}

func (hmi HTMLMenuItem) HTMLId() safehtml.Identifier {
	return shared.MakeSafeId(fmt.Sprintf("%s-%s",
		hmi.Menu.HTMLId(),
		hmi.LocalId,
	))
}

func (hmi HTMLMenuItem) SubmenuURL() safehtml.URL {
	return shared.MakeSafeURL(hmi.MenuType().Params(shared.ParamsArgs{
		Equates:  "--",
		Combines: "/",
	}))
}

func (hmi HTMLMenuItem) Slug() safehtml.URL {
	return hmi.SubmenuURL()
}
func (hmi HTMLMenuItem) LinksQueryParams() safehtml.URL {
	return hmi.ItemURL()
}

func (hmi HTMLMenuItem) ItemURL() safehtml.URL {
	return shared.MakeSafeURL(hmi.MenuType().Params(shared.ParamsArgs{
		Equates:  "=",
		Combines: "&",
	}))
}

func newHTMLMenuItem(mi model.MenuItem, args *HTMLMenuItemArgs) HTMLMenuItem {
	return HTMLMenuItem{}.Renew(mi, args)
}

func (hmi HTMLMenuItem) Renew(mi model.MenuItem, args *HTMLMenuItemArgs) HTMLMenuItem {
	if args == nil {
		shared.Panicf("RenewWithArgs: args must not be nil")
	}
	hmi = zeroStateHTMLMenuItem
	hmi.Label = shared.MakeSafeHTML(mi.Label)
	hmi.LocalId = strings.ToLower(mi.LocalId)
	hmi.Menu = args.Menu
	//mt,err := shared.MenuTypeByParentTypeAndMenuName(shared.GroupTypeMenuType, args.LocalId)
	mt, err := shared.MenuTypeByParentTypeAndMenuName(hmi.Menu.MenuType(), mi.LocalId)
	if err != nil {
		shared.Panicf(err.Error())
	}
	hmi.Type = mt
	if args.IconState == ZeroStateIcon {
		args.IconState = CollapsedIcon
	}
	hmi.IconState = args.IconState
	return hmi
}

func (hmi HTMLMenuItem) IsIconBlank() bool {
	return hmi.IconState == BlankIcon
}

func (hmi HTMLMenuItem) IsTopLevelMenu() bool {
	return hmi.Menu.Level() == 0
}

func (hmi HTMLMenuItem) NotTopLevelMenu() bool {
	return hmi.Menu.Level() != 0
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
	var htmlItems []HTMLMenuItem

	hr.HTTPStatus = http.StatusOK

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
	htmlItems = shared.ConvertSlice(items.Items, func(item model.MenuItem) HTMLMenuItem {
		return newHTMLMenuItem(item, &HTMLMenuItemArgs{
			Menu: *p.Menu,
		})
	})
	hr.HTML, err = menuTemplate.ExecuteToHTML(HTMLMenu{
		MenuItems: htmlItems,
	})
end:
	if err != nil {
		hr.HTTPStatus = http.StatusInternalServerError
	}
	return hr, err
}

//func getMenuItemsForType(ctx Context, host, key string) (items []HTMLMenuItem, err error) {
//	var keys []string
//	var gt sqlc.GroupType
//	var gs []sqlc.Group
//	var db *storage.NestedDBTX
//
//	db = storage.GetNestedDBTX(v.DataStore)
//	err = db.Exec(func(dbtx *storage.NestedDBTX) (err error) {
//		q := v.Queries(dbtx)
//		switch keys[1] {
//		case shared.GroupTypeMenuType: // Group Type
//			gs, err = q.ListGroupsByType(ctx, sqlc.ListGroupsByTypeParams{
//				Type:           strings.ToUpper(keys[2]),
//				GroupsArchived: storage.NotArchived,
//				GroupsDeleted:  storage.NotDeleted,
//			})
//			if err != nil {
//				goto end
//			}
//		}
//		err = nil
//		items = func(ctx Context, host string, gt groupType, gs []sqlc.Group) []HTMLMenuItem {
//			var menuItems []HTMLMenuItem
//
//			// Instantiate new menu
//			// Groups are level == 1, aka children of Group Types where level == 0
//			m := newHTMLMenu(host, 1)
//
//			args := MenuItemArgs{
//				IconState: BlankIcon,
//			}
//			menuItems = make([]HTMLMenuItem, len(gs)+1)
//			menuItems[0] = newHTMLMenuItemWithArgs(&m, model.MenuItem{
//				LocalId: "none",
//				Label:   fmt.Sprintf("<No %s>", gt.Plural),
//			}, args)
//			groups, err := model.LoadGroups(ctx,model.GroupsParams{
//				Host:       shared.NewHost(host),
//				GroupType:  gt.Type,
//			})
//			for i, g := range gs {
//				grp := model.NewGroup(g)
//				menuItems[i+1] = newHTMLMenuItemWithArgs(&m, model.MenuItem{
//					LocalId: strings.ToLower(grp.Type),
//					Label:   grp.Name,
//				}, args)
//			}
//			return menuItems
//
//		}(ctx, host, gt, gs)
//	end:
//		return err
//	})
//	if err != nil {
//		goto end
//	}
//end:
//	return items, err
//}
//
