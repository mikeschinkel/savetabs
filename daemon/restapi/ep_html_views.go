package restapi

import (
	"context"
	"net/http"

	"github.com/mikeschinkel/savetabs/daemon/guard"
)

func (a *API) GetHtmlContextMenuContextMenuTypeId(w http.ResponseWriter, r *http.Request, contextMenuType ContextMenuType, id IdParameter) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		return guard.GetContextMenuHTML(ctx, guard.ContextMenuArgs{
			Host:            r.Host,
			ContextMenuType: contextMenuType,
			DBId:            id,
		})
	})
}
func (a *API) GetHtmlContextMenuContextMenuTypeIdRenameForm(w http.ResponseWriter, r *http.Request, contextMenuType ContextMenuType, id IdParameter) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		return guard.GetContextMenuRenameFormHTML(ctx, guard.ContextMenuArgs{
			Host:            r.Host,
			ContextMenuType: contextMenuType,
			DBId:            id,
		})
	})
}

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		return guard.GetSubmenuHTML(ctx, r.Host, menuItem)
	})
}

func (a *API) GetHtmlLinkset(w http.ResponseWriter, r *http.Request, params GetHtmlLinksetParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		// TODO: Implement validation for these filters before passing them on
		return guard.GetLinksetHTML(ctx, guard.LinksetParams{
			RequestURI: r.RequestURI,
			Host:       r.Host,
		})
	})
}

func (a *API) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		return guard.GetMenuHTML(ctx, guard.MenuParams{
			Host: r.Host,
		})
	})
}

func (a *API) GetHtmlError(w http.ResponseWriter, r *http.Request, params GetHtmlErrorParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		return guard.GetErrorHTML(guard.ErrorParams{
			Host: r.Host,
			Msg:  *params.Err,
		})
	})
}
func (a *API) GetHtmlAlert(w http.ResponseWriter, r *http.Request, params GetHtmlAlertParams) {
	a.sendView(context.Background(), w, r, func(ctx Context) (guard.HTMLResponse, error) {
		return guard.GetAlertHTML(ctx, guard.AlertParams{
			Type: string(*params.Typ),
			Msg:  *params.Msg,
			Host: r.Host,
		})
	})
}
