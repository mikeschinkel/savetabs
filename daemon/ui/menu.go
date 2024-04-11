package ui

import (
	"bytes"
	"context"
	"fmt"
)

type menu struct {
	apiURL    string
	MenuItems []menuItem
}

func (m menu) HTMLMenuURL() string {
	return fmt.Sprintf("%s/html/menu", m.apiURL)
}

var menuTemplate = GetTemplate("menu")

func MenuHTML(host string) (html []byte, err error) {
	var out bytes.Buffer
	var items []menuItem

	gts, err := queries.ListGroupsType(context.Background())
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
	return html, err
}
