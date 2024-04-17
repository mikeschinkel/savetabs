package restapi

import (
	"context"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetHtmlBrowse(w http.ResponseWriter, r *http.Request) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return ui.GetBrowseHTML(ctx, r.Host)
	})
}
