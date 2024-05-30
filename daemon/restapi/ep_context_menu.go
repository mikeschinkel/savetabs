package restapi

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"savetabs/guard"
	"savetabs/shared"
)

func (a *API) PutContextMenuContextMenuTypeIdName(w http.ResponseWriter, r *http.Request, contextMenuType ContextMenuType, id IdParameter) {
	var uvs url.Values
	var name string
	var args guard.ContextMenuItemArgs
	var hr guard.HTMLResponse
	var merged bool

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendHTMLError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	uvs, err = url.ParseQuery(string(body))
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	name = strings.TrimSpace(uvs.Get("name"))
	if name == "" {
		a.sendHTMLError(w, r, http.StatusBadRequest, "Name value cannot be empty")
		return
	}

	slog.Debug("PUT name:", "name", name, "content_menu", contextMenuType, "db_id", id)

	args = guard.ContextMenuItemArgs{
		ContextMenuArgs: guard.ContextMenuArgs{
			Host:            r.Host,
			ContextMenuType: contextMenuType,
			DBId:            id,
		},
		Name: name,
	}

	merged, err = guard.UpdateContextMenuItemName(ctx, args)
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	if merged {
		args.Name = fmt.Sprintf("Merged with %s", args.Name)
	}

	hr.HTML = shared.MakeSafeHTML(args.Name)
	a.sendHTML(w, hr.HTML)
}
