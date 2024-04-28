package ui

import (
	"fmt"
	"net/http"

	"github.com/google/safehtml"
)

type menu struct {
	apiURL    string
	MenuItems []menuItem
}

func (m menu) HTMLMenuURL() string {
	return fmt.Sprintf("%s/html/menu", m.apiURL)
}

func (m menu) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", m.apiURL)
}

var menuTemplate = GetTemplate("menu")

func (v *Views) GetMenuHTML(ctx Context, host string) (html safehtml.HTML, status int, err error) {
	var items []menuItem

	gts, err := v.Queries.ListGroupsType(ctx)
	if err != nil {
		goto end
	}
	items = menuItemsFromListGroupTypesRows(host, gts)
	html, err = menuTemplate.ExecuteToHTML(menu{
		apiURL:    makeURL(host),
		MenuItems: items,
	})
	if err != nil {
		goto end
	}
end:
	return html, http.StatusInternalServerError, err
}
