package ui

import (
	"fmt"
	"net/http"

	"github.com/google/safehtml"
	"savetabs/sqlc"
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
	var gts []sqlc.ListGroupsTypeRow

	db := sqlc.GetNestedDBTX(v.DataStore)
	err = db.Exec(func(dbtx sqlc.DBTX) (err error) {
		gts, err = v.Queries(dbtx).ListGroupsType(ctx)
		return err
	})
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
