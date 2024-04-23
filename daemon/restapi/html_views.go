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

func (a *API) GetLinks(w http.ResponseWriter, r *http.Request, params GetLinksParams) {
	sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetLinksHTML(ctx, r.Host, params, r.URL.RawQuery)
	})
}

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuHTML(ctx, r.Host)
	})
}
