package restapi

import (
	"context"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return ui.GetMenuHTML(ctx, r.Host)
	})
}
