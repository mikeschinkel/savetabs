package ui

import (
	"fmt"
	"net/http"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

type HTMLMenu struct {
	apiURL    safehtml.URL
	Level     int
	MenuType  *shared.MenuType
	MenuItems []HTMLMenuItem
}

func NewHTMLMenu(apiURL safehtml.URL, mt *shared.MenuType, level int) HTMLMenu {
	return HTMLMenu{}.Renew(apiURL, mt, level)
}

var pristineMenu HTMLMenu

func (m HTMLMenu) Renew(apiURL safehtml.URL, mt shared.MenuType, level int) HTMLMenu {
	m = pristineMenu
	m.apiURL = apiURL
	m.Level = level
	m.MenuType = mt
	m.MenuItems = make([]HTMLMenuItem, 0)
	return m
}

func (m HTMLMenu) HTMLMenuURL() string {
	return fmt.Sprintf("%s/html/menu", m.apiURL)
}

func (m HTMLMenu) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", m.apiURL)
}

var menuTemplate = GetTemplate("menu")

type HTMLMenuParams struct {
	Host shared.Host
}

func GetMenuHTML(ctx Context, p HTMLMenuParams) (hr HTMLResponse, err error) {
	var menu model.Menu

	hr.HTTPStatus = http.StatusOK

	menu, err = model.MenuLoad(ctx, model.MenuParams{
		Type: shared.GroupTypeMenuType,
	})

	hm := HTMLMenu{
		apiURL:    shared.MakeSafeURL(p.Host.URL()),
		MenuItems: make([]HTMLMenuItem, len(menu.Items)),
	}
	for i, item := range menu.Items {
		item.MenuItemable = menu
		hm.MenuItems[i] = hm.MenuItems[i].Renew(&hm, item)
	}

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
