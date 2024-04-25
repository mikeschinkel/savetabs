package restapi

import (
	"context"
	"errors"
	"net/http"
)

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	a.sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuItemHTML(ctx, r.Host, menuItem)
	})
}

func (a *API) GetHtmlLinkset(w http.ResponseWriter, r *http.Request, params GetHtmlLinksetParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetLinkSetHTML(ctx, r.Host, params, r.URL.RawQuery)
	})
}

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	a.sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetMenuHTML(ctx, r.Host)
	})
}

func (a *API) GetHtmlError(w http.ResponseWriter, r *http.Request, params GetHtmlErrorParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return a.Views.GetErrorHTML(ctx, errors.New(*params.Err))
	})
}
