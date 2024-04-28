package restapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/safehtml"
	"savetabs/ui"
)

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	a.sendView(context.Background(), w, r, func(ctx Context) (safehtml.HTML, int, error) {
		return a.Views.GetMenuItemHTML(ctx, r.Host, menuItem)
	})
}

func (a *API) GetHtmlLinkset(w http.ResponseWriter, r *http.Request, params GetHtmlLinksetParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) (safehtml.HTML, int, error) {
		return a.Views.GetLinkSetHTML(ctx, r.Host, r.RequestURI, params)
	})
}

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	a.sendView(context.Background(), w, r, func(ctx Context) (safehtml.HTML, int, error) {
		return a.Views.GetMenuHTML(ctx, r.Host)
	})
}

func (a *API) GetHtmlError(w http.ResponseWriter, r *http.Request, params GetHtmlErrorParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) (safehtml.HTML, int, error) {
		return a.Views.GetErrorHTML(ctx, errors.New(*params.Err))
	})
}
func (a *API) GetHtmlAlert(w http.ResponseWriter, r *http.Request, params GetHtmlAlertParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) (safehtml.HTML, int, error) {
		return a.Views.GetAlertHTML(ctx, ui.AlertType(*params.Typ), ui.Message{Text: *params.Msg})
	})
}
