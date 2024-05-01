package restapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/sqlc"
	"savetabs/storage"
	"savetabs/ui"
)

func (a *API) PostHtmlLinkset(w http.ResponseWriter, r *http.Request) {
	var msg string
	var status int
	var params ui.FilterGetter

	ctx := context.TODO()

	err := r.ParseForm()
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	_, ok := r.Form["link_id"]
	if !ok {
		a.sendError(w, r, http.StatusBadRequest, ErrNoLinkIdsSubmitted.Error())
		return
	}
	linkIds := make([]int64, len(r.Form["link_id"]))
	for i, id := range r.Form["link_id"] {
		linkIds[i], err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid, expected integer", "link_id", id)
		}
	}
	queryJSONs, ok := r.Form["query_json"]
	if !ok {
		queryJSONs = []string{"{}"}
	}
	queryJSON := strings.Join(queryJSONs, "")

	ds := sqlc.GetDatastore()
	db := sqlc.GetNestedDBTX(ds)
	err = db.Exec(func(dbtx sqlc.DBTX) error {
		params, err = GetHtmlLinksetParamsFromJSON(queryJSON)
		if err != nil {
			slog.Error(err.Error())
			params, _ = GetHtmlLinksetParamsFromJSON("{}")
		}
		msg, err = storage.UpsertLinkSet(ctx, ds, linkSetAction{
			Action:  shared.ActionType(r.Form.Get("action")),
			LinkIds: linkIds,
		})
		switch {
		case err == nil:
			var linksHTML safehtml.HTML
			linksHTML, status, err = a.Views.GetLinkSetHTML(ctx, r.Host, r.RequestURI, params)
			if err != nil {
				if status == 0 {
					status = http.StatusInternalServerError
				}
				a.sendError(w, r, status, err.Error())
			}
			linkURLs, err := ds.Queries(dbtx).GetLinkURLs(ctx, linkIds)
			if err != nil {
				slog.Error("Failed to get link URLs for %v", linkIds)
			}
			if len(linkURLs) > 3 {
				linkURLs = linkURLs[:4]
				linkURLs[3] = "..."
			}
			alertHTML, _, _ := a.Views.GetOOBAlertHTML(ctx, ui.SuccessAlert, ui.Message{
				Text:  msg,
				Items: linkURLs,
			})
			a.sendHTML(w, safehtml.HTMLConcat(linksHTML, alertHTML))
			goto end
		case errors.Is(err, ErrFailedToUnmarshal):
			a.sendError(w, r, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrFailedUpsertLinks):
			// TODO: Break out errors into different status for different reasons
			fallthrough
		default:
			a.sendError(w, r, http.StatusInternalServerError, err.Error())
		}
	end:
		return err
	})
}

var _ storage.LinkSetActionGetter = (*linkSetAction)(nil)

type linkSetAction struct {
	Action  shared.ActionType
	LinkIds []int64
}

func (l linkSetAction) GetAction() shared.ActionType {
	return l.Action
}

func (a linkSetAction) GetLinkIds() (ids []int64, err error) {
	return a.LinkIds, err
}
