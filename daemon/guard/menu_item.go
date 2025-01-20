package guard

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/ui"
)

type MenuItem struct {
	value string
}

// matchMenuItemIds matcbes values menu item and also capturing the prefix and
// the localId into a 2 element slice, e.g.
//
//	result[0] => prefix, e.g. `gt`, grp`
//	result[1] => localId, e.g, `tabgroup` & `category`; `golang` and `htmx`, etc.
var matchMenuItemIds = regexp.MustCompile(fmt.Sprintf(`^(%s)-(.+)$`, shared.MenuTypesForRegexp()))

// GetSubmenuHTML return HTMX-flavored HTML for a single menu item
func GetSubmenuHTML(ctx Context, host string, menuItem string) (_ HTMLResponse, err error) {
	var apiURL safehtml.URL
	var mt *shared.MenuType

	hr := ui.NewHTMLResponse()
	// TODO: Review `id`, `key`, `menuItem` etc. semantics
	if !matchMenuItemIds.MatchString(menuItem) {
		err = errors.Join(ErrInvalidMenuItemFormat, fmt.Errorf(`menu_item=%s`, menuItem))
		hr.StatusCode = http.StatusBadRequest
		goto end
	}
	mt, err = shared.MenuTypeByName(menuItem)
	if err != nil {
		goto end
	}
	apiURL = shared.MakeSafeURL(shared.NewHost(host).URL())
	hr, err = ui.GetSubmenuHTML(ctx, ui.SubmenuHTMLArgs{
		MenuType: mt,
		Menu: ui.NewHTMLMenu(ui.HTMLMenuArgs{
			APIURL:   apiURL,
			MenuType: mt,
		}),
	})
end:
	return HTMLResponse{hr}, err
}
