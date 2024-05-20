package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

var _ shared.Menu = (*HTMLMenu)(nil)

type HTMLMenu struct {
	apiURL    safehtml.URL
	Type      *shared.MenuType
	MenuItems []HTMLMenuItem
	level     int
	parent    *HTMLMenu
}

func (hm HTMLMenu) HTMLId() (id safehtml.Identifier) {
	if hm.level == 0 {
		id = safehtml.IdentifierFromConstant("mi")
		goto end
	}
	id = shared.MakeSafeIdf("mi-%s", strings.Join(hm.Type.Slice(), "-"))
end:
	return id
}

func (hm HTMLMenu) MenuType() *shared.MenuType {
	return hm.Type
}

func (hm HTMLMenu) Parent() shared.Menu {
	return hm.parent
}

func (hm HTMLMenu) ItemURL() safehtml.URL {
	panic("ItemURL() should not be called for HTMLMenu")
}

func (hm HTMLMenu) SubmenuURL() safehtml.URL {
	panic("SubmenuURL() should not be called for HTMLMenu")
}

func (hm HTMLMenu) Level() int {
	return hm.level
}

type HTMLMenuArgs struct {
	APIURL safehtml.URL
	Type   *shared.MenuType
}

func NewHTMLMenu(args HTMLMenuArgs) *HTMLMenu {
	return HTMLMenu{}.Renew(args)
}

var zeroStateMenu HTMLMenu

func (hm HTMLMenu) Renew(args HTMLMenuArgs) *HTMLMenu {
	hm = zeroStateMenu
	hm.apiURL = args.APIURL
	hm.MenuItems = make([]HTMLMenuItem, 0)
	hm.Type = args.Type
	if hm.Type == nil {
		shared.Panicf("NewHTMLMenu() or HTMLMenu.Renew() called with `nil` Type")
	}
	hm.level = hm.Type.Level()
	if hm.level > 0 {
		hm.parent = NewHTMLMenu(HTMLMenuArgs{
			APIURL: hm.apiURL,
			Type:   hm.Type.Parent,
		})
	}
	return &hm
}

func (hm HTMLMenu) HTMLMenuURL() string {
	return fmt.Sprintf("%s/html/menu", hm.apiURL)
}

func (hm HTMLMenu) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", hm.apiURL)
}

var menuTemplate = GetTemplate("menu")

type HTMLMenuParams struct {
	Host shared.Host
}

func GetMenuHTML(ctx Context, p HTMLMenuParams) (hr HTMLResponse, err error) {
	var menu *model.Menu

	hr.HTTPStatus = http.StatusOK

	menu, err = model.MenuLoad(ctx, model.MenuParams{
		Type: shared.GroupTypeMenuType,
	})

	hm := NewHTMLMenu(HTMLMenuArgs{
		APIURL: shared.MakeSafeURL(p.Host.URL()),
		Type:   menu.Type,
	})

	hm.MenuItems = shared.ConvertSlice(menu.Items, func(item model.MenuItem) HTMLMenuItem {
		return newHTMLMenuItem(item, &HTMLMenuItemArgs{
			IconState: CollapsedIcon,
			Menu:      hm,
		})
	})

	hr.HTML, err = menuTemplate.ExecuteToHTML(hm)
	if err != nil {
		hr.HTTPStatus = http.StatusInternalServerError
		goto end
	}
end:
	return hr, err
}

//func menuItemsFromListGroupTypesRows(host string, gtrs []sqlc.ListGroupsTypeRow) []menuItem {
//	var menuItems []menuItem
//
//	cnt := len(gtrs)
//
//	// No need to show invalid as a group type if
//	// there are no resources of that type
//	invalid := -1
//	for i, gtr := range gtrs {
//		if gtr.LinkCount != 0 {
//			continue
//		}
//		if gtr.Type != "I" {
//			continue
//		}
//		cnt--
//		invalid = i
//		break
//	}
//	menuItems = make([]menuItem, cnt)
//	for i, gtr := range gtrs {
//		if i == invalid {
//			continue
//		}
//		src := newGroupTypeFromListGroupsTypeRow(gtr)
//		menuItems[i] = newMenuItem(src, host, gtr.Plural.String)
//	}
//	menuItems = append(menuItems,
//		newHTMLMenuItemWithArgs(allLinks{}, host, "All Links", menuItemArgs{
//			SummaryClass: topLevelSummaryClass,
//			IconState:    BlankIcon,
//		}),
//	)
//	return menuItems
//}
//
//
