package ui

import (
	"fmt"
	"net/http"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

var _ shared.MenuItemParent = (*HTMLMenu)(nil)

type HTMLMenu struct {
	apiURL    safehtml.URL
	menuType  *shared.MenuType
	parent    *HTMLMenu
	MenuItems []HTMLMenuItem
}

func NewHTMLMenu(args HTMLMenuArgs) *HTMLMenu {
	return HTMLMenu{}.Renew(args)
}

var zeroStateMenu HTMLMenu

func (hm HTMLMenu) IsTopLevelMenu() bool {
	return hm.menuType.Level() <= 1
}

func (hm HTMLMenu) Renew(args HTMLMenuArgs) *HTMLMenu {
	hm = zeroStateMenu
	hm.apiURL = args.APIURL
	hm.MenuItems = make([]HTMLMenuItem, 0)
	hm.menuType = args.MenuType
	if hm.menuType == nil {
		shared.Panicf("NewHTMLMenu() or HTMLMenu.Renew() called with `nil` menuType")
	}
	return &hm
}

func (hm HTMLMenu) APIURL() safehtml.URL {
	return hm.apiURL
}

func (hm HTMLMenu) HTMLId() (id safehtml.Identifier) {
	return safehtml.IdentifierFromConstant("mi")
}

func (hm HTMLMenu) MenuType() *shared.MenuType {
	return hm.menuType
}

func (hm HTMLMenu) Parent() shared.MenuItemParent {
	return hm.parent
}

func (hm HTMLMenu) Level() int {
	return 0
}

type HTMLMenuArgs struct {
	APIURL   safehtml.URL
	MenuType *shared.MenuType
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
	var hm *HTMLMenu

	menuType := shared.GroupTypeMenuType

	hr = NewHTMLResponse()

	menu, err = model.LoadMenu(ctx, model.MenuParams{
		Type: menuType,
	})
	if err != nil {
		hr.SetCode(http.StatusInternalServerError)
		goto end
	}

	hm = NewHTMLMenu(HTMLMenuArgs{
		APIURL:   shared.MakeSafeURL(p.Host.URL()),
		MenuType: menuType,
	})

	hm.MenuItems = shared.ConvertSlice(menu.Items, func(item model.MenuItem) HTMLMenuItem {
		return newHTMLMenuItem(item, &HTMLMenuItemArgs{
			Parent:   hm,
			MenuType: menuType,
		})
	})

	hr.HTML, err = menuTemplate.ExecuteToHTML(hm)
	if err != nil {
		hr.SetCode(http.StatusInternalServerError)
		goto end
	}
end:
	return hr, err
}
