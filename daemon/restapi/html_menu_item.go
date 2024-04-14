package restapi

import (
	"context"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, error) {
		return ui.GetMenuItemHTML(ctx, r.Host, menuItem)
	})
}
