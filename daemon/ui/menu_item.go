package ui

import (
	"net/http"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

//goland:noinspection GoUnusedGlobalVariable

type HTMLMenuItem struct {
	HTMLId           safehtml.Identifier
	Label            safehtml.HTML
	LinksQueryParams safehtml.URL
	Slug             safehtml.URL
	Menu             *HTMLMenu
	MenuItemArgs
}

type MenuItemArgs struct {
	IconState IconState
}

func newHTMLMenuItem(menu *HTMLMenu, mi model.MenuItem) HTMLMenuItem {
	return newHTMLMenuItemWithArgs(menu, mi, nil)
}

func newHTMLMenuItemWithArgs(menu *HTMLMenu, mi model.MenuItem, args *MenuItemArgs) HTMLMenuItem {
	return HTMLMenuItem{}.RenewWithArgs(menu, mi, args)
}

var pristineHTMLMenuItem = HTMLMenuItem{
	MenuItemArgs: MenuItemArgs{
		IconState: CollapsedIcon,
	},
}

func (hmi HTMLMenuItem) Renew(menu *HTMLMenu, mi model.MenuItem) HTMLMenuItem {
	return hmi.RenewWithArgs(menu, mi, nil)
}

func (hmi HTMLMenuItem) RenewWithArgs(menu *HTMLMenu, mi model.MenuItem, args *MenuItemArgs) HTMLMenuItem {
	hmi = pristineHTMLMenuItem
	hmi.HTMLId = mi.HTMLId()
	hmi.Label = shared.MakeSafeHTML(mi.Label)
	hmi.Menu = menu
	hmi.Slug = mi.SubmenuURL()
	hmi.LinksQueryParams = mi.ItemURL()
	if args != nil {
		hmi.MenuItemArgs = *args
	}
	return hmi
}

func (hmi HTMLMenuItem) IsIconBlank() bool {
	return hmi.IconState == BlankIcon
}

func (hmi HTMLMenuItem) IsTopLevelMenu() bool {
	return hmi.Menu.Level == 0
}

func (hmi HTMLMenuItem) NotTopLevelMenu() bool {
	return hmi.Menu.Level != 0
}

type MenuItemHTMLParams struct {
	Menu     *HTMLMenu
	MenuType *shared.MenuType
	ItemType string
}

// GetMenuItemHTML responds to HTTP GET request with an text/html response
// containing the HTMX=flavored HTML for a menu item, which also includes its
// children.
func GetMenuItemHTML(ctx Context, p MenuItemHTMLParams) (hr HTMLResponse, err error) {
	var items model.MenuItems
	var htmlItems []HTMLMenuItem
	var m HTMLMenu

	hr.HTTPStatus = http.StatusOK

	if p.Menu == nil {
		panic("ERROR: A nil HTMLMenu was passed to ui.GetMenuItemHTML().")
	}
	m = *p.Menu

	items, err = model.LoadMenuItems(ctx, model.LoadMenuItemParams{
		MenuType: p.MenuType,
		Menu:     model.NewMenu(p.MenuType),
	})
	if err != nil {
		goto end
	}
	htmlItems = make([]HTMLMenuItem, len(items.Items))
	for i, item := range items.Items {
		htmlItems[i] = htmlItems[i].Renew(&m, item)
	}
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
