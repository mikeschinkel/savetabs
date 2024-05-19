package guard

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/ui"
)

type MenuItem struct {
	value string
}

func NewMenuItem(v string) MenuItem {
	return MenuItem{value: v}
}

// matchMenuItemIds matcbes values menu item and also capturing the prefix and
// the localId into a 2 element slice, e.g.
//
//	result[0] => prefix, e.g. `gt`, grp`
//	result[1] => localId, e.g, `tabgroup` & `category`; `golang` and `htmx`, etc.
var matchMenuItemIds = regexp.MustCompile(fmt.Sprintf(`^(%s)-(.+)$`, shared.MenuTypesForRegexp()))

// GetMenuItemHTML return HTMX-flavored HTML for a single menu item
func GetMenuItemHTML(ctx Context, host string, menuItem MenuItem) (_ HTMLResponse, err error) {
	var hr ui.HTMLResponse
	var apiURL safehtml.URL
	var mt shared.MenuType

	// TODO: Review `id`, `key`, `menuItem` etc. semantics
	id := menuItem.value
	if !matchMenuItemIds.MatchString(id) {
		err = errors.Join(ErrInvalidKeyFormat, fmt.Errorf(`id=%s`, id))
		hr.HTTPStatus = http.StatusBadRequest
		goto end
	}
	mt, err = shared.MenuTypeByValue(id)
	if err != nil {
		goto end
	}
	apiURL = shared.MakeSafeURL(shared.NewHost(host).URL())
	hr, err = ui.GetMenuItemHTML(ctx, ui.MenuItemHTMLParams{
		Menu:     shared.Ptr(ui.NewHTMLMenu(apiURL, &shared.GroupTypeMenuType, 1)), // TODO: Need to make 2nd two params dyanmic.
		MenuType: mt,
		ItemType: mt.Leaf(),
	})
end:
	return HTMLResponse(hr), err
}
