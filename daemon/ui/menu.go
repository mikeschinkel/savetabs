package ui

import (
	"bytes"
	"fmt"
	"net/http"
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

func (v *Views) GetMenuHTML(ctx Context, host string) (html []byte, status int, err error) {
	var out bytes.Buffer
	var items []menuItem

	gts, err := v.Queries.ListGroupsType(ctx)
	if err != nil {
		goto end
	}
	items = menuItemsFromListGroupTypesRows(host, gts)
	err = menuTemplate.Execute(&out, menu{
		apiURL:    makeURL(host),
		MenuItems: items,
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, http.StatusInternalServerError, err
}
