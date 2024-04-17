package restapi

import (
	"context"
	"net/http"
)

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuItemHTML(ctx, r.Host, menuItem)
	})
}

func (a *API) GetLinks(w http.ResponseWriter, r *http.Request, params GetLinksParams) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetLinksHTML(ctx, r.Host, params)
	})
}

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuHTML(ctx, r.Host)
	})
}
