package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

var _ shared.MenuItemable = (*HTMLMenuItem)(nil)

type HTMLMenuItem struct {
	shared.MenuItemable
	LocalId   string
	Label     safehtml.HTML
	Type      *shared.MenuType
	IconState IconState
}

type HTMLMenuItemArgs struct {
	shared.MenuItemable
}

var zeroStateHTMLMenuItem HTMLMenuItem

func (hmi HTMLMenuItem) MenuType() *shared.MenuType {
	return hmi.Type
}

func (hmi HTMLMenuItem) HTMLId() safehtml.Identifier {
	return shared.MakeSafeId(fmt.Sprintf("%s-%s",
		hmi.MenuItemable.HTMLId(),
		hmi.LocalId,
	))
}

func (hmi HTMLMenuItem) SubmenuURL() safehtml.URL {
	var u string
	if hmi.IsLeaf() {
		u = "#"
		goto end
	}
	u = hmi.MenuType().Params(shared.ParamsArgs{
		Equates:  "--",
		Combines: "/",
	})
end:
	return shared.MakeSafeURL(u)
}

func (hmi HTMLMenuItem) LinksQuery() safehtml.URL {
	mt := hmi.MenuType()
	return shared.MakeSafeURLf("?%s", mt.Params(shared.ParamsArgs{
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
	hmi.MenuItemable = args.MenuItemable
	pmt := hmi.MenuItemable.MenuType()
	mt, err := shared.MenuTypeByParentTypeAndMenuName(pmt, mi.LocalId)
	if err != nil {
		mt = shared.CloneLeafMenuType()
		mt.Parent = pmt
		mt.SetName(hmi.LocalId)
	}
	hmi.Type = mt
	return hmi
}

func (hmi HTMLMenuItem) IsLeaf() bool {
	return hmi.MenuType().IsLeaf
}

func (hmi HTMLMenuItem) IsTopLevelMenu() bool {
	return hmi.MenuItemable.Level() <= 1
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
			MenuItemable: *p.Menu,
		})
	})
	hr.HTML, err = menuTemplate.ExecuteToHTML(p.Menu)
end:
	if err != nil {
		hr.SetCode(http.StatusInternalServerError)
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
