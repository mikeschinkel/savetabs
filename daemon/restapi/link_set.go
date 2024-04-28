package restapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/storage"
	"savetabs/ui"
)

func (a *API) PostHtmlLinkset(w http.ResponseWriter, r *http.Request) {
	var msg string
	var status int

	ctx := context.TODO()

	err := r.ParseForm()
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	linkIds, ok := r.Form["link_id"]
	if !ok {
		a.sendError(w, r, http.StatusBadRequest, ErrNoLinkIdsSubmitted.Error())
		return
	}
	msg, err = storage.UpsertLinkSet(ctx, linkSetAction{
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
		alertHTML, _, _ := a.Views.GetOOBAlertHTML(ctx, ui.SuccessAlert, ui.Message{
			Text:  msg,
			Items: a.urlsForMsg(ctx, linkIds),
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
}

var _ storage.LinkSetActionGetter = (*linkSetAction)(nil)

type linkSetAction struct {
	Action  shared.ActionType
	LinkIds []string
}

func (l linkSetAction) GetAction() shared.ActionType {
	return l.Action
}

func (a linkSetAction) GetLinkIds() (ids []int64, err error) {
	var linkId int64

	ids = make([]int64, len(a.LinkIds))
	for i, id := range a.LinkIds {
		linkId, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			goto end
		}
		ids[i] = linkId
	}
end:
	return ids, err
}
