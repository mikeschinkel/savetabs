package restapi

import (
	"context"
	"net/http"
)

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuItemHTML(ctx, r.Host, menuItem)
	})
}

func (a *API) GetHtmlLinkset(w http.ResponseWriter, r *http.Request, params GetHtmlLinksetParams) {
	sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetLinkSetHTML(ctx, r.Host, params, r.URL.RawQuery)
	})
}

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuHTML(ctx, r.Host)
	})
}
